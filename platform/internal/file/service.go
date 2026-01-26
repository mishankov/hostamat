package file

import (
	"archive/zip"
	"context"
	"errors"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/platforma-dev/platforma/log"
	"github.com/platforma-dev/platforma/queue"
)

type repository interface {
	ObjectExists(ctx context.Context, bucketName, fileName string) (bool, error)
	GetObject(ctx context.Context, bucketName, fileName string) (io.Reader, error)
	PutObject(ctx context.Context, bucketName, fileName string, file io.Reader, size int64) error
	DeleteFromFolder(ctx context.Context, bucketName, folderName string) error
}

type Service struct {
	repo                  repository
	bucketName            string
	deleteFolderProcessor *queue.Processor[string]
}

func NewService(repo repository, bucketName string) *Service {

	service := &Service{
		repo:       repo,
		bucketName: bucketName,
	}

	deleteFolderQueue := queue.NewChanQueue[string](5, time.Second)
	service.deleteFolderProcessor = queue.New(queue.HandlerFunc[string](func(ctx context.Context, folder string) {
		err := service.DeleteFromFolder(ctx, folder)
		if err != nil {
			log.ErrorContext(ctx, "failed to delete files from folder", "folder", folder, "error", err)
		}
	}), deleteFolderQueue, 2, 3*time.Second)

	return service
}

func (s *Service) GetFile(ctx context.Context, path string) (io.Reader, string, error) {
	path = strings.TrimPrefix(path, "/")

	log.DebugContext(ctx, "get file", "path", path)

	// If path has a known file extension, try exact match only
	if isAllowedExtension(path) {
		reader, err := s.repo.GetObject(ctx, s.bucketName, path)
		return reader, path, err
	}

	// Try resolution in order: exact path, path.html, path/index.html
	pathsToTry := []string{
		path,
		path + ".html",
		path + "/index.html",
	}

	for _, p := range pathsToTry {
		exists, err := s.repo.ObjectExists(ctx, s.bucketName, p)
		if err != nil {
			continue
		}
		if exists {
			reader, err := s.repo.GetObject(ctx, s.bucketName, p)
			return reader, p, err
		}
	}

	// If none found, return error for the original path
	reader, err := s.repo.GetObject(ctx, s.bucketName, path)
	return reader, path, err
}

func (s *Service) UploadFile(ctx context.Context, path, filename string, file io.Reader, size int64) ([]File, []File, error) {
	var uploaded, skipped []File

	switch filepath.Ext(filename) {
	case ".zip":
		var err error
		uploaded, skipped, err = s.uploadZipArchive(ctx, path, file, size)
		if err != nil {
			return nil, nil, err
		}
	default:
		if err := s.uploadSingleFile(ctx, path, filename, file, size); err != nil {
			if errors.Is(err, ErrInvalidFileType) {
				skipped = []File{{filename, 0}}
			}
			return nil, nil, err
		} else {
			uploaded = []File{{filename, int(size)}}
		}
	}

	return uploaded, skipped, nil
}

func (s *Service) DeleteFromFolderAsync(ctx context.Context, folder string) error {
	return s.deleteFolderProcessor.Enqueue(ctx, folder)
}

func (s *Service) DeleteFromFolder(ctx context.Context, folder string) error {
	err := s.repo.DeleteFromFolder(ctx, s.bucketName, folder)
	if err != nil {
		return err
	}

	log.InfoContext(ctx, "files deleted", "folder", folder)
	return nil
}

func (s *Service) uploadSingleFile(ctx context.Context, folder, filename string, file io.Reader, size int64) error {
	if !isAllowedExtension(filename) {
		return ErrInvalidFileType
	}
	return s.repo.PutObject(ctx, s.bucketName, folder+"/"+filename, file, size)
}

func (s *Service) uploadZipArchive(ctx context.Context, folder string, file io.Reader, size int64) ([]File, []File, error) {
	fileReaderAt, ok := file.(io.ReaderAt)
	if !ok {
		return nil, nil, errors.New("zip file must be seekable")
	}

	zipReader, err := zip.NewReader(fileReaderAt, size)
	if err != nil {
		return nil, nil, err
	}

	var uploadedFiles, skippedFiles []File

	for _, f := range zipReader.File {
		if f.FileInfo().IsDir() {
			continue
		}

		filename := strings.Join(strings.Split(f.Name, "/")[1:], "/")
		if !isAllowedFile(filename) {
			skippedFiles = append(skippedFiles, File{filename, 0})
			continue
		}

		fileReader, err := f.Open()
		if err != nil {
			skippedFiles = append(skippedFiles, File{filename, 0})
			continue
		}

		err = s.uploadSingleFile(ctx, folder, filename, fileReader, int64(f.UncompressedSize64))
		fileReader.Close()
		if err != nil {
			skippedFiles = append(skippedFiles, File{filename, 0})
			continue
		}

		uploadedFiles = append(uploadedFiles, File{filename, int(f.UncompressedSize64)})
	}

	return uploadedFiles, skippedFiles, nil
}

func isAllowedFile(filename string) bool {
	if filename == "" || strings.HasPrefix(filename, ".") || strings.Contains(filename, "/.") || strings.Contains(filename, "node_modules") {
		return false
	}
	return isAllowedExtension(filename)
}

func isAllowedExtension(filename string) bool {
	ext := strings.ToLower(strings.TrimPrefix(strings.ToLower(filepath.Ext(filename)), "."))
	switch ext {
	case "html", "css", "js",
		"png", "jpg", "jpeg", "svg", "gif", "ico", "webp",
		"woff2", "woff", "ttf", "eot":
		return true
	default:
		return false
	}
}
