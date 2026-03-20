package userauth

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/platforma-dev/platforma/log"
)

type logoutService interface {
	DeleteSession(ctx context.Context, sessionID string) error
	CookieName() string
}

type LogoutHandler struct {
	service logoutService
}

func NewLogoutHandler(service logoutService) *LogoutHandler {
	return &LogoutHandler{service: service}
}

func (h *LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ev := log.EventFromContext(r.Context())

	if r.Method != http.MethodPost {

		ev.AddStep(slog.LevelWarn, "userauth.Logout.method_not_allowed")

		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie(h.service.CookieName())
	if err == nil {
		if err := h.service.DeleteSession(r.Context(), cookie.Value); err != nil {
			http.Error(w, "failed to logout", http.StatusInternalServerError)
			return
		}
	} else {
		ev.AddStep(slog.LevelInfo, "userauth.Logout.missing_cookie")
	}

	http.SetCookie(w, &http.Cookie{
		Name:     h.service.CookieName(),
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   requestIsSecure(r),
	})

	w.WriteHeader(http.StatusOK)
}
