package feedback

import (
	"context"
	"database/sql"
	"io/fs"
	"time"

	"hostamat/internal/feedback/migrations"
)

const tableName = "feedback"

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

func (r *Repository) Create(ctx context.Context, item *Feedback) error {
	const query = `
		INSERT INTO ` + tableName + ` (
			id,
			user_id,
			username,
			email,
			url,
			message,
			created_at
		)
		VALUES (
			:id,
			:user_id,
			:username,
			:email,
			:url,
			:message,
			:created_at
		)
	`

	_, err := r.db.NamedExecContext(ctx, query, item)
	return err
}

func (r *Repository) GetUserContext(ctx context.Context, userID string) (*UserContext, error) {
	var user UserContext
	const query = `
		SELECT id, username, COALESCE(email, '') AS email
		FROM users
		WHERE id = $1
	`
	if err := r.db.GetContext(ctx, &user, query, userID); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) ListRecent(ctx context.Context, limit int) ([]Feedback, error) {
	var items []Feedback
	const query = `
		SELECT id, user_id, username, email, url, message, created_at
		FROM ` + tableName + `
		ORDER BY created_at DESC
		LIMIT $1
	`
	if err := r.db.SelectContext(ctx, &items, query, limit); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *Repository) CountByUserSince(ctx context.Context, userID string, since time.Time) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count, "SELECT COUNT(*) FROM "+tableName+" WHERE user_id = $1 AND created_at >= $2", userID, since)
	return count, err
}

func (r *Repository) DeleteByUser(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM "+tableName+" WHERE user_id = $1", userID)
	return err
}
