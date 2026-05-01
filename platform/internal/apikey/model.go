package apikey

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
	platformauth "github.com/platforma-dev/platforma/auth"
)

const (
	ScopeSitesRead    Scope = "sites:read"
	ScopeSitesWrite   Scope = "sites:write"
	ScopeSitesUpload  Scope = "sites:upload"
	ScopeDomainsRead  Scope = "domains:read"
	ScopeDomainsWrite Scope = "domains:write"
)

const (
	maxNameLength       = 80
	maxActiveKeys       = 20
	lastUsedUpdateAfter = 5 * time.Minute
)

type Scope string

type APIKey struct {
	ID           string         `db:"id" json:"id"`
	UserID       string         `db:"user_id" json:"userId"`
	Name         string         `db:"name" json:"name"`
	KeyPrefix    string         `db:"key_prefix" json:"keyPrefix"`
	SecretDigest string         `db:"secret_digest" json:"-"`
	Scopes       pq.StringArray `db:"scopes" json:"scopes"`
	ExpiresAt    sql.NullTime   `db:"expires_at" json:"-"`
	LastUsedAt   sql.NullTime   `db:"last_used_at" json:"-"`
	RevokedAt    sql.NullTime   `db:"revoked_at" json:"-"`
	CreatedAt    time.Time      `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time      `db:"updated_at" json:"updatedAt"`
}

type Owner struct {
	ID       string              `db:"id"`
	Username string              `db:"username"`
	Created  time.Time           `db:"created"`
	Updated  time.Time           `db:"updated"`
	Status   platformauth.Status `db:"status"`
}

type CreateRequest struct {
	Name      string     `json:"name"`
	Scopes    []Scope    `json:"scopes"`
	ExpiresAt *time.Time `json:"expiresAt"`
}

type UpdateRequest struct {
	Name string `json:"name"`
}

type PublicKeyMetadata struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	KeyPrefix  string     `json:"keyPrefix"`
	Scopes     []Scope    `json:"scopes"`
	ExpiresAt  *time.Time `json:"expiresAt"`
	LastUsedAt *time.Time `json:"lastUsedAt"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	Status     string     `json:"status"`
}

type CreateResponse struct {
	Key    PublicKeyMetadata `json:"key"`
	Secret string            `json:"secret"`
}

type AuthenticatedKey struct {
	Key   *APIKey
	Owner *Owner
}

func (k *APIKey) Metadata(now time.Time) PublicKeyMetadata {
	return PublicKeyMetadata{
		ID:         k.ID,
		Name:       k.Name,
		KeyPrefix:  k.KeyPrefix,
		Scopes:     scopesFromStrings(k.Scopes),
		ExpiresAt:  nullTimePtr(k.ExpiresAt),
		LastUsedAt: nullTimePtr(k.LastUsedAt),
		CreatedAt:  k.CreatedAt,
		UpdatedAt:  k.UpdatedAt,
		Status:     k.Status(now),
	}
}

func (k *APIKey) Status(now time.Time) string {
	if k.RevokedAt.Valid {
		return "revoked"
	}

	if k.ExpiresAt.Valid && !k.ExpiresAt.Time.After(now) {
		return "expired"
	}

	return "active"
}

func (k *APIKey) HasScope(scope Scope) bool {
	for _, candidate := range k.Scopes {
		if candidate == string(scope) {
			return true
		}
	}

	return false
}

func nullTimePtr(value sql.NullTime) *time.Time {
	if !value.Valid {
		return nil
	}

	return &value.Time
}

func scopesFromStrings(values []string) []Scope {
	scopes := make([]Scope, 0, len(values))
	for _, value := range values {
		scopes = append(scopes, Scope(value))
	}

	return scopes
}

func stringsFromScopes(values []Scope) []string {
	scopes := make([]string, 0, len(values))
	for _, value := range values {
		scopes = append(scopes, string(value))
	}

	return scopes
}
