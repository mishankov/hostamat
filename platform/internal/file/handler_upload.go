package file

import (
	"encoding/json"
	"net/http"
)

type UploadHandler struct {
	service *Service
}

func NewUploadHandler(service *Service) *UploadHandler {
	return &UploadHandler{
		service: service,
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

	link, uploaded, skipped, err := h.service.UploadFile(ctx, handler.Filename, file, handler.Size)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := struct {
		Link          string   `json:"link"`
		UploadedFiles []string `json:"uploadedFiles"`
		SkippedFiles  []string `json:"skippedFiles"`
	}{
		Link:          link,
		UploadedFiles: uploaded,
		SkippedFiles:  skipped,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
