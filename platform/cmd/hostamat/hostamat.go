package main

import (
	"context"
	"hostamat/internal/adapters/caddy"
	"hostamat/internal/adapters/config"
	"hostamat/internal/adapters/s3"
	"hostamat/internal/adapters/timeweb"
	"hostamat/internal/handlers"
	"hostamat/internal/middlewares"
	"hostamat/pkg/platforma/application"
	"hostamat/pkg/platforma/httpclient"
	"hostamat/pkg/platforma/httpserver"
	"hostamat/pkg/platforma/httpserver/handlegroup"
	"hostamat/pkg/platforma/log"
	"io"
	"net/http"
	"os"
	"strings"
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
	caddyAPI := caddy.New(httpClient, conf.Caddy.BaseURL, conf.Caddy.ServerHost)

	// S3 client
	s3Client, err := s3.New(conf.S3.BaseURL, conf.S3.Login, conf.S3.Password, false)
	if err != nil {
		logger.Error("failed to create s3 client", "error", err)
		os.Exit(1)
	}

	// Handlers
	uploadFileHandler := handlers.NewUploadFileHandler(logger, timewebAPI, caddyAPI, s3Client, conf.S3.Bucket, conf.BaseDomain)

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
		"": http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.InfoContext(ctx, "request", "host", r.Host, "uri", r.RequestURI)

			subdomain := strings.Split(r.Host, ".")[0]

			file, err := s3Client.GetObject(ctx, conf.S3.Bucket, subdomain+"/index.html")
			if err != nil {
				logger.ErrorContext(ctx, "error getting file", "error", err)
			}

			n, err := io.Copy(w, file)
			logger.InfoContext(ctx, "file copy", "bytes", n, "error", err)
		}),
	})

	// Application
	app := application.New()

	app.OnStart(application.RunnerFunc(func(ctx context.Context) error {
		return s3Client.MakeBucket(ctx, conf.S3.Bucket)
	}), application.StartupTaskConfig{Name: "createBucket", AbortOnError: false})

	app.Register("api", api)

	if err := app.Run(ctx); err != nil {
		logger.Error("Applicaton finished with error", "error", err)
	}
}
