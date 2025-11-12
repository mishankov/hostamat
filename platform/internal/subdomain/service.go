package subdomain

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sync"

	"github.com/platforma-dev/platforma/auth"

	"github.com/google/uuid"
)

type repository interface {
	Get(ctx context.Context, id string) (*Subdomain, error)
	Create(ctx context.Context, subdomain *Subdomain) error
	Delete(ctx context.Context, id string) error
}

type webserver interface {
	AddSubdomain(ctx context.Context, baseDomain, subdomain, shorthandId string) error
	GetHostsFromShorthandId(ctx context.Context, shorthandId string) ([]string, error)
	RewriteHostsInShorthandId(ctx context.Context, shorthandId string, hosts []string) error
}

type Service struct {
	caddyMutex  sync.Mutex
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

	name := generateName()

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

	s.caddyMutex.Lock()
	defer s.caddyMutex.Unlock()

	err = s.webserver.AddSubdomain(ctx, s.baseDomain, subdomain.Name, s.shorthandId)
	if err != nil {
		return nil, fmt.Errorf("failed to add sudomain to webserver: %w", err)
	}

	return subdomain, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	subdomain, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	s.caddyMutex.Lock()
	defer s.caddyMutex.Unlock()

	hosts, err := s.webserver.GetHostsFromShorthandId(ctx, s.shorthandId)
	if err != nil {
		return err
	}

	var newHosts []string
	for _, host := range hosts {
		if host != subdomain.FullDomain {
			newHosts = append(newHosts, host)
		}
	}

	err = s.webserver.RewriteHostsInShorthandId(ctx, s.shorthandId, newHosts)
	if err != nil {
		return err
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete subdomain from database: %w", err)
	}

	return nil
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

// name has this pattern: <adjective>-<noun>-<5 random numbers>
func generateName() string {
	adjective := adjectives[rand.Intn(len(adjectives))]
	noun := nouns[rand.Intn(len(nouns))]
	number := 0
	for i := range 5 {
		number += (rand.Intn(9) + 1) * int(math.Pow10(i))
	}

	return fmt.Sprintf("%s-%s-%d", adjective, noun, number)
}
