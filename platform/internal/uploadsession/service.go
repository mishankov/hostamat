package uploadsession

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/platforma-dev/platforma/auth"
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

	err = s.repo.UpdateStatus(ctx, us)
	if err != nil {
		return err
	}

	err = s.repo.DeleteFilesForUploadSession(ctx, us.Id)
	if err != nil {
		return err
	}

	return nil
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

func (s *Service) GetSize(ctx context.Context, uploadSessionId string) (int, error) {
	return s.repo.GetFilesSizeForUploadSession(ctx, uploadSessionId)
}

func (s *Service) CreateFile(ctx context.Context, uploadSessionId, filename string, size int, uploaded bool) error {
	return s.repo.CreateFile(ctx, File{
		Id:              uuid.NewString(),
		UploadSessionId: uploadSessionId,
		Name:            filename,
		Uploaded:        uploaded,
		Size:            size,
	})
}

func (s *Service) GetFiles(ctx context.Context, uploadSessionId string) ([]File, error) {
	return s.repo.GetFiles(ctx, uploadSessionId)
}
