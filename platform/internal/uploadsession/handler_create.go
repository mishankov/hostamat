package uploadsession

import (
	"encoding/json"
	"net/http"
)

type createUploadSessionService interface {
	Create() (*UploadSession, error)
}

type CreateUploadSessionHandler struct {
	service createUploadSessionService
}

func NewCreateUploadSessionHandler(service createUploadSessionService) *CreateUploadSessionHandler {
	return &CreateUploadSessionHandler{
		service: service,
	}
}

func (h *CreateUploadSessionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, err := h.service.Create()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(session)
}
