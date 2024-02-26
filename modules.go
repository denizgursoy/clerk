package main

import (
	"context"

	"github.com/denizgursoy/clerk/pkg/v1/config"
	"github.com/denizgursoy/clerk/pkg/v1/repository"
	"github.com/denizgursoy/clerk/pkg/v1/server"
	"github.com/denizgursoy/clerk/pkg/v1/usecases"
	"github.com/rs/zerolog/log"
	"go.etcd.io/etcd/client/v3"
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

func StartBalance(lc fx.Lifecycle, m usecases.MemberUseCase) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				m.TriggerBalance()
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			m.StopBalance()

			return nil
		},
	})
}

func StopETCDClientOnShutDown(lc fx.Lifecycle, c *clientv3.Client) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return repository.Stop(c)
		},
	})
}
