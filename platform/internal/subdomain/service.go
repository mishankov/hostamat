package subdomain

import (
	"context"
	"fmt"
)

type repository interface {
	Get(ctx context.Context, id string) (*Subdomain, error)
	Create(ctx context.Context, subdomain *Subdomain) error
}

type Service struct {
	repo repository
}

func NewService(repo repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Get(ctx context.Context, id string) (*Subdomain, error) {
	subdomain, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get subdomain: %w", err)
	}
	return subdomain, nil
}

func (s *Service) Create(ctx context.Context, subdomain *Subdomain) error {
	if subdomain.Name == "" {
		return fmt.Errorf("subdomain name cannot be empty")
	}

	err := s.repo.Create(ctx, subdomain)
	if err != nil {
		return fmt.Errorf("failed to create subdomain: %w", err)
	}

	return nil
}
