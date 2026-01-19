package site

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
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

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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

	site, err := h.service.Create(ctx, requestBody.Name)
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
			ErrorMessage: "Ошибка сервера при создании сайта",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(site)
}
