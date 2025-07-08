package file

import (
	"archive/zip"
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"hostamat/internal/adapters/caddy"
	"io"
	"path/filepath"
	"strings"
)

type repository interface {
	ObjectExists(ctx context.Context, fileName string) (bool, error)
	GetObject(ctx context.Context, fileName string) (io.Reader, error)
	PutObject(ctx context.Context, fileName string, file io.Reader, size int64) error
}

type Service struct {
	repo        repository
	bucketName  string
	caddyAPI    *caddy.CaddyAPI
	baseDomain  string
	shorthandId string
}

func NewService(repo repository, bucketName string, caddyAPI *caddy.CaddyAPI, baseDomain, shorthandId string) *Service {
	return &Service{
		repo:        repo,
		bucketName:  bucketName,
		caddyAPI:    caddyAPI,
		baseDomain:  baseDomain,
		shorthandId: shorthandId,
	}
}

func (s *Service) GetFile(ctx context.Context, fileName string) (io.Reader, error) {
	return s.repo.GetObject(ctx, fileName)
}

func (s *Service) FileExists(ctx context.Context, fileName string) (bool, error) {
	return s.repo.ObjectExists(ctx, fileName)
}

func (s *Service) UploadSingleFile(ctx context.Context, subdomain, filename string, file io.Reader, size int64) error {
	if !isAllowedExtension(filename) {
		return ErrInvalidFileType
	}
	return s.repo.PutObject(ctx, subdomain+"/"+filename, file, size)
}

func (s *Service) UploadZipArchive(ctx context.Context, subdomain string, zipFile io.ReaderAt, size int64) ([]string, []string, error) {
	zipReader, err := zip.NewReader(zipFile, size)
	if err != nil {
		return nil, nil, err
	}

	var uploadedFiles []string
	var skippedFiles []string

	for _, f := range zipReader.File {
		if f.FileInfo().IsDir() {
			continue
		}

		filename := strings.Join(strings.Split(f.Name, "/")[1:], "/")
		if !isAllowedFile(filename) {
			skippedFiles = append(skippedFiles, filename)
			continue
		}

		fileReader, err := f.Open()
		if err != nil {
			skippedFiles = append(skippedFiles, filename)
			continue
		}

		err = s.UploadSingleFile(ctx, subdomain, filename, fileReader, int64(f.UncompressedSize64))
		fileReader.Close()
		if err != nil {
			skippedFiles = append(skippedFiles, filename)
			continue
		}

		uploadedFiles = append(uploadedFiles, filename)
	}

	return uploadedFiles, skippedFiles, nil
}

func (s *Service) GenerateSubdomain() string {
	return randomBase64String(5)
}

func (s *Service) CreateSubdomain(ctx context.Context, subdomain string) error {
	return s.caddyAPI.AddSubdomain(ctx, s.baseDomain, subdomain, s.shorthandId)
}

func (s *Service) GetContentType(filename string) string {
	switch {
	case strings.HasSuffix(filename, ".html"):
		return "text/html"
	case strings.HasSuffix(filename, ".css"):
		return "text/css"
	case strings.HasSuffix(filename, ".js"):
		return "application/javascript"
	case strings.HasSuffix(filename, ".png"):
		return "image/png"
	case strings.HasSuffix(filename, ".jpg"), strings.HasSuffix(filename, ".jpeg"):
		return "image/jpeg"
	case strings.HasSuffix(filename, ".svg"):
		return "image/svg+xml"
	default:
		return "application/octet-stream"
	}
}

func isAllowedFile(filename string) bool {
	if filename == "" || strings.HasPrefix(filename, ".") || strings.Contains(filename, "/.") {
		return false
	}
	return isAllowedExtension(filename)
}

func isAllowedExtension(filename string) bool {
	ext := strings.ToLower(strings.TrimPrefix(strings.ToLower(filepath.Ext(filename)), "."))
	switch ext {
	case "html", "css", "js", "png", "jpg", "jpeg", "svg":
		return true
	default:
		return false
	}
}

func randomBase64String(length int) string {
	bytes := make([]byte, length+1)
	rand.Read(bytes)
	return strings.ToLower(base64.URLEncoding.EncodeToString(bytes)[:length])
}

var (
	ErrInvalidFileType = errors.New("invalid file type")
)
