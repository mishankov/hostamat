package site

import (
	"time"

	"hostamat/internal/formcollection"
	"hostamat/internal/uploadsession"
)

const (
	maxSites          = 5
	siteNameMaxLength = 25
	maxFilesSize      = 20 * 1024 * 1024
)

const (
	SiteWarningSeverityWarning = "warning"
	SiteWarningMissingHomePage = "missing_home_page"
)

type Site struct {
	Id        string        `db:"id" json:"id"`
	Name      string        `db:"name" json:"name"`
	User      string        `db:"user" json:"user"`
	Subdomain string        `db:"subdomain" json:"subdomain"`
	Domain    string        `db:"domain" json:"domain"`
	CreatedAt time.Time     `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time     `db:"updated_at" json:"updatedAt"`
	Warnings  []SiteWarning `db:"-" json:"warnings"`
}

type SiteWarning struct {
	Site      string    `db:"site" json:"-"`
	Code      string    `db:"code" json:"code"`
	Severity  string    `db:"severity" json:"severity"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

type UploadedPage struct {
	FileName string `json:"fileName"`
	Path     string `json:"path"`
	URL      string `json:"url"`
}

type UploadSummary struct {
	UploadedFiles  []uploadsession.File             `json:"uploadedFiles"`
	SkippedFiles   []uploadsession.File             `json:"skippedFiles"`
	AvailablePages []UploadedPage                   `json:"availablePages"`
	HasHomePage    bool                             `json:"hasHomePage"`
	FormDiscovery  *formcollection.DiscoverySummary `json:"formDiscovery,omitempty"`
}
