package site

import (
	"context"
	"errors"
	"hostamat/internal/file"
	"hostamat/internal/subdomain"
	"hostamat/internal/uploadsession"
	"io"
	"strings"
	"time"

	"github.com/platforma-dev/platforma/auth"

	"github.com/google/uuid"
)

type repository interface {
	Get(ctx context.Context, id string) (*Site, error)
	GetByDomain(ctx context.Context, domain string) (*Site, error)
	GetList(ctx context.Context, userId string) ([]*Site, error)
	Create(ctx context.Context, site *Site) error
	Delete(ctx context.Context, id string) error
}

type subdomainService interface {
	Create(ctx context.Context) (*subdomain.Subdomain, error)
	Delete(ctx context.Context, id string) error
}

type Service struct {
	repository           repository
	subdomainService     subdomainService
	fileService          *file.Service
	uploadSessionService *uploadsession.Service
}

func NewService(repository repository, subdomainService subdomainService, fileService *file.Service, uploadSessionService *uploadsession.Service) *Service {
	return &Service{
		repository:           repository,
		subdomainService:     subdomainService,
		fileService:          fileService,
		uploadSessionService: uploadSessionService,
	}
}

func (s *Service) Create(ctx context.Context, name string) (*Site, error) {
	var err error

	// GetList also checks access so no need to check it further...
	sites, err := s.GetList(ctx)
	if err != nil {
		return nil, err
	}

	if len(sites) >= maxSites {
		return nil, errors.New("max sites limit reached")
	}

	// ... but access is still checked beacuse we need user info
	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, errors.New("access denied")
	}

	subdomain, err := s.subdomainService.Create(ctx)
	if err != nil {
		return nil, err
	}

	site := &Site{
		Id:        uuid.NewString(),
		Name:      name,
		User:      user.ID,
		Subdomain: subdomain.Id,
		Domain:    subdomain.FullDomain,
		CreatedAt: time.Now(),
	}

	err = s.repository.Create(ctx, site)
	if err != nil {
		return nil, err
	}

	return site, nil
}

func (s *Service) Get(ctx context.Context, id string) (*Site, error) {
	site, err := s.repository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if ok := s.checkAccess(ctx, site); !ok {
		return nil, errors.New("access denied")
	}

	return site, nil
}

func (s *Service) GetIdByDomain(ctx context.Context, domain string) (string, error) {
	site, err := s.repository.GetByDomain(ctx, domain)
	if err != nil {
		return "", err
	}

	return site.Id, nil
}

func (s *Service) GetByDomain(ctx context.Context, domain string) (*Site, error) {
	site, err := s.repository.GetByDomain(ctx, domain)
	if err != nil {
		return nil, err
	}

	if ok := s.checkAccess(ctx, site); !ok {
		return nil, errors.New("access denied")
	}

	return site, nil
}

func (s *Service) GetByDomainNoAuth(ctx context.Context, domain string) (*Site, error) {
	return s.repository.GetByDomain(ctx, domain)
}

func (s *Service) GetList(ctx context.Context) ([]*Site, error) {
	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, errors.New("access denied")
	}

	return s.repository.GetList(ctx, user.ID)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	site, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	err = s.subdomainService.Delete(ctx, site.Subdomain)
	if err != nil {
		return err
	}

	err = s.DeleteFiles(ctx, site.Id)
	if err != nil {
		return err
	}

	// `s.Get` already checks permission, so can delete without additional checks
	return s.repository.Delete(ctx, id)
}

func (s *Service) DeleteFiles(ctx context.Context, siteId string) error {
	return s.fileService.DeleteFromFolder(ctx, siteId)
}

func (s *Service) CreateUploadSession(ctx context.Context, siteId string) (*uploadsession.UploadSession, error) {
	// Get site to check access
	site, err := s.Get(ctx, siteId)
	if err != nil {
		return nil, err
	}

	err = s.DeleteFiles(ctx, site.Domain)
	if err != nil {
		return nil, err
	}

	return s.uploadSessionService.Create(ctx, site.Id)
}

func (s *Service) CloseUploadSession(ctx context.Context, uploadSessionId string) error {
	us, err := s.uploadSessionService.Get(ctx, uploadSessionId)
	if err != nil {
		return err
	}

	_, err = s.Get(ctx, us.Site)
	if err != nil {
		return err
	}

	return s.uploadSessionService.Close(ctx, us.Id)
}

func (s *Service) GetFile(ctx context.Context, host, filePath string) (io.Reader, error) {
	site, err := s.GetByDomainNoAuth(ctx, host)
	if err != nil {
		return nil, err
	}

	filePath = strings.Trim(filePath, "/")

	return s.fileService.GetFile(ctx, site.Id+"/"+filePath)
}

func (s *Service) UploadFile(ctx context.Context, uploadSessionId, filename string, file io.Reader, size int64) ([]string, []string, error) {
	us, err := s.uploadSessionService.Get(ctx, uploadSessionId)
	if err != nil {
		return nil, nil, err
	}

	_, err = s.Get(ctx, us.Site)
	if err != nil {
		return nil, nil, err
	}

	return s.fileService.UploadFile(ctx, us.Site, filename, file, size)
}

func (s *Service) checkAccess(ctx context.Context, site *Site) bool {
	user := auth.UserFromContext(ctx)
	if user == nil {
		return false
	}

	if site.User == user.ID {
		return true
	}

	return false
}
