package feedback

import "github.com/platforma-dev/platforma/httpserver"

type Domain struct {
	Repository   *Repository
	Service      *Service
	HandlerGroup *httpserver.HandlerGroup
}

func (d *Domain) GetRepository() any {
	return d.Repository
}

func New(db db) *Domain {
	repository := NewRepository(db)
	service := NewService(repository)

	handlerGroup := httpserver.NewHandlerGroup()
	handlerGroup.Handle("POST /", NewCreateHandler(service))

	return &Domain{
		Repository:   repository,
		Service:      service,
		HandlerGroup: handlerGroup,
	}
}
