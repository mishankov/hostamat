package subdomain

import (
	"context"
	"database/sql"
	"fmt"
	"hostamat/internal/subdomain/internal/migrations"

	"github.com/mishankov/platforma/database"
)

const (
	tableName = "subdomains"
)

type db interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}

type Repository struct {
	db db
}

func NewRepository(db db) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Schema() ([]database.Migration, database.Schema) {
	return migrations.Migrations, migrations.Schema
}

func (r *Repository) Get(ctx context.Context, id string) (*Subdomain, error) {
	const query = `SELECT * FROM ` + tableName + ` WHERE id = $1`

	var subdomain Subdomain
	err := r.db.GetContext(ctx, &subdomain, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get subdomain: %w", err)
	}

	return &subdomain, nil
}

func (r *Repository) Create(ctx context.Context, subdomain *Subdomain) error {
	const query = `INSERT INTO ` + tableName + ` (id, "user", name, full_domain) VALUES (:id, :user, :name, :full_domain)`

	_, err := r.db.NamedExecContext(ctx, query, subdomain)
	if err != nil {
		return fmt.Errorf("failed to create subdomain: %w", err)
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	const query = `DELETE FROM ` + tableName + ` WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete subdomain: %w", err)
	}

	return nil
}
