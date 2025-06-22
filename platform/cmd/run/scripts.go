package main

import (
	"context"
	"hostamat/internal/adapters/config"
	"hostamat/internal/adapters/timeweb"
	"hostamat/pkg/platforma/httpclient"
	"hostamat/pkg/platforma/log"
	"log/slog"
	"os"
	"strings"
)

func main() {
	command := os.Args[1]

	switch command {
	case "flush_subdomains":
		flushSubdomains()
	}
}

func flushSubdomains() {
	ctx := context.Background()
	conf := config.New()
	logger := log.New(conf.Logger.Type, slog.LevelInfo, nil)
	httpClient := httpclient.New(logger)
	timewebAPI := timeweb.New(httpClient, conf.Timeweb.BaseURL, conf.Timeweb.APIToken)

	domain, err := timewebAPI.GetDomain(ctx, conf.Timeweb.BaseDomain)
	if err != nil {
		logger.Error("Error getting domain", "error", err)
	}

	logger.Info("Subdomains amount", "amount", len(domain.Subdomains))

	if len(domain.Subdomains) == 0 {
		return
	}

	logger.Info("Deleting subdomains...")

	for i, subdomain := range domain.Subdomains {
		logger.Info("Deleting subdomain", "index", i)

		err := timewebAPI.DeleteSubdomain(ctx, domain.Fqdn, strings.Split(subdomain.Fqdn, ".")[0])
		if err != nil {
			logger.Error("Error deletting domain", "error", err)
		}
	}

	domain, err = timewebAPI.GetDomain(ctx, conf.Timeweb.BaseDomain)
	if err != nil {
		logger.Error("Error getting domain", "error", err)
	}
	logger.Info("Subdomains amount", "amount", len(domain.Subdomains))
}
