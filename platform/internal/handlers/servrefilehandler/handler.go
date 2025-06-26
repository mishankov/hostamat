package servrefilehandler

import (
	"context"
	"io"
	"net/http"
	"strings"
)

type logger interface {
	DebugContext(ctx context.Context, msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
}

type s3 interface {
	ObjectExists(ctx context.Context, bucketName, fileName string) (bool, error)
	GetObject(ctx context.Context, bucketName, fileName string) (io.Reader, error)
}

type GetFileHandler struct {
	logger     logger
	s3         s3
	bucketName string
}

func New(logger logger, s3 s3, bucketName string) *GetFileHandler {
	return &GetFileHandler{
		logger:     logger,
		s3:         s3,
		bucketName: bucketName,
	}
}

func (h *GetFileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	subdomain := strings.Split(r.Host, ".")[0]
	filePath := r.URL.Path
	if filePath == "/" {
		filePath = "/index.html"
	}

	fileName := subdomain + filePath
	exists, err := h.s3.ObjectExists(ctx, h.bucketName, fileName)
	if err != nil {
		h.logger.ErrorContext(ctx, "error checking file existence", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	file, err := h.s3.GetObject(ctx, h.bucketName, fileName)
	if err != nil {
		h.logger.ErrorContext(ctx, "error getting file", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", contentTypeFromFileName(filePath))

	n, err := io.Copy(w, file)
	if err != nil {
		h.logger.ErrorContext(ctx, "error copying file", "error", err)
		return
	}
	h.logger.DebugContext(ctx, "file served", "bytes", n)
}

func contentTypeFromFileName(fileName string) string {
	contentType := "application/octet-stream"
	if strings.HasSuffix(fileName, ".html") {
		contentType = "text/html"
	} else if strings.HasSuffix(fileName, ".css") {
		contentType = "text/css"
	} else if strings.HasSuffix(fileName, ".js") {
		contentType = "application/javascript"
	} else if strings.HasSuffix(fileName, ".png") {
		contentType = "image/png"
	} else if strings.HasSuffix(fileName, ".jpg") || strings.HasSuffix(fileName, ".jpeg") {
		contentType = "image/jpeg"
	} else if strings.HasSuffix(fileName, ".svg") {
		contentType = "image/svg+xml"
	}

	return contentType
}
