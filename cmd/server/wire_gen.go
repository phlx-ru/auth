// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"auth/internal/biz"
	"auth/internal/clients"
	"auth/internal/conf"
	"auth/internal/data"
	"auth/internal/pkg/metrics"
	"auth/internal/server"
	"auth/internal/service"
	"context"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
)

// Injectors from wire.go:

// wireData init database
func wireData(confData *conf.Data, logger log.Logger) (data.Database, func(), error) {
	database, cleanup, err := data.NewData(confData, logger)
	if err != nil {
		return nil, nil, err
	}
	return database, func() {
		cleanup()
	}, nil
}

// wireApp init kratos application.
func wireApp(contextContext context.Context, database data.Database, auth *conf.Auth, confServer *conf.Server, notifications clients.Notifications, metricsMetrics metrics.Metrics, logger log.Logger) (*kratos.App, error) {
	userRepo := data.NewUserRepo(database, logger, metricsMetrics)
	sessionRepo := data.NewSessionRepo(database, logger, metricsMetrics)
	codeRepo := data.NewCodeRepo(database, logger, metricsMetrics)
	historyRepo := data.NewHistoryRepo(database, logger, metricsMetrics)
	authUsecase := biz.NewAuthUsecase(userRepo, sessionRepo, codeRepo, historyRepo, notifications, metricsMetrics, logger)
	authService := service.NewAuthService(authUsecase, metricsMetrics, logger)
	userUsecase := biz.NewUserUsecase(userRepo, metricsMetrics, logger)
	userService := service.NewUserService(userUsecase, metricsMetrics, logger)
	grpcServer := server.NewGRPCServer(confServer, authService, userService)
	httpServer := server.NewHTTPServer(confServer, auth, authService, userService, metricsMetrics, logger)
	app := newApp(contextContext, logger, grpcServer, httpServer)
	return app, nil
}