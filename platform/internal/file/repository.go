package file

import (
	"context"
	"io"

	"github.com/platforma-dev/platforma/log"
)

type s3 interface {
	MakeBucket(ctx context.Context, bucketName string) error
	ObjectExists(ctx context.Context, bucketName, fileName string) (bool, error)
	GetObject(ctx context.Context, bucketName, fileName string) (io.Reader, error)
	PutObject(ctx context.Context, bucketName, fileName string, file io.Reader, size int64) error
	CleanFolder(ctx context.Context, bucketName, folderName string) error
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
	log.DebugContext(ctx, "get object from s3", "bucket", bucketName, "file", fileName)

	exist, err := r.ObjectExists(ctx, bucketName, fileName)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, ErrFileNotFound
	}

	f, err := r.s3Client.GetObject(ctx, bucketName, fileName)
	return f, err
}

func (r *Repository) PutObject(ctx context.Context, bucketName, fileName string, file io.Reader, size int64) error {
	return r.s3Client.PutObject(ctx, bucketName, fileName, file, size)
}

func (r *Repository) DeleteFromFolder(ctx context.Context, bucketName, folderName string) error {
	return r.s3Client.CleanFolder(ctx, bucketName, folderName)
}
