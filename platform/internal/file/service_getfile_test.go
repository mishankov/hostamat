package file_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"github.com/platforma-dev/platforma/log"
	"hostamat/internal/file"
)

func TestGetFileResolvesHTMLVariantWithoutExistenceChecks(t *testing.T) {
	t.Parallel()

	repo := &getFileRepoMock{
		getObjectResults: map[string]getObjectResult{
			"site/about": {
				err: file.ErrFileNotFound,
			},
			"site/about.html": {
				reader: bytes.NewBufferString("ok"),
			},
		},
	}

	service := file.NewService(repo, "bucket")

	reader, resolvedPath, err := service.GetFile(testContext(), "site/about")
	if err != nil {
		t.Fatalf("GetFile() error = %v", err)
	}

	if resolvedPath != "site/about.html" {
		t.Fatalf("GetFile() resolvedPath = %q, want %q", resolvedPath, "site/about.html")
	}

	if reader == nil {
		t.Fatal("GetFile() reader = nil")
	}

	if repo.objectExistsCalls != 0 {
		t.Fatalf("ObjectExists() calls = %d, want 0", repo.objectExistsCalls)
	}

	if got := repo.getObjectCalls["site/about"]; got != 1 {
		t.Fatalf("GetObject(%q) calls = %d, want 1", "site/about", got)
	}

	if got := repo.getObjectCalls["site/about.html"]; got != 1 {
		t.Fatalf("GetObject(%q) calls = %d, want 1", "site/about.html", got)
	}
}

func TestGetFileReturnsNotFoundWithoutRetryingOriginalPath(t *testing.T) {
	t.Parallel()

	repo := &getFileRepoMock{
		getObjectResults: map[string]getObjectResult{
			"site/missing": {
				err: file.ErrFileNotFound,
			},
			"site/missing.html": {
				err: file.ErrFileNotFound,
			},
			"site/missing/index.html": {
				err: file.ErrFileNotFound,
			},
		},
	}

	service := file.NewService(repo, "bucket")

	_, resolvedPath, err := service.GetFile(testContext(), "site/missing")
	if !errors.Is(err, file.ErrFileNotFound) {
		t.Fatalf("GetFile() error = %v, want ErrFileNotFound", err)
	}

	if resolvedPath != "site/missing" {
		t.Fatalf("GetFile() resolvedPath = %q, want %q", resolvedPath, "site/missing")
	}

	if got := repo.getObjectCalls["site/missing"]; got != 1 {
		t.Fatalf("GetObject(%q) calls = %d, want 1", "site/missing", got)
	}
}

func TestGetFileReturnsStorageErrorsImmediately(t *testing.T) {
	t.Parallel()

	repo := &getFileRepoMock{
		getObjectResults: map[string]getObjectResult{
			"site/about": {
				err: errors.New("s3 unavailable"),
			},
		},
	}

	service := file.NewService(repo, "bucket")

	_, _, err := service.GetFile(testContext(), "site/about")
	if err == nil || err.Error() != "s3 unavailable" {
		t.Fatalf("GetFile() error = %v, want s3 unavailable", err)
	}

	if len(repo.getObjectCalls) != 1 {
		t.Fatalf("GetFile() probed %d paths, want 1", len(repo.getObjectCalls))
	}
}

type getObjectResult struct {
	reader io.Reader
	err    error
}

type getFileRepoMock struct {
	getObjectResults  map[string]getObjectResult
	getObjectCalls    map[string]int
	objectExistsCalls int
}

func (m *getFileRepoMock) ObjectExists(ctx context.Context, bucketName, fileName string) (bool, error) {
	m.objectExistsCalls++
	return false, nil
}

func (m *getFileRepoMock) GetObject(ctx context.Context, bucketName, fileName string) (io.Reader, error) {
	if m.getObjectCalls == nil {
		m.getObjectCalls = make(map[string]int)
	}
	m.getObjectCalls[fileName]++

	result, ok := m.getObjectResults[fileName]
	if !ok {
		return nil, file.ErrFileNotFound
	}

	return result.reader, result.err
}

func (m *getFileRepoMock) PutObject(ctx context.Context, bucketName, fileName string, fileReader io.Reader, size int64) error {
	return nil
}

func (m *getFileRepoMock) DeleteFromFolder(ctx context.Context, bucketName, folderName string) error {
	return nil
}

func testContext() context.Context {
	return context.WithValue(context.Background(), log.WideEventKey, log.NewEvent("test"))
}
