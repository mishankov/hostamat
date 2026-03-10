package userauth

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/platforma-dev/platforma/database"
)

type db interface {
	NamedExecContext(ctx context.Context, query string, arg any) (sql.Result, error)
	GetContext(ctx context.Context, dest any, query string, args ...any) error
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type Repository struct {
	db db
}

func NewRepository(db db) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Migrations() []database.Migration {
	return []database.Migration{
		{
			ID: "init",
			Up: `CREATE TABLE IF NOT EXISTS users (
				id VARCHAR(255) PRIMARY KEY,
				username VARCHAR(255),
				password TEXT,
				salt TEXT,
				created TIMESTAMP,
				updated TIMESTAMP,
				status VARCHAR(50),
				provider VARCHAR(50) NOT NULL DEFAULT 'legacy',
				provider_user_id VARCHAR(255)
			)`,
			Down: "DROP TABLE users",
		},
		{
			ID: "add_provider_columns",
			Up: `ALTER TABLE users ADD COLUMN IF NOT EXISTS provider VARCHAR(50) NOT NULL DEFAULT 'legacy';
UPDATE users SET provider = 'legacy' WHERE provider IS NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS provider_user_id VARCHAR(255);`,
			Down: `ALTER TABLE users DROP COLUMN IF EXISTS provider_user_id;
ALTER TABLE users DROP COLUMN IF EXISTS provider;`,
		},
		{
			ID:   "drop_username_unique",
			Up:   `ALTER TABLE users DROP CONSTRAINT IF EXISTS users_username_key;`,
			Down: `ALTER TABLE users ADD CONSTRAINT users_username_key UNIQUE (username);`,
		},
		{
			ID:   "add_provider_user_id_unique_index",
			Up:   `CREATE UNIQUE INDEX IF NOT EXISTS users_provider_provider_user_id_idx ON users (provider, provider_user_id) WHERE provider_user_id IS NOT NULL;`,
			Down: `DROP INDEX IF EXISTS users_provider_provider_user_id_idx;`,
		},
	}
}

func (r *Repository) Get(ctx context.Context, id string) (*User, error) {
	var user User
	if err := r.db.GetContext(ctx, &user, `SELECT * FROM users WHERE id = $1`, id); err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}

func (r *Repository) GetByProviderUserID(ctx context.Context, provider Provider, providerUserID string) (*User, error) {
	var user User
	if err := r.db.GetContext(ctx, &user, `SELECT * FROM users WHERE provider = $1 AND provider_user_id = $2`, provider, providerUserID); err != nil {
		return nil, fmt.Errorf("failed to get user by provider: %w", err)
	}

	return &user, nil
}

func (r *Repository) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (id, username, password, salt, created, updated, status, provider, provider_user_id)
		VALUES (:id, :username, :password, :salt, :created, :updated, :status, :provider, :provider_user_id)
	`
	if _, err := r.db.NamedExecContext(ctx, query, user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *Repository) UpdateUsername(ctx context.Context, id, username string) error {
	query := `
		UPDATE users
		SET username = $1, updated = CURRENT_TIMESTAMP
		WHERE id = $2
	`
	if _, err := r.db.ExecContext(ctx, query, username, id); err != nil {
		return fmt.Errorf("failed to update username: %w", err)
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	if _, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
