package site

import (
	"github.com/mishankov/platforma/httpserver"
)

type Domain struct {
	Repository  *Repository
	Service     *Service
	HandleGroup *httpserver.HandlerGroup
}

func (d *Domain) GetRepository() any {
	return d.Repository
}

func New(db db, subdomainService subdomainService) *Domain {
	repository := NewRepository(db)
	service := NewService(repository, subdomainService)

	createHandle := NewCreateSiteHandler(service)
	listHandler := NewListSiteHandler(service)
	deleteHandler := NewDeleteSiteHandler(service)

	handlerGroup := httpserver.NewHandlerGroup()
	handlerGroup.Handle("POST /", createHandle)
	handlerGroup.Handle("GET /", listHandler)
	handlerGroup.Handle("DELETE /", deleteHandler)

	return &Domain{
		Repository:  repository,
		Service:     service,
		HandleGroup: handlerGroup,
	}
}
