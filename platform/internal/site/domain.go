package site

import (
	"hostamat/internal/file"
	"hostamat/internal/uploadsession"
	"net/http"

	"github.com/platforma-dev/platforma/httpserver"
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

func New(db db, subdomainService subdomainService, customDomainService customDomainService, fileService *file.Service, uploadSessionService *uploadsession.Service, formCollectionService formCollectionService) *Domain {
	repository := NewRepository(db)
	service := NewService(repository, subdomainService, customDomainService, fileService, uploadSessionService, formCollectionService)

	siteHandlerGroup := httpserver.NewHandlerGroup()
	siteHandlerGroup.Handle("GET /", NewListSiteHandler(service))
	siteHandlerGroup.Handle("PUT /", NewUpdateSiteHandler(service))
	siteHandlerGroup.Handle("POST /", NewCreateSiteHandler(service))
	siteHandlerGroup.Handle("DELETE /", NewDeleteSiteHandler(service))

	fileHandlerGroup := httpserver.NewHandlerGroup()
	fileHandlerGroup.Handle("POST /", NewUploadFileHandler(service))

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
