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
			domain VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			version INT NOT NULL
		)
	`
	_, err := r.db.Exec(query)
	return err
}

func (r *Repository) Create(ctx context.Context, site *Site) error {
	query := `
		INSERT INTO ` + tableName + ` (id, subdomain, domain, created_at, updated_at, version)
		VALUES (:id, :subdomain, :domain, :created_at, :updated_at, :version)
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

func (r *Repository) GetByDomain(ctx context.Context, domain string) (*Site, error) {
	var site Site
	err := r.db.GetContext(ctx, &site, "SELECT * FROM "+tableName+" WHERE domain = $1", domain)
	if err != nil {
		return nil, err
	}
	return &site, nil
}

func (r *Repository) GetList(ctx context.Context) ([]*Site, error) {
	var sites []*Site
	err := r.db.SelectContext(ctx, &sites, "SELECT * FROM "+tableName)
	if err != nil {
		return nil, err
	}
	return sites, nil
}

func (r *Repository) Update(ctx context.Context, site *Site) error {
	query := `
		UPDATE ` + tableName + `
		SET subdomain = :subdomain, domain = :domain, updated_at = :updated_at, version = :version
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, site)
	return err
}
