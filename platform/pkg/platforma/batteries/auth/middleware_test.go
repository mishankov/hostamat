package auth_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"hostamat/pkg/platforma/batteries/auth"
	"hostamat/pkg/platforma/batteries/session"
	"hostamat/pkg/platforma/batteries/user"
)

func TestAuthenticationMiddleware_ValidSession(t *testing.T) {
	sessionSvc := &mockSessionService{
		sessions: map[string]*session.Session{
			"valid-session-id": {ID: "valid-session-id", User: "user-id"},
		},
	}
	userSvc := &mockUserService{
		users: map[string]*user.User{
			"user-id": {ID: "user-id", Username: "testuser"},
		},
	}
	middleware := auth.NewAuthenticationMiddleware(sessionSvc, userSvc)

	handler := middleware.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "valid-session-id"})
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestAuthenticationMiddleware_NoSessionCookie(t *testing.T) {
	sessionSvc := &mockSessionService{}
	userSvc := &mockUserService{}
	middleware := auth.NewAuthenticationMiddleware(sessionSvc, userSvc)

	handler := middleware.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called when authentication fails")
	}))

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}
}

func TestAuthenticationMiddleware_InvalidSession(t *testing.T) {
	sessionSvc := &mockSessionService{
		sessions: map[string]*session.Session{},
	}
	userSvc := &mockUserService{}
	middleware := auth.NewAuthenticationMiddleware(sessionSvc, userSvc)

	handler := middleware.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called when authentication fails")
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "invalid-session-id"})
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}
}

func TestAuthenticationMiddleware_SessionServiceError(t *testing.T) {
	sessionSvc := &mockSessionService{
		error: errors.New("database error"),
	}
	userSvc := &mockUserService{}
	middleware := auth.NewAuthenticationMiddleware(sessionSvc, userSvc)

	handler := middleware.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called when authentication fails")
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "session-id"})
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}
}

func TestAuthenticationMiddleware_UserServiceError(t *testing.T) {
	sessionSvc := &mockSessionService{
		sessions: map[string]*session.Session{
			"valid-session-id": {ID: "valid-session-id", User: "user-id"},
		},
	}
	userSvc := &mockUserService{
		error: errors.New("user not found"),
	}
	middleware := auth.NewAuthenticationMiddleware(sessionSvc, userSvc)

	handler := middleware.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called when authentication fails")
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "valid-session-id"})
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}
}

type mockSessionService struct {
	sessions map[string]*session.Session
	error    error
}

func (m *mockSessionService) Get(ctx context.Context, id string) (*session.Session, error) {
	if m.error != nil {
		return nil, m.error
	}
	if session, ok := m.sessions[id]; ok {
		return session, nil
	}
	return nil, errors.New("session not found")
}

func (m *mockSessionService) Create(ctx context.Context, session *session.Session) error {
	return m.error
}

func (m *mockSessionService) Delete(ctx context.Context, id string) error {
	return m.error
}

type mockUserService struct {
	users map[string]*user.User
	error error
}

func (m *mockUserService) Get(ctx context.Context, id string) (*user.User, error) {
	if m.error != nil {
		return nil, m.error
	}
	if user, ok := m.users[id]; ok {
		return user, nil
	}
	return nil, errors.New("user not found")
}
