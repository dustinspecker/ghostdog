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
	testConfig := NewTest()

	if testConfig.Config.Fs.Name() != "MemMapFS" {
		t.Errorf("expected fs to be MemMapFS, but got %s", testConfig.Config.Fs.Name())
	}

	if testConfig.Config.LogCtx.Level != log.DebugLevel {
		t.Errorf("expected log level to be debug, but got %s", testConfig.Config.LogCtx.Level)
	}
}

func TestNewTestLogHandler(t *testing.T) {
	testConfig := NewTest()

	testConfig.Config.LogCtx.Debug("test")

	if len(testConfig.LogHandler.Entries) != 1 {
		t.Fatalf("expected to have 1 log entry: %v", testConfig.LogHandler.Entries)
	}

	entry := testConfig.LogHandler.Entries[0]
	if entry.Message != "test" {
		t.Errorf("test message should be saved in entries, but got %s", entry.Message)
	}
}
