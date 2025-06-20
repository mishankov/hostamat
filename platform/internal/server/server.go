package server

import (
	"context"
	"hostamat/pkg/platforma/httpserver/handlegroup"
	"net/http"
)

type logger interface {
	Info(msg string, args ...any)
	DebugContext(ctx context.Context, msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
}

type Server struct {
	port              string
	uploadFileHandler http.Handler
	logger            logger
}

func New(port string, logger logger, uploadFileHandler http.Handler) *Server {
	return &Server{port: port, logger: logger, uploadFileHandler: uploadFileHandler}
}

func (s *Server) Run() {
	server := handlegroup.New()

	server.Use(&RequestLogger{logger: s.logger})

	server.Handle("/uploadFile", s.uploadFileHandler)

	server.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	s.logger.Info("Starting server at: http://localhost:" + s.port)
	if err := http.ListenAndServe(":"+s.port, server); err != nil {
		s.logger.Info("server stopped", "error", err)
	}
}
