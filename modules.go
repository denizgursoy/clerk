package main

import (
	"context"

	"github.com/denizgursoy/clerk/internal/config"
	"github.com/denizgursoy/clerk/internal/server"
	"go.uber.org/fx"
)

func StartHTTPServer(lc fx.Lifecycle, s server.Server, c config.Config) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return server.Start(s, c)
		},
	})
}
