package customdomain

import (
	"context"
	"time"

	"github.com/platforma-dev/platforma/application"
	"github.com/platforma-dev/platforma/httpserver"
	"github.com/platforma-dev/platforma/scheduler"
)

type Domain struct {
	Repository       *Repository
	Service          *Service
	HandlerGroup     *httpserver.HandlerGroup
	CleanupScheduler *scheduler.Scheduler
}

func (d *Domain) GetRepository() any {
	return d.Repository
}

func New(db db, siteRepo siteRepository, hostSync hostSync, dnsCfg DNSConfig) *Domain {
	repository := NewRepository(db)
	service := NewService(repository, siteRepo, hostSync, netResolver{}, dnsCfg)
	cleanupScheduler, err := scheduler.New("@hourly", application.RunnerFunc(func(ctx context.Context) error {
		_, err := service.DeleteStaleUnverified(ctx, 48*time.Hour)
		return err
	}))
	if err != nil {
		panic(err)
	}

	handlerGroup := httpserver.NewHandlerGroup()
	handlerGroup.Handle("GET /", NewListHandler(service))
	handlerGroup.Handle("POST /", NewCreateHandler(service))
	handlerGroup.Handle("PUT /", NewUpdateHandler(service))
	handlerGroup.Handle("POST /verify", NewVerifyHandler(service))
	handlerGroup.Handle("DELETE /", NewDeleteHandler(service))

	return &Domain{
		Repository:       repository,
		Service:          service,
		HandlerGroup:     handlerGroup,
		CleanupScheduler: cleanupScheduler,
	}
}
