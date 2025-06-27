package httpserver

import (
	"context"
	"errors"
	"hostamat/pkg/platforma/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type HttpServer struct {
	mux             *http.ServeMux
	port            string
	shutdownTimeout time.Duration
	middlewares     []Middleware
}

func New(port string, shutdownTimeout time.Duration, handleGroups map[string]http.Handler) *HttpServer {
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

func (s *HttpServer) Run(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":" + s.port,
		Handler: s.mux,
	}

	go func() {
		log.InfoContext(ctx, "starting http server", "address", server.Addr)

		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.ErrorContext(ctx, "HTTP server error", "error", err)
		}
		log.InfoContext(ctx, "stopped serving new connections.")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(ctx, s.shutdownTimeout)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.ErrorContext(ctx, "HTTP shutdown error", "error", err)
		return err
	}
	log.InfoContext(ctx, "graceful shutdown complete.")

	return nil
}
