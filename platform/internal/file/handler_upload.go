package file

import (
	"archive/zip"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"hostamat/internal/adapters/caddy"
)

type serviceUpload interface {
	UploadFile(ctx context.Context, fileName string, file io.Reader, size int64) error
	FileExists(ctx context.Context, fileName string) (bool, error)
}

type UploadHandler struct {
	service     serviceUpload
	caddyAPI    *caddy.CaddyAPI
	baseDomain  string
	shorthandId string
}

func NewUploadHandler(service serviceUpload, caddyAPI *caddy.CaddyAPI, baseDomain, shorthandId string) *UploadHandler {
	return &UploadHandler{
		service:     service,
		caddyAPI:    caddyAPI,
		baseDomain:  baseDomain,
		shorthandId: shorthandId,
	}
}

func (h *UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseMultipartForm(10 << 20) // 10 MB

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	subdomain := randomBase64String(5)
	response := struct {
		Link          string   `json:"link"`
		UploadedFiles []string `json:"uploadedFiles"`
		SkippedFiles  []string `json:"skippedFiles"`
		Errors        []string `json:"errors"`
	}{}

	fileExtension := strings.Split(handler.Filename, ".")[strings.Count(handler.Filename, ".")]
	switch fileExtension {
	case "html":
		err = h.service.UploadFile(ctx, subdomain+"/"+handler.Filename, file, handler.Size)
		if err != nil {
			http.Error(w, "Error uploading file", http.StatusInternalServerError)
			return
		}
		response.UploadedFiles = append(response.UploadedFiles, handler.Filename)
	case "zip":
		zipFile, err := zip.NewReader(file, handler.Size)
		if err != nil {
			http.Error(w, "Error reading zip file", http.StatusBadRequest)
			return
		}

		for _, f := range zipFile.File {
			if f.FileInfo().IsDir() {
				continue
			}

			fileName := strings.Join(strings.Split(f.Name, "/")[1:], "/")
			if fileName == "" || strings.HasPrefix(fileName, ".") || strings.Contains(fileName, "/.") {
				response.SkippedFiles = append(response.SkippedFiles, fileName)
				continue
			}

			fileReader, err := f.Open()
			if err != nil {
				response.Errors = append(response.Errors, fileName+": "+err.Error())
				continue
			}

			err = h.service.UploadFile(ctx, subdomain+"/"+fileName, fileReader, int64(f.UncompressedSize64))
			fileReader.Close()
			if err != nil {
				response.Errors = append(response.Errors, fileName+": "+err.Error())
				continue
			}

			response.UploadedFiles = append(response.UploadedFiles, fileName)
		}
	default:
		http.Error(w, "Only html and zip files are allowed", http.StatusBadRequest)
		return
	}

	if err := h.caddyAPI.AddSubdomain(ctx, h.baseDomain, subdomain, h.shorthandId); err != nil {
		response.Errors = append(response.Errors, "Failed to create subdomain: "+err.Error())
	}

	response.Link = subdomain + "." + h.baseDomain
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func randomBase64String(length int) string {
	bytes := make([]byte, length+1)
	rand.Read(bytes)
	return strings.ToLower(base64.URLEncoding.EncodeToString(bytes)[:length])
}
