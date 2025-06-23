package main

import (
	"context"
	"hostamat/internal/adapters/caddy"
	"hostamat/internal/adapters/config"
	"hostamat/internal/adapters/timeweb"
	"hostamat/internal/handlers"
	"hostamat/internal/middlewares"
	"hostamat/pkg/platforma/application"
	"hostamat/pkg/platforma/httpclient"
	"hostamat/pkg/platforma/httpserver"
	"hostamat/pkg/platforma/httpserver/handlegroup"
	"hostamat/pkg/platforma/log"
	"net/http"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	ctx := context.Background()

	// Basic components
	conf := config.New()
	logger := log.New(conf.Logger.Type, conf.Logger.Level, nil)
	httpClient := httpclient.New(logger)

	// Setting default logger for platforma pkg
	log.SetDefault(logger)

	// API adapters
	timewebAPI := timeweb.New(httpClient, conf.Timeweb.BaseURL, conf.Timeweb.APIToken)
	caddyAPI := caddy.New(httpClient, conf.Caddy.BaseURL, conf.Caddy.ServerHost, conf.BaseDomain)

	s3Client, err := minio.New(conf.S3.BaseURL, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.S3.Login, conf.S3.Password, ""),
		Secure: false,
	})

	if err != nil {
		logger.Error("failed to create s3 client", "error", err)
	}

	// Handlers
	uploadFileHandler := handlers.NewUploadFileHandler(logger, timewebAPI, caddyAPI, s3Client, conf.S3.Bucket, conf.Timeweb.BaseDomain)

	// Middlewares
	requestLogger := middlewares.NewRequestLogger(logger)

	// Handle groups
	fileHandlers := handlegroup.New()
	fileHandlers.Use(requestLogger)
	fileHandlers.Handle("/upload", uploadFileHandler)

	// Services
	api := httpserver.New(conf.Port, map[string]http.Handler{
		"/files": fileHandlers,
		"/healthcheck": http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		}),
	})

	// testScheduler := scheduler.New(3*time.Second, func(ctx context.Context) error {
	// 	logger.DebugContext(ctx, "hey")
	// 	return nil
	// })

	// Application
	app := application.New()
	app.Register("API", api)
	// app.Register("Sheduler", testScheduler)

	if err := app.Run(ctx); err != nil {
		logger.Error("Applicaton finished with error", "error", err)
	}
}
