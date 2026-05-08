package formcollection

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"

	"hostamat/internal/formcollection/migrations"

	"github.com/lib/pq"
)

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

func (r *Repository) UserOwnsSite(ctx context.Context, userID, siteID string) (bool, error) {
	var count int
	err := r.db.GetContext(ctx, &count, `SELECT COUNT(*) FROM sites WHERE id = $1 AND "user" = $2`, siteID, userID)
	return count > 0, err
}

func (r *Repository) CreateDiscovery(ctx context.Context, discovery *Discovery, forms []DiscoveredForm) error {
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO form_discoveries (id, site, created_at)
		VALUES (:id, :site, :created_at)
	`, discovery)
	if err != nil {
		return err
	}

	for _, form := range forms {
		if _, err := r.db.NamedExecContext(ctx, `
			INSERT INTO form_definitions (
				id,
				discovery_id,
				site,
				file_name,
				page_path,
				form_key,
				label,
				original_action,
				method,
				fields,
				supported,
				unsupported_reason,
				active,
				public_token,
				created_at,
				updated_at
			)
			VALUES (
				:id,
				:discovery_id,
				:site,
				:file_name,
				:page_path,
				:form_key,
				:label,
				:original_action,
				:method,
				:fields,
				:supported,
				:unsupported_reason,
				:active,
				:public_token,
				:created_at,
				:updated_at
			)
		`, &form); err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) DisableActiveBySite(ctx context.Context, siteID string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE form_definitions
		SET active = FALSE, updated_at = NOW()
		WHERE site = $1 AND active = TRUE
	`, siteID)
	return err
}

func (r *Repository) GetLatestDiscoveryID(ctx context.Context, siteID string) (string, error) {
	var id string
	err := r.db.GetContext(ctx, &id, `
		SELECT id
		FROM form_discoveries
		WHERE site = $1
		ORDER BY created_at DESC
		LIMIT 1
	`, siteID)
	return id, err
}

func (r *Repository) ListDefinitions(ctx context.Context, userID, siteID string) ([]DiscoveredForm, error) {
	args := []any{userID}
	query := `
		SELECT fd.*
		FROM form_definitions fd
		INNER JOIN sites s ON s.id = fd.site
		WHERE s."user" = $1
	`
	if siteID != "" {
		args = append(args, siteID)
		query += ` AND fd.site = $2`
	}
	query += ` ORDER BY fd.created_at DESC, fd.page_path ASC`

	var forms []DiscoveredForm
	err := r.db.SelectContext(ctx, &forms, query, args...)
	if err != nil {
		return nil, err
	}
	return forms, nil
}

func (r *Repository) GetDefinitionsForEnable(ctx context.Context, userID, siteID, discoveryID string, ids []string) ([]DiscoveredForm, error) {
	if len(ids) == 0 {
		return []DiscoveredForm{}, nil
	}

	query := `
		SELECT fd.*
		FROM form_definitions fd
		INNER JOIN sites s ON s.id = fd.site
		WHERE s."user" = $1
			AND fd.site = $2
			AND fd.discovery_id = $3
			AND fd.id = ANY($4)
	`
	var forms []DiscoveredForm
	err := r.db.SelectContext(ctx, &forms, query, userID, siteID, discoveryID, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	return forms, nil
}

func (r *Repository) ActivateDefinition(ctx context.Context, id, token string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE form_definitions
		SET active = TRUE, public_token = $2, updated_at = NOW()
		WHERE id = $1
	`, id, token)
	return err
}

func (r *Repository) GetActiveByToken(ctx context.Context, token string) (*DiscoveredForm, error) {
	var form DiscoveredForm
	err := r.db.GetContext(ctx, &form, `
		SELECT *
		FROM form_definitions
		WHERE public_token = $1
	`, token)
	if err != nil {
		return nil, err
	}
	return &form, nil
}

func (r *Repository) CreateSubmission(ctx context.Context, submission *Submission) error {
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO form_submissions (id, site, form, page_path, fields, created_at)
		VALUES (:id, :site, :form, :page_path, :fields, :created_at)
	`, submission)
	return err
}

func (r *Repository) ListSubmissions(ctx context.Context, userID string, filters SubmissionFilters) ([]Submission, error) {
	args := []any{userID}
	query := `
		SELECT fs.id,
		       fs.site,
		       fs.form,
		       fs.page_path,
		       fs.fields,
		       fs.created_at,
		       s.name AS site_name,
		       fd.label AS form_label,
		       fd.active AS active,
		       NOT fd.active AS historical
		FROM form_submissions fs
		INNER JOIN sites s ON s.id = fs.site
		INNER JOIN form_definitions fd ON fd.id = fs.form
		WHERE s."user" = $1
	`
	if filters.SiteID != "" {
		args = append(args, filters.SiteID)
		query += fmt.Sprintf(` AND fs.site = $%d`, len(args))
	}
	if filters.FormID != "" {
		args = append(args, filters.FormID)
		query += fmt.Sprintf(` AND fs.form = $%d`, len(args))
	}
	query += ` ORDER BY fs.created_at DESC`

	var submissions []Submission
	err := r.db.SelectContext(ctx, &submissions, query, args...)
	if err != nil {
		return nil, err
	}
	return submissions, nil
}

func (r *Repository) DeleteSubmission(ctx context.Context, userID, submissionID string) (bool, error) {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM form_submissions fs
		USING sites s
		WHERE fs.site = s.id
			AND s."user" = $1
			AND fs.id = $2
	`, userID, submissionID)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	return rows > 0, err
}

func (r *Repository) DeleteBySite(ctx context.Context, siteID string) error {
	if _, err := r.db.ExecContext(ctx, `DELETE FROM form_submissions WHERE site = $1`, siteID); err != nil {
		return err
	}
	if _, err := r.db.ExecContext(ctx, `DELETE FROM form_definitions WHERE site = $1`, siteID); err != nil {
		return err
	}
	_, err := r.db.ExecContext(ctx, `DELETE FROM form_discoveries WHERE site = $1`, siteID)
	return err
}

func (r *Repository) DeleteByUser(ctx context.Context, userID string) error {
	if _, err := r.db.ExecContext(ctx, `
		DELETE FROM form_submissions
		WHERE site IN (SELECT id FROM sites WHERE "user" = $1)
	`, userID); err != nil {
		return err
	}
	if _, err := r.db.ExecContext(ctx, `
		DELETE FROM form_definitions
		WHERE site IN (SELECT id FROM sites WHERE "user" = $1)
	`, userID); err != nil {
		return err
	}
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM form_discoveries
		WHERE site IN (SELECT id FROM sites WHERE "user" = $1)
	`, userID)
	return err
}
