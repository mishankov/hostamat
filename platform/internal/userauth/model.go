package userauth

import (
	"time"

	platformauth "github.com/platforma-dev/platforma/auth"
)

type Provider string

const (
	ProviderLegacy Provider = "legacy"
	ProviderYandex Provider = "yandex"
)

type User struct {
	ID             string              `db:"id"               json:"id"`
	Username       string              `db:"username"         json:"username"`
	Password       string              `db:"password"         json:"-"`
	Salt           string              `db:"salt"             json:"-"`
	Created        time.Time           `db:"created"          json:"created"`
	Updated        time.Time           `db:"updated"          json:"updated"`
	Status         platformauth.Status `db:"status"           json:"status"`
	Provider       Provider            `db:"provider"         json:"provider"`
	ProviderUserID *string             `db:"provider_user_id" json:"providerUserId"`
}
