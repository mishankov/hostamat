package subdomain

type Domain struct {
	Repository *Repository
	Service    *Service
}

func (d *Domain) GetRepository() any {
	return d.Repository
}

func New(db db, hostSync hostSync, baseDomain string) *Domain {
	repository := NewRepository(db)
	service := NewService(repository, hostSync, baseDomain)

	return &Domain{
		Repository: repository,
		Service:    service,
	}
}
