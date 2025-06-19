package httpserver

import (
	"hostamat/pkg/platforma/httpserver/middleware"
	"net/http"
	"slices"
)

type h struct {
	pattern string
	handler http.Handler
}

type ServicePackage struct {
	middlewares []middleware.Middleware
	handlers    []h
}

func New() *ServicePackage {
	return &ServicePackage{}
}

func (sp *ServicePackage) Use(middleware middleware.Middleware) {
	sp.middlewares = append(sp.middlewares, middleware)
}

func (sp *ServicePackage) UseFunc(middlewareFunc func(http.Handler) http.Handler) {
	sp.middlewares = append(sp.middlewares, middleware.MiddlewareFunc(middlewareFunc))
}

func (sp *ServicePackage) Handle(pattern string, handler http.Handler) {
	sp.handlers = append(sp.handlers, h{pattern, handler})
}

func (sp *ServicePackage) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	sp.handlers = append(sp.handlers, h{pattern, http.HandlerFunc(handler)})
}

func (sp *ServicePackage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server := http.NewServeMux()
	for _, handler := range sp.handlers {
		finalHandler := handler.handler
		for _, middleware := range slices.Backward(sp.middlewares) {
			finalHandler = middleware.Wrap(finalHandler)
		}

		server.Handle(handler.pattern, finalHandler)
	}

	server.ServeHTTP(w, r)
}
