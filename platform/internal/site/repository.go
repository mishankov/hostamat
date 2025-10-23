package site

import (
	"context"
	"database/sql"
	"hostamat/internal/site/internal/migrations"

	"github.com/mishankov/platforma/database"
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

func (r *Repository) Schema() ([]database.Migration, database.Schema) {
	return migrations.Migrations, migrations.Schema
}

func (r *Repository) Create(ctx context.Context, site *Site) error {
	query := `
		INSERT INTO ` + tableName + ` (id, "user", subdomain, domain, created_at, updated_at, version)
		VALUES (:id, :user, :subdomain, :domain, :created_at, :updated_at, :version)
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
	err := r.db.SelectContext(ctx, &sites, "SELECT * FROM "+tableName+" WHERE \"user\" = $1", userId)
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

func (r *Repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM "+tableName+" WHERE id = $1", id)
	return err
}
