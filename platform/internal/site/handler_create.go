package site

import (
	"context"
	"encoding/json"
	"net/http"
)

type createSiteService interface {
	Create(context.Context, string) (*Site, error)
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

	var requestBody struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if requestBody.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	site, err := h.service.Create(ctx, requestBody.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(site)
}
