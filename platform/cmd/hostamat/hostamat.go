package main

import (
	"context"
	"hostamat/internal/adapters/caddy"
	"hostamat/internal/adapters/config"
	"hostamat/internal/adapters/s3"
	"hostamat/internal/infra/db"
	"hostamat/internal/uploadsession"

	"hostamat/internal/handlers"
	"hostamat/internal/handlers/servefilehandler"
	"hostamat/internal/handlers/uploadfilehandler"
	"hostamat/internal/middlewares"
	"hostamat/internal/schedulers/cleanup"
	"hostamat/pkg/platforma/application"
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

	// Upload Session
	uploadSessionRepo := uploadsession.NewRepository(db)
	uploadSessionService := uploadsession.NewService(uploadSessionRepo)
	uploadSessionHandlerCreate := uploadsession.NewCreateUploadSessionHandler(uploadSessionService)

	// Middlewares
	requestLogger := middlewares.NewRequestLogger(logger)

	// Handlers
	healthcheckHandler := handlers.NewHealthcheckHandler()
	uploadFileHandler := uploadfilehandler.New(logger, caddyAPI, s3Client, conf.S3.Bucket, conf.PublicBaseDomain, conf.Caddy.ShorthandId)
	serveFileHandler := servefilehandler.New(logger, s3Client, conf.S3.Bucket)

	// Handle groups
	filesAPI := httpserver.NewHandlerGroup()
	filesAPI.Handle("/upload", uploadFileHandler)

	// Services
	httpServer := httpserver.New(logger, conf.Port, 10*time.Second)

	httpServer.Handle("/healthcheck", healthcheckHandler)
	httpServer.HandleGroup("/files/api", filesAPI)
	httpServer.Handle("/uploadsession", uploadSessionHandlerCreate)
	// HandleGroup used here because client would request files by additional path
	// TODO: should probably document it in platforma
	httpServer.HandleGroup(conf.FileServerPath, serveFileHandler)

	httpServer.Use(requestLogger)

	// Cleanup service to remove files from S3 and subdomains from Caddy
	cleanupService := cleanup.New(s3Client, caddyAPI, conf.S3.Bucket, conf.Caddy.ShorthandId)
	cleanupJob := scheduler.New(logger, 30*time.Minute, cleanupService)

	// Application
	app := application.New()

	// Migrate DB
	app.OnStart(application.RunnerFunc(func(ctx context.Context) error {
		return uploadSessionRepo.Migrate() //TODO: add ctx
	}), application.StartupTaskConfig{Name: "migrateDB", AbortOnError: true})

	// Ensure that bucket exists before starting application
	app.OnStart(application.RunnerFunc(func(ctx context.Context) error {
		return s3Client.MakeBucket(ctx, conf.S3.Bucket)
	}), application.StartupTaskConfig{Name: "createBucket", AbortOnError: true})

	// Ensure that caddy has desired configuration
	app.OnStart(application.RunnerFunc(func(ctx context.Context) error {
		err := caddyAPI.CheckConfigWithId(ctx, conf.Caddy.ShorthandId)
		if err == nil {
			return nil
		}

		err = caddyAPI.CreateBaseConfig(ctx, conf.FileServerHost, conf.FileServerPath, conf.Caddy.ShorthandId)

		return nil
	}), application.StartupTaskConfig{Name: "prepareCaddy", AbortOnError: true})

	app.Register("http", httpServer)
	app.Register("cleanup", cleanupJob)

	if err := app.Run(ctx); err != nil {
		logger.Error("Applicaton finished with error", "error", err)
	}
}
