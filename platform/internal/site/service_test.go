package site_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"hostamat/internal/site"
)

func TestService_Create(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		mockRepo := newMockRepository()
		service := site.NewService(mockRepo)

		testDomain := "example.com"
		createdSite, err := service.Create(context.Background(), testDomain)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		if createdSite.Id == "" {
			t.Error("Expected non-empty site ID")
		}
		if createdSite.Domain != testDomain {
			t.Errorf("Expected domain %q, got %q", testDomain, createdSite.Domain)
		}
		if createdSite.CreatedAt.IsZero() {
			t.Error("Expected CreatedAt to be set")
		}

		// Verify the site was stored
		_, err = mockRepo.Get(context.Background(), createdSite.Id)
		if err != nil {
			t.Errorf("Site not found in repository: %v", err)
		}
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := &errorProneMockRepository{
			shouldError: true,
		}
		service := site.NewService(mockRepo)

		_, err := service.Create(context.Background(), "example.com")
		if err == nil {
			t.Error("Expected error from repository")
		}
	})
}

func TestService_Get(t *testing.T) {
	t.Run("existing site", func(t *testing.T) {
		mockRepo := newMockRepository()
		service := site.NewService(mockRepo)

		// Create test site directly in mock
		testSite := &site.Site{
			Id:        "test-get-id",
			Domain:    "example.com",
			CreatedAt: time.Now(),
		}
		if err := mockRepo.Create(context.Background(), testSite); err != nil {
			t.Fatalf("Failed to setup test site: %v", err)
		}

		// Test Get
		got, err := service.Get(context.Background(), "test-get-id")
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}

		if got.Id != testSite.Id {
			t.Errorf("Expected ID %q, got %q", testSite.Id, got.Id)
		}
		if got.Domain != testSite.Domain {
			t.Errorf("Expected domain %q, got %q", testSite.Domain, got.Domain)
		}
	})

	t.Run("non-existent site", func(t *testing.T) {
		mockRepo := newMockRepository()
		service := site.NewService(mockRepo)

		_, err := service.Get(context.Background(), "non-existent-id")
		if err == nil {
			t.Error("Expected error for non-existent site")
		}
	})
}

type mockRepository struct {
	mu    sync.Mutex
	sites map[string]*site.Site
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		sites: make(map[string]*site.Site),
	}
}

func (m *mockRepository) Get(ctx context.Context, id string) (*site.Site, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	s, exists := m.sites[id]
	if !exists {
		return nil, errors.New("not found")
	}
	return s, nil
}

func (m *mockRepository) Create(ctx context.Context, site *site.Site) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.sites[site.Id]; exists {
		return errors.New("already exists")
	}
	m.sites[site.Id] = site
	return nil
}

type errorProneMockRepository struct {
	mockRepository
	shouldError bool
}

func (m *errorProneMockRepository) Create(ctx context.Context, site *site.Site) error {
	if m.shouldError {
		return errors.New("forced repository error")
	}
	return m.mockRepository.Create(ctx, site)
}
