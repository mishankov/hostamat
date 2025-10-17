package subdomain

type Domain struct {
	Repository *Repository
	Service    *Service
}

func (d *Domain) GetRepository() any {
	return d.Repository
}

func New(db db, webserver webserver, baseDomain, shorthandId string) *Domain {
	repository := NewRepository(db)
	service := NewService(repository, webserver, baseDomain, shorthandId)

	return &Domain{
		Repository: repository,
		Service:    service,
	}
}
