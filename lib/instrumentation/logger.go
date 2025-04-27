package instrumentation

import (
	"log/slog"
	"os"
)

type Logger = slog.Logger

const (
	LOG_LEVEL_DEBUG string = "debug"
	LOG_LEVEL_INFO  string = "info"
	LOG_LEVEL_WARN  string = "warn"
	LOG_LEVEL_ERROR string = "error"
)

func New(lvl string) *Logger {
	opts := &slog.HandlerOptions{
		Level: parseLevel(lvl),
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))

	slog.SetDefault(logger)

	return logger
}

func parseLevel(lvl string) slog.Level {
	switch lvl {
	case LOG_LEVEL_DEBUG:
		return slog.LevelDebug
	case LOG_LEVEL_INFO:
		return slog.LevelInfo
	case LOG_LEVEL_WARN:
		return slog.LevelWarn
	case LOG_LEVEL_ERROR:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
