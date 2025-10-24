package uploadsession

import (
	"context"

	"github.com/mishankov/platforma/database"
)

type Repository struct {
	db *database.Database
}

func NewRepository(db *database.Database) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, uploadSession *UploadSession) error {
	_, err := r.db.NamedExecContext(ctx, `INSERT INTO uploadsessions (id, "user", site, created_at, updated_at, status) VALUES (:id, :user, :site, :created_at, :updated_at, :status)`, uploadSession)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateStatus(ctx context.Context, uploadSession *UploadSession) error {
	_, err := r.db.NamedExecContext(ctx, `UPDATE uploadsessions SET status = :status, updated_at = :updated_at WHERE id = :id`, uploadSession)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Schema() ([]database.Migration, database.Schema) {
	return []database.Migration{}, database.Schema{Queries: []string{
		`CREATE TABLE IF NOT EXISTS uploadsessions (
			id VARCHAR(255) PRIMARY KEY,
			"user" VARCHAR(255) NOT NULL,
			site VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			status VARCHAR(255) NOT NULL
		)`}}
}
