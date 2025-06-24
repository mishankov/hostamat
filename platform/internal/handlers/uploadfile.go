package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"hostamat/internal/adapters/caddy"
	"hostamat/internal/adapters/timeweb"
	"io"
	"log"
	"net/http"
)

type logger interface {
	DebugContext(ctx context.Context, msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
}

type s3 interface {
	PutObject(ctx context.Context, bucketName, fileName string, file io.Reader, fileSize int64) error
}

type UploadFileHandler struct {
	logger     logger
	timewebAPI *timeweb.TimewebAPI
	caddyAPI   *caddy.CaddyAPI
	s3         s3
	s3Bucket   string
	baseDomain string
}

func NewUploadFileHandler(logger logger, timewebAPI *timeweb.TimewebAPI, caddyAPI *caddy.CaddyAPI, s3 s3, s3Bucket string, baseDomain string) *UploadFileHandler {
	return &UploadFileHandler{
		logger:     logger,
		timewebAPI: timewebAPI,
		caddyAPI:   caddyAPI,
		s3:         s3,
		s3Bucket:   s3Bucket,
		baseDomain: baseDomain,
	}
}

func (h *UploadFileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseMultipartForm(10 << 20) // 10 MB

	file, handler, err := r.FormFile("file")
	if err != nil {
		h.logger.ErrorContext(ctx, "error retrieving file", "error", err)
		http.Error(w, "Error retrieving file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	subdomain := randomBase64String(5)

	err = h.timewebAPI.AddSubdomain(ctx, h.baseDomain, subdomain)
	if err != nil {
		h.logger.ErrorContext(ctx, "error in timeweb add subdomain", "error", err)
	}

	err = h.caddyAPI.AddSubdomain(ctx, subdomain)
	if err != nil {
		h.logger.ErrorContext(ctx, "error in caddy add subdomain", "error", err)
	}

	err = h.s3.PutObject(ctx, h.s3Bucket, subdomain+"/"+handler.Filename, file, handler.Size)
	if err != nil {
		log.Fatal("Error saving file to minio", err)
	}

	fmt.Fprint(w, subdomain+"."+h.baseDomain)
}

func randomBase64String(length int) string {
	bytes := make([]byte, length+1)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatal(err)
	}

	result := base64.URLEncoding.EncodeToString(bytes)

	return result[:len(result)-1]
}
