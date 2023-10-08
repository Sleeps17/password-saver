package logger

import (
	"io"
	"log/slog"
)

const (
	LocalEnv = "local"
	DevEnv   = "dev"
	ProdEnv  = "prod"
)

func Setup(env string, w io.Writer) *slog.Logger {

	var logger *slog.Logger

	switch env {
	case LocalEnv:
		logger = slog.New(
			slog.NewTextHandler(w, &slog.HandlerOptions{AddSource: true, Level: slog.LevelInfo}),
		)
	case DevEnv:
		logger = slog.New(
			slog.NewJSONHandler(w, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case ProdEnv:
		logger = slog.New(
			slog.NewJSONHandler(w, &slog.HandlerOptions{Level: slog.LevelError}),
		)
	default:
		logger = slog.New(
			slog.NewTextHandler(w, &slog.HandlerOptions{AddSource: true, Level: slog.LevelInfo}),
		)
	}

	return logger
}
