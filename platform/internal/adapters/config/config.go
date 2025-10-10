package config

import (
	"log/slog"
	"os"
)

type caddy struct {
	BaseURL     string
	ShorthandId string // https://caddyserver.com/docs/api#using-id-in-json
}

type s3 struct {
	BaseURL  string
	Login    string
	Password string
	Bucket   string
}

type db struct {
	Host     string
	Port     string
	DB       string
	User     string
	Password string
}

type logger struct {
	Level slog.Level

	// text or json
	Type string
}

type Config struct {
	Port             string // Port to run http server
	PublicBaseDomain string // Base domain for uploaded sites
	FileServerHost   string // Internal host for file server
	FileServerPath   string
	Logger           logger
	Caddy            caddy
	S3               s3
	DB               db
}

func New() *Config {
	return &Config{
		Port:             os.Getenv("HOSTAMAT_PORT"),
		PublicBaseDomain: os.Getenv("HOSTAMAT_DOMAIN"),
		FileServerHost:   os.Getenv("HOSTAMAT_FILESERVER_HOST"),
		FileServerPath:   os.Getenv("HOSTAMAT_FILESERVER_PATH"),
		Caddy: caddy{
			BaseURL:     os.Getenv("HOSTAMAT_CADDY_URL"),
			ShorthandId: os.Getenv("HOSTAMAT_CADDY_SHORTHAND_ID"),
		},
		S3: s3{
			BaseURL:  os.Getenv("HOSTAMAT_S3_URL"),
			Login:    os.Getenv("HOSTAMAT_S3_LOGIN"),
			Password: os.Getenv("HOSTAMAT_S3_PASSWORD"),
			Bucket:   os.Getenv("HOSTAMAT_S3_BUCKET"),
		},
		DB: db{
			Host:     os.Getenv("HOSTAMAT_POSTGRES_HOST"),
			Port:     os.Getenv("HOSTAMAT_POSTGRES_Port"),
			DB:       os.Getenv("HOSTAMAT_POSTGRES_DB"),
			User:     os.Getenv("HOSTAMAT_POSTGRES_USER"),
			Password: os.Getenv("HOSTAMAT_POSTGRES_PASSWORD"),
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
