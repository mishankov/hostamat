package handlers

import (
	"context"
	"encoding/json"
	"hostamat/pkg/platforma/application"
	"net/http"
)

type healther interface {
	Health(context.Context) *application.ApplicationHealth
}

type HealthCheckHandler struct {
	app healther
}

func NewHealthCheckHandler(app healther) *HealthCheckHandler {
	return &HealthCheckHandler{app: app}
}

func (h *HealthCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	health := h.app.Health(r.Context())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(health)
}
