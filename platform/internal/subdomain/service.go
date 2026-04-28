package subdomain

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"math/rand"

	"github.com/platforma-dev/platforma/auth"
	"github.com/platforma-dev/platforma/log"

	"github.com/google/uuid"
)

type repository interface {
	Get(ctx context.Context, id string) (*Subdomain, error)
	Create(ctx context.Context, subdomain *Subdomain) error
	Delete(ctx context.Context, id string) error
	ListHosts(ctx context.Context) ([]string, error)
}

type hostSync interface {
	Add(ctx context.Context, host string) error
	Remove(ctx context.Context, host string) error
}

type Service struct {
	repo       repository
	hostSync   hostSync
	baseDomain string
}

func NewService(repo repository, hostSync hostSync, baseDomain string) *Service {
	return &Service{
		repo:       repo,
		hostSync:   hostSync,
		baseDomain: baseDomain,
	}
}

func (s *Service) Get(ctx context.Context, id string) (*Subdomain, error) {
	ev := log.EventFromContext(ctx)

	ev.AddStep(slog.LevelInfo, "subdomain.Get")
	ev.AddAttrs(map[string]any{
		"subdomain.id": id,
	})

	subdomain, err := s.repo.Get(ctx, id)

	if err != nil {
		err := fmt.Errorf("failed to get subdomain: %w", err)
		ev.AddError(err)
		return nil, err
	}

	ev.AddAttrs(map[string]any{
		"subdomain.name":        subdomain.Name,
		"subdomain.full_domain": subdomain.FullDomain,
		"user.id":               subdomain.User,
	})

	if !s.checkAccess(ctx, subdomain) {
		ev.AddStep(slog.LevelWarn, "subdomain.Get.access_denied")
		return nil, errors.New("access denied")
	}

	return subdomain, nil
}

func (s *Service) Create(ctx context.Context) (*Subdomain, error) {
	ev := log.EventFromContext(ctx)

	ev.AddStep(slog.LevelInfo, "subdomain.Create")

	var err error

	user := auth.UserFromContext(ctx)
	if user == nil {
		ev.AddAttrs(map[string]any{
			"user.accessDenied": true,
		})
		return nil, errors.New("access denied")
	}

	ev.AddAttrs(map[string]any{
		"user.id": user.ID,
	})

	name := generateName()

	subdomain := &Subdomain{
		Id:         uuid.NewString(),
		User:       user.ID,
		Name:       name,
		FullDomain: name + "." + s.baseDomain,
	}

	ev.AddAttrs(map[string]any{
		"subdomain.id":          subdomain.Id,
		"subdomain.name":        subdomain.Name,
		"subdomain.full_domain": subdomain.FullDomain,
	})

	err = s.repo.Create(ctx, subdomain)
	if err != nil {
		err := fmt.Errorf("failed to create subdomain: %w", err)
		ev.AddError(err)
		return nil, err
	}

	err = s.hostSync.Add(ctx, subdomain.FullDomain)
	if err != nil {
		err := fmt.Errorf("failed to add sudomain to webserver: %w", err)
		ev.AddError(err)
		return nil, err
	}

	return subdomain, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	ev := log.EventFromContext(ctx)

	ev.AddStep(slog.LevelInfo, "subdomain.Delete")
	ev.AddAttrs(map[string]any{
		"subdomain.id": id,
	})

	subdomain, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	ev.AddAttrs(map[string]any{
		"subdomain.name":        subdomain.Name,
		"subdomain.full_domain": subdomain.FullDomain,
	})

	err = s.hostSync.Remove(ctx, subdomain.FullDomain)
	if err != nil {

		ev.AddError(err)

		return err
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {

		ev.AddError(err)

		return fmt.Errorf("failed to delete subdomain from database: %w", err)
	}

	return nil
}

func (s *Service) ListHosts(ctx context.Context) ([]string, error) {
	return s.repo.ListHosts(ctx)
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
