package main

import (
	"github.com/denizgursoy/clerk/internal/config"
	"github.com/denizgursoy/clerk/internal/server"
	"github.com/denizgursoy/clerk/internal/usecases"
	"go.uber.org/fx"
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
