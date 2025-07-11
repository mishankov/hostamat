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

}

type mockS3 struct {
	makeBucketCallCount   int
	objectExistsCallCount int
	getObjectCallCount    int
	putObjectCallCount    int
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
