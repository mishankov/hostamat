package main

import (
	"context"
	"hostamat/internal/adapters/config"
	"hostamat/internal/adapters/timeweb"
	"hostamat/internal/infra/logger"
	"hostamat/pkg/platforma/httpclient"
	"log"
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
	logger := logger.New(conf.Logger.Type, slog.LevelInfo)
	httpClient := httpclient.New(logger)
	timewebAPI := timeweb.New(httpClient, conf.Timeweb.BaseURL, conf.Timeweb.APIToken)

	domain, err := timewebAPI.GetDomain(ctx, conf.Timeweb.BaseDomain)
	if err != nil {
		log.Fatal("Error getting domain:", err)
	}

	log.Println("Subdomains amount:", len(domain.Subdomains))

	if len(domain.Subdomains) == 0 {
		return
	}

	log.Println("Deleting subdomains...")

	for i, subdomain := range domain.Subdomains {
		log.Println("Deleting subdomain", i)

		err := timewebAPI.DeleteSubdomain(ctx, domain.Fqdn, strings.Split(subdomain.Fqdn, ".")[0])
		if err != nil {
			log.Fatal("Error deleting domain:", err)
		}
	}

	domain, err = timewebAPI.GetDomain(ctx, conf.Timeweb.BaseDomain)
	if err != nil {
		log.Fatal("Error getting domain:", err)
	}
	log.Println("Subdomains amount:", len(domain.Subdomains))
}
