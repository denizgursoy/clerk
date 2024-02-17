package server

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"partitioner/internal/config"
)

func GetHTTPServer() *echo.Echo {
	return echo.New()
}

func Start(e *echo.Echo, cfg config.Config) error {
	if err := e.Start(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		log.Err(err).Msgf("could not start the HTTP server on port %d", cfg.Port)

		return err
	}

	return nil
}

func Stop(e *echo.Echo) error {
	if err := e.Close(); err != nil {
		log.Err(err).Msg("could not stop the HTTP server")

		return err
	}

	return nil
}
