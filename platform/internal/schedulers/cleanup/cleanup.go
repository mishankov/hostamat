package cleanup

import (
	"context"
	"fmt"
	"hostamat/internal/adapters/caddy"
	"hostamat/internal/adapters/s3"
)

type Cleanup struct {
	s3          *s3.S3
	bucketName  string
	shorthandId string
	caddy       *caddy.CaddyAPI
}

func New(s3 *s3.S3, caddy *caddy.CaddyAPI, bucketName, shorthandId string) *Cleanup {
	return &Cleanup{s3: s3, bucketName: bucketName, caddy: caddy, shorthandId: shorthandId}
}

func (c *Cleanup) Run(ctx context.Context) error {
	// Recreate bucket from S3
	if err := c.s3.CleanBucket(ctx, c.bucketName); err != nil {
		return fmt.Errorf("failed to delete S3 bucket: %w", err)
	}

	// TODO: Delete all subdomains from Caddy
	if err := c.caddy.ClearSubdomains(ctx, c.shorthandId); err != nil {
		return fmt.Errorf("failed to delete Caddy subdomains: %w", err)
	}

	return nil
}
