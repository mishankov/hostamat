package file

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	"hostamat/internal/adapters/caddy"
)

type UploadHandler struct {
	service     *Service
	caddyAPI    *caddy.CaddyAPI
	baseDomain  string
	shorthandId string
}

func NewUploadHandler(service *Service, caddyAPI *caddy.CaddyAPI, baseDomain, shorthandId string) *UploadHandler {
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

	subdomain := h.service.GenerateSubdomain()
	response := struct {
		Link          string   `json:"link"`
		UploadedFiles []string `json:"uploadedFiles"`
		SkippedFiles  []string `json:"skippedFiles"`
		Errors        []string `json:"errors"`
	}{}

	var uploaded, skipped []string
	switch filepath.Ext(handler.Filename) {
	case ".html":
		if err := h.service.UploadSingleFile(ctx, subdomain, handler.Filename, file, handler.Size); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		uploaded = []string{handler.Filename}
	case ".zip":
		uploaded, skipped, err = h.service.UploadZipArchive(ctx, subdomain, file, handler.Size)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, "Only html and zip files are allowed", http.StatusBadRequest)
		return
	}

	if err := h.caddyAPI.AddSubdomain(ctx, h.baseDomain, subdomain, h.shorthandId); err != nil {
		response.Errors = append(response.Errors, "Failed to create subdomain: "+err.Error())
	}

	response.Link = subdomain + "." + h.baseDomain
	response.UploadedFiles = uploaded
	response.SkippedFiles = skipped

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
