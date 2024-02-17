package usecases

import "github.com/rs/zerolog/log"

func StartApplication() error {
	log.Info().Msg("starting the server")

	return nil
}
