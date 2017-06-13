package op

import (
	"testing"

	"github.com/rs/zerolog"
)

func Test_Console(t *testing.T) {
	log := zerolog.New(NewConsole())

	log.Debug().Msg("debug")
	log.Info().Msg("info")
	log.Warn().Msg("warn")
	log.Error().Msg("error")
	// log.Fatal().Msg("Fatal")
	// log.Panic().Msg("Panic")
}
