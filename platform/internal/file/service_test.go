package file_test

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"hostamat/internal/file"
	"hostamat/internal/site"
	"io"
	"slices"
	"testing"
)

func TestService_GetFile(t *testing.T) {
	repo := &mockRepo{
		storage: map[string]map[string]mockFile{
			"bucket": {
				"site_id/index.html": {
					reader: &bytes.Buffer{},
					size:   3,
				},
			},
		},
	}
	ss := &mockSiteService{
		sites: map[string]*site.Site{
			"my.domain": {
				Id:        "site_id",
				Subdomain: "my",
				Domain:    "my.domain",
			},
		},
	}

	s := file.NewService(repo, ss, "bucket")

	file, err := s.GetFile(context.TODO(), "index.html", "my.domain")
	if err != nil {
		t.Fatalf("error is not expected, got: %s", err)
	}

	if file == nil {
		t.Fatalf("file was not expected to be nil")
	}
}

func TestService_UploadFile(t *testing.T) {
	repo := &mockRepo{
		storage: map[string]map[string]mockFile{
			"bucket": {},
		},
	}
	ss := &mockSiteService{
		sites: map[string]*site.Site{
			"my.domain": {
				Id:        "site_id",
				Subdomain: "my",
				Domain:    "my.domain",
			},
		},
	}

	s := file.NewService(repo, ss, "bucket")

	t.Run("single file upload success", func(t *testing.T) {
		uploaded, skipped, err := s.UploadFile(context.TODO(), "index.html", &bytes.Buffer{}, 3, "site_id")

		if err != nil {
			t.Fatalf("error was not expected, got: %s", err)
		}

		if len(uploaded) != 1 {
			t.Fatalf("expected to upload 1 file, got: %d", len(uploaded))
		}

		if uploaded[0] != "index.html" {
			t.Fatalf("expected to upload index.html, got: %s", uploaded[0])
		}

		if len(skipped) != 0 {
			t.Fatalf("expected to skip no files, got: %d", len(skipped))
		}
	})

	t.Run("zip file upload success", func(t *testing.T) {
		zipFile, size := createTestZip(t, map[string]string{
			"archive/index.html": "<html></html>",
			"archive/style.css":  "body {}",
		})

		uploaded, skipped, err := s.UploadFile(context.TODO(), "archive.zip", zipFile, size, "site_id")
		if err != nil {
			t.Fatalf("error was not expected, got: %s", err)
		}

		if len(uploaded) != 2 {
			t.Errorf("expected to upload 2 files, got: %d", len(uploaded))
		}

		for _, name := range []string{"index.html", "style.css"} {
			if !slices.Contains(uploaded, name) {
				t.Errorf("%s expected to be uploaded, got: %s", name, uploaded)
			}
		}

		if len(skipped) != 0 {
			t.Errorf("expected to skip no files, got: %d, slice: %s", len(skipped), skipped)
		}
	})

}

type mockFile struct {
	reader io.Reader
	size   int64
}

type mockRepo struct {
	storage map[string]map[string]mockFile
}

func (m *mockRepo) ObjectExists(ctx context.Context, bucketName, fileName string) (bool, error) {
	bucket, ok := m.storage[bucketName]
	if !ok {
		return false, nil
	}

	_, ok = bucket[fileName]
	return ok, nil
}

func (m *mockRepo) GetObject(ctx context.Context, bucketName, fileName string) (io.Reader, error) {
	bucket, ok := m.storage[bucketName]
	if !ok {
		return nil, errors.New("no bucket")
	}

	file, ok := bucket[fileName]
	if !ok {
		return nil, errors.New("no file")
	}

	return file.reader, nil
}

func (m *mockRepo) PutObject(ctx context.Context, bucketName, fileName string, file io.Reader, size int64) error {
	f := mockFile{
		reader: file,
		size:   size,
	}

	_, ok := m.storage[bucketName]
	if !ok {
		m.storage[bucketName] = make(map[string]mockFile)
	}

	m.storage[bucketName][fileName] = f

	return nil
}

func (m *mockRepo) DeleteFromFolder(ctx context.Context, bucketName, folderName string) error {
	delete(m.storage[bucketName], folderName)
	return nil
}

type mockSiteService struct {
	sites map[string]*site.Site
}

func (m *mockSiteService) GetByDomain(ctx context.Context, domain string) (*site.Site, error) {
	return m.sites[domain], nil
}

func createTestZip(t *testing.T, files map[string]string) (io.Reader, int64) {
	t.Helper()
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	for name, content := range files {
		f, err := w.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		_, err = f.Write([]byte(content))
		if err != nil {
			t.Fatal(err)
		}
	}

	err := w.Close()
	if err != nil {
		t.Fatal(err)
	}

	return bytes.NewReader(buf.Bytes()), int64(buf.Len())
}
