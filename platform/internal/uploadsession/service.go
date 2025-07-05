package uploadsession

import (
	"time"

	"github.com/google/uuid"
)

type repository interface {
	Get(id string) (*UploadSession, error)
	Create(session *UploadSession) error
	Update(session *UploadSession) error
}

type Service struct {
	repository repository
}

func NewService(repository repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Create() (*UploadSession, error) {
	us := &UploadSession{
		Id:        uuid.NewString(),
		CreatedAt: time.Now(),
		Status:    Open,
		ExpiresAt: time.Now().Add(time.Minute * 5), // TODO: make this configurable
	}

	err := s.repository.Create(us)
	if err != nil {
		return nil, err
	}

	return us, nil
}

func (s *Service) Get(id string) (*UploadSession, error) {
	return s.repository.Get(id)
}

func (s *Service) Cancel(id string) error {
	us, err := s.repository.Get(id)
	if err != nil {
		return err
	}

	us.Status = Cancelled
	us.ResolvedAt = time.Now()
	return s.repository.Update(us)
}

func (s *Service) Finish(id string) error {
	us, err := s.repository.Get(id)
	if err != nil {
		return err
	}

	us.Status = Finished
	us.ResolvedAt = time.Now()
	return s.repository.Update(us)
}
