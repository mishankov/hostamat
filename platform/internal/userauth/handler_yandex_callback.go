package userauth

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"net/http"
	"time"
)

type callbackYandexAuthService interface {
	CreateSessionFromYandexCode(ctx context.Context, code string) (string, error)
	CookieName() string
	OAuthStateCookieName() string
	OAuthStateCookiePath() string
}

type YandexCallbackHandler struct {
	service callbackYandexAuthService
}

func NewYandexCallbackHandler(service callbackYandexAuthService) *YandexCallbackHandler {
	return &YandexCallbackHandler{service: service}
}

func (h *YandexCallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	if query.Get("error") != "" {
		clearCookie(w, h.service.OAuthStateCookieName(), h.service.OAuthStateCookiePath(), requestIsSecure(r))
		writeRedirectTarget(w, "/login?error=oauth_denied")
		return
	}

	if !isValidState(r, h.service.OAuthStateCookieName(), query.Get("state")) {
		clearCookie(w, h.service.OAuthStateCookieName(), h.service.OAuthStateCookiePath(), requestIsSecure(r))
		writeRedirectTarget(w, "/login?error=oauth_failed")
		return
	}

	sessionID, err := h.service.CreateSessionFromYandexCode(r.Context(), query.Get("code"))
	if err != nil {
		clearCookie(w, h.service.OAuthStateCookieName(), h.service.OAuthStateCookiePath(), requestIsSecure(r))
		writeRedirectTarget(w, "/login?error=oauth_failed")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     h.service.CookieName(),
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   requestIsSecure(r),
	})
	clearCookie(w, h.service.OAuthStateCookieName(), h.service.OAuthStateCookiePath(), requestIsSecure(r))

	writeRedirectTarget(w, "/dashboard")
}

func isValidState(r *http.Request, cookieName, state string) bool {
	if state == "" {
		return false
	}

	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return false
	}

	return subtle.ConstantTimeCompare([]byte(cookie.Value), []byte(state)) == 1
}

func clearCookie(w http.ResponseWriter, name, path string, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     path,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   secure,
	})
}

func writeRedirectTarget(w http.ResponseWriter, redirectTo string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(struct {
		RedirectTo string `json:"redirectTo"`
	}{
		RedirectTo: redirectTo,
	})
}
