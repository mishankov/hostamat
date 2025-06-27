package httpserver

import (
	"net/http"
)

// HandleGroup represents a group of HTTP handlers that share common middlewares.
type HandleGroup struct {
	mux         *http.ServeMux
	middlewares []Middleware
}

// NewHandleGroup creates a new HandleGroup with an initialized http.ServeMux.
func NewHandleGroup() *HandleGroup {
	return &HandleGroup{mux: http.NewServeMux()}
}

// Use adds a middleware to the HandleGroup's middleware chain.
func (hg *HandleGroup) Use(middleware Middleware) {
	hg.middlewares = append(hg.middlewares, middleware)
}

// UseFunc adds a function as a middleware to the HandleGroup's middleware chain.
func (hg *HandleGroup) UseFunc(middlewareFunc func(http.Handler) http.Handler) {
	hg.middlewares = append(hg.middlewares, MiddlewareFunc(middlewareFunc))
}

// Handle registers an http.Handler for the given pattern, applying all
// middlewares in the HandleGroup's chain.
func (hg *HandleGroup) Handle(pattern string, handler http.Handler) {
	hg.mux.Handle(pattern, wrapHandlerInMiddleware(handler, hg.middlewares))
}

// HandleFunc registers an http.HandlerFunc for the given pattern, applying all
// middlewares in the HandleGroup's chain.
func (hg *HandleGroup) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	hg.mux.Handle(pattern, wrapHandlerInMiddleware(http.HandlerFunc(handler), hg.middlewares))
}

// ServeHTTP implements the http.Handler interface, allowing HandleGroup to
// be used as an HTTP handler itself.
func (hg *HandleGroup) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hg.mux.ServeHTTP(w, r)
}
