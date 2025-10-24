package site

import (
	"encoding/json"
	"net/http"
)

type UploadFileHandler struct {
	service *Service
}

func NewUploadFileHandler(service *Service) *UploadFileHandler {
	return &UploadFileHandler{
		service: service,
	}
}
func (h *UploadFileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	siteId := r.Header.Get("H-Site-Id")
	if siteId == "" {
		http.Error(w, "Missing site id", http.StatusBadRequest)
		return
	}

	r.ParseMultipartForm(10 << 20) // 10 MB
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	uploaded, skipped, err := h.service.fileService.UploadFile(ctx, siteId, handler.Filename, file, handler.Size)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := struct {
		UploadedFiles []string `json:"uploadedFiles"`
		SkippedFiles  []string `json:"skippedFiles"`
	}{
		UploadedFiles: uploaded,
		SkippedFiles:  skipped,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
