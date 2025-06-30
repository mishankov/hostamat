package main

import (
	"context"
	"hostamat/internal/adapters/caddy"
	"hostamat/internal/adapters/config"
	"hostamat/internal/adapters/s3"

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
	"net/http"
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
	caddyAPI := caddy.New(httpClient, conf.Caddy.BaseURL)

	// S3 client
	s3Client, err := s3.New(conf.S3.BaseURL, conf.S3.Login, conf.S3.Password, false)
	if err != nil {
		logger.Error("failed to create s3 client", "error", err)
		os.Exit(1)
	}

	// Middlewares
	requestLogger := middlewares.NewRequestLogger(logger)

	// Handlers
	healthcheckHandler := handlers.NewHealthcheckHandler()
	uploadFileHandler := uploadfilehandler.New(logger, caddyAPI, s3Client, conf.S3.Bucket, conf.PublicBaseDomain, conf.FileServerHost, conf.FileServerPath)
	serveFileHandler := servefilehandler.New(logger, s3Client, conf.S3.Bucket)

	// Handle groups
	filesAPI := httpserver.NewHandleGroup()
	filesAPI.Handle("/upload", uploadFileHandler)

	// Services
	httpServer := httpserver.New(logger, conf.Port, 10*time.Second, map[string]http.Handler{
		"/healthcheck":      healthcheckHandler,
		"/files/api":        filesAPI,
		conf.FileServerPath: serveFileHandler,
	})

	httpServer.Use(requestLogger)

	cleanupService := cleanup.New(s3Client, conf.S3.Bucket, caddyAPI)
	_ = scheduler.New(logger, 30*time.Second, cleanupService)

	// Application
	app := application.New()

	// Ensure that bucket exists before starting application
	app.OnStart(application.RunnerFunc(func(ctx context.Context) error {
		return s3Client.MakeBucket(ctx, conf.S3.Bucket)
	}), application.StartupTaskConfig{Name: "createBucket", AbortOnError: false})

	app.Register("http", httpServer)
	// app.Register("cleanup", cleanupJob)

	if err := app.Run(ctx); err != nil {
		logger.Error("Applicaton finished with error", "error", err)
	}
}
