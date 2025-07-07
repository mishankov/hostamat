package httpserver

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type logger interface {
	InfoContext(ctx context.Context, msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
}

type HttpServer struct {
	logger          logger
	mux             *http.ServeMux
	port            string
	shutdownTimeout time.Duration
	middlewares     []Middleware
}

func New(logger logger, port string, shutdownTimeout time.Duration) *HttpServer {
	return &HttpServer{logger: logger, mux: http.NewServeMux(), port: port, shutdownTimeout: shutdownTimeout}
}

func (s *HttpServer) Handle(pattern string, handler http.Handler) {
	s.mux.Handle(pattern, handler)
}

func (s *HttpServer) HandleFunc(pattern string, handlerFunc http.HandlerFunc) {
	s.mux.HandleFunc(pattern, handlerFunc)
}

func (s *HttpServer) HandleGroup(pattern string, handler http.Handler) {
	s.mux.Handle(pattern+"/", http.StripPrefix(pattern, handler))
}

func (s *HttpServer) Use(middleware Middleware) {
	s.middlewares = append(s.middlewares, middleware)
}

func (s *HttpServer) UseFunc(middlewareFunc func(http.Handler) http.Handler) {
	s.middlewares = append(s.middlewares, MiddlewareFunc(middlewareFunc))
}

func (s *HttpServer) Run(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":" + s.port,
		Handler: wrapHandlerInMiddleware(s.mux, s.middlewares),
	}

	go func() {
		s.logger.InfoContext(ctx, "starting http server", "address", server.Addr)

		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			s.logger.ErrorContext(ctx, "HTTP server error", "error", err)
		}
		s.logger.InfoContext(ctx, "stopped serving new connections.")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(ctx, s.shutdownTimeout)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		s.logger.ErrorContext(ctx, "HTTP shutdown error", "error", err)
		return err
	}
	s.logger.InfoContext(ctx, "graceful shutdown complete.")

	return nil
}
