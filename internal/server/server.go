package server

import (
	"fmt"
	"net"

	"github.com/denizgursoy/clerk/internal/config"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type Server struct {
	UnimplementedConsumerServiceServer
}

func GetHTTPServer() Server {
	return Server{}
}

func Start(s Server, cfg config.Config) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Info().Msgf("failed to listen on port %d: %v", cfg.Port, err)

		return err
	}

	server := grpc.NewServer()
	RegisterConsumerServiceServer(server, &s)
	log.Printf("gRPC Server listening at %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		log.Info().Msgf("failed to serve: %v", err)

		return err
	}

	return nil
}

func Stop(s Server) error {
	return nil
}
