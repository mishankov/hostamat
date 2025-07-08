package file

import (
	"context"
	"io"
	"strings"
)

type repository interface {
	ObjectExists(ctx context.Context, fileName string) (bool, error)
	GetObject(ctx context.Context, fileName string) (io.Reader, error)
	PutObject(ctx context.Context, fileName string, file io.Reader, size int64) error
}

type Service struct {
	repo       repository
	bucketName string
}

func NewService(repo repository, bucketName string) *Service {
	return &Service{
		repo:       repo,
		bucketName: bucketName,
	}
}

func (s *Service) GetFile(ctx context.Context, fileName string) (io.Reader, error) {
	return s.repo.GetObject(ctx, fileName)
}

func (s *Service) UploadFile(ctx context.Context, fileName string, file io.Reader, size int64) error {
	return s.repo.PutObject(ctx, fileName, file, size)
}

func (s *Service) FileExists(ctx context.Context, fileName string) (bool, error) {
	return s.repo.ObjectExists(ctx, fileName)
}

func contentTypeFromFileName(fileName string) string {
	switch {
	case strings.HasSuffix(fileName, ".html"):
		return "text/html"
	case strings.HasSuffix(fileName, ".css"):
		return "text/css"
	case strings.HasSuffix(fileName, ".js"):
		return "application/javascript"
	case strings.HasSuffix(fileName, ".png"):
		return "image/png"
	case strings.HasSuffix(fileName, ".jpg"), strings.HasSuffix(fileName, ".jpeg"):
		return "image/jpeg"
	case strings.HasSuffix(fileName, ".svg"):
		return "image/svg+xml"
	default:
		return "application/octet-stream"
	}
}
