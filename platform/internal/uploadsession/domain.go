package uploadsession

import (
	"github.com/jmoiron/sqlx"
)

type Domain struct {
	Repository *Repository
	Service    *Service
}

func (d *Domain) GetRepository() any {
	return d.Repository
}

func New(db *sqlx.DB) *Domain {
	repository := NewRepository(db)
	service := NewService(repository)

	return &Domain{
		Repository: repository,
		Service:    service,
	}
}
