package uploadsession_test

import (
	"errors"
	"hostamat/internal/uploadsession"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateUploadSessionHandler_ServeHTTP(t *testing.T) {
	t.Run("invalid method", func(t *testing.T) {
		handler := uploadsession.NewCreateUploadSessionHandler(nil)
		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
		}
	})

	t.Run("service error", func(t *testing.T) {
		handler := uploadsession.NewCreateUploadSessionHandler(&mockService{
			createFn: func() (*uploadsession.UploadSession, error) {
				return nil, errors.New("service error")
			},
		})

		req := httptest.NewRequest("POST", "/", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})

	t.Run("success", func(t *testing.T) {
		handler := uploadsession.NewCreateUploadSessionHandler(&mockService{
			createFn: func() (*uploadsession.UploadSession, error) {
				return &uploadsession.UploadSession{
					Id: "1",
				}, nil
			},
		})

		req := httptest.NewRequest("POST", "/", nil)
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
	})
}

type mockService struct {
	createFn func() (*uploadsession.UploadSession, error)
	finishFn func(string) error
}

func (m *mockService) Create() (*uploadsession.UploadSession, error) {
	return m.createFn()
}

func (m *mockService) Finish(id string) error {
	if m.finishFn != nil {
		return m.finishFn(id)
	}
	return nil
}
