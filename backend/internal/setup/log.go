package setup

import (
	"log/slog"
	"os"
)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

func NewLogger() Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}
