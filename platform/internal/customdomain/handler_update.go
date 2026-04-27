package customdomain

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/platforma-dev/platforma/log"
)

type updateDomainService interface {
	UpdateSite(context.Context, string, string) (*DomainRecord, error)
}

type UpdateHandler struct {
	service updateDomainService
}

func NewUpdateHandler(service updateDomainService) *UpdateHandler {
	return &UpdateHandler{service: service}
}

func (h *UpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ev := log.EventFromContext(ctx)

	if r.Method != http.MethodPut {
		ev.AddStep(slog.LevelWarn, "customdomain.Update.method_not_allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	domainID := r.Header.Get("H-Domain-Id")
	if domainID == "" {
		ev.AddStep(slog.LevelWarn, "customdomain.Update.missing_domain_id")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(handlerError{
			ErrorCode:    errorCodeMissingDomainID,
			ErrorMessage: "Отсутствует ID домена",
		})
		return
	}

	var requestBody struct {
		SiteID string `json:"siteId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		ev.AddStep(slog.LevelWarn, "customdomain.Update.invalid_body")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(handlerError{
			ErrorCode:    errorCodeInvalidBody,
			ErrorMessage: "Неверный формат запроса",
		})
		return
	}

	record, err := h.service.UpdateSite(ctx, domainID, requestBody.SiteID)
	if err != nil {
		if errors.Is(err, errMissingSiteID) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(handlerError{
				ErrorCode:    errorCodeMissingSiteID,
				ErrorMessage: err.Error(),
			})
			return
		}
		if errors.Is(err, errNotFound) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(handlerError{
				ErrorCode:    errorCodeNotFound,
				ErrorMessage: err.Error(),
			})
			return
		}
		if errors.Is(err, errAccessDenied) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(handlerError{
				ErrorCode:    errorCodeAccessDenied,
				ErrorMessage: err.Error(),
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(handlerError{
			ErrorCode:    errorCodeOther,
			ErrorMessage: "Ошибка сервера при изменении домена",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(record)
}
