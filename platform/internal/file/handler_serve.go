package file

import (
	"context"
	"io"
	"net/http"
	"strings"
)

type serviceServe interface {
	GetFile(ctx context.Context, fileName string) (io.Reader, error)
	FileExists(ctx context.Context, fileName string) (bool, error)
}

type ServeHandler struct {
	service serviceServe
}

func NewServeHandler(service serviceServe) *ServeHandler {
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
	exists, err := h.service.FileExists(ctx, fileName)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	file, err := h.service.GetFile(ctx, fileName)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", contentTypeFromFileName(filePath))
	io.Copy(w, file)
}
