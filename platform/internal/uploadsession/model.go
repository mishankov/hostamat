package uploadsession

import "time"

type UploadSessionStatus string

const (
	UploadSessionStatusOpen   UploadSessionStatus = "open"
	UploadSessionStatusClosed UploadSessionStatus = "closed"
)

type UploadSession struct {
	Id        string              `db:"id" json:"id"`
	User      string              `db:"user" json:"user"`
	Site      string              `db:"site" json:"site"`
	CreatedAt time.Time           `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time           `db:"updated_at" json:"updatedAt"`
	Status    UploadSessionStatus `db:"status" json:"status"`
}

type File struct {
	Id              string `db:"id" json:"id"`
	UploadSessionId string `db:"upload_session" json:"-"`
	Name            string `db:"name" json:"name"`
	Uploaded        bool   `db:"uploaded" json:"uploaded"`
	Size            int    `db:"size" json:"size"`
}
