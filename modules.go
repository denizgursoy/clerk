package main

import (
	"context"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"partitioner/internal/config"
	"partitioner/internal/server"
)

func StartHTTPServer(lc fx.Lifecycle, e *echo.Echo, c config.Config) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return server.Start(e, c)
		},
		OnStop: func(ctx context.Context) error {
			return server.Stop(e)
		},
	})
}
