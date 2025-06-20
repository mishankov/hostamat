package server

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type RequestLogger struct {
	logger logger
}

func (m *RequestLogger) Wrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "traceId", uuid.NewString())

		m.logger.DebugContext(ctx, "service invoke started", "path", r.URL.Path, "method", r.Method)
		h.ServeHTTP(w, r)
		m.logger.DebugContext(ctx, "service invoke finished")
	})
}
