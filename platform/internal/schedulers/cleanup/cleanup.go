package cleanup

import (
	"context"
	"fmt"
	"hostamat/internal/adapters/caddy"
	"hostamat/internal/adapters/s3"
	"hostamat/internal/adapters/timeweb"
	"strings"
)

type Cleanup struct {
	s3             *s3.S3
	bucketName     string
	webserver      *caddy.CaddyAPI
	domainprovider *timeweb.TimewebAPI
}

func New(s3 *s3.S3, bucketName string, webserver *caddy.CaddyAPI, domainprovider *timeweb.TimewebAPI) *Cleanup {
	return &Cleanup{s3: s3, bucketName: bucketName, webserver: webserver, domainprovider: domainprovider}
}

func (c *Cleanup) Run(ctx context.Context) error {
	// Recreate bucket from S3
	if err := c.s3.CleanBucket(ctx, c.bucketName); err != nil {
		return fmt.Errorf("failed to delete S3 bucket: %w", err)
	}

	// Get and delete all subdomains from Timeweb
	domain, err := c.domainprovider.GetDomain(ctx, c.bucketName)
	if err != nil {
		return fmt.Errorf("failed to get domain info: %w", err)
	}

	for _, subdomain := range domain.Subdomains {
		subdomainName := strings.TrimSuffix(subdomain.Fqdn, "."+c.bucketName)
		if err := c.domainprovider.DeleteSubdomain(ctx, c.bucketName, subdomainName); err != nil {
			return fmt.Errorf("failed to delete subdomain %s: %w", subdomainName, err)
		}
	}

	// TODO: Delete all subdomains from Caddy

	return nil
}
