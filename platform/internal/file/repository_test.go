package file_test

import (
	"context"
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

	if mock.makeBucketCallCount != 1 {
		t.Fatalf("MakeBucket expected to be called once, got: %d", mock.makeBucketCallCount)
	}
}

func TestRepository_ObjectExists(t *testing.T) {
	mock := &mockS3{}
	r := file.NewRepository(mock)

	_, err := r.ObjectExists(context.Background(), "test-bucket", "test-file")
	if err != nil {
		t.Fatalf("ObjectExists() unexpected error: %v", err)
	}

	if mock.objectExistsCallCount != 1 {
		t.Errorf("ObjectExists() expected to call s3.ObjectExists once, got %d calls", mock.objectExistsCallCount)
	}
}

func TestRepository_GetObject(t *testing.T) {
	mock := &mockS3{}
	r := file.NewRepository(mock)

	_, err := r.GetObject(context.Background(), "test-bucket", "test-file")
	if err != nil {
		t.Fatalf("GetObject() unexpected error: %v", err)
	}

	if mock.getObjectCallCount != 1 {
		t.Errorf("GetObject() expected to call s3.GetObject once, got %d calls", mock.getObjectCallCount)
	}
}

func TestRepository_PutObject(t *testing.T) {
	mock := &mockS3{}
	r := file.NewRepository(mock)

	err := r.PutObject(context.Background(), "test-bucket", "test-file", nil, 0)
	if err != nil {
		t.Fatalf("PutObject() unexpected error: %v", err)
	}

	if mock.putObjectCallCount != 1 {
		t.Errorf("PutObject() expected to call s3.PutObject once, got %d calls", mock.putObjectCallCount)
	}
}

type mockS3 struct {
	makeBucketCallCount   int
	objectExistsCallCount int
	getObjectCallCount    int
	putObjectCallCount    int
	cleanFolderCallCount  int
}

func (m *mockS3) MakeBucket(ctx context.Context, bucketName string) error {
	m.makeBucketCallCount++
	return nil
}

func (m *mockS3) ObjectExists(ctx context.Context, bucketName, fileName string) (bool, error) {
	m.objectExistsCallCount++
	return false, nil
}

func (m *mockS3) GetObject(ctx context.Context, bucketName, fileName string) (io.Reader, error) {
	m.getObjectCallCount++
	return nil, nil
}

func (m *mockS3) PutObject(ctx context.Context, bucketName, fileName string, file io.Reader, size int64) error {
	m.putObjectCallCount++
	return nil
}

func (m *mockS3) CleanFolder(ctx context.Context, bucketName, folderName string) error {
	m.cleanFolderCallCount++
	return nil
}
