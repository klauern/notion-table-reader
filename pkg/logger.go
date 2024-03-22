package pkg

import (
	"log/slog"
	"os"
	"strings"
)

var (
	Log      *slog.Logger
	LogLevel slog.Level = slog.LevelInfo
)

func init() {
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
	Log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     LogLevel,
	}))
}
