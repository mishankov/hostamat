package file

import (
	"context"
	"io"
	"net/http"
	"strings"
)

type serveService interface {
	GetFile(ctx context.Context, fileName string) (io.Reader, error)
	GetContentType(filePath string) string
}

type ServeHandler struct {
	service serveService
}

func NewServeHandler(service serveService) *ServeHandler {
	return &ServeHandler{
		service: service,
	}
}

func (h *ServeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	subdomain := strings.Split(r.Host, ".")[0]
	filePath := r.URL.Path
	if filePath == "/" {
		filePath = "/index.html"
	}

	fileName := subdomain + filePath
	file, err := h.service.GetFile(ctx, fileName)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", h.service.GetContentType(filePath))
	io.Copy(w, file)
}
