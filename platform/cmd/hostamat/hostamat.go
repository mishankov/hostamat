package main

import (
	"context"
	"hostamat/internal/adapters/caddy"
	"hostamat/internal/adapters/config"
	"hostamat/internal/adapters/s3"
	"hostamat/internal/file"
	"hostamat/internal/infra/db"
	"hostamat/internal/middlewares"
	"hostamat/internal/schedulers/cleanup"
	"hostamat/internal/site"
	"hostamat/internal/subdomain"
	"hostamat/pkg/platforma/application"
	"hostamat/pkg/platforma/batteries/healthcheck"
	"hostamat/pkg/platforma/batteries/traceid"
	"hostamat/pkg/platforma/httpclient"
	"hostamat/pkg/platforma/httpserver"
	"hostamat/pkg/platforma/log"
	"hostamat/pkg/platforma/scheduler"
	"os"
	"time"
)

func main() {
	ctx := context.Background()

	// Basic components
	conf := config.New()
	logger := log.New(os.Stdout, conf.Logger.Type, conf.Logger.Level, nil)
	httpClient := httpclient.New(logger)

	logger.DebugContext(ctx, "config", "config", conf)

	// Setting default logger for platforma pkg
	log.SetDefault(logger)

	// API adapters
	caddyAPI := caddy.New(httpClient, conf.Caddy.BaseURL, "srv0") // Default server name

	// S3 client
	s3Client, err := s3.New(conf.S3.BaseURL, conf.S3.Login, conf.S3.Password, false)
	if err != nil {
		logger.Error("failed to create s3 client", "error", err)
		os.Exit(1)
	}

	// Database
	db, err := db.CreateDB("hostamat.db")
	if err != nil {
		logger.Error("failed to create database", "error", err)
		os.Exit(1)
	}

	// Middlewares
	traceId := traceid.NewTraceId(nil, "Hostamat-Trace-Id")
	requestLogger := middlewares.NewRequestLogger(logger)

	// Domain service setup
	subdomainRepo := subdomain.NewRepository(db)
	subdomainService := subdomain.NewService(subdomainRepo, caddyAPI, conf.PublicBaseDomain, conf.Caddy.ShorthandId)

	// Site service setup
	siteRepo := site.NewRepository(db)
	siteService := site.NewService(siteRepo, subdomainService)
	siteCreateHandle := site.NewCreateSiteHandler(siteService)
	siteListHandler := site.NewListSiteHandler(siteService)
	siteDeleteHandler := site.NewDeleteSiteHandler(siteService)

	siteHandlerGroup := httpserver.NewHandlerGroup()
	siteHandlerGroup.Handle("POST /", siteCreateHandle)
	siteHandlerGroup.Handle("GET /", siteListHandler)
	siteHandlerGroup.Handle("DELETE /", siteDeleteHandler)

	// File service setup
	fileRepo := file.NewRepository(s3Client)
	fileService := file.NewService(fileRepo, siteService, conf.S3.Bucket)
	serveFileHandler := file.NewServeHandler(fileService)
	uploadFileHandler := file.NewUploadHandler(fileService)
	deleteHandler := file.NewDeleteHandler(fileService)

	// Handle groups
	filesAPI := httpserver.NewHandlerGroup()
	filesAPI.Handle("POST /", uploadFileHandler)
	filesAPI.Handle("DELETE /", deleteHandler)

	// Services
	httpServer := httpserver.New(conf.Port, 10*time.Second)

	httpServer.Handle("/files/api", filesAPI)
	httpServer.Handle("/sites/api", siteHandlerGroup)
	// HandleGroup used here because client would request files by additional path
	httpServer.HandleGroup(conf.FileServerPath, serveFileHandler)

	httpServer.Use(traceId, requestLogger)

	// Cleanup service to remove files from S3 and subdomains from Caddy
	cleanupService := cleanup.New(s3Client, caddyAPI, conf.S3.Bucket, conf.Caddy.ShorthandId)
	cleanupJob := scheduler.New(30*time.Minute, cleanupService)

	// Application
	app := application.New()
	healthcheckHandler := healthcheck.NewHealthCheckHandler(app)
	httpServer.Handle("/health", healthcheckHandler)

	// Migrate DB
	app.OnStartFunc(subdomainRepo.Migrate, application.StartupTaskConfig{Name: "migrateSubdomain", AbortOnError: true})
	app.OnStartFunc(siteRepo.Migrate, application.StartupTaskConfig{Name: "migrateSite", AbortOnError: true})
	app.OnStartFunc(func(ctx context.Context) error {
		return fileRepo.Migrate(ctx, conf.S3.Bucket)
	}, application.StartupTaskConfig{Name: "migrateFile", AbortOnError: true})

	// Ensure that caddy has desired configuration
	app.OnStart(application.RunnerFunc(func(ctx context.Context) error {
		err := caddyAPI.CheckConfigWithId(ctx, conf.Caddy.ShorthandId)
		if err == nil {
			return nil
		}

		return caddyAPI.CreateBaseConfig(ctx, conf.FileServerHost, conf.FileServerPath, conf.Caddy.ShorthandId)
	}), application.StartupTaskConfig{Name: "prepareCaddy", AbortOnError: true})

	app.Register("http", httpServer)
	app.Register("cleanup", cleanupJob)

	if err := app.Run(ctx); err != nil {
		logger.Error("Applicaton finished with error", "error", err)
	}
}
