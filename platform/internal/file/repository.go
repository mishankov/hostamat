package file

import (
	"context"
	"io"
)

type s3 interface {
	MakeBucket(ctx context.Context, bucketName string) error
	ObjectExists(ctx context.Context, bucketName, fileName string) (bool, error)
	GetObject(ctx context.Context, bucketName, fileName string) (io.Reader, error)
	PutObject(ctx context.Context, bucketName, fileName string, file io.Reader, size int64) error
}

type Repository struct {
	s3Client s3
}

func NewRepository(s3Client s3) *Repository {
	return &Repository{
		s3Client: s3Client,
	}
}

func (r *Repository) Migrate(ctx context.Context, bucketName string) error {
	return r.s3Client.MakeBucket(ctx, bucketName)
}

func (r *Repository) ObjectExists(ctx context.Context, bucketName, fileName string) (bool, error) {
	return r.s3Client.ObjectExists(ctx, bucketName, fileName)
}

func (r *Repository) GetObject(ctx context.Context, bucketName, fileName string) (io.Reader, error) {
	return r.s3Client.GetObject(ctx, bucketName, fileName)
}

func (r *Repository) PutObject(ctx context.Context, bucketName, fileName string, file io.Reader, size int64) error {
	return r.s3Client.PutObject(ctx, bucketName, fileName, file, size)
}
