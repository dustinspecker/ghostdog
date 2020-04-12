package config

import (
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/memory"
	"github.com/spf13/afero"
)

type Config struct {
	Fs               afero.Fs
	LogCtx           *log.Entry
	WorkingDirectory string
}

func New(logLevel string) (Config, error) {
	// TODO: how to test working directory
	// TODO: how to test when os.Getwd returns an error
	cwd, err := os.Getwd()
	if err != nil {
		return Config{}, err
	}

	parsedLogLevel, err := log.ParseLevel(logLevel)
	if err != nil {
		return Config{}, err
	}

	log.SetLevel(parsedLogLevel)

	// TODO: how to unit test the handler is CLI
	log.SetHandler(cli.New(os.Stderr))

	// TODO: how to unit test these fields
	logCtx := log.WithFields(log.Fields{
		"app": "ghostdog",
	})

	return Config{
		Fs:               afero.NewOsFs(),
		LogCtx:           logCtx,
		WorkingDirectory: cwd,
	}, nil
}

func NewTest() Config {
	log.SetLevel(log.DebugLevel)
	log.SetHandler(memory.New())

	logCtx := log.WithFields(log.Fields{
		"app": "ghostdog-test",
	})

	return Config{
		Fs:               afero.NewMemMapFs(),
		LogCtx:           logCtx,
		WorkingDirectory: ".",
	}
}
