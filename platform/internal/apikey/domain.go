package apikey

import "github.com/platforma-dev/platforma/httpserver"

type Domain struct {
	Repository   *Repository
	Service      *Service
	HandlerGroup *httpserver.HandlerGroup
}

func New(db db, pepper string) *Domain {
	repository := NewRepository(db)
	service := NewService(repository, pepper)

	handlerGroup := httpserver.NewHandlerGroup()
	handlerGroup.Handle("GET /", NewListHandler(service))
	handlerGroup.Handle("POST /", NewCreateHandler(service))
	handlerGroup.Handle("PATCH /", NewUpdateHandler(service))
	handlerGroup.Handle("DELETE /", NewRevokeHandler(service))

	return &Domain{
		Repository:   repository,
		Service:      service,
		HandlerGroup: handlerGroup,
	}
}

func (d *Domain) GetRepository() any {
	return d.Repository
}
