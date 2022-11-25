package server

import (
	v1 "auth/api/auth/v1"
	"auth/internal/conf"
	"auth/internal/middlewares"
	internalJWT "auth/internal/pkg/jwt"
	"auth/internal/pkg/metrics"
	"auth/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(
	c *conf.Server,
	a *conf.Auth,
	auth *service.AuthService,
	user *service.UserService,
	metric metrics.Metrics,
	logger log.Logger,
) *http.Server {
	var opts = []http.ServerOption{
		http.Timeout(c.Http.Timeout.AsDuration()),
		http.Middleware(
			middlewares.Duration(metric, logger),
			recovery.Recovery(),
			jwt.Server(internalJWT.Check(a.Jwt.Secret)),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterAuthHTTPServer(srv, auth)
	v1.RegisterUserHTTPServer(srv, user)
	return srv
}
