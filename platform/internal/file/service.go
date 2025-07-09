package file

import (
	"archive/zip"
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"path/filepath"
	"strings"
)

type repository interface {
	ObjectExists(ctx context.Context, bucketName, fileName string) (bool, error)
	GetObject(ctx context.Context, bucketName, fileName string) (io.Reader, error)
	PutObject(ctx context.Context, bucketName, fileName string, file io.Reader, size int64) error
}

type domainService interface {
	AddSubdomain(ctx context.Context, baseDomain, subdomain, shorthandId string) error
}

type Service struct {
	domainService domainService
	repo          repository
	bucketName    string
	baseDomain    string
	shorthandId   string
}

func NewService(repo repository, bucketName string, domainService domainService, baseDomain, shorthandId string) *Service {
	return &Service{
		repo:          repo,
		bucketName:    bucketName,
		domainService: domainService,
		baseDomain:    baseDomain,
		shorthandId:   shorthandId,
	}
}

func (s *Service) GetFile(ctx context.Context, fileName string) (io.Reader, error) {
	return s.repo.GetObject(ctx, s.bucketName, fileName)
}

func (s *Service) UploadFile(ctx context.Context, filename string, file io.Reader, size int64) (string, []string, []string, error) {
	subdomain := s.GenerateSubdomain()
	var uploaded, skipped []string

	switch filepath.Ext(filename) {
	case ".html":
		if err := s.uploadSingleFile(ctx, subdomain, filename, file, size); err != nil {
			return "", nil, nil, err
		}
		uploaded = []string{filename}
	case ".zip":
		var err error
		uploaded, skipped, err = s.uploadZipArchive(ctx, subdomain, file, size)
		if err != nil {
			return "", nil, nil, err
		}
	default:
		return "", nil, nil, ErrInvalidFileType
	}

	if err := s.CreateSubdomain(ctx, subdomain); err != nil {
		return "", nil, nil, err
	}

	return subdomain + "." + s.baseDomain, uploaded, skipped, nil
}

func (s *Service) uploadSingleFile(ctx context.Context, subdomain, filename string, file io.Reader, size int64) error {
	if !isAllowedExtension(filename) {
		return ErrInvalidFileType
	}
	return s.repo.PutObject(ctx, s.bucketName, subdomain+"/"+filename, file, size)
}

func (s *Service) uploadZipArchive(ctx context.Context, subdomain string, file io.Reader, size int64) ([]string, []string, error) {
	fileReaderAt, ok := file.(io.ReaderAt)
	if !ok {
		return nil, nil, errors.New("zip file must be seekable")
	}

	zipReader, err := zip.NewReader(fileReaderAt, size)
	if err != nil {
		return nil, nil, err
	}

	var uploadedFiles, skippedFiles []string

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

		err = s.uploadSingleFile(ctx, subdomain, filename, fileReader, int64(f.UncompressedSize64))
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
	return s.domainService.AddSubdomain(ctx, s.baseDomain, subdomain, s.shorthandId)
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
