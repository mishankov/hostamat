package file_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"hostamat/internal/file"
)

type mockDeleteService struct {
	deleteFunc func(ctx context.Context, folder string) error
}

func (m *mockDeleteService) DeleteFromFolder(ctx context.Context, folder string) error {
	return m.deleteFunc(ctx, folder)
}

func TestDeleteHandler_ServeHTTP(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		folderParam  string
		mockError    error
		expectedCode int
	}{
		{
			name:         "successful delete",
			method:       http.MethodDelete,
			folderParam:  "test-folder",
			expectedCode: http.StatusNoContent,
		},
		{
			name:         "wrong method",
			method:       http.MethodGet,
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "missing folder param",
			method:       http.MethodDelete,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "delete error",
			method:       http.MethodDelete,
			folderParam:  "test-folder",
			mockError:    errors.New("delete failed"),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &mockDeleteService{
				deleteFunc: func(ctx context.Context, folder string) error {
					if tt.mockError != nil {
						return tt.mockError
					}
					return nil
				},
			}

			handler := file.NewDeleteHandler(service)
			req := httptest.NewRequest(tt.method, "/delete?folder="+tt.folderParam, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			resp := w.Result()
			if resp.StatusCode != tt.expectedCode {
				t.Errorf("expected status %d, got %d", tt.expectedCode, resp.StatusCode)
			}
		})
	}
}
