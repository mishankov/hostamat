package config

import (
	"log/slog"
	"os"
)

type timeweb struct {
	APIToken   string
	BaseURL    string
	BaseDomain string
}

type caddy struct {
	BaseURL    string
	ServerHost string
}

type minio struct {
	BaseURL  string
	Login    string
	Password string
	Bucket   string
}

type logger struct {
	Level slog.Level

	// text or json
	Type string
}

type Config struct {
	Port       string
	BaseDomain string
	Logger     logger
	Caddy      caddy
	Timeweb    timeweb
	Minio      minio
}

func New() *Config {
	return &Config{
		Port:       os.Getenv("HOSTAMAT_PORT"),
		BaseDomain: os.Getenv("HOSTAMAT_DOMAIN"),
		Caddy: caddy{
			BaseURL:    os.Getenv("HOSTAMAT_CADDY_URL"),
			ServerHost: os.Getenv("HOSTAMAT_CADDY_HOST"),
		},
		Timeweb: timeweb{
			APIToken:   os.Getenv("HOSTAMAT_TIMEWEB_TOKEN"),
			BaseURL:    os.Getenv("HOSTAMAT_TIMEWEB_URL"),
			BaseDomain: os.Getenv("HOSTAMAT_TIMEWEB_DOMAIN"),
		},
		Minio: minio{
			BaseURL:  os.Getenv("HOSTAMAT_MINIO_URL"),
			Login:    os.Getenv("HOSTAMAT_MINIO_LOGIN"),
			Password: os.Getenv("HOSTAMAT_MINIO_PASSWORD"),
			Bucket:   os.Getenv("HOSTAMAT_MINIO_BUCKET"),
		},
		Logger: logger{
			Level: func() slog.Level {
				switch os.Getenv("HOSTAMAT_LOG_LEVEL") {
				case "debug":
					return -4
				case "info":
					return 0
				}
				return 0
			}(),
			Type: os.Getenv("HOSTAMAT_LOG_TYPE"),
		},
	}
}
