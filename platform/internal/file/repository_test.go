package file_test

import (
	"context"
	"errors"
	"hostamat/internal/file"
	"io"
	"testing"
)

func TestRepository_Migrate(t *testing.T) {
	mock := &mockS3{}
	r := file.NewRepository(mock)

	err := r.Migrate(context.Background(), "")

	if err != nil {
		t.Fatalf("error expected to be nil, got: %s", err)
	}
}

func TestRepository_ObjectExists(t *testing.T) {
	mock := &mockS3{}
	r := file.NewRepository(mock)

	_, err := r.ObjectExists(context.Background(), "test-bucket", "test-file")
	if err != nil {
		t.Fatalf("ObjectExists() unexpected error: %v", err)
	}
}

func TestRepository_GetObject(t *testing.T) {
	mock := &mockS3{}
	r := file.NewRepository(mock)

	_, err := r.GetObject(context.Background(), "test-bucket", "test-file")
	if err == nil {
		t.Fatalf("GetObject() expected error when file does not exist, got nil")
	}
}

func TestRepository_PutObject(t *testing.T) {
	mock := &mockS3{}
	r := file.NewRepository(mock)

	err := r.PutObject(context.Background(), "test-bucket", "test-file", nil, 0)
	if err != nil {
		t.Fatalf("PutObject() unexpected error: %v", err)
	}
}

type mockS3 struct {
}

func (m *mockS3) MakeBucket(ctx context.Context, bucketName string) error {
	return nil
}

func (m *mockS3) ObjectExists(ctx context.Context, bucketName, fileName string) (bool, error) {
	return false, nil
}

func (m *mockS3) GetObject(ctx context.Context, bucketName, fileName string) (io.Reader, error) {
	return nil, errors.New("file not found")
}

func (m *mockS3) PutObject(ctx context.Context, bucketName, fileName string, file io.Reader, size int64) error {
	return nil
}

func (m *mockS3) CleanFolder(ctx context.Context, bucketName, folderName string) error {
	return nil
}
