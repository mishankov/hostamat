package file

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
)

type Repository struct {
	client     *minio.Client
	bucketName string
}

func NewRepository(client *minio.Client, bucketName string) *Repository {
	return &Repository{
		client:     client,
		bucketName: bucketName,
	}
}

func (r *Repository) ObjectExists(ctx context.Context, fileName string) (bool, error) {
	_, err := r.client.StatObject(ctx, r.bucketName, fileName, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *Repository) GetObject(ctx context.Context, fileName string) (io.Reader, error) {
	return r.client.GetObject(ctx, r.bucketName, fileName, minio.GetObjectOptions{})
}

func (r *Repository) PutObject(ctx context.Context, fileName string, file io.Reader, size int64) error {
	_, err := r.client.PutObject(ctx, r.bucketName, fileName, file, size, minio.PutObjectOptions{})
	return err
}
