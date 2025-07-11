package site_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"hostamat/internal/site"
)

func TestCreateSiteHandler_ServeHTTP(t *testing.T) {
	t.Run("invalid method", func(t *testing.T) {
		handler := site.NewCreateSiteHandler(&mockServiceCreate{
			createFn: func(subdomain string) (*site.Site, error) {
				return &site.Site{
					Id:        "1",
					Subdomain: subdomain,
					Domain:    subdomain + ".example.com",
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

	t.Run("invalid request body", func(t *testing.T) {
		handler := site.NewCreateSiteHandler(nil)
		req := httptest.NewRequest("POST", "/", strings.NewReader("invalid json"))
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("service error", func(t *testing.T) {
		handler := site.NewCreateSiteHandler(&mockServiceCreate{
			createFn: func(domain string) (*site.Site, error) {
				return nil, errors.New("service error")
			},
		})

		req := httptest.NewRequest("POST", "/", strings.NewReader(`{"subdomain":"example"}`))
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})

	t.Run("success", func(t *testing.T) {
		handler := site.NewCreateSiteHandler(&mockServiceCreate{
			createFn: func(domain string) (*site.Site, error) {
				return &site.Site{
					Id:        "1",
					Subdomain: domain,
				}, nil
			},
		})

		req := httptest.NewRequest("POST", "/", strings.NewReader(`{"subdomain":"example"}`))
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
		}

		if w.Header().Get("Content-Type") != "application/json" {
			t.Errorf("Expected content type %s, got %s", "application/json", w.Header().Get("Content-Type"))
		}

		if !strings.Contains(w.Body.String(), `"id":"1"`) {
			t.Errorf("Expected body to contain %s, got %s", `"id":"1"`, w.Body.String())
		}
		if !strings.Contains(w.Body.String(), `"subdomain":"example"`) {
			t.Errorf("Expected body to contain %s, got %s", `"subdomain":"example"`, w.Body.String())
		}
	})
}

type mockServiceCreate struct {
	createFn func(string) (*site.Site, error)
}

func (m *mockServiceCreate) Create(domain string) (*site.Site, error) {
	return m.createFn(domain)
}
