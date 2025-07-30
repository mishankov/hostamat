package uploadsession_test

import (
	"context"
	"errors"
	"testing"

	"hostamat/internal/uploadsession"

	"github.com/google/uuid"
)

type mockRepo struct {
	sessions map[string]*uploadsession.UploadSession
}

func (m *mockRepo) Create(ctx context.Context, session *uploadsession.UploadSession) error {
	if _, exists := m.sessions[session.Id]; exists {
		return errors.New("duplicate id")
	}
	m.sessions[session.Id] = session
	return nil
}

func (m *mockRepo) Get(ctx context.Context, id string) (*uploadsession.UploadSession, error) {
	session, exists := m.sessions[id]
	if !exists {
		return nil, errors.New("not found")
	}
	return session, nil
}

func (m *mockRepo) Update(ctx context.Context, session *uploadsession.UploadSession) error {
	if _, exists := m.sessions[session.Id]; !exists {
		return errors.New("not found")
	}
	m.sessions[session.Id] = session
	return nil
}

func TestService_Create(t *testing.T) {
	repo := &mockRepo{sessions: make(map[string]*uploadsession.UploadSession)}
	service := uploadsession.NewService(repo)

	t.Run("success", func(t *testing.T) {
		siteId := "test-site-id"
		session, err := service.Create(context.Background(), siteId)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		if session.SiteId != siteId {
			t.Errorf("expected siteId %q, got %q", siteId, session.SiteId)
		}
		if session.Status != uploadsession.StatusOpen {
			t.Errorf("expected status %q, got %q", uploadsession.StatusOpen, session.Status)
		}
		if _, err := uuid.Parse(session.Id); err != nil {
			t.Errorf("expected valid UUID, got %q", session.Id)
		}
	})

	t.Run("empty siteId", func(t *testing.T) {
		_, err := service.Create(context.Background(), "")
		if err == nil {
			t.Error("expected error for empty siteId")
		}
	})
}

func TestService_Close(t *testing.T) {
	repo := &mockRepo{sessions: make(map[string]*uploadsession.UploadSession)}
	service := uploadsession.NewService(repo)

	// Create test session
	session, err := service.Create(context.Background(), "test-site-id")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("success", func(t *testing.T) {
		if err := service.Close(context.Background(), session.Id); err != nil {
			t.Fatalf("Close() error = %v", err)
		}

		// Verify session was closed
		updated, err := repo.Get(context.Background(), session.Id)
		if err != nil {
			t.Fatal(err)
		}
		if updated.Status != uploadsession.StatusClosed {
			t.Errorf("expected status %q, got %q", uploadsession.StatusClosed, updated.Status)
		}
	})

	t.Run("already closed", func(t *testing.T) {
		if err := service.Close(context.Background(), session.Id); err != nil {
			t.Errorf("Close() on already closed session should be no-op, got error: %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		err := service.Close(context.Background(), "nonexistent")
		if err == nil {
			t.Error("expected error for nonexistent session")
		}
	})
}
