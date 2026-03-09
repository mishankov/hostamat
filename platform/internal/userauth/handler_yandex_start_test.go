package userauth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockStartYandexAuthService struct {
	redirectURL string
	state       string
	err         error
}

func (m *mockStartYandexAuthService) StartYandexAuth() (string, string, error) {
	return m.redirectURL, m.state, m.err
}

func (m *mockStartYandexAuthService) OAuthStateCookieName() string {
	return oauthStateCookieName
}

func (m *mockStartYandexAuthService) OAuthStateCookiePath() string {
	return "/api/auth/yandex/callback"
}

func TestYandexStartHandler_ReturnsAuthURLAndSetsCookie(t *testing.T) {
	handler := NewYandexStartHandler(&mockStartYandexAuthService{
		redirectURL: "https://oauth.yandex.ru/authorize?state=test-state",
		state:       "test-state",
	})

	req := httptest.NewRequest(http.MethodGet, "/auth/yandex/start", nil)
	req.Header.Set("X-Forwarded-Proto", "https")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	resp := rec.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var payload struct {
		AuthURL string `json:"authUrl"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if payload.AuthURL != "https://oauth.yandex.ru/authorize?state=test-state" {
		t.Fatalf("unexpected auth url: %s", payload.AuthURL)
	}

	cookies := resp.Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}

	if cookies[0].Name != oauthStateCookieName {
		t.Fatalf("unexpected cookie name: %s", cookies[0].Name)
	}

	if cookies[0].Path != "/api/auth/yandex/callback" {
		t.Fatalf("unexpected cookie path: %s", cookies[0].Path)
	}

	if !cookies[0].Secure {
		t.Fatalf("expected secure cookie")
	}
}
