//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"context"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	"auth/internal/biz"
	"auth/internal/conf"
	"auth/internal/data"
	"auth/internal/pkg/metrics"
	"auth/internal/server"
	"auth/internal/service"
)

// wireData init database
func wireData(*conf.Data, log.Logger) (data.Database, func(), error) {
	panic(wire.Build(data.ProviderDataSet))
}

// wireApp init kratos application.
func wireApp(context.Context, data.Database, *conf.Server, metrics.Metrics, log.Logger) (
	*kratos.App,
	error,
) {
	panic(wire.Build(server.ProviderSet, data.ProviderRepoSet, biz.ProviderSet, service.ProviderSet, newApp))
}
