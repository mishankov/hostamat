package site

import (
	"hostamat/internal/file"
	"hostamat/internal/uploadsession"
	"net/http"

	"github.com/mishankov/platforma/httpserver"
)

type Domain struct {
	Repository                *Repository
	Service                   *Service
	SiteHandlerGroup          *httpserver.HandlerGroup
	FileHandlerGroup          *httpserver.HandlerGroup
	FileServeHandler          http.Handler
	UploadSessionHandlerGroup *httpserver.HandlerGroup
}

func (d *Domain) GetRepository() any {
	return d.Repository
}

func New(db db, subdomainService subdomainService, fileService *file.Service, uploadSessionService *uploadsession.Service) *Domain {
	repository := NewRepository(db)
	service := NewService(repository, subdomainService, fileService, uploadSessionService)

	createHandle := NewCreateSiteHandler(service)
	listHandler := NewListSiteHandler(service)
	deleteHandler := NewDeleteSiteHandler(service)

	siteHandlerGroup := httpserver.NewHandlerGroup()
	siteHandlerGroup.Handle("POST /", createHandle)
	siteHandlerGroup.Handle("GET /", listHandler)
	siteHandlerGroup.Handle("DELETE /", deleteHandler)

	fileHandlerGroup := httpserver.NewHandlerGroup()
	fileHandlerGroup.Handle("POST /", NewUploadFileHandler(service))
	fileHandlerGroup.Handle("DELETE /", NewDeleteFileHandler(service))

	uploadSessionHandlerGroup := httpserver.NewHandlerGroup()
	uploadSessionHandlerGroup.Handle("POST /", NewCreateUploadSessionHandler(service))
	uploadSessionHandlerGroup.Handle("DELETE /", NewCloseUploadSessionHandler(service))

	return &Domain{
		Repository:                repository,
		Service:                   service,
		SiteHandlerGroup:          siteHandlerGroup,
		FileHandlerGroup:          fileHandlerGroup,
		FileServeHandler:          NewServeFileHandler(service),
		UploadSessionHandlerGroup: uploadSessionHandlerGroup,
	}
}
