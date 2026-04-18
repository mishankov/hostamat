package customdomain

import (
	"context"
	"net/http"
)

type deleteService interface {
	Delete(context.Context, string) error
}

type DeleteHandler struct {
	service deleteService
}

func NewDeleteHandler(service deleteService) *DeleteHandler {
	return &DeleteHandler{service: service}
}

func (h *DeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	domainID := r.Header.Get("H-Domain-Id")
	if err := h.service.Delete(r.Context(), domainID); err != nil {
		respondWithError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
