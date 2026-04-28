package customdomain

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type listService interface {
	List(context.Context) ([]DomainRecord, error)
}

type ListHandler struct {
	service listService
}

func NewListHandler(service listService) *ListHandler {
	return &ListHandler{service: service}
}

func (h *ListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	domains, err := h.service.List(r.Context())
	if err != nil {
		respondWithError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(struct {
		Domains []DomainRecord `json:"domains"`
	}{
		Domains: domains,
	})
}

func respondWithError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	switch {
	case errors.Is(err, errAccessDenied):
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(handlerError{
			ErrorCode:    errorCodeAccessDenied,
			ErrorMessage: err.Error(),
		})
	case errors.Is(err, errMissingHostname):
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(handlerError{
			ErrorCode:    errorCodeMissingHostname,
			ErrorMessage: err.Error(),
		})
	case errors.Is(err, errInvalidHostname), errors.Is(err, errApexHostnameRejected):
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(handlerError{
			ErrorCode:    errorCodeInvalidHostname,
			ErrorMessage: err.Error(),
		})
	case errors.Is(err, errMissingSiteID):
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(handlerError{
			ErrorCode:    errorCodeMissingSiteID,
			ErrorMessage: err.Error(),
		})
	case errors.Is(err, errMissingDomainID):
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(handlerError{
			ErrorCode:    errorCodeMissingDomainID,
			ErrorMessage: err.Error(),
		})
	case errors.Is(err, errHostnameTaken):
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(handlerError{
			ErrorCode:    errorCodeHostnameTaken,
			ErrorMessage: err.Error(),
		})
	case errors.Is(err, errNotFound):
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(handlerError{
			ErrorCode:    errorCodeNotFound,
			ErrorMessage: err.Error(),
		})
	default:
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(handlerError{
			ErrorCode:    errorCodeOther,
			ErrorMessage: "Ошибка сервера при работе с доменами",
		})
	}
}
