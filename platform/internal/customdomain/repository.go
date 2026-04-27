package customdomain

import (
	"context"
	"database/sql"
	"io/fs"

	"hostamat/internal/customdomain/migrations"
	"hostamat/internal/site"
)

const tableName = "custom_domains"

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
	return &Repository{db: db}
}

func (r *Repository) Migrations() fs.FS {
	return migrations.FS
}

func (r *Repository) Create(ctx context.Context, domain *CustomDomain) error {
	const query = `
		INSERT INTO ` + tableName + ` (
			id,
			site,
			"user",
			hostname,
			domain_type,
			status,
			verification_token,
			last_verification_error,
			last_checked_at,
			activated_at,
			created_at,
			updated_at
		)
		VALUES (
			:id,
			:site,
			:user,
			:hostname,
			:domain_type,
			:status,
			:verification_token,
			:last_verification_error,
			:last_checked_at,
			:activated_at,
			:created_at,
			:updated_at
		)
	`

	_, err := r.db.NamedExecContext(ctx, query, domain)
	return err
}

func (r *Repository) Get(ctx context.Context, id string) (*CustomDomain, error) {
	var domain CustomDomain
	err := r.db.GetContext(ctx, &domain, "SELECT * FROM "+tableName+" WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	return &domain, nil
}

func (r *Repository) GetActiveByHostname(ctx context.Context, hostname string) (*CustomDomain, error) {
	var domain CustomDomain
	err := r.db.GetContext(ctx, &domain, "SELECT * FROM "+tableName+" WHERE hostname = $1 AND status = $2", hostname, DomainStatusActive)
	if err != nil {
		return nil, err
	}

	return &domain, nil
}

func (r *Repository) GetBySite(ctx context.Context, siteID string) ([]*CustomDomain, error) {
	var domains []*CustomDomain
	err := r.db.SelectContext(ctx, &domains, "SELECT * FROM "+tableName+" WHERE site = $1 ORDER BY created_at", siteID)
	if err != nil {
		return nil, err
	}

	return domains, nil
}

func (r *Repository) GetListByUser(ctx context.Context, userID string) ([]*DomainWithSite, error) {
	const query = `
		SELECT
			d.id,
			d.site,
			d."user",
			d.hostname,
			d.domain_type,
			d.status,
			d.verification_token,
			d.last_verification_error,
			d.last_checked_at,
			d.activated_at,
			d.created_at,
			d.updated_at,
			s.name AS site_name,
			s.domain AS site_domain
		FROM ` + tableName + ` d
		JOIN sites s ON s.id = d.site
		WHERE d."user" = $1
		ORDER BY s.created_at, d.created_at
	`

	var domains []*DomainWithSite
	err := r.db.SelectContext(ctx, &domains, query, userID)
	if err != nil {
		return nil, err
	}

	return domains, nil
}

func (r *Repository) Update(ctx context.Context, domain *CustomDomain) error {
	const query = `
		UPDATE ` + tableName + `
		SET
			site = :site,
			status = :status,
			last_verification_error = :last_verification_error,
			last_checked_at = :last_checked_at,
			activated_at = :activated_at,
			updated_at = :updated_at
		WHERE id = :id
	`

	_, err := r.db.NamedExecContext(ctx, query, domain)
	return err
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM "+tableName+" WHERE id = $1", id)
	return err
}

func (r *Repository) ResolveSiteByHostname(ctx context.Context, hostname string) (*site.Site, error) {
	const query = `
		SELECT s.*
		FROM ` + tableName + ` d
		JOIN sites s ON s.id = d.site
		WHERE d.hostname = $1 AND d.status = $2
	`

	var resolved site.Site
	err := r.db.GetContext(ctx, &resolved, query, hostname, DomainStatusActive)
	if err != nil {
		return nil, err
	}

	return &resolved, nil
}

func (r *Repository) ListActiveHosts(ctx context.Context) ([]string, error) {
	var hosts []string
	err := r.db.SelectContext(ctx, &hosts, "SELECT hostname FROM "+tableName+" WHERE status = $1 ORDER BY hostname", DomainStatusActive)
	if err != nil {
		return nil, err
	}

	return hosts, nil
}
