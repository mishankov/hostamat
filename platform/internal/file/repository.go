package file

import (
	"context"
	"io"

	"hostamat/internal/adapters/s3"
)

type Repository struct {
	s3Client *s3.S3
}

func NewRepository(s3Client *s3.S3) *Repository {
	return &Repository{
		s3Client: s3Client,
	}
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
