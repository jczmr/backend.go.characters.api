package logger

import (
	"log/slog"
	"os"
)

func NewSlogLogger() *slog.Logger {
	//TODO could be read from .env to better handle when DEBUG like logs are required
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo, // Set your desired log level
	}))
}
