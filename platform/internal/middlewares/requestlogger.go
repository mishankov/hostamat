package middlewares

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
)

type logger interface {
	DebugContext(ctx context.Context, msg string, args ...any)
}

type responseWrapper struct {
	http.ResponseWriter
	statusCode int
	headers    http.Header
	body       []byte
}

func (rw *responseWrapper) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWrapper) Write(body []byte) (int, error) {
	rw.body = body
	return rw.ResponseWriter.Write(body)
}

type RequestLogger struct {
	logger logger
}

func NewRequestLogger(logger logger) *RequestLogger {
	return &RequestLogger{logger: logger}
}

// TODO: wrap ResponseWriter to log response code and headers
func (m *RequestLogger) Wrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ww := &responseWrapper{ResponseWriter: w}

		headerAttrs := make([]slog.Attr, 0, len(r.Header))
		for k, v := range r.Header {
			headerAttrs = append(headerAttrs, slog.String(k, strings.Join(v, ", ")))
		}
		m.logger.DebugContext(
			ctx, "handler invoke started",
			"path", r.URL.Path, "method", r.Method, "remote", r.RemoteAddr, "headers",
			slog.Attr{Key: "headers", Value: slog.GroupValue(headerAttrs...)},
		)

		h.ServeHTTP(ww, r)

		m.logger.DebugContext(
			ctx, "handler invoke finished",
			"status", ww.statusCode, "headers", ww.headers, "body", string(ww.body),
		)
	})
}
