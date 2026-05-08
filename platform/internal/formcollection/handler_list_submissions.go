package formcollection

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type listSubmissionsService interface {
	ListSubmissions(ctx context.Context, filters SubmissionFilters) ([]Submission, error)
}

type ListSubmissionsHandler struct {
	service listSubmissionsService
}

func NewListSubmissionsHandler(service listSubmissionsService) *ListSubmissionsHandler {
	return &ListSubmissionsHandler{service: service}
}

func (h *ListSubmissionsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	submissions, err := h.service.ListSubmissions(r.Context(), SubmissionFilters{
		SiteID: r.URL.Query().Get("siteId"),
		FormID: r.URL.Query().Get("formId"),
	})
	if err != nil {
		if errors.Is(err, errAccessDenied) {
			respondJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Нужно войти в аккаунт")
			return
		}
		respondJSONError(w, http.StatusInternalServerError, "OTHER", "Не удалось загрузить ответы")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(submissions)
}
