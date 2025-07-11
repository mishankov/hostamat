package file_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"platform/internal/file"
	"strings"
	"testing"
)

type mockUploadService struct {
	uploadFunc func(ctx context.Context, filename string, file io.Reader, size int64, siteId string) ([]string, []string, error)
}

func (m *mockUploadService) UploadFile(ctx context.Context, filename string, file io.Reader, size int64, siteId string) ([]string, []string, error) {
	return m.uploadFunc(ctx, filename, file, size, siteId)
}

func TestUploadHandler_ServeHTTP(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		siteIdHeader   string
		mockResponse   []string
		mockSkipped    []string
		mockError      error
		expectedCode   int
		expectedBody   map[string][]string
	}{
		{
			name:          "successful upload",
			method:        http.MethodPost,
			siteIdHeader:  "test-site",
			mockResponse:  []string{"file1.txt"},
			mockSkipped:   []string{"file2.txt"},
			expectedCode:  http.StatusOK,
			expectedBody: map[string][]string{
				"uploadedFiles": {"file1.txt"},
				"skippedFiles":  {"file2.txt"},
			},
		},
		{
			name:          "wrong method",
			method:        http.MethodGet,
			expectedCode:  http.StatusMethodNotAllowed,
		},
		{
			name:          "missing site id",
			method:        http.MethodPost,
			expectedCode:  http.StatusBadRequest,
		},
		{
			name:          "upload error",
			method:        http.MethodPost,
			siteIdHeader:  "test-site",
			mockError:     errors.New("upload failed"),
			expectedCode:  http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &mockUploadService{
				uploadFunc: func(ctx context.Context, filename string, file io.Reader, size int64, siteId string) ([]string, []string, error) {
					if tt.mockError != nil {
						return nil, nil, tt.mockError
					}
					return tt.mockResponse, tt.mockSkipped, nil
				},
			}

			handler := file.NewUploadHandler(service)
			
			// Create multipart form with test file
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			part, _ := writer.CreateFormFile("file", "test.txt")
			part.Write([]byte("test content"))
			writer.Close()

			req := httptest.NewRequest(tt.method, "/upload", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			if tt.siteIdHeader != "" {
				req.Header.Set("H-Site-Id", tt.siteIdHeader)
			}

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			resp := w.Result()
			if resp.StatusCode != tt.expectedCode {
				t.Errorf("expected status %d, got %d", tt.expectedCode, resp.StatusCode)
			}

			if tt.expectedBody != nil {
				var result struct {
					UploadedFiles []string `json:"uploadedFiles"`
					SkippedFiles  []string `json:"skippedFiles"`
				}
				if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
					t.Errorf("failed to decode response: %v", err)
				}

				if !equalStringSlices(result.UploadedFiles, tt.expectedBody["uploadedFiles"]) {
					t.Errorf("expected uploaded files %v, got %v", tt.expectedBody["uploadedFiles"], result.UploadedFiles)
				}

				if !equalStringSlices(result.SkippedFiles, tt.expectedBody["skippedFiles"]) {
					t.Errorf("expected skipped files %v, got %v", tt.expectedBody["skippedFiles"], result.SkippedFiles)
				}
			}
		})
	}
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
