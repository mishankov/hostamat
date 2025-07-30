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

	folder := r.URL.Query().Get("folder")
	if folder == "" {
		http.Error(w, "Missing folder parameter", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteFromFolder(ctx, folder); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
