package log

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

var Logger = slog.Default()

type contextKey string

const (
	TraceIdKey     contextKey = "traceId"
	ServiceNameKey contextKey = "serviceName"
	SartupTaskKey  contextKey = "satrtupTask"
)

// SetDefault sets the default logger used by the package-level logging functions.
func SetDefault(l *slog.Logger) {
	Logger = l
}

type contextHandler struct {
	slog.Handler
	additionKeys []fmt.Stringer
}

// Handle processes the log record by adding context values before passing it to the underlying handler.
func (h *contextHandler) Handle(ctx context.Context, r slog.Record) error {
	if traceId, ok := ctx.Value(TraceIdKey).(string); ok {
		r.AddAttrs(slog.String("traceId", traceId))
	}

	if serviceName, ok := ctx.Value(ServiceNameKey).(string); ok {
		r.AddAttrs(slog.String("serviceName", serviceName))
	}

	if startupTask, ok := ctx.Value(SartupTaskKey).(string); ok {
		r.AddAttrs(slog.String("startupTask", startupTask))
	}

	for _, key := range h.additionKeys {
		if value, ok := ctx.Value(key).(string); ok {
			r.AddAttrs(slog.String(key.String(), value))
		}
	}

	return h.Handler.Handle(ctx, r)
}

// New creates a new slog.Logger with the specified type (json/text), log level, and additional context keys to include.
func New(loggerType string, level slog.Level, contextKeys []fmt.Stringer) *slog.Logger {
	if loggerType == "json" {
		return slog.New(&contextHandler{slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}), contextKeys})
	}

	return slog.New(&contextHandler{slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}), contextKeys})
}

// Debug logs a message at Debug level.
func Debug(msg string, args ...any) {
	Logger.Debug(msg, args...)
}

// DebugContext logs a message at Debug level with context.
func DebugContext(ctx context.Context, msg string, args ...any) {
	Logger.DebugContext(ctx, msg, args...)
}

// Info logs a message at Info level.
func Info(msg string, args ...any) {
	Logger.Info(msg, args...)
}

// InfoContext logs a message at Info level with context.
func InfoContext(ctx context.Context, msg string, args ...any) {
	Logger.InfoContext(ctx, msg, args...)
}

// Warn logs a message at Warn level.
func Warn(msg string, args ...any) {
	Logger.Warn(msg, args...)
}

// WarnContext logs a message at Warn level with context.
func WarnContext(ctx context.Context, msg string, args ...any) {
	Logger.WarnContext(ctx, msg, args...)
}

// Error logs a message at Error level.
func Error(msg string, args ...any) {
	Logger.Error(msg, args...)
}

// ErrorContext logs a message at Error level with context.
func ErrorContext(ctx context.Context, msg string, args ...any) {
	Logger.ErrorContext(ctx, msg, args...)
}
