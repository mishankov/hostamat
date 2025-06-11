package handlers

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"hostamat/internal/adapters/caddy"
	"hostamat/internal/adapters/timeweb"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type logger interface {
	DebugContext(ctx context.Context, msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
}

type UploadFileHandler struct {
	logger     logger
	timewebAPI *timeweb.TimewebAPI
	caddyAPI   *caddy.CaddyAPI
	baseDomain string
}

func NewUploadFileHandler(logger logger, timewebAPI *timeweb.TimewebAPI, caddyAPI *caddy.CaddyAPI, baseDomain string) *UploadFileHandler {
	return &UploadFileHandler{logger: logger, timewebAPI: timewebAPI, caddyAPI: caddyAPI, baseDomain: baseDomain}
}

func (ufh *UploadFileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(r.Context(), "traceId", uuid.NewString())

	ufh.logger.DebugContext(ctx, "upload file handler started")

	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseMultipartForm(10 << 20) // 10 MB

	file, handler, err := r.FormFile("file")
	if err != nil {
		ufh.logger.ErrorContext(ctx, "error retrieving file", "error", err)
		http.Error(w, "Error retrieving file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	fmt.Fprintf(w, "Received file: %s\n", handler.Filename)
	fmt.Fprintf(w, "File size: %d\n", handler.Size)

	var buf = bytes.Buffer{}
	if _, err := io.Copy(&buf, file); err != nil {
		ufh.logger.ErrorContext(ctx, "error saving file", "error", err)
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	subdomain := randomBase64String(5)

	err = ufh.timewebAPI.AddSubdomain(ctx, ufh.baseDomain, subdomain)
	if err != nil {
		ufh.logger.ErrorContext(ctx, "error in timeweb add subdomain", "error", err)
	}

	err = ufh.caddyAPI.AddSubdomain(ctx, subdomain)
	if err != nil {
		ufh.logger.ErrorContext(ctx, "error in caddy add subdomain", "error", err)
	}

	fmt.Fprintf(w, "File saved successfully.")

	ufh.logger.DebugContext(ctx, "upload file handler finished")
}

func randomBase64String(length int) string {
	bytes := make([]byte, length+1)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatal(err)
	}

	result := base64.URLEncoding.EncodeToString(bytes)

	return result[:len(result)-1]
}
