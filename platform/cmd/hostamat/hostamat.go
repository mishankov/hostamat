package main

import (
	"context"
	"fmt"
	"hostamat/internal/adapters/caddy"
	"hostamat/internal/adapters/config"
	"hostamat/internal/adapters/s3"
	"hostamat/internal/adapters/yandexoauth"
	"hostamat/internal/customdomain"
	"hostamat/internal/feedback"
	"hostamat/internal/file"
	"hostamat/internal/formcollection"
	"hostamat/internal/hostsync"
	"hostamat/internal/site"
	"hostamat/internal/subdomain"
	"hostamat/internal/uploadsession"
	"hostamat/internal/user"
	"hostamat/internal/userauth"
	"os"
	"strings"
	"time"

	"hostamat/internal/httpclient"

	"github.com/platforma-dev/platforma/application"
	"github.com/platforma-dev/platforma/database"
	"github.com/platforma-dev/platforma/httpserver"
	"github.com/platforma-dev/platforma/log"
	"github.com/platforma-dev/platforma/session"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	ctx := context.Background()
	conf := config.New()

	// Logging
	logger := log.NewWideEventLogger(
		os.Stdout,
		log.NewDefaultSampler(1*time.Second, 400, .25),
		"json",
		nil,
	)
	loggerMiddleware := log.NewWideEventMiddleware(logger, "http_server", nil)

	log.SetDefault(logger)

	// Basic components
	httpClient := httpclient.New(10 * time.Second)

	logger.DebugContext(ctx, "config", "config", conf)

	// Setting default logger for platforma pkg

	// API adapters
	caddyAPI := caddy.New(httpClient, conf.Caddy.BaseURL, "srv0") // Default server name
	hostSyncService := hostsync.NewService(caddyAPI, conf.Caddy.ShorthandId)

	// S3 client
	s3Client, err := s3.New(conf.S3.BaseURL, conf.S3.Login, conf.S3.Password, false)
	if err != nil {
		logger.Error("failed to create s3 client", "error", err)
		os.Exit(1)
	}

	// Database
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", conf.DB.User, conf.DB.Password, conf.DB.Host, conf.DB.Port, conf.DB.DB)
	db, err := database.New(connStr)
	if err != nil {
		logger.Error("failed to create database", "error", err)
		os.Exit(1)
	}
	conn := db.Connection()
	defer conn.Close()

	// Domains
	siteRepository := site.NewRepository(conn)
	subdomainDomain := subdomain.New(conn, hostSyncService, conf.PublicBaseDomain)
	fileDomain := file.New(s3Client, conf.S3.Bucket)
	uploadSessionDomain := uploadsession.New(conn)
	customDomain := customdomain.New(conn, siteRepository, hostSyncService, customdomain.DNSConfig{
		VerificationPrefix: conf.CustomDomains.VerificationPrefix,
		CNAMETarget:        conf.CustomDomains.CNAMETarget,
	})
	feedbackDomain := feedback.New(conn)
	formCollectionDomain := formcollection.New(conn, fileDomain.Service)
	siteDomain := site.New(conn, subdomainDomain.Service, customDomain.Service, fileDomain.Service, uploadSessionDomain.Service, formCollectionDomain.Service)
	sessionDomain := session.New(conn)
	userDomain := user.New(siteDomain.Service, uploadSessionDomain.Service, feedbackDomain.Service, formCollectionDomain.Service)
	yandexClient := yandexoauth.New(
		httpClient,
		conf.Auth.YandexAuthorizeURL,
		conf.Auth.YandexTokenURL,
		conf.Auth.YandexUserInfoURL,
		conf.Auth.YandexClientID,
		conf.Auth.YandexClientSecret,
		conf.Auth.YandexRedirectURI,
	)
	authDomain := userauth.New(conn, sessionDomain.Service, yandexClient, "hsessionid", conf.Auth.YandexRedirectURI, userDomain.CleanupQueue)

	// Middlewares
	recoverMiddleware := httpserver.NewRecoverMiddleware()
	traceId := log.NewTraceIDMiddleware(nil, "Hostamat-Trace-Id")

	// API with authentication
	apiGroup := httpserver.NewHandlerGroup()
	apiGroup.Use(authDomain.Middleware)

	apiGroup.Mount("/files", siteDomain.FileHandlerGroup)
	apiGroup.Mount("/sites", siteDomain.SiteHandlerGroup)
	apiGroup.Mount("/domains", customDomain.HandlerGroup)
	apiGroup.Mount("/feedback", feedbackDomain.HandlerGroup)
	apiGroup.Mount("/forms", formCollectionDomain.HandlerGroup)
	apiGroup.Mount("/uploadsessions", siteDomain.UploadSessionHandlerGroup)

	// HTTP server
	httpServer := httpserver.New(conf.Port, 10*time.Second)

	httpServer.Mount("/api", apiGroup)
	httpServer.Mount("/auth", authDomain.HandleGroup)
	httpServer.Mount("/forms", formCollectionDomain.PublicHandlerGroup)
	if conf.FileServerPath != "" {
		httpServer.Mount(strings.TrimRight(conf.FileServerPath, "/")+"/forms", formCollectionDomain.PublicHandlerGroup)
	}
	// Mount used here because client would request files by additional path
	httpServer.Mount(conf.FileServerPath, siteDomain.FileServeHandler)

	httpServer.Handle("/metrics", promhttp.Handler())

	httpServer.Use(recoverMiddleware, traceId, loggerMiddleware)

	// Application
	app := application.New()
	healthcheckHandler := application.NewHealthCheckHandler(app)
	httpServer.Handle("/health", healthcheckHandler)

	// Migrate File storage
	app.OnStartFunc(func(ctx context.Context) error {
		return fileDomain.Repository.Migrate(ctx, conf.S3.Bucket)
	}, application.StartupTaskConfig{Name: "migrateFile", AbortOnError: true})

	// Ensure that caddy has desired configuration
	app.OnStart(application.RunnerFunc(func(ctx context.Context) error {
		err := caddyAPI.CheckConfigWithId(ctx, conf.Caddy.ShorthandId)
		if err == nil {
			return nil
		}

		return caddyAPI.CreateBaseConfig(ctx, conf.FileServerHost, conf.FileServerPath, conf.Caddy.ShorthandId)
	}), application.StartupTaskConfig{Name: "prepareCaddy", AbortOnError: true})
	app.OnStart(application.RunnerFunc(func(ctx context.Context) error {
		generatedHosts, err := subdomainDomain.Service.ListHosts(ctx)
		if err != nil {
			return err
		}

		customHosts, err := customDomain.Service.ListActiveHosts(ctx)
		if err != nil {
			return err
		}

		hosts := append(generatedHosts, customHosts...)

		return hostSyncService.Sync(ctx, hosts)
	}), application.StartupTaskConfig{Name: "syncHosts", AbortOnError: true})

	app.RegisterDatabase("default", db)

	app.RegisterDomain("subdomain", "default", subdomainDomain)
	app.RegisterDomain("customdomain", "default", customDomain)
	app.RegisterDomain("feedback", "default", feedbackDomain)
	app.RegisterDomain("formcollection", "default", formCollectionDomain)
	app.RegisterDomain("site", "default", siteDomain)
	app.RegisterDomain("file", "", fileDomain)
	app.RegisterDomain("auth", "default", authDomain)
	app.RegisterDomain("session", "default", sessionDomain)
	app.RegisterDomain("uploadsession", "default", uploadSessionDomain)

	app.RegisterService("http", httpServer)
	app.RegisterService("deleteFolderQueue", fileDomain.DeleteFolderQueue)
	app.RegisterService("userCleanupQueue", userDomain.CleanupQueue)
	app.RegisterService("staleCustomDomainScheduler", customDomain.CleanupScheduler)

	if err := app.Run(ctx); err != nil {
		logger.Error("Applicaton finished with error", "error", err)
	}
}
