package uploadsession_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"hostamat/internal/uploadsession"
)

func TestService_Create(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		mockRepo := newMockRepository()
		service := uploadsession.NewService(mockRepo)

		session, err := service.Create()
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		if session.Id == "" {
			t.Error("Expected non-empty session ID")
		}
		if session.Status != uploadsession.Open {
			t.Errorf("Expected status %q, got %q", uploadsession.Open, session.Status)
		}
		if time.Until(session.ExpiresAt) < time.Minute*4 || time.Until(session.ExpiresAt) > time.Minute*5 {
			t.Errorf("Expected expires_at to be ~5 minutes from now, got %v", session.ExpiresAt)
		}

		// Verify the session was stored
		_, err = mockRepo.Get(session.Id)
		if err != nil {
			t.Errorf("Session not found in repository: %v", err)
		}
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := &mockRepository{
			sessions: make(map[string]*uploadsession.UploadSession),
		}
		// Force Create to fail
		mockRepo.Create = func(*uploadsession.UploadSession) error {
			return errors.New("forced repository error")
		}

		service := uploadsession.NewService(mockRepo)
		_, err := service.Create()
		if err == nil {
			t.Error("Expected error from repository")
		}
	})
}

type mockRepository struct {
	mu       sync.Mutex
	sessions map[string]*uploadsession.UploadSession
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		sessions: make(map[string]*uploadsession.UploadSession),
	}
}

func (m *mockRepository) Get(id string) (*uploadsession.UploadSession, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, exists := m.sessions[id]
	if !exists {
		return nil, errors.New("not found")
	}
	return session, nil
}

func (m *mockRepository) Create(session *uploadsession.UploadSession) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.sessions[session.Id]; exists {
		return errors.New("already exists")
	}
	m.sessions[session.Id] = session
	return nil
}

func (m *mockRepository) Update(session *uploadsession.UploadSession) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.sessions[session.Id]; !exists {
		return errors.New("not found")
	}
	m.sessions[session.Id] = session
	return nil
}
