package subdomain

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	tableName = "subdomains"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Migrate() error {
	query := `
		CREATE TABLE IF NOT EXISTS ` + tableName + ` (
			id VARCHAR(255) PRIMARY KEY,
			name VARCHAR(255) NOT NULL
		)
	`
	_, err := r.db.Exec(query)
	return err
}

func (r *Repository) Get(ctx context.Context, id string) (*Subdomain, error) {
	const query = `SELECT id, name FROM ` + tableName + ` WHERE id = $1`

	var subdomain Subdomain
	err := r.db.GetContext(ctx, &subdomain, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get subdomain: %w", err)
	}

	return &subdomain, nil
}

func (r *Repository) Create(ctx context.Context, subdomain *Subdomain) error {
	const query = `INSERT INTO ` + tableName + ` (id, name) VALUES (:id, :name)`

	_, err := r.db.NamedExecContext(ctx, query, subdomain)
	if err != nil {
		return fmt.Errorf("failed to create subdomain: %w", err)
	}

	return nil
}
