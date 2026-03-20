package userauth

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/platforma-dev/platforma/log"
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
	ev := log.EventFromContext(r.Context())

	ev.AddAttrs(map[string]any{
		"auth.provider": "yandex",
	})

	if r.Method != http.MethodGet {

		ev.AddStep(slog.LevelWarn, "userauth.YandexStart.method_not_allowed")

		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	redirectURL, state, err := h.service.StartYandexAuth()
	if err != nil {

		ev.AddError(err)

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
	err = json.NewEncoder(w).Encode(struct {
		AuthURL string `json:"authUrl"`
	}{
		AuthURL: redirectURL,
	})
	if err != nil {
		ev.AddError(err)
	}
}

func requestIsSecure(r *http.Request) bool {
	if r.TLS != nil {
		return true
	}

	return strings.EqualFold(strings.TrimSpace(strings.Split(r.Header.Get("X-Forwarded-Proto"), ",")[0]), "https")
}
