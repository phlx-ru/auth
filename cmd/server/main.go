package main

import (
	"context"
	"flag"
	"os"
	"path"

	"auth/internal/clients"
	"auth/internal/conf"
	pkgConfig "auth/internal/pkg/config"
	"auth/internal/pkg/logger"
	"auth/internal/pkg/metrics"
	"auth/internal/pkg/runtime"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/joho/godotenv"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name = `auth_server`
	// Version is the version of the compiled software.
	Version = `1.1.1`
	// flagconf is the config flag.
	flagconf string
	// dotenv is loaded from config path .env file
	dotenv string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "./configs", "config path, eg: -conf config.yaml")
	flag.StringVar(&dotenv, "dotenv", ".env", ".env file, eg: -dotenv .env")
}

func newApp(ctx context.Context, logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Context(ctx),
		kratos.Server(
			gs,
			hs,
		),
	)
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	flag.Parse()

	var err error

	ctx := context.Background()

	envPath := path.Join(flagconf, dotenv)
	err = godotenv.Overload(envPath)
	if err != nil {
		return err
	}

	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
		config.WithDecoder(pkgConfig.EnvReplaceDecoder),
	)
	defer func() {
		_ = c.Close()
	}()

	if err = c.Load(); err != nil {
		return err
	}

	var bc conf.Bootstrap
	if err = c.Scan(&bc); err != nil {
		return err
	}

	logs := logger.New(id, Name, Version, bc.Log.Level)
	logHelper := logger.NewHelper(logs, "scope", "server")

	metric, err := metrics.New(bc.Metrics.Address, Name, bc.Metrics.Mute)
	if err != nil {
		return err
	}
	defer metric.Close()
	metric.Increment("starts.count")

	database, cleanup, err := wireData(bc.Data, logs)
	if err != nil {
		return err
	}
	defer cleanup()

	go database.CollectDatabaseMetrics(ctx, metric)
	go runtime.CollectGoMetrics(ctx, metric)

	if err = database.Prepare(ctx, bc.Data.Database.Migrate); err != nil {
		return err
	}

	n := bc.Client.Grpc.Notifications
	notificationsClient, err := clients.NewNotifications(ctx, n.Endpoint, n.Timeout.AsDuration(), metric, logs)

	app, err := wireApp(ctx, database, bc.Auth, bc.Server, notificationsClient, metric, logs)
	if err != nil {
		panic(err)
	}

	// start and wait for stop signal
	if err = app.Run(); err != nil {
		return err
	}

	logHelper.Info("app terminates")

	return nil
}
