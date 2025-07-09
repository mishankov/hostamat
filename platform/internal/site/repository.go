package site

import (
	"context"

	"github.com/jmoiron/sqlx"
)

const (
	tableName = "sites"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Migrate() error {
	query := `
		CREATE TABLE IF NOT EXISTS ` + tableName + ` (
			id VARCHAR(255) PRIMARY KEY,
			subdomain VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL
		)
	`
	_, err := r.db.Exec(query)
	return err
}

func (r *Repository) Create(ctx context.Context, site *Site) error {
	query := `
		INSERT INTO ` + tableName + ` (id, subdomain, created_at)
		VALUES (:id, :subdomain, :created_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, site)
	return err
}

func (r *Repository) Get(ctx context.Context, id string) (*Site, error) {
	var site Site
	err := r.db.GetContext(ctx, &site, "SELECT * FROM "+tableName+" WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &site, nil
}
