package userauth

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	platformauth "github.com/platforma-dev/platforma/auth"
	"github.com/platforma-dev/platforma/log"
)

type getUserService interface {
	GetFromSession(ctx context.Context, sessionID string) (*platformauth.User, error)
	GetProfileFromSession(ctx context.Context, sessionID string) (*User, error)
	CookieName() string
}

type GetHandler struct {
	service getUserService
}

func NewGetHandler(service getUserService) *GetHandler {
	return &GetHandler{service: service}
}

func (h *GetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ev := log.EventFromContext(r.Context())

	if r.Method != http.MethodGet {

		ev.AddStep(slog.LevelWarn, "userauth.Get.method_not_allowed")

		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie(h.service.CookieName())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := h.service.GetProfileFromSession(r.Context(), cookie.Value)
	if errors.Is(err, platformauth.ErrUserNotFound) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, "failed to get user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	resp := struct {
		Username string   `json:"username"`
		Email    string   `json:"email"`
		Provider Provider `json:"provider"`
	}{
		Username: user.Username,
		Email:    stringValue(user.Email),
		Provider: ProviderYandex,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {

		ev.AddError(err)

		log.ErrorContext(r.Context(), "failed to encode response", "error", err)
	}
}
