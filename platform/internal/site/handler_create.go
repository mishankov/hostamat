package site

import (
	"encoding/json"
	"net/http"
)

type createSiteService interface {
	Create(domain string) (*Site, error)
}

type CreateSiteHandler struct {
	service createSiteService
}

func NewCreateSiteHandler(service createSiteService) *CreateSiteHandler {
	return &CreateSiteHandler{
		service: service,
	}
}

func (h *CreateSiteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Domain string `json:"domain"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	site, err := h.service.Create(request.Domain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(site)
}
