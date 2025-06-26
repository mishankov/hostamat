package httpserver

import (
	"context"
	"hostamat/pkg/platforma/log"
	"net/http"
)

type HttpServer struct {
	mux         *http.ServeMux
	port        string
	middlewares []Middleware
}

func New(port string, handleGroups map[string]http.Handler) *HttpServer {
	mux := http.NewServeMux()
	for pattern, handleGroup := range handleGroups {
		mux.Handle(pattern+"/", http.StripPrefix(pattern, handleGroup))
	}

	return &HttpServer{mux: mux, port: port}
}

func (s *HttpServer) Use(middleware Middleware) {
	s.middlewares = append(s.middlewares, middleware)
}

func (s *HttpServer) UseFunc(middlewareFunc func(http.Handler) http.Handler) {
	s.middlewares = append(s.middlewares, MiddlewareFunc(middlewareFunc))
}

// TODO: graceful shutdown
func (s *HttpServer) Run(ctx context.Context) error {
	log.InfoContext(ctx, "starting http server", "address", "http://localhost:"+s.port)
	return http.ListenAndServe(":"+s.port, wrapHandlerInMiddleware(s.mux, s.middlewares))
}
