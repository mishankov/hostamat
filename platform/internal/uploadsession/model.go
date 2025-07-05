package uploadsession

import "time"

type UploadSessionStatus string

const (
	Open      UploadSessionStatus = "open"
	Finished  UploadSessionStatus = "finished"
	Cancelled UploadSessionStatus = "cancelled"
)

type UploadSession struct {
	Id         string              `db:"id" json:"id"`
	CreatedAt  time.Time           `db:"created_at" json:"created_at"`
	Status     UploadSessionStatus `db:"status" json:"status"`
	ExpiresAt  time.Time           `db:"expires_at" json:"expires_at"`
	ResolvedAt time.Time           `db:"resolved_at" json:"resolved_at,omitempty"`
}
