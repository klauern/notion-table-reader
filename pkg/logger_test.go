package pkg

import (
	"log/slog"
	"os"
	"testing"
)

func TestSetupLogging(t *testing.T) {
	// Set up test environment
	oldLogLevel := os.Getenv("LOG_LEVEL")
	os.Setenv("LOG_LEVEL", "debug")
	defer os.Setenv("LOG_LEVEL", oldLogLevel)

	// Call the function to be tested
	SetupLogging()

	// Verify the expected log level
	expectedLogLevel := slog.LevelDebug
	if LogLevel != expectedLogLevel {
		t.Errorf("Expected log level %v, but got %v", expectedLogLevel, LogLevel)
	}

	// Verify the default logger
	defaultLogger := slog.Default()
	if defaultLogger == nil {
		t.Error("Default logger is nil")
	}
}
