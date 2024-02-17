package main

import (
	"go.uber.org/fx"
	"partitioner/internal/config"
	"partitioner/internal/server"
	"partitioner/internal/usecases"
)

func main() {
	fx.New(
		fx.Provide(
			server.GetHTTPServer,
			config.CreateConfig,
		),
		fx.Invoke(
			usecases.StartApplication,
			StartHTTPServer,
		),
	).Run()
}
