package user

import (
	"context"
	"hostamat/internal/site"
	"time"

	"github.com/platforma-dev/platforma/auth"
	"github.com/platforma-dev/platforma/log"
	"github.com/platforma-dev/platforma/queue"
)

type Site interface {
	GetId() string
}

type siteService interface {
	GetList(ctx context.Context) ([]*site.Site, error)
	Delete(ctx context.Context, id string) error
}

type uploadSessionService interface {
	DeleteByUser(ctx context.Context, userId string) error
}

type feedbackService interface {
	DeleteByUser(ctx context.Context, userId string) error
}

type formCollectionService interface {
	DeleteByUser(ctx context.Context, userId string) error
}

type Service struct {
	siteService           siteService
	uploadSessionService  uploadSessionService
	feedbackService       feedbackService
	formCollectionService formCollectionService
	cleanupProcessor      *queue.Processor[auth.UserCleanupJob]
}

func NewService(siteService siteService, uploadSessionService uploadSessionService, feedbackService feedbackService, formCollectionService formCollectionService) *Service {
	service := &Service{
		siteService:           siteService,
		uploadSessionService:  uploadSessionService,
		feedbackService:       feedbackService,
		formCollectionService: formCollectionService,
	}

	cleanupQueue := queue.NewChanQueue[auth.UserCleanupJob](5, time.Second)
	service.cleanupProcessor = queue.New(queue.HandlerFunc[auth.UserCleanupJob](func(ctx context.Context, job auth.UserCleanupJob) {
		ctx = context.WithValue(ctx, auth.UserContextKey, &auth.User{ID: job.UserID})
		err := service.cleanupUser(ctx, job.UserID)
		if err != nil {
			log.ErrorContext(ctx, "failed to cleanup user", "userId", job.UserID, "error", err)
		}
	}), cleanupQueue, 2, 3*time.Second)

	return service
}

func (s *Service) cleanupUser(ctx context.Context, userId string) error {
	sites, err := s.siteService.GetList(ctx)
	if err != nil {
		return err
	}

	if s.formCollectionService != nil {
		err = s.formCollectionService.DeleteByUser(ctx, userId)
		if err != nil {
			log.ErrorContext(ctx, "failed to delete form data during user cleanup", "userId", userId, "error", err)
			return err
		}
	}

	for _, site := range sites {
		err = s.siteService.Delete(ctx, site.Id)
		if err != nil {
			log.ErrorContext(ctx, "failed to delete site during user cleanup", "siteId", site.Id, "userId", userId, "error", err)
		}
	}

	err = s.uploadSessionService.DeleteByUser(ctx, userId)
	if err != nil {
		log.ErrorContext(ctx, "failed to delete upload sessions during user cleanup", "userId", userId, "error", err)
		return err
	}

	if s.feedbackService != nil {
		err = s.feedbackService.DeleteByUser(ctx, userId)
		if err != nil {
			log.ErrorContext(ctx, "failed to delete feedback during user cleanup", "userId", userId, "error", err)
			return err
		}
	}

	log.InfoContext(ctx, "user cleanup completed", "userId", userId)
	return nil
}
