package site

import (
	"context"
	"encoding/json"
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(handlerError{
			ErrorCode:    errorCodeMissingSiteId,
			ErrorMessage: "Отсутствует ID сайта",
		})
		return
	}

	err := h.service.Delete(ctx, siteId)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(handlerError{
			ErrorCode:    errorCodeOther,
			ErrorMessage: "Ошибка сервера при удалении сайта",
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
