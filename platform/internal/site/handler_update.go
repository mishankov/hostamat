package site

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
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

	if r.Method != http.MethodPut {
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

	var requestBody struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(handlerError{
			ErrorCode:    errorCodeInvalidBody,
			ErrorMessage: "Неверный формат запроса",
		})
		return
	}

	if requestBody.Name == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(handlerError{
			ErrorCode:    errorCodeMissingName,
			ErrorMessage: "Название обязательно",
		})
		return
	}

	err := h.service.UpdateName(ctx, siteId, requestBody.Name)
	if err != nil {
		if errors.Is(err, errLongName) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(handlerError{
				ErrorCode:    errorCodeNameTooLong,
				ErrorMessage: err.Error(),
			})
			return
		}
		if errors.Is(err, errShortName) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(handlerError{
				ErrorCode:    errorCodeNameTooShort,
				ErrorMessage: err.Error(),
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(handlerError{
			ErrorCode:    errorCodeOther,
			ErrorMessage: "Ошибка сервера при изменении сайта",
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
