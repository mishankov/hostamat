package site

import (
	"context"
	"hostamat/internal/subdomain"
	"time"

	"github.com/google/uuid"
)

type repository interface {
	Get(ctx context.Context, id string) (*Site, error)
	GetByDomain(ctx context.Context, domain string) (*Site, error)
	GetList(ctx context.Context) ([]*Site, error)
	Create(ctx context.Context, site *Site) error
	Update(ctx context.Context, site *Site) error
}

type subdomainService interface {
	Create(ctx context.Context) (*subdomain.Subdomain, error)
}

type Service struct {
	repository       repository
	subdomainService subdomainService
}

func NewService(repository repository, subdomainService subdomainService) *Service {
	return &Service{
		repository:       repository,
		subdomainService: subdomainService,
	}
}

func (s *Service) Create(ctx context.Context) (*Site, error) {
	var err error

	subdomain, err := s.subdomainService.Create(ctx)
	if err != nil {
		return nil, err
	}

	site := &Site{
		Id:        uuid.NewString(),
		Subdomain: subdomain.Id,
		Domain:    subdomain.FullDomain,
		CreatedAt: time.Now(),
	}

	err = s.repository.Create(ctx, site)
	if err != nil {
		return nil, err
	}

	return site, nil
}

func (s *Service) Get(ctx context.Context, id string) (*Site, error) {
	return s.repository.Get(ctx, id)
}

func (s *Service) GetByDomain(ctx context.Context, domain string) (*Site, error) {
	return s.repository.GetByDomain(ctx, domain)
}

func (s *Service) GetList(ctx context.Context) ([]*Site, error) {
	return s.repository.GetList(ctx)
}

func (s *Service) UpdateVersion(ctx context.Context, id string) error {
	site, err := s.repository.Get(ctx, id)
	if err != nil {
		return err
	}

	site.Version++
	site.UpdatedAt = time.Now()
	return s.repository.Update(ctx, site)
}
