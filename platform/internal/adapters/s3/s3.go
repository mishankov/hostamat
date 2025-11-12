package s3

import (
	"context"
	"fmt"
	"io"

	"github.com/platforma-dev/platforma/log"

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

func (s3 *S3) CleanBucket(ctx context.Context, bucketName string) error {
	// List all objects in the bucket
	objectsCh := s3.client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Recursive: true,
	})

	// Remove all objects
	for obj := range objectsCh {
		if obj.Err != nil {
			return fmt.Errorf("error listing objects: %w", obj.Err)
		}
		err := s3.client.RemoveObject(ctx, bucketName, obj.Key, minio.RemoveObjectOptions{})
		if err != nil {
			return fmt.Errorf("error removing object %s: %w", obj.Key, err)
		}
	}

	return nil
}

func (s3 *S3) CleanFolder(ctx context.Context, bucketName, folderName string) error {
	// List all objects in the folder
	objectsCh := s3.client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    folderName,
		Recursive: true,
	})

	// Remove all objects
	for obj := range objectsCh {
		if obj.Err != nil {
			return fmt.Errorf("error listing objects: %w", obj.Err)
		}
		err := s3.client.RemoveObject(ctx, bucketName, obj.Key, minio.RemoveObjectOptions{})
		if err != nil {
			return fmt.Errorf("error removing object %s: %w", obj.Key, err)
		}
	}

	return nil
}

func (s3 *S3) ObjectExists(ctx context.Context, bucketName, fileName string) (bool, error) {
	_, err := s3.client.StatObject(ctx, bucketName, fileName, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}

	return true, nil
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
