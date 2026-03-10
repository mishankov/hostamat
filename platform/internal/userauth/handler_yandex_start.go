package userauth

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type startYandexAuthService interface {
	StartYandexAuth() (string, string, error)
	OAuthStateCookieName() string
	OAuthStateCookiePath() string
}

type YandexStartHandler struct {
	service startYandexAuthService
}

func NewYandexStartHandler(service startYandexAuthService) *YandexStartHandler {
	return &YandexStartHandler{service: service}
}

func (h *YandexStartHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	redirectURL, state, err := h.service.StartYandexAuth()
	if err != nil {
		http.Error(w, "failed to start oauth flow", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     h.service.OAuthStateCookieName(),
		Value:    state,
		Path:     h.service.OAuthStateCookiePath(),
		Expires:  time.Now().Add(10 * time.Minute),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   requestIsSecure(r),
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(struct {
		AuthURL string `json:"authUrl"`
	}{
		AuthURL: redirectURL,
	})
}

func requestIsSecure(r *http.Request) bool {
	if r.TLS != nil {
		return true
	}

	return strings.EqualFold(strings.TrimSpace(strings.Split(r.Header.Get("X-Forwarded-Proto"), ",")[0]), "https")
}
