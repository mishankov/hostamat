package main

import (
	"hostamat/internal/adapters/config"

	"hostamat/pkg/platforma/log"
	"log/slog"
	"os"
)

func main() {
	command := os.Args[1]

	switch command {
	case "flush_subdomains":
		flushSubdomains()
	}
}

func flushSubdomains() {
	conf := config.New()
	logger := log.New(os.Stdout, conf.Logger.Type, slog.LevelInfo, nil)

	logger.Info("TODO")
}
