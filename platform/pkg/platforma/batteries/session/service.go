package session

import "context"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, session *Session) error {
	return s.repo.Create(ctx, session)
}

func (s *Service) Get(ctx context.Context, id string) (*Session, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) GetByUserId(ctx context.Context, id string) (*Session, error) {
	return s.repo.GetByUserId(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
