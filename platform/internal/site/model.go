package site

import (
	"time"
)

const (
	maxSites          = 5
	siteNameMaxLength = 25
	maxFilesSize      = 20 * 1024 * 1024
)

type Site struct {
	Id              string    `db:"id" json:"id"`
	Name            string    `db:"name" json:"name"`
	User            string    `db:"user" json:"user"`
	Subdomain       string    `db:"subdomain" json:"subdomain"`
	Domain          string    `db:"domain" json:"domain"`
	ActiveReleaseId *string   `db:"active_release_id" json:"activeReleaseId"`
	CreatedAt       time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt       time.Time `db:"updated_at" json:"updatedAt"`
}

type UploadedFile struct {
	Name     string `json:"name"`
	Uploaded bool   `json:"uploaded"`
	Size     int    `json:"size"`
}
