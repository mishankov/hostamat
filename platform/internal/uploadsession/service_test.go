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
		mockRepo := &errorProneMockRepository{
			shouldError: true,
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

type errorProneMockRepository struct {
	mockRepository
	shouldError bool
}

func (m *errorProneMockRepository) Create(session *uploadsession.UploadSession) error {
	if m.shouldError {
		return errors.New("forced repository error")
	}
	return m.mockRepository.Create(session)
}

func (m *errorProneMockRepository) Update(session *uploadsession.UploadSession) error {
	if m.shouldError {
		return errors.New("forced repository error")
	}
	return m.mockRepository.Update(session)
}

func TestService_Cancel(t *testing.T) {
	t.Run("successful cancellation", func(t *testing.T) {
		mockRepo := newMockRepository()
		service := uploadsession.NewService(mockRepo)

		// Create test session
		testSession := &uploadsession.UploadSession{
			Id:        "test-cancel-id",
			Status:    uploadsession.Open,
			CreatedAt: time.Now(),
			ExpiresAt: time.Now().Add(time.Hour),
		}
		if err := mockRepo.Create(testSession); err != nil {
			t.Fatalf("Failed to setup test session: %v", err)
		}

		// Test Cancel
		err := service.Cancel("test-cancel-id")
		if err != nil {
			t.Fatalf("Cancel failed: %v", err)
		}

		// Verify status was updated
		updated, err := mockRepo.Get("test-cancel-id")
		if err != nil {
			t.Fatalf("Failed to get updated session: %v", err)
		}
		if updated.Status != uploadsession.Cancelled {
			t.Errorf("Expected status %q, got %q", uploadsession.Cancelled, updated.Status)
		}
		if updated.ResolvedAt.IsZero() {
			t.Error("Expected ResolvedAt to be set")
		}
	})

	t.Run("non-existent session", func(t *testing.T) {
		mockRepo := newMockRepository()
		service := uploadsession.NewService(mockRepo)

		err := service.Cancel("non-existent-id")
		if err == nil {
			t.Error("Expected error for non-existent session")
		}
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := &errorProneMockRepository{
			shouldError: true,
		}
		service := uploadsession.NewService(mockRepo)

		err := service.Cancel("any-id")
		if err == nil {
			t.Error("Expected error from repository")
		}
	})
}

func TestService_Finish(t *testing.T) {
	t.Run("successful finish", func(t *testing.T) {
		mockRepo := newMockRepository()
		service := uploadsession.NewService(mockRepo)

		// Create test session
		testSession := &uploadsession.UploadSession{
			Id:        "test-finish-id",
			Status:    uploadsession.Open,
			CreatedAt: time.Now(),
			ExpiresAt: time.Now().Add(time.Hour),
		}
		if err := mockRepo.Create(testSession); err != nil {
			t.Fatalf("Failed to setup test session: %v", err)
		}

		// Test Finish
		err := service.Finish("test-finish-id")
		if err != nil {
			t.Fatalf("Finish failed: %v", err)
		}

		// Verify status was updated
		updated, err := mockRepo.Get("test-finish-id")
		if err != nil {
			t.Fatalf("Failed to get updated session: %v", err)
		}
		if updated.Status != uploadsession.Finished {
			t.Errorf("Expected status %q, got %q", uploadsession.Finished, updated.Status)
		}
		if updated.ResolvedAt.IsZero() {
			t.Error("Expected ResolvedAt to be set")
		}
	})

	t.Run("non-existent session", func(t *testing.T) {
		mockRepo := newMockRepository()
		service := uploadsession.NewService(mockRepo)

		err := service.Finish("non-existent-id")
		if err == nil {
			t.Error("Expected error for non-existent session")
		}
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := &errorProneMockRepository{
			shouldError: true,
		}
		service := uploadsession.NewService(mockRepo)

		err := service.Finish("any-id")
		if err == nil {
			t.Error("Expected error from repository")
		}
	})
}

func TestService_Get(t *testing.T) {
	t.Run("existing session", func(t *testing.T) {
		mockRepo := newMockRepository()
		service := uploadsession.NewService(mockRepo)

		// Create test session directly in mock
		testSession := &uploadsession.UploadSession{
			Id:        "test-get-id",
			CreatedAt: time.Now(),
			Status:    uploadsession.Open,
			ExpiresAt: time.Now().Add(time.Hour),
		}
		if err := mockRepo.Create(testSession); err != nil {
			t.Fatalf("Failed to setup test session: %v", err)
		}

		// Test Get
		got, err := service.Get("test-get-id")
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}

		if got.Id != testSession.Id {
			t.Errorf("Expected ID %q, got %q", testSession.Id, got.Id)
		}
		if got.Status != testSession.Status {
			t.Errorf("Expected status %q, got %q", testSession.Status, got.Status)
		}
	})

	t.Run("non-existent session", func(t *testing.T) {
		mockRepo := newMockRepository()
		service := uploadsession.NewService(mockRepo)

		_, err := service.Get("non-existent-id")
		if err == nil {
			t.Error("Expected error for non-existent session")
		}
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := &errorProneMockRepository{
			shouldError: true,
		}

		service := uploadsession.NewService(mockRepo)
		_, err := service.Get("any-id")
		if err == nil {
			t.Error("Expected error from repository")
		}
	})
}
