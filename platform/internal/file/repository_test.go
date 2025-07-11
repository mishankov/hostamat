package file_test

import (
	"context"
	"hostamat/internal/file"
	"testing"
)

func Testrepository_Migrate(t *testing.T) {
	r := file.NewRepository(&mockS3{})

}

type mockS3 struct {
	makeBucketCallCount int
}

func (m *mockS3) MakeBucket(ctx context.Context, bucketName string) error {
	m.makeBucketCallCount++
	return nil
}
