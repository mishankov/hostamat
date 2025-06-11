package logger

import (
	"context"
	"log/slog"
	"os"
)

type contextHandler struct {
	slog.Handler
}

// Handle overrides the default Handle method to add context values.

func (h *contextHandler) Handle(ctx context.Context, r slog.Record) error {
	if traceId, ok := ctx.Value("traceId").(string); ok {
		r.AddAttrs(slog.String("traceId", traceId))
	}

	return h.Handler.Handle(ctx, r)
}

func New(loggerType string, level slog.Level) *slog.Logger {
	if loggerType == "json" {
		return slog.New(&contextHandler{slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})})
	}

	return slog.New(&contextHandler{slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})})
}
