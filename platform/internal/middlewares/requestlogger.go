package middlewares

import (
	"context"
	"hostamat/pkg/platforma/log"
	"net/http"

	"github.com/google/uuid"
)

type logger interface {
	DebugContext(ctx context.Context, msg string, args ...any)
}

type RequestLogger struct {
	logger logger
}

func NewRequestLogger(logger logger) *RequestLogger {
	return &RequestLogger{logger: logger}
}

func (m *RequestLogger) Wrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), log.TraceIdKey, uuid.NewString())
		r = r.WithContext(ctx)

		m.logger.DebugContext(ctx, "service invoke started", "path", r.URL.Path, "method", r.Method, "remote", r.RemoteAddr, "headers", r.Header)
		h.ServeHTTP(w, r)
		m.logger.DebugContext(ctx, "service invoke finished")
	})
}
