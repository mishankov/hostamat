package log

import (
	"context"
	"log/slog"
)

var Logger = slog.Default()

func SetDefault(l *slog.Logger) {
	Logger = l
}

// TODO: more methods
func InfoContext(ctx context.Context, msg string, args ...any) {
	Logger.InfoContext(ctx, msg, args...)
}

func ErrorContext(ctx context.Context, msg string, args ...any) {
	Logger.ErrorContext(ctx, msg, args...)
}
