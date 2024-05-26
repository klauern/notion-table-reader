package pkg

import (
	"log/slog"
	"os"
	"strings"
)

var LogLevel slog.Level = slog.LevelInfo

func SetupLogging() {
	switch level := os.Getenv("LOG_LEVEL"); strings.ToLower(level) {
	case "debug":
		LogLevel = slog.LevelDebug
	case "warn":
		LogLevel = slog.LevelWarn
	case "error", "fatal":
		LogLevel = slog.LevelError
	default:
		LogLevel = slog.LevelInfo
	}
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     LogLevel,
	}))
	slog.SetDefault(log)
}
