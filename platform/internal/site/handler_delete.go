package site

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/platforma-dev/platforma/log"
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
	ev := log.EventFromContext(ctx)

	if r.Method != http.MethodDelete {

		ev.AddStep(slog.LevelWarn, "site.Delete.method_not_allowed")

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	siteId := r.Header.Get("H-Site-Id")
	if siteId == "" {

		ev.AddStep(slog.LevelWarn, "site.Delete.missing_site_id")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(handlerError{
			ErrorCode:    errorCodeMissingSiteId,
			ErrorMessage: "Отсутствует ID сайта",
		})
		if err != nil {
			ev.AddError(err)
		}
		return
	}

	ev.AddAttrs(map[string]any{
		"site.id": siteId,
	})

	err := h.service.Delete(ctx, siteId)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode(handlerError{
			ErrorCode:    errorCodeOther,
			ErrorMessage: "Ошибка сервера при удалении сайта",
		})
		if err != nil {
			ev.AddError(err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
