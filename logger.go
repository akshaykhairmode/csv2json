package main

import (
	"os"

	"github.com/rs/zerolog"
)

func GetLogger(verbose bool) *zerolog.Logger {

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()

	if verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	return &logger

}
