package uploadsession

import "github.com/jmoiron/sqlx"

const (
	tableName = "upload_sessions"
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
			created_at TIMESTAMP NOT NULL,
			status VARCHAR(255) NOT NULL,
			expires_at TIMESTAMP NOT NULL,
			resolved_at TIMESTAMP
		)
	`
	_, err := r.db.Exec(query)
	return err
}

func (r *Repository) Get(id string) (*UploadSession, error) {
	var session UploadSession
	err := r.db.Get(&session, "SELECT * FROM "+tableName+" WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *Repository) Create(session *UploadSession) error {
	query := `
		INSERT INTO ` + tableName + ` (id, created_at, status, expires_at, resolved_at)
		VALUES (:id, :created_at, :status, :expires_at, :resolved_at)
	`
	_, err := r.db.NamedExec(query, session)
	return err
}

func (r *Repository) Update(session *UploadSession) error {
	query := `
		UPDATE ` + tableName + `
		SET status = :status, resolved_at = :resolved_at
		WHERE id = :id
	`
	_, err := r.db.NamedExec(query, session)
	return err
}
