package site_test

import (
	"context"
	"encoding/json"
	"errors"
	"hostamat/internal/site"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestListSiteHandler_ServeHTTP(t *testing.T) {
	t.Run("invalid method", func(t *testing.T) {
		handler := site.NewListSiteHandler(nil)
		req := httptest.NewRequest("POST", "/", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
		}
	})

	t.Run("service error", func(t *testing.T) {
		handler := site.NewListSiteHandler(&mockServiceList{
			getListFn: func(context.Context) ([]*site.Site, error) {
				return nil, errors.New("service error")
			},
		})

		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})

	t.Run("success", func(t *testing.T) {
		handler := site.NewListSiteHandler(&mockServiceList{
			getListFn: func(_ context.Context) ([]*site.Site, error) {
				return []*site.Site{
					{
						Id:        "1",
						Subdomain: "site1",
						Domain:    "site1.example.com",
						CreatedAt: time.Now(),
					},
					{
						Id:        "2",
						Subdomain: "site2",
						Domain:    "site2.example.com",
						CreatedAt: time.Now(),
					},
				}, nil
			},
		})

		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		if w.Header().Get("Content-Type") != "application/json" {
			t.Errorf("Expected content type %s, got %s", "application/json", w.Header().Get("Content-Type"))
		}

		var response []site.Site
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if len(response) != 2 {
			t.Fatalf("Expected 2 sites, got %d", len(response))
		}

		if response[0].Id != "1" || response[0].Subdomain != "site1" || response[0].Domain != "site1.example.com" {
			t.Errorf("Expected site1 with id=1, subdomain=site1, domain=site1.example.com, got %+v", response[0])
		}
		if response[1].Id != "2" || response[1].Subdomain != "site2" || response[1].Domain != "site2.example.com" {
			t.Errorf("Expected site2 with id=2, subdomain=site2, domain=site2.example.com, got %+v", response[1])
		}
	})
}

type mockServiceList struct {
	getListFn func(context.Context) ([]*site.Site, error)
}

func (m *mockServiceList) GetList(ctx context.Context) ([]*site.Site, error) {
	return m.getListFn(ctx)
}
