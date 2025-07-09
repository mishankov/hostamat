package uploadsession

import (
	"encoding/json"
	"net/http"
)

type closeUploadSessionService interface {
	Get(id string) (*UploadSession, error)
	Finish(id string) error
}

type CloseUploadSessionHandler struct {
	service closeUploadSessionService
}

func NewCloseUploadSessionHandler(service closeUploadSessionService) *CloseUploadSessionHandler {
	return &CloseUploadSessionHandler{
		service: service,
	}
}

func (h *CloseUploadSessionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionID := r.Header.Get("H-Upload-Session")
	if sessionID == "" {
		http.Error(w, "Missing H-Upload-Session header", http.StatusBadRequest)
		return
	}

	err := h.service.Finish(sessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updatedSession, err := h.service.Get(sessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedSession)
}
