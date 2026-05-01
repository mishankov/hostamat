package site

import (
	"time"

	"hostamat/internal/uploadsession"
)

const (
	maxSites          = 5
	siteNameMaxLength = 25
	maxFilesSize      = 20 * 1024 * 1024
)

type Site struct {
	Id        string    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	User      string    `db:"user" json:"user"`
	Subdomain string    `db:"subdomain" json:"subdomain"`
	Domain    string    `db:"domain" json:"domain"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

type UploadedPage struct {
	FileName string `json:"fileName"`
	Path     string `json:"path"`
	URL      string `json:"url"`
}

type UploadSummary struct {
	UploadedFiles  []uploadsession.File `json:"uploadedFiles"`
	SkippedFiles   []uploadsession.File `json:"skippedFiles"`
	AvailablePages []UploadedPage       `json:"availablePages"`
	HasHomePage    bool                 `json:"hasHomePage"`
}
