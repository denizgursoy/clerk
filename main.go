package main

import (
	"github.com/denizgursoy/clerk/pkg/v1/config"
	"github.com/denizgursoy/clerk/pkg/v1/handlers"
	"github.com/denizgursoy/clerk/pkg/v1/repository"
	"github.com/denizgursoy/clerk/pkg/v1/server"
	"github.com/denizgursoy/clerk/pkg/v1/usecases"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(
			server.NewGRPCServer,
			config.CreateConfig,
			fx.Annotate(usecases.NewMemberUserCase, fx.As(new(usecases.MemberUseCase))),
			fx.Annotate(repository.NewMemberETCDRepository, fx.As(new(usecases.MemberRepository))),
		),
		fx.Invoke(
			handlers.NewMemberGRPCHandler,
			StartGRPCServer,
			StartBalance,
		),
	).Run()
}
