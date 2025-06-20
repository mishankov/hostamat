package handlegroup

import (
	"hostamat/pkg/platforma/httpserver/middleware"
	"net/http"
	"slices"
)

type HandleGroup struct {
	mux         *http.ServeMux
	middlewares []middleware.Middleware
}

func New() *HandleGroup {
	return &HandleGroup{mux: http.NewServeMux()}
}

func (hg *HandleGroup) Use(middleware middleware.Middleware) {
	hg.middlewares = append(hg.middlewares, middleware)
}

func (hg *HandleGroup) UseFunc(middlewareFunc func(http.Handler) http.Handler) {
	hg.middlewares = append(hg.middlewares, middleware.MiddlewareFunc(middlewareFunc))
}

func (hg *HandleGroup) Handle(pattern string, handler http.Handler) {
	hg.mux.Handle(pattern, hg.wrapHandlerInMiddleware(handler))
}

func (hg *HandleGroup) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	hg.mux.Handle(pattern, hg.wrapHandlerInMiddleware(http.HandlerFunc(handler)))
}

func (hg *HandleGroup) wrapHandlerInMiddleware(handler http.Handler) http.Handler {
	finalHandler := handler
	for _, middleware := range slices.Backward(hg.middlewares) {
		finalHandler = middleware.Wrap(finalHandler)
	}

	return finalHandler
}

func (hg *HandleGroup) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hg.mux.ServeHTTP(w, r)
}
