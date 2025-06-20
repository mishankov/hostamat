package httpserver

import (
	"log"
	"net/http"
)

type HttpServer struct {
	mux  *http.ServeMux
	port string
}

func New(port string, handleGroups map[string]http.Handler) *HttpServer {
	mux := http.NewServeMux()
	for pattern, handleGroup := range handleGroups {
		mux.Handle(pattern+"/", http.StripPrefix(pattern, handleGroup))
	}

	return &HttpServer{mux: mux, port: port}
}

func (s *HttpServer) Run() error {
	log.Println("starting http server at http://localhost:" + s.port)
	return http.ListenAndServe(":"+s.port, s.mux)
}
