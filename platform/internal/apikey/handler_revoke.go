package apikey

import "net/http"

type RevokeHandler struct {
	service *Service
}

func NewRevokeHandler(service *Service) *RevokeHandler {
	return &RevokeHandler{service: service}
}

func (h *RevokeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	keyID := r.Header.Get("H-API-Key-Id")
	if keyID == "" {
		keyID = r.PathValue("keyId")
	}

	if err := h.service.Revoke(r.Context(), keyID); err != nil {
		respondDashboardError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
