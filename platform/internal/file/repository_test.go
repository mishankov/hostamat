package file_test

import (
	"bytes"
	"context"
	"errors"
	"hostamat/internal/file"
	"io"
	"testing"
)

func TestRepository_Migrate(t *testing.T) {
	r := file.NewRepository(&mockS3{})

	err := r.Migrate(context.Background(), "")

	if err != nil {
		t.Fatalf("error expected to be nil, got: %s", err)
	}
}

type mockS3 struct {
	makeBucketCallCount   int
	objectExistsCallCount int
	getObjectCallCount    int
	putObjectCallCount    int

	objects map[string]map[string][]byte // bucket -> key -> content
}

func (m *mockS3) MakeBucket(ctx context.Context, bucketName string) error {
	m.makeBucketCallCount++
	if m.objects == nil {
		m.objects = make(map[string]map[string][]byte)
	}
	m.objects[bucketName] = make(map[string][]byte)
	return nil
}

func (m *mockS3) ObjectExists(ctx context.Context, bucketName, fileName string) (bool, error) {
	m.objectExistsCallCount++
	if bucket, ok := m.objects[bucketName]; ok {
		_, exists := bucket[fileName]
		return exists, nil
	}
	return false, nil
}

func (m *mockS3) GetObject(ctx context.Context, bucketName, fileName string) (io.Reader, error) {
	m.getObjectCallCount++
	if bucket, ok := m.objects[bucketName]; ok {
		if content, exists := bucket[fileName]; exists {
			return bytes.NewReader(content), nil
		}
	}
	return nil, errors.New("object not found")
}

func (m *mockS3) PutObject(ctx context.Context, bucketName, fileName string, file io.Reader, size int64) error {
	m.putObjectCallCount++
	if _, ok := m.objects[bucketName]; !ok {
		return errors.New("bucket does not exist")
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	m.objects[bucketName][fileName] = content
	return nil
}
