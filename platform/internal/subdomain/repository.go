package subdomain

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	tableName = "subdomains"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Migrate() error {
	query := `
		CREATE TABLE IF NOT EXISTS ` + tableName + ` (
			id VARCHAR(255) PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE,
			full_domain VARCHAR(255) NOT NULL
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
		return nil, fmt.Errorf("failed to get subdomain: %w", err)
	}

	return &subdomain, nil
}

func (r *Repository) Create(ctx context.Context, subdomain *Subdomain) error {
	const query = `INSERT INTO ` + tableName + ` (id, name, full_domain) VALUES (:id, :name, :full_domain)`

	_, err := r.db.NamedExecContext(ctx, query, subdomain)
	if err != nil {
		return fmt.Errorf("failed to create subdomain: %w", err)
	}

	return nil
}
