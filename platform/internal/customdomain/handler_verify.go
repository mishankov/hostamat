package customdomain

import (
	"context"
	"encoding/json"
	"net/http"
)

type verifyService interface {
	Verify(context.Context, string) (*DomainRecord, error)
}

type VerifyHandler struct {
	service verifyService
}

func NewVerifyHandler(service verifyService) *VerifyHandler {
	return &VerifyHandler{service: service}
}

func (h *VerifyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	domainID := r.Header.Get("H-Domain-Id")
	record, err := h.service.Verify(r.Context(), domainID)
	if err != nil {
		respondWithError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(record)
}
