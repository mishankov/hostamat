package httpserver

import (
	"net/http"
)

type HandleGroup struct {
	mux         *http.ServeMux
	middlewares []Middleware
}

func NewHandleGroup() *HandleGroup {
	return &HandleGroup{mux: http.NewServeMux()}
}

func (hg *HandleGroup) Use(middleware Middleware) {
	hg.middlewares = append(hg.middlewares, middleware)
}

func (hg *HandleGroup) UseFunc(middlewareFunc func(http.Handler) http.Handler) {
	hg.middlewares = append(hg.middlewares, MiddlewareFunc(middlewareFunc))
}

func (hg *HandleGroup) Handle(pattern string, handler http.Handler) {
	hg.mux.Handle(pattern, wrapHandlerInMiddleware(handler, hg.middlewares))
}

func (hg *HandleGroup) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	hg.mux.Handle(pattern, wrapHandlerInMiddleware(http.HandlerFunc(handler), hg.middlewares))
}

func (hg *HandleGroup) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hg.mux.ServeHTTP(w, r)
}
