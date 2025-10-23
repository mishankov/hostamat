package uploadsession

import "hostamat/pkg/platforma/database"

type Domain struct {
	Repository *Repository
	Service    *Service
}

func (d *Domain) GetRepository() any {
	return d.Repository
}

func New(db database.Database) *Domain {
	repository := NewRepository(db)
	service := NewService(repository)

	return &Domain{
		Repository: repository,
		Service:    service,
	}
}
