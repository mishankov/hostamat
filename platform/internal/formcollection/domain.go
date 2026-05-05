package formcollection

import "github.com/platforma-dev/platforma/httpserver"

type Domain struct {
	Repository         *Repository
	Service            *Service
	HandlerGroup       *httpserver.HandlerGroup
	PublicHandlerGroup *httpserver.HandlerGroup
}

func (d *Domain) GetRepository() any {
	return d.Repository
}

func New(db db, fileService fileService) *Domain {
	repository := NewRepository(db)
	service := NewService(repository, fileService)

	handlerGroup := httpserver.NewHandlerGroup()
	handlerGroup.Handle("GET /", NewListDefinitionsHandler(service))
	handlerGroup.Handle("POST /enable", NewEnableHandler(service))
	handlerGroup.Handle("GET /submissions", NewListSubmissionsHandler(service))
	handlerGroup.Handle("DELETE /submissions", NewDeleteSubmissionHandler(service))

	publicHandlerGroup := httpserver.NewHandlerGroup()
	publicHandlerGroup.Handle("POST /submit/{token}", NewSubmitHandler(service))

	return &Domain{
		Repository:         repository,
		Service:            service,
		HandlerGroup:       handlerGroup,
		PublicHandlerGroup: publicHandlerGroup,
	}
}
