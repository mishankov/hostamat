package apikey

import (
	"encoding/json"
	"net/http"
)

type ListHandler struct {
	service *Service
}

func NewListHandler(service *Service) *ListHandler {
	return &ListHandler{service: service}
}

func (h *ListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	keys, err := h.service.List(r.Context())
	if err != nil {
		respondDashboardError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(struct {
		Keys []PublicKeyMetadata `json:"keys"`
	}{Keys: keys})
}
