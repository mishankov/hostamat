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
)

func SetDefault(l *slog.Logger) {
	Logger = l
}

type contextHandler struct {
	slog.Handler
	additionKeys []fmt.Stringer
}

// Handle overrides the default Handle method to add context values.
func (h *contextHandler) Handle(ctx context.Context, r slog.Record) error {
	if traceId, ok := ctx.Value(TraceIdKey).(string); ok {
		r.AddAttrs(slog.String("traceId", traceId))
	}

	if serviceName, ok := ctx.Value(ServiceNameKey).(string); ok {
		r.AddAttrs(slog.String("serviceName", serviceName))
	}

	for _, key := range h.additionKeys {
		if value, ok := ctx.Value(key).(string); ok {
			r.AddAttrs(slog.String(key.String(), value))
		}
	}

	return h.Handler.Handle(ctx, r)
}

func New(loggerType string, level slog.Level, contextKeys []fmt.Stringer) *slog.Logger {
	if loggerType == "json" {
		return slog.New(&contextHandler{slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}), contextKeys})
	}

	return slog.New(&contextHandler{slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}), contextKeys})
}

// TODO: more methods
func InfoContext(ctx context.Context, msg string, args ...any) {
	Logger.InfoContext(ctx, msg, args...)
}

func ErrorContext(ctx context.Context, msg string, args ...any) {
	Logger.ErrorContext(ctx, msg, args...)
}
