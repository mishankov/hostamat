package site

import (
	_ "embed"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/platforma-dev/platforma/log"
)

//go:embed 404.html
var notFoundHTML []byte

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
	ev := log.EventFromContext(ctx)

	ev.AddAttrs(map[string]any{
		"site.domain": r.Host,
		"file.path":   r.URL.Path,
	})

	filePath := r.URL.Path

	file, resolvedPath, err := h.service.GetFile(ctx, r.Host, filePath)
	if err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		_, err = w.Write(notFoundHTML)
		if err != nil {
			ev.AddError(err)
		}
		return
	}

	// Close the file if it implements io.Closer
	if closer, ok := file.(io.Closer); ok {
		defer closer.Close()
	}

	w.Header().Set("Content-Type", h.getContentType(resolvedPath))

	ev.AddAttrs(map[string]any{
		"file.resolved_path": resolvedPath,
	})

	if _, err := io.Copy(w, file); err != nil {

		ev.AddError(err)

		log.ErrorContext(ctx, "failed to copy file to response", "error", err, "path", resolvedPath)
	}
}

func (h *ServeFileHandler) getContentType(filename string) string {
	ext := filepath.Ext(filename)
	contentType := mime.TypeByExtension(ext)

	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Add charset for text-based content types
	if strings.HasPrefix(contentType, "text/") ||
		contentType == "application/javascript" ||
		contentType == "application/json" {
		contentType += "; charset=utf-8"
	}

	return contentType
}
