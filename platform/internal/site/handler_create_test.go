package site_test

import (
	"encoding/json"
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
					Domain:    domain + ".example.com",
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

		var response site.Site
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if response.Id != "1" {
			t.Errorf("Expected id %q, got %q", "1", response.Id)
		}
		if response.Subdomain != "example" {
			t.Errorf("Expected subdomain %q, got %q", "example", response.Subdomain)
		}
		if response.Domain != "example.example.com" {
			t.Errorf("Expected domain %q, got %q", "example.example.com", response.Domain)
		}
	})
}

type mockServiceCreate struct {
	createFn func(string) (*site.Site, error)
}

func (m *mockServiceCreate) Create(domain string) (*site.Site, error) {
	return m.createFn(domain)
}
