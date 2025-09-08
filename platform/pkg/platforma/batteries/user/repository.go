package user

import (
	"context"
	"database/sql"
	"hostamat/pkg/platforma/database"
)

type db interface {
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type Repository struct {
	db db
}

func NewRepository(db db) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Schema() ([]database.Migration, database.Schema) {
	return []database.Migration{}, database.Schema{
		Queries: []string{
			`
			CREATE TABLE IF NOT EXISTS users (
				id TEXT PRIMARY KEY,
				username TEXT UNIQUE,
				password TEXT,
				salt TEXT,
				created TIMESTAMP,
				updated TIMESTAMP,
				status TEXT
			)
			`,
		},
	}
}

func (r *Repository) Get(ctx context.Context, id string) (*User, error) {
	var user User
	err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (id, username, password, salt, created, updated, status)
		VALUES (:id, :username, :password, :salt, :created, :updated, :status)
	`
	_, err := r.db.NamedExecContext(ctx, query, user)
	return err
}
