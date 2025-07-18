package file

import (
	"context"
	"io"
	"net/http"
	"strings"
)

type serveService interface {
	GetFile(ctx context.Context, fileName string, domain string) (io.Reader, error)
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

	filePath := r.URL.Path
	if filePath == "/" {
		filePath = "/index.html"
	}

	file, err := h.service.GetFile(ctx, filePath, r.Host)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", h.getContentType(filePath))
	io.Copy(w, file)
}

func (h *ServeHandler) getContentType(filename string) string {
	switch {
	case strings.HasSuffix(filename, ".html"):
		return "text/html"
	case strings.HasSuffix(filename, ".css"):
		return "text/css"
	case strings.HasSuffix(filename, ".js"):
		return "application/javascript"
	case strings.HasSuffix(filename, ".png"):
		return "image/png"
	case strings.HasSuffix(filename, ".jpg"), strings.HasSuffix(filename, ".jpeg"):
		return "image/jpeg"
	case strings.HasSuffix(filename, ".svg"):
		return "image/svg+xml"
	default:
		return "application/octet-stream"
	}
}
