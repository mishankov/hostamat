package site

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type repository interface {
	Get(ctx context.Context, id string) (*Site, error)
	GetList(ctx context.Context) ([]*Site, error)
	Create(ctx context.Context, site *Site) error
}

type Service struct {
	repository repository
	baseDomain string
}

func NewService(repository repository, baseDomain string) *Service {
	return &Service{
		repository: repository,
		baseDomain: baseDomain,
	}
}

func (s *Service) Create(ctx context.Context, subdomain string) (*Site, error) {
	site := &Site{
		Id:        uuid.NewString(),
		Subdomain: subdomain,
		Domain:    subdomain + "." + s.baseDomain,
		CreatedAt: time.Now(),
	}

	err := s.repository.Create(ctx, site)
	if err != nil {
		return nil, err
	}

	return site, nil
}

func (s *Service) Get(ctx context.Context, id string) (*Site, error) {
	return s.repository.Get(ctx, id)
}

func (s *Service) GetList(ctx context.Context) ([]*Site, error) {
	return s.repository.GetList(ctx)
}
