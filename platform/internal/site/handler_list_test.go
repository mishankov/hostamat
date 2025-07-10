package site_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"hostamat/internal/site"
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
			getListFn: func() ([]*site.Site, error) {
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
			getListFn: func() ([]*site.Site, error) {
				return []*site.Site{
					{
						Id:        "1",
						Subdomain: "site1",
					},
					{
						Id:        "2",
						Subdomain: "site2",
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

		body := w.Body.String()
		if !strings.Contains(body, `"id":"1"`) || !strings.Contains(body, `"subdomain":"site1"`) {
			t.Errorf("Expected body to contain site1, got %s", body)
		}
		if !strings.Contains(body, `"id":"2"`) || !strings.Contains(body, `"subdomain":"site2"`) {
			t.Errorf("Expected body to contain site2, got %s", body)
		}
	})
}

type mockServiceList struct {
	getListFn func() ([]*site.Site, error)
}

func (m *mockServiceList) GetList() ([]*site.Site, error) {
	return m.getListFn()
}
