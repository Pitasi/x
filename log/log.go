package log

import (
	"log/slog"
	"os"
)

var std *slog.Logger

func init() {
	std = defaultLogger()
	slog.SetDefault(std)
}

func defaultLogger() *slog.Logger {
	return slog.New(getHandler())
}

func getHandler() slog.Handler {
	format := os.Getenv("LOG_FORMAT")
	addSourceConfig := os.Getenv("LOG_ADD_SOURCE")
	levelConfig := os.Getenv("LOG_LEVEL")

	addSource := addSourceConfig != "false" && addSourceConfig != "FALSE"

	var level slog.Leveler
	switch levelConfig {
	case "error", "ERROR", "err", "ERR":
		level = slog.LevelError
	case "warning", "WARNING", "warn", "WARN":
		level = slog.LevelWarn
	case "debug", "DEBUG", "dbg", "DBG":
		level = slog.LevelDebug
	case "info", "INFO":
		level = slog.LevelInfo
	default:
		level = slog.LevelInfo
	}

	switch format {
	case "JSON", "json":
		return slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level:     level,
			AddSource: addSource,
		})
	default:
		return slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level:     level,
			AddSource: addSource,
		})
	}
}

// Module returns a [slog.Logger] with a "module" attribute set to the given value.
// If the value is empty, the standard logger will be returned.
func Module(name string) *slog.Logger {
	if name == "" {
		return std
	}
	return std.With(slog.String("module", name))
}
