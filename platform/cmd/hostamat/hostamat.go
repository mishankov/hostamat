package main

import (
	"hostamat/internal/adapters/caddy"
	"hostamat/internal/adapters/config"
	"hostamat/internal/adapters/timeweb"
	"hostamat/internal/handlers"
	"hostamat/internal/infra/httpclient"
	"hostamat/internal/infra/logger"
	"hostamat/internal/server"
)

func main() {
	conf := config.New()

	logger := logger.New(conf.Logger.Type, conf.Logger.Level)

	httpClient := httpclient.New(logger)

	timewebAPI := timeweb.New(httpClient, conf.Timeweb.BaseURL, conf.Timeweb.APIToken)
	caddyAPI := caddy.New(httpClient, conf.Caddy.BaseURL, conf.Caddy.ServerHost, conf.BaseDomain)

	uploadFileHandler := handlers.NewUploadFileHandler(logger, timewebAPI, caddyAPI, conf.Timeweb.BaseDomain)

	server := server.New(conf.Port, logger, uploadFileHandler)

	server.Run()
}
