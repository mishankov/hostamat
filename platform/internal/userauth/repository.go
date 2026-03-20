package userauth

import (
	"context"
	"database/sql"
	"fmt"
	"hostamat/internal/userauth/migrations"
	"io/fs"
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

func (r *Repository) Migrations() fs.FS {
	return migrations.FS
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
