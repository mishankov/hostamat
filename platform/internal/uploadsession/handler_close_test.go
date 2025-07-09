package uploadsession_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"hostamat/internal/uploadsession"
)

func TestCloseUploadSessionHandler_ServeHTTP(t *testing.T) {
	t.Run("invalid method", func(t *testing.T) {
		handler := uploadsession.NewCloseUploadSessionHandler(nil)
		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
		}
	})

	t.Run("missing header", func(t *testing.T) {
		handler := uploadsession.NewCloseUploadSessionHandler(nil)
		req := httptest.NewRequest("POST", "/", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("service error", func(t *testing.T) {
		handler := uploadsession.NewCloseUploadSessionHandler(&mockService{
			finishFn: func(id string) error {
				return errors.New("service error")
			},
		})

		req := httptest.NewRequest("POST", "/", nil)
		req.Header.Set("H-Upload-Session", "test-id")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})

	t.Run("success", func(t *testing.T) {
		handler := uploadsession.NewCloseUploadSessionHandler(&mockService{
			finishFn: func(id string) error {
				return nil
			},
		})

		req := httptest.NewRequest("POST", "/", nil)
		req.Header.Set("H-Upload-Session", "test-id")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		if w.Header().Get("Content-Type") != "application/json" {
			t.Errorf("Expected content type %s, got %s", "application/json", w.Header().Get("Content-Type"))
		}

		if !strings.Contains(w.Body.String(), `"status":"finished"`) {
			t.Errorf("Expected body to contain %s, got %s", `"status":"finished"`, w.Body.String())
		}
	})
}

type mockService struct {
	finishFn func(string) error
}

func (m *mockService) Finish(id string) error {
	return m.finishFn(id)
}
