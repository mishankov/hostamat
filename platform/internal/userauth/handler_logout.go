package userauth

import (
	"context"
	"net/http"
	"time"
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
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie(h.service.CookieName())
	if err == nil {
		if err := h.service.DeleteSession(r.Context(), cookie.Value); err != nil {
			http.Error(w, "failed to logout", http.StatusInternalServerError)
			return
		}
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
