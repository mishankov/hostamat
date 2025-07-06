package uploadsession_test

import (
	"errors"
	"sync"

	"hostamat/internal/uploadsession"
)

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
