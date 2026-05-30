package logger

import (
	"log/slog"
	"os"
)

const (
	LocalEnv = "local"
	DevEnv   = "development"
	ProdEnv  = "production"
)

func SetUpLogger(env string) *slog.Logger {

	var logger *slog.Logger
	switch env {
	case LocalEnv:
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case DevEnv:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case ProdEnv:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return logger
}
