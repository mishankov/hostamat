package formcollection

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type listDefinitionsService interface {
	ListDefinitions(ctx context.Context, siteID string) ([]DiscoveredForm, error)
}

type ListDefinitionsHandler struct {
	service listDefinitionsService
}

func NewListDefinitionsHandler(service listDefinitionsService) *ListDefinitionsHandler {
	return &ListDefinitionsHandler{service: service}
}

func (h *ListDefinitionsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	forms, err := h.service.ListDefinitions(r.Context(), r.URL.Query().Get("siteId"))
	if err != nil {
		if errors.Is(err, errAccessDenied) {
			respondJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Нужно войти в аккаунт")
			return
		}
		respondJSONError(w, http.StatusInternalServerError, "OTHER", "Не удалось загрузить формы")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(forms)
}
