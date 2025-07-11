package file_test

import (
	"context"
	"errors"
	"hostamat/internal/file"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServeHandler_ServeHTTP(t *testing.T) {
	tests := []struct {
		name         string
		host         string
		path         string
		mockResponse io.Reader
		mockError    error
		expectedCode int
		expectedType string
		expectedBody string
	}{
		{
			name:         "successful html file",
			host:         "test.example.com",
			path:         "/index.html",
			mockResponse: strings.NewReader("<html></html>"),
			expectedCode: http.StatusOK,
			expectedType: "text/html",
			expectedBody: "<html></html>",
		},
		{
			name:         "root path redirects to index.html",
			host:         "test.example.com",
			path:         "/",
			mockResponse: strings.NewReader("<html></html>"),
			expectedCode: http.StatusOK,
			expectedType: "text/html",
		},
		{
			name:         "css file",
			host:         "test.example.com",
			path:         "/style.css",
			mockResponse: strings.NewReader("body{}"),
			expectedCode: http.StatusOK,
			expectedType: "text/css",
		},
		{
			name:         "file not found",
			host:         "test.example.com",
			path:         "/missing.html",
			mockError:    errors.New("not found"),
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &mockServeService{
				getFileFunc: func(ctx context.Context, fileName string, domain string) (io.Reader, error) {
					if tt.mockError != nil {
						return nil, tt.mockError
					}
					return tt.mockResponse, nil
				},
			}

			handler := file.NewServeHandler(service)
			req := httptest.NewRequest("GET", tt.path, nil)
			req.Host = tt.host
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			resp := w.Result()
			if resp.StatusCode != tt.expectedCode {
				t.Errorf("expected status %d, got %d", tt.expectedCode, resp.StatusCode)
			}

			if tt.expectedType != "" {
				contentType := resp.Header.Get("Content-Type")
				if contentType != tt.expectedType {
					t.Errorf("expected content type %s, got %s", tt.expectedType, contentType)
				}
			}

			if tt.expectedBody != "" {
				body, _ := io.ReadAll(resp.Body)
				if string(body) != tt.expectedBody {
					t.Errorf("expected body %s, got %s", tt.expectedBody, string(body))
				}
			}
		})
	}
}

type mockServeService struct {
	getFileFunc func(ctx context.Context, fileName string, domain string) (io.Reader, error)
}

func (m *mockServeService) GetFile(ctx context.Context, fileName string, domain string) (io.Reader, error) {
	return m.getFileFunc(ctx, fileName, domain)
}
