package cleanup

import (
	"context"
	"fmt"
	"hostamat/internal/adapters/caddy"
	"hostamat/internal/adapters/s3"
)

type Cleanup struct {
	s3         *s3.S3
	bucketName string
	webserver  *caddy.CaddyAPI
}

func New(s3 *s3.S3, bucketName string, webserver *caddy.CaddyAPI) *Cleanup {
	return &Cleanup{s3: s3, bucketName: bucketName, webserver: webserver}
}

func (c *Cleanup) Run(ctx context.Context) error {
	// Recreate bucket from S3
	if err := c.s3.CleanBucket(ctx, c.bucketName); err != nil {
		return fmt.Errorf("failed to delete S3 bucket: %w", err)
	}

	// TODO: Delete all subdomains from Caddy

	return nil
}
