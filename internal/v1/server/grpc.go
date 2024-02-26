package server

import (
	"fmt"
	"net"

	"github.com/denizgursoy/clerk/internal/v1/config"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func NewGRPCServer() *grpc.Server {
	return grpc.NewServer()
}

func Start(cfg config.Config, srv *grpc.Server) error {
	// add a listener address
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Error().Err(err).Int("port", cfg.Port).Msg("could not listen port")

		return fmt.Errorf("could not bind tcp: %w", err)
	}
	if err = srv.Serve(lis); err != nil {
		log.Error().Err(err).Int("port", cfg.Port).Msg("could not start server")

		return fmt.Errorf("could not start grpc server: %w", err)
	}

	return nil
}

func Stop(srv *grpc.Server) {
	srv.GracefulStop()
}
