package feedback

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/platforma-dev/platforma/auth"
)

type createService interface {
	Create(ctx context.Context, message, url string) (*Feedback, error)
}

type CreateHandler struct {
	service createService
}

func NewCreateHandler(service createService) *CreateHandler {
	return &CreateHandler{service: service}
}

func (h *CreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MaxRequestBodyBytes)

	var requestBody struct {
		Message string `json:"message"`
		URL     string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		respondWithError(w, http.StatusBadRequest, errorCodeInvalidBody, "Неверный формат запроса")
		return
	}

	item, err := h.service.Create(r.Context(), requestBody.Message, requestBody.URL)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrUserNotFound):
			respondWithError(w, http.StatusUnauthorized, errorCodeUnauthorized, "Нужно войти в аккаунт")
		case errors.Is(err, errMissingMessage):
			respondWithError(w, http.StatusBadRequest, errorCodeMissingMessage, err.Error())
		case errors.Is(err, errLongMessage):
			respondWithError(w, http.StatusBadRequest, errorCodeMessageTooLong, err.Error())
		case errors.Is(err, errLongURL):
			respondWithError(w, http.StatusBadRequest, errorCodeURLTooLong, err.Error())
		case errors.Is(err, errMissingEmail):
			respondWithError(w, http.StatusBadRequest, errorCodeMissingEmail, err.Error())
		case errors.Is(err, errRateLimited):
			respondWithError(w, http.StatusTooManyRequests, errorCodeRateLimited, err.Error())
		default:
			respondWithError(w, http.StatusInternalServerError, errorCodeOther, "Не удалось отправить сообщение")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(struct {
		ID string `json:"id"`
	}{
		ID: item.ID,
	})
}

func respondWithError(w http.ResponseWriter, status int, code errorCode, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(handlerError{
		ErrorCode:    code,
		ErrorMessage: message,
	})
}
