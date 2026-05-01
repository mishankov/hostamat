package apikey

import (
	"encoding/json"
	"errors"
	"net/http"
)

type dashboardError struct {
	ErrorCode    ErrorCode `json:"errorCode"`
	ErrorMessage string    `json:"errorMessage"`
}

func respondDashboardError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	switch {
	case errors.Is(err, ErrInvalidBody), errors.Is(err, ErrMissingName), errors.Is(err, ErrNameTooLong), errors.Is(err, ErrMissingScopes), errors.Is(err, ErrInvalidScope), errors.Is(err, ErrInvalidExpiration):
		w.WriteHeader(http.StatusBadRequest)
	case errors.Is(err, ErrActiveKeyLimit):
		w.WriteHeader(http.StatusConflict)
	case errors.Is(err, ErrUnknownKey):
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	_ = json.NewEncoder(w).Encode(dashboardError{
		ErrorCode:    CodeForError(err),
		ErrorMessage: err.Error(),
	})
}
