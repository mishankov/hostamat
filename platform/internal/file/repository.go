package file

import (
	"context"
	"io"

	"hostamat/internal/adapters/s3"
)

type Repository struct {
	s3Client   *s3.S3
	bucketName string
}

func NewRepository(s3Client *s3.S3, bucketName string) *Repository {
	return &Repository{
		s3Client:   s3Client,
		bucketName: bucketName,
	}
}

func (r *Repository) ObjectExists(ctx context.Context, fileName string) (bool, error) {
	return r.s3Client.ObjectExists(ctx, r.bucketName, fileName)
}

func (r *Repository) GetObject(ctx context.Context, fileName string) (io.Reader, error) {
	return r.s3Client.GetObject(ctx, r.bucketName, fileName)
}

func (r *Repository) PutObject(ctx context.Context, fileName string, file io.Reader, size int64) error {
	return r.s3Client.PutObject(ctx, r.bucketName, fileName, file, size)
}
