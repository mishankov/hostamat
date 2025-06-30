package uploadfilehandler

import (
	"archive/zip"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"hostamat/internal/adapters/caddy"

	"io"
	"log"
	"net/http"
	"strings"
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
	logger                logger
	caddyAPI              *caddy.CaddyAPI
	s3                    s3
	s3Bucket              string
	baseDomain            string
	target                string
	fileHandlerPathPrefix string
}

func New(logger logger, caddyAPI *caddy.CaddyAPI, s3 s3, s3Bucket string, baseDomain string, target string, fileHandlerPathPrefix string) *UploadFileHandler {
	return &UploadFileHandler{
		logger:                logger,
		caddyAPI:              caddyAPI,
		s3:                    s3,
		s3Bucket:              s3Bucket,
		baseDomain:            baseDomain,
		target:                target,
		fileHandlerPathPrefix: fileHandlerPathPrefix,
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

	fileExtension := strings.Split(handler.Filename, ".")[strings.Count(handler.Filename, ".")]
	switch fileExtension {
	case "html":
		err = h.s3.PutObject(ctx, h.s3Bucket, subdomain+"/"+handler.Filename, file, handler.Size)
		if err != nil {
			h.logger.ErrorContext(ctx, "error in s3 put object", "error", err)
			http.Error(w, "Error uploading file", http.StatusInternalServerError)
			return
		}
	case "zip":
		zipFile, err := zip.NewReader(file, handler.Size)
		if err != nil {
			h.logger.ErrorContext(ctx, "error in zip new reader", "error", err)
			http.Error(w, "Error uploading file", http.StatusInternalServerError)
			return
		}
		for _, file := range zipFile.File {
			zippedFile, err := file.Open()
			if err != nil {
				h.logger.ErrorContext(ctx, "error in zip file open", "error", err)
				http.Error(w, "Error uploading file", http.StatusInternalServerError)
				return
			}

			fileName := strings.Join(strings.Split(file.Name, "/")[1:], "/")

			h.logger.DebugContext(ctx, "current file", "fileName", fileName, "fileSize", file.UncompressedSize64)

			if fileName == "" || file.FileInfo().IsDir() {
				continue
			}

			err = h.s3.PutObject(ctx, h.s3Bucket, subdomain+"/"+fileName, zippedFile, int64(file.UncompressedSize64))
			if err != nil {
				h.logger.ErrorContext(ctx, "error in s3 put object", "error", err, "filName", fileName, "fileSize", file.UncompressedSize64)
				http.Error(w, "Error uploading file", http.StatusInternalServerError)
				return
			}
		}
	default:
		http.Error(w, "Only html and zip files are allowed", http.StatusBadRequest)
		return
	}

	err = h.caddyAPI.AddSubdomain(ctx, h.baseDomain, subdomain, h.target, h.fileHandlerPathPrefix)
	if err != nil {
		h.logger.ErrorContext(ctx, "error in caddy add subdomain", "error", err)
	}

	fmt.Fprint(w, subdomain+"."+h.baseDomain)
}

func randomBase64String(length int) string {
	bytes := make([]byte, length+1)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatal(err)
	}

	result := base64.URLEncoding.EncodeToString(bytes)

	return strings.ToLower(result[:len(result)-1])
}
