package file_test

import (
	"bytes"
	"context"
	"errors"
	"hostamat/internal/file"
	"hostamat/internal/site"
	"io"
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

type mockSiteService struct {
	sites map[string]*site.Site
}

func (m *mockSiteService) GetByDomain(ctx context.Context, domain string) (*site.Site, error) {
	return m.sites[domain], nil
}
