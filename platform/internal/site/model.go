package site

import (
	"time"
)

type Site struct {
	Id        string    `db:"id" json:"id"`
	Subdomain string    `db:"subdomain" json:"subdomain"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}
