package site

import (
	"encoding/json"
	"net/http"

	"github.com/platforma-dev/platforma/log"
)

type CloseUploadSessionHandler struct {
	service *Service
}

func NewCloseUploadSessionHandler(service *Service) *CloseUploadSessionHandler {
	return &CloseUploadSessionHandler{
		service: service,
	}
}

func (h *CloseUploadSessionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	uploadSessionId := r.Header.Get("H-UploadSession-Id")
	if uploadSessionId == "" {
		http.Error(w, "Missing upload session id", http.StatusBadRequest)
		return
	}

	files, err := h.service.CloseUploadSession(r.Context(), uploadSessionId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(files); err != nil {
		log.ErrorContext(r.Context(), "failed to encode body", "error", err.Error())
		return
	}
}
