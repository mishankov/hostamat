package site_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"hostamat/internal/site"
	"hostamat/internal/subdomain"
)

func TestService_Create(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		mockRepo := newMockRepository()
		service := site.NewService(mockRepo, &mockDomainService{})

		createdSite, err := service.Create(context.Background())
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		if createdSite.Id == "" {
			t.Error("Expected non-empty site ID")
		}
		if createdSite.Subdomain != "id" {
			t.Errorf("Expected subdomain %q, got %q", "id", createdSite.Subdomain)
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
		service := site.NewService(mockRepo, &mockDomainService{})

		_, err := service.Create(context.Background())
		if err == nil {
			t.Error("Expected error from repository")
		}
	})
}

func TestService_Get(t *testing.T) {
	t.Run("existing site", func(t *testing.T) {
		mockRepo := newMockRepository()
		service := site.NewService(mockRepo, &mockDomainService{})

		// Create test site directly in mock
		testSite := &site.Site{
			Id:        "test-get-id",
			Subdomain: "example",
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
		if got.Subdomain != testSite.Subdomain {
			t.Errorf("Expected subdomain %q, got %q", testSite.Subdomain, got.Subdomain)
		}
	})

	t.Run("non-existent site", func(t *testing.T) {
		mockRepo := newMockRepository()
		service := site.NewService(mockRepo, &mockDomainService{})

		_, err := service.Get(context.Background(), "non-existent-id")
		if err == nil {
			t.Error("Expected error for non-existent site")
		}
	})
}

func TestService_GetList(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := newMockRepository()
		service := site.NewService(mockRepo, &mockDomainService{})

		// Create test sites directly in mock
		testSites := []*site.Site{
			{
				Id:        "1",
				Subdomain: "site1",
				CreatedAt: time.Now(),
			},
			{
				Id:        "2",
				Subdomain: "site2",
				CreatedAt: time.Now(),
			},
		}

		for _, s := range testSites {
			if err := mockRepo.Create(context.Background(), s); err != nil {
				t.Fatalf("Failed to setup test site: %v", err)
			}
		}

		// Test GetList
		sites, err := service.GetList(context.Background())
		if err != nil {
			t.Fatalf("GetList failed: %v", err)
		}

		if len(sites) != len(testSites) {
			t.Fatalf("Expected %d sites, got %d", len(testSites), len(sites))
		}
	})

	t.Run("empty list", func(t *testing.T) {
		mockRepo := newMockRepository()
		service := site.NewService(mockRepo, &mockDomainService{})

		sites, err := service.GetList(context.Background())
		if err != nil {
			t.Fatalf("GetList failed: %v", err)
		}

		if len(sites) != 0 {
			t.Errorf("Expected empty list, got %d items", len(sites))
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

func (m *mockRepository) GetByDomain(ctx context.Context, domain string) (*site.Site, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, site := range m.sites {
		if site.Domain == domain {
			return site, nil
		}
	}
	return nil, nil
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

func (m *mockRepository) GetList(ctx context.Context) ([]*site.Site, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	sites := make([]*site.Site, 0, len(m.sites))
	for _, s := range m.sites {
		sites = append(sites, s)
	}
	return sites, nil
}

func (m *mockRepository) Update(ctx context.Context, site *site.Site) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.sites[site.Id]; !exists {
		return errors.New("not found")
	}
	m.sites[site.Id] = site
	return nil
}

func (m *mockRepository) Delete(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.sites[id]; !exists {
		return errors.New("not found")
	}
	delete(m.sites, id)
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

type mockDomainService struct{}

func (s *mockDomainService) Create(_ context.Context) (*subdomain.Subdomain, error) {
	return &subdomain.Subdomain{
		Id:         "id",
		Name:       "name",
		FullDomain: "name.example.com",
	}, nil
}
