package userauth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockCallbackYandexAuthService struct {
	sessionID string
	err       error
}

func (m *mockCallbackYandexAuthService) CreateSessionFromYandexCode(_ context.Context, _ string) (string, error) {
	if m.err != nil {
		return "", m.err
	}

	return m.sessionID, nil
}

func (m *mockCallbackYandexAuthService) CookieName() string {
	return "hsessionid"
}

func (m *mockCallbackYandexAuthService) OAuthStateCookieName() string {
	return oauthStateCookieName
}

func (m *mockCallbackYandexAuthService) OAuthStateCookiePath() string {
	return "/api/auth/yandex/callback"
}

func TestYandexCallbackHandler_ReturnsDeniedRedirectTarget(t *testing.T) {
	handler := NewYandexCallbackHandler(&mockCallbackYandexAuthService{})

	req := httptest.NewRequest(http.MethodGet, "/auth/yandex/callback?error=access_denied", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	resp := rec.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var payload struct {
		RedirectTo string `json:"redirectTo"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if payload.RedirectTo != "/login?error=oauth_denied" {
		t.Fatalf("unexpected redirect target: %s", payload.RedirectTo)
	}
}

func TestYandexCallbackHandler_ReturnsFailedRedirectTargetForInvalidState(t *testing.T) {
	handler := NewYandexCallbackHandler(&mockCallbackYandexAuthService{})

	req := httptest.NewRequest(http.MethodGet, "/auth/yandex/callback?state=wrong&code=abc", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	resp := rec.Result()
	var payload struct {
		RedirectTo string `json:"redirectTo"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if payload.RedirectTo != "/login?error=oauth_failed" {
		t.Fatalf("unexpected redirect target: %s", payload.RedirectTo)
	}
}

func TestYandexCallbackHandler_ReturnsDashboardRedirectAndSessionCookie(t *testing.T) {
	handler := NewYandexCallbackHandler(&mockCallbackYandexAuthService{sessionID: "session-1"})

	req := httptest.NewRequest(http.MethodGet, "/auth/yandex/callback?state=ok&code=abc", nil)
	req.AddCookie(&http.Cookie{Name: oauthStateCookieName, Value: "ok"})
	req.Header.Set("X-Forwarded-Proto", "https")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	resp := rec.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var payload struct {
		RedirectTo string `json:"redirectTo"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if payload.RedirectTo != "/dashboard" {
		t.Fatalf("unexpected redirect target: %s", payload.RedirectTo)
	}

	cookies := resp.Cookies()
	if len(cookies) != 2 {
		t.Fatalf("expected 2 cookies, got %d", len(cookies))
	}

	foundSessionCookie := false
	for _, cookie := range cookies {
		if cookie.Name == "hsessionid" && cookie.Value == "session-1" && cookie.Secure {
			foundSessionCookie = true
		}
	}

	if !foundSessionCookie {
		t.Fatalf("expected secure session cookie to be set")
	}
}

func TestYandexCallbackHandler_ReturnsFailedRedirectTargetWhenExchangeFails(t *testing.T) {
	handler := NewYandexCallbackHandler(&mockCallbackYandexAuthService{err: errors.New("boom")})

	req := httptest.NewRequest(http.MethodGet, "/auth/yandex/callback?state=ok&code=abc", nil)
	req.AddCookie(&http.Cookie{Name: oauthStateCookieName, Value: "ok"})
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	resp := rec.Result()
	var payload struct {
		RedirectTo string `json:"redirectTo"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if payload.RedirectTo != "/login?error=oauth_failed" {
		t.Fatalf("unexpected redirect target: %s", payload.RedirectTo)
	}
}
