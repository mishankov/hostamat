package site

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type repository interface {
	Get(ctx context.Context, id string) (*Site, error)
	Create(ctx context.Context, site *Site) error
}

type Service struct {
	repository repository
}

func NewService(repository repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Create(ctx context.Context, subdomain string) (*Site, error) {
	site := &Site{
		Id:        uuid.NewString(),
		Subdomain: subdomain,
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
