package file

import (
	"net/http"

	"github.com/mishankov/platforma/httpserver"
)

type Domain struct {
	Repository        *Repository
	Service           *Service
	ServeHandle       http.Handler
	ChangeHandleGroup *httpserver.HandlerGroup
}

func (d *Domain) GetRepository() any {
	return d.Repository
}

func New(s3Client s3, siteService siteService, bucketName string) *Domain {
	repository := NewRepository(s3Client)
	service := NewService(repository, siteService, bucketName)

	serveFileHandler := NewServeHandler(service)
	uploadFileHandler := NewUploadHandler(service)
	deleteHandler := NewDeleteHandler(service)

	// Handle groups
	filesAPI := httpserver.NewHandlerGroup()
	filesAPI.Handle("POST /", uploadFileHandler)
	filesAPI.Handle("DELETE /", deleteHandler)

	return &Domain{
		Repository:        repository,
		Service:           service,
		ServeHandle:       serveFileHandler,
		ChangeHandleGroup: filesAPI,
	}
}
