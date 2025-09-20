package site_test

import (
	"context"
	"encoding/json"
	"errors"
	"hostamat/internal/site"
	"hostamat/pkg/platforma/batteries/auth"
	"hostamat/pkg/platforma/batteries/user"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateSiteHandler_ServeHTTP(t *testing.T) {
	t.Run("invalid method", func(t *testing.T) {
		handler := site.NewCreateSiteHandler(&mockServiceCreate{
			createFn: func(_ context.Context) (*site.Site, error) {
				return &site.Site{
					Id:        "1",
					Subdomain: "my",
					Domain:    "my.example.com",
				}, nil
			},
		})
		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
		}
	})

	t.Run("service error", func(t *testing.T) {
		handler := site.NewCreateSiteHandler(&mockServiceCreate{
			createFn: func(_ context.Context) (*site.Site, error) {
				return nil, errors.New("service error")
			},
		})

		req := withUserRequest("POST", "/", `{"subdomain":"example"}`, "test-user")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})

	t.Run("success", func(t *testing.T) {
		handler := site.NewCreateSiteHandler(&mockServiceCreate{
			createFn: func(_ context.Context) (*site.Site, error) {
				return &site.Site{
					Id:        "1",
					User:      "test-user",
					Subdomain: "my",
					Domain:    "my.example.com",
				}, nil
			},
		})

		req := withUserRequest("POST", "/", `{"subdomain":"example"}`, "test-user")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
		}

		if w.Header().Get("Content-Type") != "application/json" {
			t.Errorf("Expected content type %s, got %s", "application/json", w.Header().Get("Content-Type"))
		}

		var response site.Site
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if response.Id != "1" {
			t.Errorf("Expected id %q, got %q", "1", response.Id)
		}
		if response.User != "test-user" {
			t.Errorf("Expected user %q, got %q", "test-user", response.User)
		}
		if response.Subdomain != "my" {
			t.Errorf("Expected subdomain %q, got %q", "my", response.Subdomain)
		}
		if response.Domain != "my.example.com" {
			t.Errorf("Expected domain %q, got %q", "my.example.com", response.Domain)
		}
	})
}

func withUserRequest(method, url, body string, userID string) *http.Request {
	req := httptest.NewRequest(method, url, strings.NewReader(body))
	ctx := context.WithValue(req.Context(), auth.UserContextKey, &user.User{
		ID: userID,
	})
	return req.WithContext(ctx)
}

type mockServiceCreate struct {
	createFn func(context.Context) (*site.Site, error)
}

func (m *mockServiceCreate) Create(ctx context.Context) (*site.Site, error) {
	return m.createFn(ctx)
}
