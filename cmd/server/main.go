package main

import (
	"context"
	"flag"
	"os"
	"path"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/joho/godotenv"

	"auth/internal/clients"
	"auth/internal/conf"
	pkgConfig "auth/internal/pkg/config"
	"auth/internal/pkg/logger"
	"auth/internal/pkg/metrics"
	"auth/internal/pkg/runtime"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name = `auth-server`
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

func initSentry(dsn, env string) (flush func(), err error) {
	err = sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		Environment:      env,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		return nil, err
	}
	return func() {
		sentry.Flush(2 * time.Second)
	}, nil
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	flag.Parse()

	var err error

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	if bc.Sentry.Dsn != "" {
		flush, errSentry := initSentry(bc.Sentry.Dsn, bc.Env)
		if errSentry != nil {
			return errSentry
		}
		defer flush()
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
	if err != nil {
		return err
	}

	app, err := wireApp(ctx, database, bc.Auth, bc.Server, notificationsClient, metric, logs)
	if err != nil {
		return err
	}

	// start and wait for stop signal
	if err = app.Run(); err != nil {
		return err
	}

	logHelper.Info("app terminates")

	return nil
}
