package site

import (
	"time"
)

const maxSites = 5

type Site struct {
	Id        string    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	User      string    `db:"user" json:"user"`
	Subdomain string    `db:"subdomain" json:"subdomain"`
	Domain    string    `db:"domain" json:"domain"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
	Version   int       `db:"version" json:"version"`
}
