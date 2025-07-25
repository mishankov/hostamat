package site

import (
	"time"
)

type Site struct {
	Id        string    `db:"id" json:"id"`
	Subdomain string    `db:"subdomain" json:"subdomain"`
	Domain    string    `db:"domain" json:"domain"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
	Version   int       `db:"version" json:"version"`
}
