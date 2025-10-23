package uploadsession

import "hostamat/pkg/platforma/database"

type Repository struct {
	db database.Database
}

func NewRepository(db database.Database) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Schema() ([]database.Migration, database.Schema) {
	return []database.Migration{}, database.Schema{Queries: []string{
		`CREATE TABLE IF NOT EXISTS uploadsessions (
			id VARCHAR(255) PRIMARY KEY,
			"user" VARCHAR(255) NOT NULL,
			site VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
			status VARCHAR(255) NOT NULL
		)`}}
}
