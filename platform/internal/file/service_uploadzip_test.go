package file_test

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"hostamat/internal/file"
)

func TestUploadZipArchiveRejectsExpandedSizeBeforeUploading(t *testing.T) {
	t.Parallel()

	archive, size := createZip(t, map[string]string{
		"archive/index.html": "123",
		"archive/app.js":     "456",
	})
	repo := newZipRepoMock()
	service := file.NewService(repo, "bucket")

	uploaded, skipped, err := service.UploadZipArchive(testContext(), "site/release", archive, size, 5)
	if !errors.Is(err, file.ErrUploadSizeExceeded) {
		t.Fatalf("UploadZipArchive() error = %v, want ErrUploadSizeExceeded", err)
	}

	if uploaded != nil {
		t.Fatalf("UploadZipArchive() uploaded = %v, want nil", uploaded)
	}

	if skipped != nil {
		t.Fatalf("UploadZipArchive() skipped = %v, want nil", skipped)
	}

	if len(repo.puts) != 0 {
		t.Fatalf("PutObject() calls = %d, want 0", len(repo.puts))
	}
}

func TestUploadZipArchiveDoesNotCountSkippedFilesAgainstLimit(t *testing.T) {
	t.Parallel()

	archive, size := createZip(t, map[string]string{
		"archive/index.html":        "12345",
		"archive/node_modules/a.js": "this skipped dependency is larger than the limit",
	})
	repo := newZipRepoMock()
	service := file.NewService(repo, "bucket")

	uploaded, skipped, err := service.UploadZipArchive(testContext(), "site/release", archive, size, 5)
	if err != nil {
		t.Fatalf("UploadZipArchive() error = %v", err)
	}

	if len(uploaded) != 1 || uploaded[0].Name != "index.html" || uploaded[0].Size != 5 {
		t.Fatalf("UploadZipArchive() uploaded = %v, want index.html size 5", uploaded)
	}

	if len(skipped) != 1 || skipped[0].Name != "node_modules/a.js" {
		t.Fatalf("UploadZipArchive() skipped = %v, want node_modules/a.js", skipped)
	}

	if len(repo.puts) != 1 {
		t.Fatalf("PutObject() calls = %d, want 1", len(repo.puts))
	}

	if repo.puts[0].name != "site/release/index.html" {
		t.Fatalf("PutObject() path = %q, want site/release/index.html", repo.puts[0].name)
	}
}

type zipPut struct {
	name string
	size int64
}

type zipRepoMock struct {
	puts []zipPut
}

func newZipRepoMock() *zipRepoMock {
	return &zipRepoMock{}
}

func (m *zipRepoMock) ObjectExists(ctx context.Context, bucketName, fileName string) (bool, error) {
	return false, nil
}

func (m *zipRepoMock) GetObject(ctx context.Context, bucketName, fileName string) (io.Reader, error) {
	return nil, file.ErrFileNotFound
}

func (m *zipRepoMock) PutObject(ctx context.Context, bucketName, fileName string, fileReader io.Reader, size int64) error {
	m.puts = append(m.puts, zipPut{name: fileName, size: size})
	_, err := io.Copy(io.Discard, fileReader)
	return err
}

func (m *zipRepoMock) DeleteFromFolder(ctx context.Context, bucketName, folderName string) error {
	return nil
}

func createZip(t *testing.T, files map[string]string) (io.Reader, int64) {
	t.Helper()

	buf := new(bytes.Buffer)
	writer := zip.NewWriter(buf)

	for name, content := range files {
		entry, err := writer.Create(name)
		if err != nil {
			t.Fatal(err)
		}

		if _, err := entry.Write([]byte(content)); err != nil {
			t.Fatal(err)
		}
	}

	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	return bytes.NewReader(buf.Bytes()), int64(buf.Len())
}
