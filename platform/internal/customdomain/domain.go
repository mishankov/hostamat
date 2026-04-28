package customdomain

import "github.com/platforma-dev/platforma/httpserver"

type Domain struct {
	Repository   *Repository
	Service      *Service
	HandlerGroup *httpserver.HandlerGroup
}

func (d *Domain) GetRepository() any {
	return d.Repository
}

func New(db db, siteRepo siteRepository, hostSync hostSync, dnsCfg DNSConfig) *Domain {
	repository := NewRepository(db)
	service := NewService(repository, siteRepo, hostSync, netResolver{}, dnsCfg)

	handlerGroup := httpserver.NewHandlerGroup()
	handlerGroup.Handle("GET /", NewListHandler(service))
	handlerGroup.Handle("POST /", NewCreateHandler(service))
	handlerGroup.Handle("PUT /", NewUpdateHandler(service))
	handlerGroup.Handle("POST /verify", NewVerifyHandler(service))
	handlerGroup.Handle("DELETE /", NewDeleteHandler(service))

	return &Domain{
		Repository:   repository,
		Service:      service,
		HandlerGroup: handlerGroup,
	}
}
