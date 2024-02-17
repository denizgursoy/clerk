package main

import (
	"context"

	"github.com/denizgursoy/clerk/pkg/v1/config"
	"github.com/denizgursoy/clerk/pkg/v1/server"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func StartGRPCServer(lc fx.Lifecycle, sd fx.Shutdowner, c config.Config, srv *grpc.Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := server.Start(c, srv); err != nil {
					if err := sd.Shutdown(); err != nil {
						log.Info().Msg("could not stop the app")
					}
				}

			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			go server.Stop(srv)

			return nil
		},
	})
}
