package formcollection

import (
	"context"
	"errors"
	"net/http"
)

type deleteSubmissionService interface {
	DeleteSubmission(ctx context.Context, submissionID string) error
}

type DeleteSubmissionHandler struct {
	service deleteSubmissionService
}

func NewDeleteSubmissionHandler(service deleteSubmissionService) *DeleteSubmissionHandler {
	return &DeleteSubmissionHandler{service: service}
}

func (h *DeleteSubmissionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		respondJSONError(w, http.StatusBadRequest, "MISSING_SUBMISSION_ID", "Отсутствует ID ответа")
		return
	}

	err := h.service.DeleteSubmission(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, errAccessDenied):
			respondJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Нужно войти в аккаунт")
		case errors.Is(err, errFormNotFound):
			respondJSONError(w, http.StatusNotFound, "SUBMISSION_NOT_FOUND", "Ответ не найден")
		default:
			respondJSONError(w, http.StatusInternalServerError, "OTHER", "Не удалось удалить ответ")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
