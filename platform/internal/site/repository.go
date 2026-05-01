package site

import (
	"context"
	"database/sql"
	"hostamat/internal/site/migrations"
	"io/fs"

	"github.com/lib/pq"
)

const (
	tableName = "sites"
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

func (r *Repository) Migrations() fs.FS {
	return migrations.FS
}

func (r *Repository) Create(ctx context.Context, site *Site) error {
	query := `
		INSERT INTO ` + tableName + ` (id, name, "user", subdomain, domain, created_at, updated_at)
		VALUES (:id, :name, :user, :subdomain, :domain, :created_at, :updated_at)
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

func (r *Repository) GetList(ctx context.Context, userId string) ([]*Site, error) {
	var sites []*Site
	err := r.db.SelectContext(ctx, &sites, "SELECT * FROM "+tableName+" WHERE \"user\" = $1 ORDER BY created_at", userId)
	if err != nil {
		return nil, err
	}
	return sites, nil
}

func (r *Repository) GetWarningsForSites(ctx context.Context, siteIDs []string) ([]SiteWarning, error) {
	if len(siteIDs) == 0 {
		return []SiteWarning{}, nil
	}

	var warnings []SiteWarning
	err := r.db.SelectContext(ctx, &warnings, `
		SELECT *
		FROM site_warnings
		WHERE site = ANY($1)
		ORDER BY created_at
	`, pq.Array(siteIDs))
	if err != nil {
		return nil, err
	}

	return warnings, nil
}

func (r *Repository) UpsertWarning(ctx context.Context, warning *SiteWarning) error {
	const query = `
		INSERT INTO site_warnings (
			site,
			code,
			severity,
			created_at,
			updated_at
		)
		VALUES (
			:site,
			:code,
			:severity,
			:created_at,
			:updated_at
		)
		ON CONFLICT (site, code) DO UPDATE SET
			severity = EXCLUDED.severity,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.NamedExecContext(ctx, query, warning)
	return err
}

func (r *Repository) DeleteWarning(ctx context.Context, siteID, code string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM site_warnings WHERE site = $1 AND code = $2", siteID, code)
	return err
}

func (r *Repository) Update(ctx context.Context, site *Site) error {
	query := `
		UPDATE ` + tableName + `
		SET name = :name, updated_at = :updated_at
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, site)
	return err
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	if _, err := r.db.ExecContext(ctx, "DELETE FROM site_warnings WHERE site = $1", id); err != nil {
		return err
	}

	_, err := r.db.ExecContext(ctx, "DELETE FROM "+tableName+" WHERE id = $1", id)
	return err
}
