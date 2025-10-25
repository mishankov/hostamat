package uploadsession

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mishankov/platforma/auth"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Get(ctx context.Context, id string) (*UploadSession, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) Close(ctx context.Context, uploadSessionId string) error {
	us, err := s.Get(ctx, uploadSessionId)
	if err != nil {
		return err
	}

	us.Status = UploadSessionStatusClosed
	us.UpdatedAt = time.Now()

	return s.repo.UpdateStatus(ctx, us)
}

func (s *Service) Create(ctx context.Context, siteId string) (*UploadSession, error) {
	user := auth.UserFromContext(ctx)

	uploadSession := &UploadSession{
		Id:        uuid.NewString(),
		User:      user.ID,
		Site:      siteId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Status:    UploadSessionStatusOpen,
	}

	err := s.repo.Create(ctx, uploadSession)
	if err != nil {
		return nil, err
	}

	return uploadSession, nil
}
