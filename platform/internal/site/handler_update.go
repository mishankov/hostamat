package site

import (
	"context"
	"encoding/json"
	"net/http"
)

type updateSiteService interface {
	UpdateName(ctx context.Context, id, name string) error
}

type UpdateSiteHandler struct {
	service updateSiteService
}

func NewUpdateSiteHandler(service updateSiteService) *UpdateSiteHandler {
	return &UpdateSiteHandler{
		service: service,
	}
}

func (h *UpdateSiteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	siteId := r.Header.Get("H-Site-Id")
	if siteId == "" {
		http.Error(w, "Missing site id", http.StatusBadRequest)
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

	err := h.service.UpdateName(ctx, siteId, requestBody.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
