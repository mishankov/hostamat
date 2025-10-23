package uploadsession

import "time"

type Uploadsession struct {
	Id        string    `db:"id" json:"id"`
	User      string    `db:"user" json:"user"`
	Site      string    `db:"site" json:"site"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
	Status    string    `db:"status" json:"status"`
}
