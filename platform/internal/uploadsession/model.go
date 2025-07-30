package uploadsession

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusOpen   Status = "open"
	StatusClosed Status = "closed"
)

type UploadSession struct {
	Id        string    `db:"id" json:"id"`
	SiteId    string    `db:"site_id" json:"siteId"`
	Status    Status    `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}
