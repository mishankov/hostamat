package userauth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	platformauth "github.com/platforma-dev/platforma/auth"
	"github.com/platforma-dev/platforma/log"
)

type userDeleter interface {
	DeleteUser(ctx context.Context) error
}

type DeleteHandler struct {
	service userDeleter
}

func NewDeleteHandler(service userDeleter) *DeleteHandler {
	return &DeleteHandler{service: service}
}

func (h *DeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := h.service.DeleteUser(r.Context())
	if errors.Is(err, platformauth.ErrUserNotFound) {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode("User deleted successfully"); err != nil {
		log.ErrorContext(r.Context(), "failed to encode response", "error", err)
	}
}
