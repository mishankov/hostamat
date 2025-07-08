package file

import (
	"io"
	"net/http"
	"strings"
)

type ServeHandler struct {
	service *Service
}

func NewServeHandler(service *Service) *ServeHandler {
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
