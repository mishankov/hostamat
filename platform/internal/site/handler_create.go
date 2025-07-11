package site

import (
	"context"
	"encoding/json"
	"net/http"
)

type createSiteService interface {
	Create(context.Context) (*Site, error)
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
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	site, err := h.service.Create(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(site)
}
