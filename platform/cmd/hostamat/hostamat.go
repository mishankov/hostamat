package main

import (
	"context"
	"fmt"
	"hostamat/internal/adapters/caddy"
	"hostamat/internal/adapters/config"
	"hostamat/internal/adapters/s3"
	"hostamat/internal/file"
	"hostamat/internal/middlewares"
	"hostamat/internal/site"
	"hostamat/internal/subdomain"
	"hostamat/internal/uploadsession"
	"os"
	"time"

	"hostamat/internal/httpclient"

	"github.com/platforma-dev/platforma/application"
	"github.com/platforma-dev/platforma/auth"
	"github.com/platforma-dev/platforma/database"
	"github.com/platforma-dev/platforma/httpserver"
	"github.com/platforma-dev/platforma/log"
	"github.com/platforma-dev/platforma/session"
)

func main() {
	ctx := context.Background()

	// Basic components
	conf := config.New()
	logger := log.New(os.Stdout, conf.Logger.Type, conf.Logger.Level, nil)
	httpClient := httpclient.New(10 * time.Second)

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
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", conf.DB.User, conf.DB.Password, conf.DB.Host, conf.DB.Port, conf.DB.DB)
	db, err := database.New(connStr)
	if err != nil {
		logger.Error("failed to create database", "error", err)
		os.Exit(1)
	}
	conn := db.Connection()
	defer conn.Close()

	// Domains
	subdomainDomain := subdomain.New(conn, caddyAPI, conf.PublicBaseDomain, conf.Caddy.ShorthandId)
	fileDomain := file.New(s3Client, conf.S3.Bucket)
	uploadSessionDomain := uploadsession.New(conn)
	siteDomain := site.New(conn, subdomainDomain.Service, fileDomain.Service, uploadSessionDomain.Service)
	sessionDomain := session.New(conn)
	authDomain := auth.New(conn, sessionDomain.Service, "hsessionid", nil, nil)

	// Middlewares
	recoverMiddleware := httpserver.NewRecoverMiddleware()
	traceId := httpserver.NewTraceIDMiddleware(nil, "Hostamat-Trace-Id")
	requestLogger := middlewares.NewRequestLogger(logger)

	// API with authentication
	apiGroup := httpserver.NewHandlerGroup()
	apiGroup.Use(authDomain.Middleware)

	apiGroup.Handle("/files", siteDomain.FileHandlerGroup)
	apiGroup.Handle("/sites", siteDomain.SiteHandlerGroup)
	apiGroup.Handle("/uploadsessions", siteDomain.UploadSessionHandlerGroup)

	// HTTP server
	httpServer := httpserver.New(conf.Port, 10*time.Second)

	httpServer.HandleGroup("/api", apiGroup)
	httpServer.HandleGroup("/auth", authDomain.HandleGroup)
	// HandleGroup used here because client would request files by additional path
	httpServer.HandleGroup(conf.FileServerPath, siteDomain.FileServeHandler)

	httpServer.Use(recoverMiddleware, traceId, requestLogger)

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

	app.RegisterDatabase("default", db)

	app.RegisterDomain("subdomain", "default", subdomainDomain)
	app.RegisterDomain("site", "default", siteDomain)
	app.RegisterDomain("file", "", fileDomain)
	app.RegisterDomain("auth", "default", authDomain)
	app.RegisterDomain("session", "default", sessionDomain)
	app.RegisterDomain("uploadsession", "default", uploadSessionDomain)

	app.RegisterService("http", httpServer)
	app.RegisterService("deleteFolderQueue", fileDomain.DeleteFolderQueue)

	if err := app.Run(ctx); err != nil {
		logger.Error("Applicaton finished with error", "error", err)
	}
}
