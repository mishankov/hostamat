package file

import (
	"context"
	"net/http"
)

type deleteService interface {
	DeleteFromFolder(ctx context.Context, folder string) error
}

type DeleteHandler struct {
	service deleteService
}

func NewDeleteHandler(service deleteService) *DeleteHandler {
	return &DeleteHandler{
		service: service,
	}
}

func (h *DeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	siteId := r.Header.Get("H-Site-Id")
	if siteId == "" {
		http.Error(w, "Missing H-Site-Id header", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteFromFolder(ctx, siteId); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
