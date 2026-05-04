package feedback_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"hostamat/internal/feedback"

	"github.com/platforma-dev/platforma/auth"
)

type createServiceMock struct {
	item    *feedback.Feedback
	err     error
	message string
	url     string
}

func (m *createServiceMock) Create(ctx context.Context, message, url string) (*feedback.Feedback, error) {
	m.message = message
	m.url = url
	if m.err != nil {
		return nil, m.err
	}
	return m.item, nil
}

func TestCreateHandlerReturnsCreated(t *testing.T) {
	t.Parallel()

	service := &createServiceMock{
		item: &feedback.Feedback{ID: "feedback-1"},
	}
	handler := feedback.NewCreateHandler(service)

	req := httptest.NewRequest(http.MethodPost, "/api/feedback", strings.NewReader(`{"message":"hello","url":"/dashboard/sites"}`))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d, body = %s", rec.Code, http.StatusCreated, rec.Body.String())
	}
	if service.message != "hello" || service.url != "/dashboard/sites" {
		t.Fatalf("service args = %q, %q", service.message, service.url)
	}
	if !strings.Contains(rec.Body.String(), `"id":"feedback-1"`) {
		t.Fatalf("body = %s, want feedback id", rec.Body.String())
	}
}

func TestCreateHandlerRejectsInvalidBody(t *testing.T) {
	t.Parallel()

	handler := feedback.NewCreateHandler(&createServiceMock{})
	req := httptest.NewRequest(http.MethodPost, "/api/feedback", bytes.NewBufferString(`{`))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if !strings.Contains(rec.Body.String(), "INVALID_BODY") {
		t.Fatalf("body = %s, want INVALID_BODY", rec.Body.String())
	}
}

func TestCreateHandlerRejectsUnauthenticated(t *testing.T) {
	t.Parallel()

	handler := feedback.NewCreateHandler(&createServiceMock{err: auth.ErrUserNotFound})
	req := httptest.NewRequest(http.MethodPost, "/api/feedback", strings.NewReader(`{"message":"hello"}`))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
	if !strings.Contains(rec.Body.String(), "UNAUTHORIZED") {
		t.Fatalf("body = %s, want UNAUTHORIZED", rec.Body.String())
	}
}

func TestCreateHandlerReturnsServerError(t *testing.T) {
	t.Parallel()

	handler := feedback.NewCreateHandler(&createServiceMock{err: errors.New("db down")})
	req := httptest.NewRequest(http.MethodPost, "/api/feedback", strings.NewReader(`{"message":"hello"}`))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

func TestCreateHandlerCapsRequestBody(t *testing.T) {
	t.Parallel()

	handler := feedback.NewCreateHandler(&createServiceMock{})
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/feedback",
		strings.NewReader(`{"message":"`+strings.Repeat("a", feedback.MaxRequestBodyBytes)+`"}`),
	)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}
