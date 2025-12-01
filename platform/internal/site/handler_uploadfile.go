package site

import (
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

	uploadSessionId := r.Header.Get("H-UploadSession-Id")
	if uploadSessionId == "" {
		http.Error(w, "Missing upload session id", http.StatusBadRequest)
		return
	}

	r.ParseMultipartForm(20 << 20) // 20 MB
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	err = h.service.UploadFile(ctx, uploadSessionId, handler.Filename, file, handler.Size)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
