package middlewares

import (
	"context"
	"net/http"
	"strings"
)

type logger interface {
	DebugContext(ctx context.Context, msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
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

func (m *RequestLogger) Wrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/metrics" {
			ww := &responseWrapper{ResponseWriter: w}

			m.logger.DebugContext(
				r.Context(), "handler invoke started",
				"path", r.URL.Path, "method", r.Method, "remote", r.RemoteAddr, "headers", maskedHeaders(r.Header),
			)

			h.ServeHTTP(ww, r)

			if ww.statusCode >= 500 {
				m.logger.ErrorContext(
					r.Context(), "handler invoke finished with error",
					"path", r.URL.Path, "method", r.Method, "remote", r.RemoteAddr, "request_headers", maskedHeaders(r.Header),
					"status", ww.statusCode, "response_headers", maskedHeaders(ww.headers), "response_body", string(ww.body),
				)
			} else {
				m.logger.DebugContext(
					r.Context(), "handler invoke finished",
					"status", ww.statusCode, "headers", maskedHeaders(ww.headers), "body", string(ww.body),
				)
			}

			return
		}

		h.ServeHTTP(w, r)
	})
}

func maskedHeaders(headers http.Header) http.Header {
	newHeaders := http.Header{}
	for hn, hvs := range headers {
		headerName := strings.ToLower(hn)
		if headerName == "authorization" {
			for _, hv := range hvs {
				if strings.HasPrefix(strings.ToLower(hv), "basic ") || strings.HasPrefix(strings.ToLower(hv), "bearer ") {
					newHeaders.Add(hn, strings.Split(hv, " ")[0]+" ***")
					continue
				}
				newHeaders.Add(hn, hv)
			}
		} else if headerName == "cookie" || headerName == "set-cookie" {
			for range hvs {
				newHeaders.Add(hn, "***")
			}
		} else {
			for _, hv := range hvs {
				newHeaders.Add(hn, hv)
			}
		}
	}

	return newHeaders
}
