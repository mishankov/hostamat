package file

type Domain struct {
	Repository *Repository
	Service    *Service
}

func (d *Domain) GetRepository() any {
	return d.Repository
}

func New(s3Client s3, bucketName string) *Domain {
	repository := NewRepository(s3Client)
	service := NewService(repository, bucketName)

	return &Domain{
		Repository: repository,
		Service:    service,
	}
}
