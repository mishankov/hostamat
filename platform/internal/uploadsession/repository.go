package uploadsession

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	tableName = "upload_sessions"
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
			site_id VARCHAR(255) NOT NULL,
			status VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := r.db.Exec(query)
	return err
}

func (r *Repository) Create(ctx context.Context, session *UploadSession) error {
	const query = `INSERT INTO ` + tableName + ` (id, site_id, status) VALUES (:id, :site_id, :status)`

	_, err := r.db.NamedExecContext(ctx, query, session)
	if err != nil {
		return fmt.Errorf("failed to create upload session: %w", err)
	}

	return nil
}

func (r *Repository) Get(ctx context.Context, id string) (*UploadSession, error) {
	const query = `SELECT * FROM ` + tableName + ` WHERE id = $1`

	var session UploadSession
	err := r.db.GetContext(ctx, &session, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get upload session: %w", err)
	}

	return &session, nil
}

func (r *Repository) Update(ctx context.Context, session *UploadSession) error {
	query := `UPDATE ` + tableName + ` SET status = :status, updated_at = :updated_at WHERE id = :id`

	_, err := r.db.NamedExecContext(ctx, query, session)
	if err != nil {
		return fmt.Errorf("failed to update upload session status: %w", err)
	}

	return nil
}
