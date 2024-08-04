package logger

import (
	"os"

	"github.com/rs/zerolog"
)

var L zerolog.Logger

type Opts struct {
	Level      zerolog.Level
	TimeFormat string
}

func InitLogger(opts Opts) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(opts.Level)

	L = zerolog.New(os.Stdout)

	L.Info().Msgf("Logger initialized with level %d", opts.Level)
}
