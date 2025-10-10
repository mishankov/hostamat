package main

import (
	"context"
	"fmt"
	"hostamat/internal/adapters/caddy"
	"hostamat/internal/adapters/config"
	"hostamat/internal/adapters/s3"
	"hostamat/internal/file"
	"hostamat/internal/middlewares"
	"hostamat/internal/schedulers/cleanup"
	"hostamat/internal/site"
	"hostamat/internal/subdomain"
	"hostamat/pkg/platforma/application"
	"hostamat/pkg/platforma/batteries/auth"
	"hostamat/pkg/platforma/batteries/healthcheck"
	"hostamat/pkg/platforma/batteries/session"
	"hostamat/pkg/platforma/batteries/traceid"
	"hostamat/pkg/platforma/batteries/user"
	"hostamat/pkg/platforma/database"
	"hostamat/pkg/platforma/httpclient"
	"hostamat/pkg/platforma/httpserver"
	"hostamat/pkg/platforma/log"
	"hostamat/pkg/platforma/scheduler"
	"os"
	"time"

	_ "github.com/lib/pq"
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
	db, err := database.New("postgres", connStr)
	if err != nil {
		logger.Error("failed to create database", "error", err)
		os.Exit(1)
	}

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

	// Session service setup
	sessionRepo := session.NewRepository(db)
	sessionService := session.NewService(sessionRepo)

	// User service setup
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo, sessionService, "hsessionid")
	registerHandler := user.NewRegisterHandler(userService)
	loginHandler := user.NewLoginHandler(userService)
	logoutHandler := user.NewLogoutHandler(userService)
	getUserHandler := user.NewGetHandler(userService)

	userAPI := httpserver.NewHandlerGroup()
	userAPI.Handle("POST /register", registerHandler)
	userAPI.Handle("POST /login", loginHandler)
	userAPI.Handle("POST /logout", logoutHandler)
	userAPI.Handle("GET /me", getUserHandler)

	// Middlewares
	recoverMiddleware := httpserver.NewRecoverMiddleware()
	traceId := traceid.NewTraceId(nil, "Hostamat-Trace-Id")
	requestLogger := middlewares.NewRequestLogger(logger)
	authMiddleware := auth.NewAuthenticationMiddleware(userService)

	// HTTP server
	httpServer := httpserver.New(conf.Port, 10*time.Second)

	apiGroup := httpserver.NewHandlerGroup()
	apiGroup.Use(authMiddleware)

	apiGroup.Handle("/files", filesAPI)
	apiGroup.Handle("/sites", siteHandlerGroup)

	httpServer.HandleGroup("/api", apiGroup)
	httpServer.HandleGroup("/auth", userAPI)
	// HandleGroup used here because client would request files by additional path
	httpServer.HandleGroup(conf.FileServerPath, serveFileHandler)

	httpServer.Use(recoverMiddleware, traceId, requestLogger)

	// Cleanup service to remove files from S3 and subdomains from Caddy
	cleanupService := cleanup.New(s3Client, caddyAPI, conf.S3.Bucket, conf.Caddy.ShorthandId)
	cleanupJob := scheduler.New(30*time.Minute, cleanupService)

	// Application
	app := application.New()
	healthcheckHandler := healthcheck.NewHealthCheckHandler(app)
	httpServer.Handle("/health", healthcheckHandler)

	// Migrate File storage
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

	app.RegisterDatabase("default", db)
	app.RegisterRepository("default", "site", siteRepo)
	app.RegisterRepository("default", "subdomain", subdomainRepo)
	app.RegisterRepository("default", "user", userRepo)
	app.RegisterRepository("default", "session", sessionRepo)

	app.RegisterService("http", httpServer)
	app.RegisterService("cleanup", cleanupJob)

	if err := app.Run(ctx); err != nil {
		logger.Error("Applicaton finished with error", "error", err)
	}
}
