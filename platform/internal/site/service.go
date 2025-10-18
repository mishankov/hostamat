package site

import (
	"context"
	"errors"
	"hostamat/internal/subdomain"
	"hostamat/pkg/platforma/auth"
	"time"

	"github.com/google/uuid"
)

type repository interface {
	Get(ctx context.Context, id string) (*Site, error)
	GetByDomain(ctx context.Context, domain string) (*Site, error)
	GetList(ctx context.Context, userId string) ([]*Site, error)
	Create(ctx context.Context, site *Site) error
	Update(ctx context.Context, site *Site) error
	Delete(ctx context.Context, id string) error
}

type subdomainService interface {
	Create(ctx context.Context) (*subdomain.Subdomain, error)
	Delete(ctx context.Context, id string) error
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

	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, errors.New("access denied")
	}

	subdomain, err := s.subdomainService.Create(ctx)
	if err != nil {
		return nil, err
	}

	site := &Site{
		Id:        uuid.NewString(),
		User:      user.ID,
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
	site, err := s.repository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if ok := s.checkAccess(ctx, site); !ok {
		return nil, errors.New("access denied")
	}

	return site, nil
}

func (s *Service) GetIdByDomain(ctx context.Context, domain string) (string, error) {
	site, err := s.repository.GetByDomain(ctx, domain)
	if err != nil {
		return "", err
	}

	return site.Id, nil
}

func (s *Service) GetByDomain(ctx context.Context, domain string) (*Site, error) {
	site, err := s.repository.GetByDomain(ctx, domain)
	if err != nil {
		return nil, err
	}

	if ok := s.checkAccess(ctx, site); !ok {
		return nil, errors.New("access denied")
	}

	return site, nil
}

func (s *Service) GetList(ctx context.Context) ([]*Site, error) {
	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, errors.New("access denied")
	}

	return s.repository.GetList(ctx, user.ID)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	site, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	err = s.subdomainService.Delete(ctx, site.Subdomain)
	if err != nil {
		return err
	}

	// `s.Get` already checks permission, so can delete without additional checks
	return s.repository.Delete(ctx, id)
}

func (s *Service) checkAccess(ctx context.Context, site *Site) bool {
	user := auth.UserFromContext(ctx)
	if user == nil {
		return false
	}

	if site.User == user.ID {
		return true
	}

	return false
}
