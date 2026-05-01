package site

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"hostamat/internal/file"
	"hostamat/internal/subdomain"
	"hostamat/internal/uploadsession"
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/platforma-dev/platforma/auth"
	"github.com/platforma-dev/platforma/log"

	"github.com/google/uuid"
)

type repository interface {
	Get(ctx context.Context, id string) (*Site, error)
	GetByDomain(ctx context.Context, domain string) (*Site, error)
	GetList(ctx context.Context, userId string) ([]*Site, error)
	Create(ctx context.Context, site *Site) error
	Update(ctx context.Context, site *Site) error
	Delete(ctx context.Context, id string) error
}

type subdomainService interface {
	Create(ctx context.Context) (*subdomain.Subdomain, error)
	Delete(ctx context.Context, id string) error
}

type customDomainService interface {
	ResolveSiteByHost(ctx context.Context, host string) (*Site, error)
	DeleteBySite(ctx context.Context, siteID string) error
}

type Service struct {
	repository           repository
	subdomainService     subdomainService
	customDomainService  customDomainService
	fileService          *file.Service
	uploadSessionService *uploadsession.Service
}

func NewService(repository repository, subdomainService subdomainService, customDomainService customDomainService, fileService *file.Service, uploadSessionService *uploadsession.Service) *Service {
	return &Service{
		repository:           repository,
		subdomainService:     subdomainService,
		customDomainService:  customDomainService,
		fileService:          fileService,
		uploadSessionService: uploadSessionService,
	}
}

func (s *Service) Create(ctx context.Context, name string) (*Site, error) {
	ev := log.EventFromContext(ctx)

	ev.AddStep(slog.LevelInfo, "site.Create")

	var err error

	// GetList also checks access so no need to check it further...
	sites, err := s.GetList(ctx)
	if err != nil {
		return nil, err
	}

	if len(sites) >= maxSites {

		ev.AddStep(slog.LevelWarn, "site.Create.max_sites_reached")

		return nil, errors.New("max sites limit reached")
	}

	// ... but access is still checked beacuse we need user info
	user := auth.UserFromContext(ctx)
	if user == nil {
		ev.AddAttrs(map[string]any{
			"user.accessDenied": true,
		})

		return nil, errors.New("access denied")
	}

	ev.AddAttrs(map[string]any{
		"user.id": user.ID,
	})

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

	err = s.validate(ctx, site)
	if err != nil {

		ev.AddAttrs(map[string]any{
			"site.id":               site.Id,
			"site.domain":           site.Domain,
			"subdomain.id":          subdomain.Id,
			"subdomain.name":        subdomain.Name,
			"subdomain.full_domain": subdomain.FullDomain,
		})
		ev.AddStep(slog.LevelWarn, "site.Create.invalid_name")

		return nil, err
	}

	err = s.repository.Create(ctx, site)
	if err != nil {

		ev.AddError(err)

		return nil, err
	}

	ev.AddAttrs(map[string]any{
		"site.id":               site.Id,
		"site.domain":           site.Domain,
		"subdomain.id":          subdomain.Id,
		"subdomain.name":        subdomain.Name,
		"subdomain.full_domain": subdomain.FullDomain,
	})

	return site, nil
}

func (s *Service) Get(ctx context.Context, id string) (*Site, error) {
	ev := log.EventFromContext(ctx)

	ev.AddStep(slog.LevelInfo, "site.Get")
	ev.AddAttrs(map[string]any{
		"site.id": id,
	})

	site, err := s.repository.Get(ctx, id)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			ev.AddStep(slog.LevelInfo, "site.Get.not_found")
		} else {
			ev.AddError(err)
		}

		return nil, err
	}

	ev.AddAttrs(map[string]any{
		"site.domain": site.Domain,
	})

	if ok := s.checkAccess(ctx, site); !ok {

		ev.AddStep(slog.LevelWarn, "site.Get.access_denied")

		return nil, errors.New("access denied")
	}

	return site, nil
}

func (s *Service) GetIdByDomain(ctx context.Context, domain string) (string, error) {
	ev := log.EventFromContext(ctx)

	ev.AddStep(slog.LevelInfo, "site.GetIdByDomain")
	ev.AddAttrs(map[string]any{
		"site.domain": domain,
	})

	site, err := s.resolveByHost(ctx, domain)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			ev.AddStep(slog.LevelInfo, "site.GetIdByDomain.not_found")
		} else {
			ev.AddError(err)
		}

		return "", err
	}

	ev.AddAttrs(map[string]any{
		"site.id": site.Id,
	})

	return site.Id, nil
}

func (s *Service) GetByDomain(ctx context.Context, domain string) (*Site, error) {
	ev := log.EventFromContext(ctx)

	ev.AddStep(slog.LevelInfo, "site.GetByDomain")
	ev.AddAttrs(map[string]any{
		"site.domain": domain,
	})

	site, err := s.resolveByHost(ctx, domain)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			ev.AddStep(slog.LevelInfo, "site.GetByDomain.not_found")
		} else {
			ev.AddError(err)
		}

		return nil, err
	}

	ev.AddAttrs(map[string]any{
		"site.id": site.Id,
	})

	if ok := s.checkAccess(ctx, site); !ok {

		ev.AddStep(slog.LevelWarn, "site.GetByDomain.access_denied")

		return nil, errors.New("access denied")
	}

	return site, nil
}

func (s *Service) GetByDomainNoAuth(ctx context.Context, domain string) (*Site, error) {
	ev := log.EventFromContext(ctx)

	ev.AddStep(slog.LevelInfo, "site.GetByDomainNoAuth")
	ev.AddAttrs(map[string]any{
		"site.domain": domain,
	})

	site, err := s.resolveByHost(ctx, domain)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			ev.AddStep(slog.LevelInfo, "site.GetByDomainNoAuth.not_found")
		} else {
			ev.AddError(err)
		}

		return nil, err
	}

	ev.AddAttrs(map[string]any{
		"site.id": site.Id,
	})

	return site, nil
}

func (s *Service) GetList(ctx context.Context) ([]*Site, error) {
	ev := log.EventFromContext(ctx)

	ev.AddStep(slog.LevelInfo, "site.GetList")

	user := auth.UserFromContext(ctx)
	if user == nil {
		ev.AddAttrs(map[string]any{
			"user.accessDenied": true,
		})

		return nil, errors.New("access denied")
	}

	ev.AddAttrs(map[string]any{
		"user.id": user.ID,
	})

	sites, err := s.repository.GetList(ctx, user.ID)
	if err != nil {
		ev.AddError(err)
	}

	return sites, err
}

func (s *Service) UpdateName(ctx context.Context, id, name string) error {
	ev := log.EventFromContext(ctx)

	ev.AddStep(slog.LevelInfo, "site.UpdateName")
	ev.AddAttrs(map[string]any{
		"site.id": id,
	})

	site, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	site.Name = name
	site.UpdatedAt = time.Now()

	err = s.validate(ctx, site)
	if err != nil {

		ev.AddStep(slog.LevelWarn, "site.UpdateName.invalid_name")

		return err
	}

	err = s.repository.Update(ctx, site)
	if err != nil {
		ev.AddError(err)
	}

	return err
}

func (s *Service) Delete(ctx context.Context, id string) error {
	ev := log.EventFromContext(ctx)

	ev.AddStep(slog.LevelInfo, "site.Delete")
	ev.AddAttrs(map[string]any{
		"site.id": id,
	})

	site, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	ev.AddAttrs(map[string]any{
		"site.domain":  site.Domain,
		"subdomain.id": site.Subdomain,
	})

	if s.customDomainService != nil {
		err = s.customDomainService.DeleteBySite(ctx, site.Id)
		if err != nil {
			return err
		}
	}

	err = s.subdomainService.Delete(ctx, site.Subdomain)
	if err != nil {
		return err
	}

	err = s.fileService.DeleteFromFolderAsync(ctx, site.Id)
	if err != nil {
		return err
	}

	// `s.Get` already checks permission, so can delete without additional checks
	err = s.repository.Delete(ctx, id)
	if err != nil {

		ev.AddError(err)

		return err
	}

	return nil
}

func (s *Service) resolveByHost(ctx context.Context, host string) (*Site, error) {
	if s.customDomainService != nil {
		return s.customDomainService.ResolveSiteByHost(ctx, host)
	}

	return s.repository.GetByDomain(ctx, host)
}

func (s *Service) CreateUploadSession(ctx context.Context, siteId string) (*uploadsession.UploadSession, error) {
	ev := log.EventFromContext(ctx)

	ev.AddStep(slog.LevelInfo, "site.CreateUploadSession")
	ev.AddAttrs(map[string]any{
		"site.id": siteId,
	})

	// Get site to check access
	site, err := s.Get(ctx, siteId)
	if err != nil {
		return nil, err
	}

	ev.AddAttrs(map[string]any{
		"site.domain": site.Domain,
	})

	err = s.fileService.DeleteFromFolder(ctx, site.Id)
	if err != nil {
		return nil, err
	}

	us, err := s.uploadSessionService.Create(ctx, site.Id)
	if err != nil {
		return nil, err
	}

	ev.AddAttrs(map[string]any{
		"uploadsession.id": us.Id,
	})

	return us, nil
}

func (s *Service) CloseUploadSession(ctx context.Context, uploadSessionId string) (*UploadSummary, error) {
	ev := log.EventFromContext(ctx)

	ev.AddStep(slog.LevelInfo, "site.CloseUploadSession")
	ev.AddAttrs(map[string]any{
		"uploadsession.id": uploadSessionId,
	})

	us, err := s.uploadSessionService.Get(ctx, uploadSessionId)
	if err != nil {
		return nil, err
	}

	ev.AddAttrs(map[string]any{
		"site.id": us.Site,
	})

	site, err := s.Get(ctx, us.Site)
	if err != nil {
		return nil, err
	}

	files, err := s.uploadSessionService.GetFiles(ctx, us.Id)
	if err != nil {
		return nil, err
	}

	err = s.uploadSessionService.Close(ctx, us.Id)
	if err != nil {
		return nil, err
	}

	return BuildUploadSummary(site.Domain, files), nil
}

func (s *Service) GetFile(ctx context.Context, host, filePath string) (io.Reader, string, error) {
	ev := log.EventFromContext(ctx)

	ev.AddStep(slog.LevelInfo, "site.GetFile")
	ev.AddAttrs(map[string]any{
		"site.domain": host,
		"file.path":   filePath,
	})

	site, err := s.GetByDomainNoAuth(ctx, host)
	if err != nil {
		return nil, "", err
	}

	ev.AddAttrs(map[string]any{
		"site.id": site.Id,
	})

	filePath = strings.Trim(filePath, "/")
	if filePath == "" {
		filePath = "index.html"
	}

	ev.AddAttrs(map[string]any{
		"file.path": filePath,
	})

	reader, resolvedPath, err := s.fileService.GetFile(ctx, site.Id+"/"+filePath)
	if err != nil {
		return nil, "", err
	}

	// Strip the site ID prefix from the resolved path
	resolvedPath = strings.TrimPrefix(resolvedPath, site.Id+"/")

	ev.AddAttrs(map[string]any{
		"file.resolved_path": resolvedPath,
	})

	return reader, resolvedPath, nil
}

func (s *Service) UploadFile(ctx context.Context, uploadSessionId, filename string, fileReader io.Reader, size int64) error {
	ev := log.EventFromContext(ctx)

	ev.AddStep(slog.LevelInfo, "site.UploadFile")
	ev.AddAttrs(map[string]any{
		"uploadsession.id": uploadSessionId,
		"file.name":        filename,
		"file.size":        size,
	})

	currentSize, err := s.uploadSessionService.GetSize(ctx, uploadSessionId)
	if err != nil {
		return err
	}

	if int64(currentSize)+size > maxFilesSize {

		ev.AddStep(slog.LevelWarn, "site.UploadFile.max_size_exceeded")

		return errUploadSize
	}

	us, err := s.uploadSessionService.Get(ctx, uploadSessionId)
	if err != nil {
		return err
	}

	ev.AddAttrs(map[string]any{
		"site.id": us.Site,
	})

	_, err = s.Get(ctx, us.Site)
	if err != nil {
		return err
	}

	uploaded, skipped, err := s.fileService.UploadFile(ctx, us.Site, filename, fileReader, size)
	if err != nil {
		if errors.Is(err, file.ErrInvalidFileType) {

			ev.AddStep(slog.LevelWarn, "site.UploadFile.invalid_file_type")

			return handlerError{ErrorCode: errorCodeInvalidExt, ErrorMessage: fmt.Sprintf("У файла %s недопустимое расширение", filename)}
		}

		return err
	}

	for _, file := range uploaded {
		s.uploadSessionService.CreateFile(ctx, us.Id, file.Name, file.Size, true)
	}

	for _, file := range skipped {
		s.uploadSessionService.CreateFile(ctx, us.Id, file.Name, file.Size, false)
	}

	return nil
}

func (s *Service) validate(_ context.Context, site *Site) error {
	if len(site.Name) > siteNameMaxLength {
		return errLongName
	}

	if len(site.Name) == 0 {
		return errShortName
	}

	return nil
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
