package uploadsession

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type repository interface {
	Create(ctx context.Context, session *UploadSession) error
	Get(ctx context.Context, id string) (*UploadSession, error)
	Update(ctx context.Context, session *UploadSession) error
}

type Service struct {
	repo repository
}

func NewService(repo repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, siteId string) (*UploadSession, error) {
	if siteId == "" {
		return nil, fmt.Errorf("siteId is required")
	}

	session := &UploadSession{
		Id:     uuid.NewString(),
		SiteId: siteId,
		Status: StatusOpen,
	}

	err := s.repo.Create(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to create upload session: %w", err)
	}

	return session, nil
}

func (s *Service) Close(ctx context.Context, id string) error {
	session, err := s.repo.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get upload session: %w", err)
	}

	if session.Status == StatusClosed {
		return nil // already closed
	}

	session.Status = StatusClosed
	session.UpdatedAt = time.Now()

	return s.repo.Update(ctx, session)
}
