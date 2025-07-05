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

func (r *Repository) Get(id string) (*UploadSession, error) {
	var session UploadSession
	err := r.db.Get(&session, "SELECT * FROM upload_sessions WHERE id = $1", id) //use tableName constant AI!
	if err != nil {
		return nil, err
	}
	return &session, nil
}
