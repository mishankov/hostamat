package httpserver

import (
	"net/http"
	"slices"
)

type Middleware interface {
	Wrap(http.Handler) http.Handler
}

type MiddlewareFunc func(http.Handler) http.Handler

func (f MiddlewareFunc) Wrap(h http.Handler) http.Handler {
	return f(h)
}

func wrapHandlerInMiddleware(handler http.Handler, middlewares []Middleware) http.Handler {
	finalHandler := handler
	for _, middleware := range slices.Backward(middlewares) {
		finalHandler = middleware.Wrap(finalHandler)
	}

	return finalHandler
}
