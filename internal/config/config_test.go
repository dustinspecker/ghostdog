package config

import (
	"testing"

	"github.com/apex/log"
)

func TestNew(t *testing.T) {
	config, err := New("debug")
	if err != nil {
		t.Fatalf("unexpected error getting config: %s", err)
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

func TestHasLogEntry(t *testing.T) {
	testConfig := NewTest()
	testConfig.Config.LogCtx.WithFields(log.Fields{"test": "yes", "app": "ghostdog"}).Debug("this is a test")

	tests := []struct {
		logLevel         log.Level
		substringMessage string
		fields           log.Fields
		expected         bool
	}{
		{log.DebugLevel, "this is a test", log.Fields{}, true},
		{log.InfoLevel, "this is a test", log.Fields{}, false},
		{log.DebugLevel, "ghostdog", log.Fields{}, false},
		{log.DebugLevel, "is a", log.Fields{}, true},
		{log.DebugLevel, "is a", log.Fields{"test": "yes"}, true},
		{log.DebugLevel, "is a", log.Fields{"test": "yes", "app": "ghostdog"}, true},
		{log.DebugLevel, "is a", log.Fields{"test": "yes", "app": "dog"}, true},
		{log.DebugLevel, "is a", log.Fields{"test": "yes", "app": "dustin"}, false},
		{log.DebugLevel, "is a", log.Fields{"name": "ghostdog"}, false},
	}

	for _, tt := range tests {
		if testConfig.HasLogEntry(tt.logLevel, tt.fields, tt.substringMessage) != tt.expected {
			t.Errorf("expected %v for %s and %s with %v, but got %v", tt.expected, tt.logLevel, tt.substringMessage, tt.fields, testConfig.HasLogEntry(tt.logLevel, tt.fields, tt.substringMessage))
		}
	}
}
