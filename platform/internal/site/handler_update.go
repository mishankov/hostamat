package site

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/platforma-dev/platforma/log"
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
	ev := log.EventFromContext(ctx)

	if r.Method != http.MethodPut {

		ev.AddStep(slog.LevelWarn, "site.Update.method_not_allowed")

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	siteId := r.Header.Get("H-Site-Id")
	if siteId == "" {

		ev.AddStep(slog.LevelWarn, "site.Update.missing_site_id")

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

	var requestBody struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {

		ev.AddStep(slog.LevelWarn, "site.Update.invalid_body")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(handlerError{
			ErrorCode:    errorCodeInvalidBody,
			ErrorMessage: "Неверный формат запроса",
		})
		if err != nil {
			ev.AddError(err)
		}
		return
	}

	if requestBody.Name == "" {

		ev.AddStep(slog.LevelWarn, "site.Update.missing_name")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(handlerError{
			ErrorCode:    errorCodeMissingName,
			ErrorMessage: "Название обязательно",
		})
		if err != nil {
			ev.AddError(err)
		}
		return
	}

	err := h.service.UpdateName(ctx, siteId, requestBody.Name)
	if err != nil {
		if errors.Is(err, errLongName) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(handlerError{
				ErrorCode:    errorCodeNameTooLong,
				ErrorMessage: err.Error(),
			})
			if err != nil {
				ev.AddError(err)
			}
			return
		}
		if errors.Is(err, errShortName) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(handlerError{
				ErrorCode:    errorCodeNameTooShort,
				ErrorMessage: err.Error(),
			})
			if err != nil {
				ev.AddError(err)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode(handlerError{
			ErrorCode:    errorCodeOther,
			ErrorMessage: "Ошибка сервера при изменении сайта",
		})
		if err != nil {
			ev.AddError(err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
