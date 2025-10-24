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

func (s *Service) CreateUploadSession(ctx context.Context, siteId string) (*UploadSession, error) {
	user := auth.UserFromContext(ctx)

	uploadSession := &UploadSession{
		Id:        uuid.NewString(),
		User:      user.ID,
		Site:      siteId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Status:    UploadSessionStatusNew,
	}

	err := s.repo.Create(ctx, uploadSession)
	if err != nil {
		return nil, err
	}

	return uploadSession, nil
}
