package apikey

import (
	"context"
	"database/sql"
	"hostamat/internal/apikey/migrations"
	"io/fs"
	"time"
)

const tableName = "api_keys"

type db interface {
	NamedExecContext(ctx context.Context, query string, arg any) (sql.Result, error)
	GetContext(ctx context.Context, dest any, query string, args ...any) error
	SelectContext(ctx context.Context, dest any, query string, args ...any) error
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type Repository struct {
	db db
}

func NewRepository(db db) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Migrations() fs.FS {
	return migrations.FS
}

func (r *Repository) Create(ctx context.Context, key *APIKey) error {
	const query = `
		INSERT INTO ` + tableName + ` (
			id,
			user_id,
			name,
			key_prefix,
			secret_digest,
			scopes,
			expires_at,
			last_used_at,
			revoked_at,
			created_at,
			updated_at
		)
		VALUES (
			:id,
			:user_id,
			:name,
			:key_prefix,
			:secret_digest,
			:scopes,
			:expires_at,
			:last_used_at,
			:revoked_at,
			:created_at,
			:updated_at
		)
	`

	_, err := r.db.NamedExecContext(ctx, query, key)
	return err
}

func (r *Repository) Get(ctx context.Context, id string) (*APIKey, error) {
	var key APIKey
	err := r.db.GetContext(ctx, &key, "SELECT * FROM "+tableName+" WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	return &key, nil
}

func (r *Repository) GetForUser(ctx context.Context, userID, id string) (*APIKey, error) {
	var key APIKey
	err := r.db.GetContext(ctx, &key, "SELECT * FROM "+tableName+" WHERE user_id = $1 AND id = $2", userID, id)
	if err != nil {
		return nil, err
	}

	return &key, nil
}

func (r *Repository) ListVisibleForUser(ctx context.Context, userID string) ([]*APIKey, error) {
	var keys []*APIKey
	err := r.db.SelectContext(ctx, &keys, "SELECT * FROM "+tableName+" WHERE user_id = $1 AND revoked_at IS NULL ORDER BY created_at DESC", userID)
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (r *Repository) CountActiveForUser(ctx context.Context, userID string, now time.Time) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count, "SELECT COUNT(*) FROM "+tableName+" WHERE user_id = $1 AND revoked_at IS NULL AND (expires_at IS NULL OR expires_at > $2)", userID, now)
	return count, err
}

func (r *Repository) UpdateName(ctx context.Context, key *APIKey) error {
	const query = `
		UPDATE ` + tableName + `
		SET name = :name, updated_at = :updated_at
		WHERE id = :id
	`

	_, err := r.db.NamedExecContext(ctx, query, key)
	return err
}

func (r *Repository) Revoke(ctx context.Context, id string, now time.Time) error {
	_, err := r.db.ExecContext(ctx, "UPDATE "+tableName+" SET revoked_at = $1, updated_at = $1 WHERE id = $2", now, id)
	return err
}

func (r *Repository) UpdateLastUsed(ctx context.Context, id string, now, threshold time.Time) error {
	_, err := r.db.ExecContext(ctx, "UPDATE "+tableName+" SET last_used_at = $1, updated_at = $1 WHERE id = $2 AND (last_used_at IS NULL OR last_used_at < $3)", now, id, threshold)
	return err
}

func (r *Repository) GetOwner(ctx context.Context, userID string) (*Owner, error) {
	var owner Owner
	err := r.db.GetContext(ctx, &owner, `SELECT id, username, created, updated, status FROM users WHERE id = $1`, userID)
	if err != nil {
		return nil, err
	}

	return &owner, nil
}
