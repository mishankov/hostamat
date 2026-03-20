package uploadsession

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/platforma-dev/platforma/auth"
	"github.com/platforma-dev/platforma/log"
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
	ev := log.EventFromContext(ctx)

	ev.AddStep(slog.LevelInfo, "uploadsession.Get")
	ev.AddAttrs(map[string]any{
		"uploadsession.id": id,
	})

	us, err := s.repo.Get(ctx, id)
	if err != nil {
		ev.AddError(err)
		return nil, err
	}

	ev.AddAttrs(map[string]any{
		"site.id": us.Site,
		"user.id": us.User,
	})

	return us, nil
}

func (s *Service) Close(ctx context.Context, uploadSessionId string) error {
	ev := log.EventFromContext(ctx)

	ev.AddStep(slog.LevelInfo, "uploadsession.Close")
	ev.AddAttrs(map[string]any{
		"uploadsession.id": uploadSessionId,
	})

	us, err := s.Get(ctx, uploadSessionId)
	if err != nil {
		return err
	}

	us.Status = UploadSessionStatusClosed
	us.UpdatedAt = time.Now()

	err = s.repo.UpdateStatus(ctx, us)
	if err != nil {
		ev.AddError(err)
		return err
	}

	err = s.repo.DeleteFilesForUploadSession(ctx, us.Id)
	if err != nil {
		ev.AddError(err)
		return err
	}

	return nil
}

func (s *Service) Create(ctx context.Context, siteId string) (*UploadSession, error) {
	ev := log.EventFromContext(ctx)

	ev.AddStep(slog.LevelInfo, "uploadsession.Create")
	ev.AddAttrs(map[string]any{
		"site.id": siteId,
	})

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
		ev.AddError(err)
		return nil, err
	}

	ev.AddAttrs(map[string]any{
		"uploadsession.id": uploadSession.Id,
	})

	return uploadSession, nil
}

func (s *Service) GetSize(ctx context.Context, uploadSessionId string) (int, error) {
	ev := log.EventFromContext(ctx)

	ev.AddStep(slog.LevelInfo, "uploadsession.GetSize")
	ev.AddAttrs(map[string]any{
		"uploadsession.id": uploadSessionId,
	})

	size, err := s.repo.GetFilesSizeForUploadSession(ctx, uploadSessionId)
	ev.AddError(err)

	return size, err
}

func (s *Service) CreateFile(ctx context.Context, uploadSessionId, filename string, size int, uploaded bool) error {
	ev := log.EventFromContext(ctx)

	ev.AddStep(slog.LevelInfo, "uploadsession.CreateFile")
	ev.AddAttrs(map[string]any{
		"uploadsession.id": uploadSessionId,
		"file.name":        filename,
		"file.size":        size,
	})

	err := s.repo.CreateFile(ctx, File{
		Id:              uuid.NewString(),
		UploadSessionId: uploadSessionId,
		Name:            filename,
		Uploaded:        uploaded,
		Size:            size,
	})
	ev.AddError(err)

	return err
}

func (s *Service) GetFiles(ctx context.Context, uploadSessionId string) ([]File, error) {
	ev := log.EventFromContext(ctx)

	ev.AddStep(slog.LevelInfo, "uploadsession.GetFiles")
	ev.AddAttrs(map[string]any{
		"uploadsession.id": uploadSessionId,
	})

	files, err := s.repo.GetFiles(ctx, uploadSessionId)
	ev.AddError(err)

	return files, err
}

func (s *Service) DeleteByUser(ctx context.Context, userId string) error {
	ev := log.EventFromContext(ctx)
	ev.AddStep(slog.LevelInfo, "uploadsession.DeleteByUser")
	ev.AddAttrs(map[string]any{
		"user.id": userId,
	})

	err := s.repo.DeleteByUser(ctx, userId)
	ev.AddError(err)

	return s.repo.DeleteByUser(ctx, userId)
}
