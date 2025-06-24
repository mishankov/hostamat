package s3

import (
	"context"
	"hostamat/pkg/platforma/log"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3 struct {
	client *minio.Client
}

func New(url, login, password string, ssl bool) (*S3, error) {
	client, err := minio.New(url, &minio.Options{
		Creds:  credentials.NewStaticV4(login, password, ""),
		Secure: ssl,
	})

	return &S3{client: client}, err
}

func (s3 *S3) MakeBucket(ctx context.Context, name string) error {
	err := s3.client.MakeBucket(ctx, name, minio.MakeBucketOptions{})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := s3.client.BucketExists(ctx, name)
		if errBucketExists == nil && exists {
			log.DebugContext(ctx, "Already own bucket", "bucket", name)
		} else {
			return err
		}
	}

	return nil
}

func (s3 *S3) PutObject(ctx context.Context, bucketName, fileName string, file io.Reader, fileSize int64) error {
	info, err := s3.client.PutObject(ctx, bucketName, fileName, file, fileSize, minio.PutObjectOptions{})

	log.DebugContext(ctx, "Put object info", "info", info)

	return err
}

func (s3 *S3) GetObject(ctx context.Context, bucketName, fileName string) (io.Reader, error) {
	obj, err := s3.client.GetObject(ctx, bucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	return obj, nil
}
