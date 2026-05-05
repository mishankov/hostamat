package user

import "github.com/platforma-dev/platforma/queue"
import "github.com/platforma-dev/platforma/auth"

type Domain struct {
	Service      *Service
	CleanupQueue *queue.Processor[auth.UserCleanupJob]
}

func New(siteService siteService, uploadSessionService uploadSessionService, feedbackService feedbackService, formCollectionService formCollectionService) *Domain {
	service := NewService(siteService, uploadSessionService, feedbackService, formCollectionService)

	return &Domain{
		Service:      service,
		CleanupQueue: service.cleanupProcessor,
	}
}
