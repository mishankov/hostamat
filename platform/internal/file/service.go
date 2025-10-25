package file

import (
	"archive/zip"
	"context"
	"errors"
	"io"
	"path/filepath"
	"strings"

	"github.com/mishankov/platforma/log"
)

type repository interface {
	ObjectExists(ctx context.Context, bucketName, fileName string) (bool, error)
	GetObject(ctx context.Context, bucketName, fileName string) (io.Reader, error)
	PutObject(ctx context.Context, bucketName, fileName string, file io.Reader, size int64) error
	DeleteFromFolder(ctx context.Context, bucketName, folderName string) error
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

func (s *Service) GetFile(ctx context.Context, path string) (io.Reader, error) {
	path = strings.TrimPrefix(path, "/")

	log.DebugContext(ctx, "get file", "path", path)

	return s.repo.GetObject(ctx, s.bucketName, path)
}

func (s *Service) UploadFile(ctx context.Context, path, filename string, file io.Reader, size int64) ([]string, []string, error) {
	var uploaded, skipped []string

	switch filepath.Ext(filename) {
	case ".zip":
		var err error
		uploaded, skipped, err = s.uploadZipArchive(ctx, path, file, size)
		if err != nil {
			return nil, nil, err
		}
	default:
		if err := s.uploadSingleFile(ctx, path, filename, file, size); err != nil {
			return nil, nil, err
		}
		uploaded = []string{filename}
	}

	return uploaded, skipped, nil
}

func (s *Service) DeleteFromFolder(ctx context.Context, folder string) error {
	return s.repo.DeleteFromFolder(ctx, s.bucketName, folder)
}

func (s *Service) uploadSingleFile(ctx context.Context, folder, filename string, file io.Reader, size int64) error {
	if !isAllowedExtension(filename) {
		return ErrInvalidFileType
	}
	return s.repo.PutObject(ctx, s.bucketName, folder+"/"+filename, file, size)
}

func (s *Service) uploadZipArchive(ctx context.Context, folder string, file io.Reader, size int64) ([]string, []string, error) {
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

		err = s.uploadSingleFile(ctx, folder, filename, fileReader, int64(f.UncompressedSize64))
		fileReader.Close()
		if err != nil {
			skippedFiles = append(skippedFiles, filename)
			continue
		}

		uploadedFiles = append(uploadedFiles, filename)
	}

	return uploadedFiles, skippedFiles, nil
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
