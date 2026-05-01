package apikey

import (
	"encoding/json"
	"net/http"
)

type UpdateHandler struct {
	service *Service
}

func NewUpdateHandler(service *Service) *UpdateHandler {
	return &UpdateHandler{service: service}
}

func (h *UpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	keyID := r.Header.Get("H-API-Key-Id")
	if keyID == "" {
		keyID = r.PathValue("keyId")
	}

	var request UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondDashboardError(w, ErrInvalidBody)
		return
	}

	key, err := h.service.UpdateName(r.Context(), keyID, request.Name)
	if err != nil {
		respondDashboardError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(key)
}
