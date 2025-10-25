package site

import (
	"io"
	"net/http"
	"strings"
)

type ServeFileHandler struct {
	service *Service
}

func NewServeFileHandler(service *Service) *ServeFileHandler {
	return &ServeFileHandler{
		service: service,
	}
}

func (h *ServeFileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	filePath := r.URL.Path
	if filePath == "/" {
		filePath = "/index.html"
	}

	file, err := h.service.GetFile(ctx, r.Host, filePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", h.getContentType(filePath))

	io.Copy(w, file)
}

func (h *ServeFileHandler) getContentType(filename string) string {
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
