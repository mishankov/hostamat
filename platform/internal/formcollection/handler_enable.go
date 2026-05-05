package formcollection

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type enableService interface {
	EnableCollection(ctx context.Context, siteID, discoveryID string, formIDs []string) (*DiscoverySummary, error)
}

type EnableHandler struct {
	service enableService
}

func NewEnableHandler(service enableService) *EnableHandler {
	return &EnableHandler{service: service}
}

func (h *EnableHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body struct {
		SiteID      string   `json:"siteId"`
		DiscoveryID string   `json:"discoveryId"`
		FormIDs     []string `json:"formIds"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSONError(w, http.StatusBadRequest, "INVALID_BODY", "Неверный формат запроса")
		return
	}

	summary, err := h.service.EnableCollection(r.Context(), body.SiteID, body.DiscoveryID, body.FormIDs)
	if err != nil {
		switch {
		case errors.Is(err, errAccessDenied):
			respondJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Нужно войти в аккаунт")
		case errors.Is(err, errNoFormsSelected):
			respondJSONError(w, http.StatusBadRequest, "NO_FORMS_SELECTED", "Не выбраны формы")
		case errors.Is(err, errFormNotFound):
			respondJSONError(w, http.StatusNotFound, "FORM_NOT_FOUND", "Форма не найдена")
		case errors.Is(err, errFormNotSupported):
			respondJSONError(w, http.StatusBadRequest, "FORM_NOT_SUPPORTED", "Эта форма не поддерживается")
		default:
			respondJSONError(w, http.StatusInternalServerError, "OTHER", "Не удалось включить сбор форм")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(summary)
}
