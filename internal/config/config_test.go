package config

import (
	"testing"

	"github.com/apex/log"
)

func TestNew(t *testing.T) {
	config, err := New("debug")
	if err != nil {
		t.Fatalf("unexpected error getting config: %w", err)
	}

	if config.Fs.Name() != "OsFs" {
		t.Errorf("expected fs to be OsFs, but got %s", config.Fs.Name())
	}

	if config.LogCtx.Level != log.DebugLevel {
		t.Errorf("expected log level to be debug, but got %s", config.LogCtx.Level)
	}
}

func TestNewReturnsErrorWhenLogLevelParseFails(t *testing.T) {
	_, err := New("dustin")
	if err == nil {
		t.Fatal("expected an error trying to parse a log level of dustin")
	}
}

func TestNewTest(t *testing.T) {
	config := NewTest()

	if config.Fs.Name() != "MemMapFS" {
		t.Errorf("expected fs to be MemMapFS, but got %s", config.Fs.Name())
	}

	if config.LogCtx.Level != log.DebugLevel {
		t.Errorf("expected log level to be debug, but got %s", config.LogCtx.Level)
	}
}
