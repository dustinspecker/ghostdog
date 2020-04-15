package config

import (
	"os"
	"strings"

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

type TestConfig struct {
	Config     Config
	LogHandler *memory.Handler
}

func hasFields(entry *log.Entry, fields log.Fields) bool {
	for key, fieldValue := range fields {
		entryValue, ok := entry.Fields[key]
		if !ok {
			return false
		}

		if !strings.Contains(entryValue.(string), fieldValue.(string)) {
			return false
		}
	}

	return true
}

func (testConfig TestConfig) HasLogEntry(logLevel log.Level, fields log.Fields, substringMessage string) bool {
	for _, entry := range testConfig.LogHandler.Entries {
		if entry.Level == logLevel && strings.Contains(entry.Message, substringMessage) && hasFields(entry, fields) {
			return true
		}
	}

	return false
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

func NewTest() TestConfig {
	log.SetLevel(log.DebugLevel)
	logHandler := memory.New()
	log.SetHandler(logHandler)

	logCtx := log.WithFields(log.Fields{
		"app": "ghostdog-test",
	})

	return TestConfig{
		Config: Config{
			Fs:               afero.NewMemMapFs(),
			LogCtx:           logCtx,
			WorkingDirectory: ".",
		},
		LogHandler: logHandler,
	}
}
