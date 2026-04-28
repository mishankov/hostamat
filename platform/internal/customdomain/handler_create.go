package customdomain

import (
	"context"
	"encoding/json"
	"net/http"
)

type createService interface {
	Create(context.Context, string, string) (*DomainRecord, error)
}

type CreateHandler struct {
	service createService
}

func NewCreateHandler(service createService) *CreateHandler {
	return &CreateHandler{service: service}
}

func (h *CreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		SiteID   string `json:"siteId"`
		Hostname string `json:"hostname"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(handlerError{
			ErrorCode:    errorCodeInvalidBody,
			ErrorMessage: "Неверный формат запроса",
		})
		return
	}

	record, err := h.service.Create(r.Context(), requestBody.SiteID, requestBody.Hostname)
	if err != nil {
		respondWithError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(record)
}
