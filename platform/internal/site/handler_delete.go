package site

import (
	"context"
	"net/http"
)

type deleteSiteService interface {
	Delete(context.Context, string) error
}

type DeleteSiteHandler struct {
	service deleteSiteService
}

func NewDeleteSiteHandler(service deleteSiteService) *DeleteSiteHandler {
	return &DeleteSiteHandler{
		service: service,
	}
}

func (h *DeleteSiteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	siteId := r.Header.Get("H-Site-Id")
	if siteId == "" {
		http.Error(w, "Missing site id", http.StatusBadRequest)
		return
	}

	err := h.service.Delete(ctx, siteId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
