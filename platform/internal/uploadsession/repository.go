package uploadsession

import (
	"context"
	"database/sql"
	"hostamat/internal/uploadsession/migrations"
	"io/fs"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Get(ctx context.Context, id string) (*UploadSession, error) {
	var uploadSession UploadSession
	err := r.db.GetContext(ctx, &uploadSession, `SELECT * FROM uploadsessions WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}

	return &uploadSession, nil
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

func (r *Repository) CreateFile(ctx context.Context, file File) error {
	_, err := r.db.NamedExecContext(ctx, `INSERT INTO files (id, upload_session, name, uploaded, size) VALUES (:id, :upload_session, :name, :uploaded, :size)`, file)
	return err
}

func (r *Repository) GetFilesSizeForUploadSession(ctx context.Context, uploadSessionId string) (int, error) {
	var size sql.NullInt64
	err := r.db.GetContext(ctx, &size, `SELECT SUM(size) FROM files WHERE upload_session = $1`, uploadSessionId)

	return int(size.Int64), err
}

func (r *Repository) GetFiles(ctx context.Context, uploadSessionId string) ([]File, error) {
	var files []File
	err := r.db.SelectContext(ctx, &files, `SELECT * FROM files WHERE upload_session = $1`, uploadSessionId)

	return files, err
}

func (r *Repository) DeleteFilesForUploadSession(ctx context.Context, uploadSessionId string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM files WHERE upload_session = $1`, uploadSessionId)
	return err
}

func (r *Repository) DeleteByUser(ctx context.Context, userId string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM files WHERE upload_session IN (SELECT id FROM uploadsessions WHERE "user" = $1)`, userId)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, `DELETE FROM uploadsessions WHERE "user" = $1`, userId)
	return err
}

func (r *Repository) Migrations() fs.FS {
	return migrations.FS
}
