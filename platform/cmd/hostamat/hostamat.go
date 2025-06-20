package main

import (
	"context"
	"hostamat/internal/adapters/caddy"
	"hostamat/internal/adapters/config"
	"hostamat/internal/adapters/timeweb"
	"hostamat/internal/handlers"
	"hostamat/internal/infra/logger"
	"hostamat/internal/middlewares"
	"hostamat/pkg/platforma/application"
	"hostamat/pkg/platforma/httpclient"
	"hostamat/pkg/platforma/httpserver"
	"hostamat/pkg/platforma/httpserver/handlegroup"
	"net/http"
)

func main() {
	ctx := context.Background()

	conf := config.New()

	logger := logger.New(conf.Logger.Type, conf.Logger.Level)

	httpClient := httpclient.New(logger)

	timewebAPI := timeweb.New(httpClient, conf.Timeweb.BaseURL, conf.Timeweb.APIToken)
	caddyAPI := caddy.New(httpClient, conf.Caddy.BaseURL, conf.Caddy.ServerHost, conf.BaseDomain)

	uploadFileHandler := handlers.NewUploadFileHandler(logger, timewebAPI, caddyAPI, conf.Timeweb.BaseDomain)

	// server := server.New(conf.Port, logger, uploadFileHandler)

	// server.Run()

	requestLogger := middlewares.NewRequestLogger(logger)

	fileHandlers := handlegroup.New()

	fileHandlers.Use(requestLogger)

	fileHandlers.Handle("/upload", uploadFileHandler)

	api := httpserver.New(conf.Port, map[string]http.Handler{
		"/files": fileHandlers,
		"/healthcheck": http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		}),
	})

	app := application.New()

	app.Register(api)

	logger.Error("Application finished", "error", app.Run(ctx))
}
