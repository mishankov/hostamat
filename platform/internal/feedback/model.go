package feedback

import "time"

const MaxMessageLength = 5000
const MaxURLLength = 2048
const MaxRequestBodyBytes = 32 * 1024
const MaxSubmissionsPerHour = 5

type Feedback struct {
	ID        string    `db:"id"         json:"id"`
	UserID    string    `db:"user_id"    json:"userId"`
	Username  string    `db:"username"   json:"username"`
	Email     string    `db:"email"      json:"email"`
	URL       string    `db:"url"        json:"url"`
	Message   string    `db:"message"    json:"message"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

type UserContext struct {
	ID       string `db:"id"`
	Username string `db:"username"`
	Email    string `db:"email"`
}
