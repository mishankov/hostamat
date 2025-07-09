package site

import (
	"time"
)

type Site struct {
	Id        string    `db:"id" json:"id"`
	Domain    string    `db:"domain" json:"domain"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}
