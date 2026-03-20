package site

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/platforma-dev/platforma/log"
)

type createSiteService interface {
	Create(context.Context, string) (*Site, error)
}

type CreateSiteHandler struct {
	service createSiteService
}

func NewCreateSiteHandler(service createSiteService) *CreateSiteHandler {
	return &CreateSiteHandler{
		service: service,
	}
}

func (h *CreateSiteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ev := log.EventFromContext(ctx)

	if r.Method != http.MethodPost {

		ev.AddStep(slog.LevelWarn, "site.Create.method_not_allowed")

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestBody struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {

		ev.AddStep(slog.LevelWarn, "site.Create.invalid_body")

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

		ev.AddStep(slog.LevelWarn, "site.Create.missing_name")

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

	site, err := h.service.Create(ctx, requestBody.Name)
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
			ErrorMessage: "Ошибка сервера при создании сайта",
		})
		if err != nil {
			ev.AddError(err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	ev.AddAttrs(map[string]any{
		"site.id":     site.Id,
		"site.domain": site.Domain,
	})

	err = json.NewEncoder(w).Encode(site)
	if err != nil {
		ev.AddError(err)
	}
}
