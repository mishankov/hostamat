package subdomain

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"hostamat/pkg/platforma/batteries/auth"
	"strings"

	"github.com/google/uuid"
)

type repository interface {
	Get(ctx context.Context, id string) (*Subdomain, error)
	Create(ctx context.Context, subdomain *Subdomain) error
}

type webserver interface {
	AddSubdomain(ctx context.Context, baseDomain, subdomain, shorthandId string) error
}

type Service struct {
	repo        repository
	webserver   webserver
	baseDomain  string
	shorthandId string
}

func NewService(repo repository, webserver webserver, baseDomain, shorthandId string) *Service {
	return &Service{
		repo:        repo,
		webserver:   webserver,
		baseDomain:  baseDomain,
		shorthandId: shorthandId,
	}
}

func (s *Service) Get(ctx context.Context, id string) (*Subdomain, error) {
	subdomain, err := s.repo.Get(ctx, id)

	if err != nil {
		return nil, fmt.Errorf("failed to get subdomain: %w", err)
	}

	if !s.checkAccess(ctx, subdomain) {
		return nil, errors.New("access denied")
	}

	return subdomain, nil
}

func (s *Service) Create(ctx context.Context) (*Subdomain, error) {
	var err error

	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, errors.New("access denied")
	}

	name := randomBase64String(5)

	subdomain := &Subdomain{
		Id:         uuid.NewString(),
		User:       user.ID,
		Name:       name,
		FullDomain: name + "." + s.baseDomain,
	}

	err = s.repo.Create(ctx, subdomain)
	if err != nil {
		return nil, fmt.Errorf("failed to create subdomain: %w", err)
	}

	err = s.webserver.AddSubdomain(ctx, s.baseDomain, subdomain.Name, s.shorthandId)
	if err != nil {
		return nil, fmt.Errorf("failed to add sudomain to webserver: %w", err)
	}

	return subdomain, nil
}

func (s *Service) checkAccess(ctx context.Context, subdomain *Subdomain) bool {
	user := auth.UserFromContext(ctx)
	if user == nil {
		return false
	}

	if subdomain.User == user.ID {
		return true
	}

	return false
}

func randomBase64String(length int) string {
	bytes := make([]byte, length+1)
	rand.Read(bytes)
	return strings.ToLower(base64.URLEncoding.EncodeToString(bytes)[:length])
}
