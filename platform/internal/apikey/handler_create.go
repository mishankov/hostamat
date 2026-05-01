package apikey

import (
	"encoding/json"
	"net/http"
)

type CreateHandler struct {
	service *Service
}

func NewCreateHandler(service *Service) *CreateHandler {
	return &CreateHandler{service: service}
}

func (h *CreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var request CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondDashboardError(w, ErrInvalidBody)
		return
	}

	response, err := h.service.Create(r.Context(), request)
	if err != nil {
		respondDashboardError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}
