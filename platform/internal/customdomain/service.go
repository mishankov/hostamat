package customdomain

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"sort"
	"strings"
	"time"

	"hostamat/internal/site"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/platforma-dev/platforma/auth"
	"github.com/platforma-dev/platforma/log"
	"golang.org/x/net/publicsuffix"
)

type repository interface {
	Create(ctx context.Context, domain *CustomDomain) error
	Get(ctx context.Context, id string) (*CustomDomain, error)
	GetActiveByHostname(ctx context.Context, hostname string) (*CustomDomain, error)
	GetBySite(ctx context.Context, siteID string) ([]*CustomDomain, error)
	GetListByUser(ctx context.Context, userID string) ([]*DomainWithSite, error)
	Update(ctx context.Context, domain *CustomDomain) error
	Delete(ctx context.Context, id string) error
	DeleteStaleUnverified(ctx context.Context, cutoff time.Time) (int64, error)
	ResolveSiteByHostname(ctx context.Context, hostname string) (*site.Site, error)
	ListActiveHosts(ctx context.Context) ([]string, error)
}

type siteRepository interface {
	Get(ctx context.Context, id string) (*site.Site, error)
	GetByDomain(ctx context.Context, domain string) (*site.Site, error)
	GetList(ctx context.Context, userID string) ([]*site.Site, error)
}

type hostSync interface {
	Add(ctx context.Context, host string) error
	Remove(ctx context.Context, host string) error
}

type dnsLookup interface {
	LookupTXT(ctx context.Context, name string) ([]string, error)
	LookupCNAME(ctx context.Context, host string) (string, error)
}

type DNSConfig struct {
	VerificationPrefix string
	CNAMETarget        string
}

type Service struct {
	repo     repository
	siteRepo siteRepository
	hostSync hostSync
	dns      dnsLookup
	dnsCfg   DNSConfig
}

func NewService(repo repository, siteRepo siteRepository, hostSync hostSync, dns dnsLookup, dnsCfg DNSConfig) *Service {
	return &Service{
		repo:     repo,
		siteRepo: siteRepo,
		hostSync: hostSync,
		dns:      dns,
		dnsCfg:   dnsCfg,
	}
}

func (s *Service) List(ctx context.Context) ([]DomainRecord, error) {
	ev := log.EventFromContext(ctx)
	ev.AddStep(slog.LevelInfo, "customdomain.List")

	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, errAccessDenied
	}

	sites, err := s.siteRepo.GetList(ctx, user.ID)
	if err != nil {
		ev.AddError(err)
		return nil, err
	}

	customDomains, err := s.repo.GetListByUser(ctx, user.ID)
	if err != nil {
		ev.AddError(err)
		return nil, err
	}

	records := make([]DomainRecord, 0, len(sites)+len(customDomains))
	for _, currentSite := range sites {
		records = append(records, DomainRecord{
			Id:         "generated:" + currentSite.Id,
			SiteId:     currentSite.Id,
			SiteName:   currentSite.Name,
			SiteDomain: currentSite.Domain,
			Hostname:   currentSite.Domain,
			Type:       DomainTypeGenerated,
			Status:     DomainStatusActive,
			CreatedAt:  currentSite.CreatedAt,
			UpdatedAt:  currentSite.UpdatedAt,
		})
	}

	for _, current := range customDomains {
		records = append(records, s.recordFromJoinedDomain(current))
	}

	sort.SliceStable(records, func(i, j int) bool {
		if records[i].SiteName == records[j].SiteName {
			if records[i].Type != records[j].Type {
				return records[i].Type == DomainTypeGenerated
			}

			return records[i].Hostname < records[j].Hostname
		}

		return records[i].SiteName < records[j].SiteName
	})

	return records, nil
}

func (s *Service) Create(ctx context.Context, siteID, hostname string) (*DomainRecord, error) {
	ev := log.EventFromContext(ctx)
	ev.AddStep(slog.LevelInfo, "customdomain.Create")

	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, errAccessDenied
	}

	if siteID == "" {
		return nil, errMissingSiteID
	}

	normalizedHostname, err := normalizeHostname(hostname)
	if err != nil {
		return nil, err
	}

	isApex, err := isApexHostname(normalizedHostname)
	if err != nil {
		return nil, err
	}

	if isApex {
		return nil, errApexHostnameRejected
	}

	siteRecord, err := s.siteRepo.Get(ctx, siteID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errNotFound
		}

		ev.AddError(err)
		return nil, err
	}

	if siteRecord.User != user.ID {
		return nil, errAccessDenied
	}

	if normalizedHostname == siteRecord.Domain {
		return nil, errHostnameTaken
	}

	if _, err := s.siteRepo.GetByDomain(ctx, normalizedHostname); err == nil {
		return nil, errHostnameTaken
	} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
		ev.AddError(err)
		return nil, err
	}

	if _, err := s.repo.GetActiveByHostname(ctx, normalizedHostname); err == nil {
		return nil, errHostnameTaken
	} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
		ev.AddError(err)
		return nil, err
	}

	now := time.Now()
	domain := &CustomDomain{
		Id:                    uuid.NewString(),
		Site:                  siteRecord.Id,
		User:                  user.ID,
		Hostname:              normalizedHostname,
		Type:                  DomainTypeCustom,
		Status:                DomainStatusPending,
		VerificationToken:     uuid.NewString(),
		LastVerificationError: "",
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	if err := s.repo.Create(ctx, domain); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return nil, errHostnameTaken
		}

		ev.AddError(err)
		return nil, err
	}

	record := s.recordFromDomain(domain, siteRecord)

	return &record, nil
}

func (s *Service) UpdateSite(ctx context.Context, domainID, siteID string) (*DomainRecord, error) {
	ev := log.EventFromContext(ctx)
	ev.AddStep(slog.LevelInfo, "customdomain.UpdateSite")

	if siteID == "" {
		return nil, errMissingSiteID
	}

	domain, currentSiteRecord, err := s.getOwnedDomain(ctx, domainID)
	if err != nil {
		return nil, err
	}

	newSiteRecord, err := s.siteRepo.Get(ctx, siteID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errNotFound
		}

		ev.AddError(err)
		return nil, err
	}

	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, errAccessDenied
	}

	if newSiteRecord.User != user.ID {
		return nil, errAccessDenied
	}

	if domain.Site == siteID {
		record := s.recordFromDomain(domain, currentSiteRecord)
		return &record, nil
	}

	now := time.Now()
	domain.Site = siteID
	domain.UpdatedAt = now

	if err := s.repo.Update(ctx, domain); err != nil {
		ev.AddError(err)
		return nil, err
	}

	record := s.recordFromDomain(domain, newSiteRecord)
	return &record, nil
}

func (s *Service) Verify(ctx context.Context, domainID string) (*DomainRecord, error) {
	ev := log.EventFromContext(ctx)
	ev.AddStep(slog.LevelInfo, "customdomain.Verify")

	domain, siteRecord, err := s.getOwnedDomain(ctx, domainID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	domain.LastCheckedAt = sql.NullTime{Time: now, Valid: true}
	domain.UpdatedAt = now

	if err := s.verifyOwnership(ctx, domain); err != nil {
		domain.Status = DomainStatusVerificationFailed
		domain.LastVerificationError = err.Error()
		domain.ActivatedAt = sql.NullTime{}
		if updateErr := s.repo.Update(ctx, domain); updateErr != nil {
			ev.AddError(updateErr)
			return nil, updateErr
		}

		record := s.recordFromDomain(domain, siteRecord)
		return &record, nil
	}

	if err := s.verifyRouting(ctx, domain, siteRecord); err != nil {
		domain.Status = DomainStatusVerificationFailed
		domain.LastVerificationError = err.Error()
		domain.ActivatedAt = sql.NullTime{}
		if updateErr := s.repo.Update(ctx, domain); updateErr != nil {
			ev.AddError(updateErr)
			return nil, updateErr
		}

		record := s.recordFromDomain(domain, siteRecord)
		return &record, nil
	}

	domain.Status = DomainStatusVerified
	domain.LastVerificationError = ""
	domain.ActivatedAt = sql.NullTime{}
	domain.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, domain); err != nil {
		ev.AddError(err)
		return nil, err
	}

	if err := s.hostSync.Add(ctx, domain.Hostname); err != nil {
		domain.Status = DomainStatusVerified
		domain.LastVerificationError = "Домен подтвержден, но Hostamat не смог активировать его в прокси. Повторите проверку позже."
		domain.ActivatedAt = sql.NullTime{}
		domain.UpdatedAt = time.Now()
		if updateErr := s.repo.Update(ctx, domain); updateErr != nil {
			ev.AddError(updateErr)
			return nil, updateErr
		}

		record := s.recordFromDomain(domain, siteRecord)
		return &record, nil
	}

	activatedAt := time.Now()
	domain.Status = DomainStatusActive
	domain.LastVerificationError = ""
	domain.ActivatedAt = sql.NullTime{Time: activatedAt, Valid: true}
	domain.UpdatedAt = activatedAt
	if err := s.repo.Update(ctx, domain); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return nil, errHostnameTaken
		}

		ev.AddError(err)
		return nil, err
	}

	record := s.recordFromDomain(domain, siteRecord)
	return &record, nil
}

func (s *Service) Delete(ctx context.Context, domainID string) error {
	domain, _, err := s.getOwnedDomain(ctx, domainID)
	if err != nil {
		return err
	}

	return s.deleteDomain(ctx, domain)
}

func (s *Service) DeleteBySite(ctx context.Context, siteID string) error {
	domains, err := s.repo.GetBySite(ctx, siteID)
	if err != nil {
		return err
	}

	for _, domain := range domains {
		if err := s.deleteDomain(ctx, domain); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) ResolveSiteByHost(ctx context.Context, host string) (*site.Site, error) {
	normalizedHost := normalizePublicHost(host)

	if normalizedHost == "" {
		return nil, sql.ErrNoRows
	}

	siteRecord, err := s.siteRepo.GetByDomain(ctx, normalizedHost)
	if err == nil {
		return siteRecord, nil
	}

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return s.repo.ResolveSiteByHostname(ctx, normalizedHost)
}

func (s *Service) ListActiveHosts(ctx context.Context) ([]string, error) {
	return s.repo.ListActiveHosts(ctx)
}

func (s *Service) DeleteStaleUnverified(ctx context.Context, olderThan time.Duration) (int64, error) {
	if olderThan <= 0 {
		return 0, fmt.Errorf("olderThan must be positive")
	}

	cutoff := time.Now().Add(-olderThan)
	deleted, err := s.repo.DeleteStaleUnverified(ctx, cutoff)
	if err != nil {
		return 0, err
	}

	log.InfoContext(ctx, "deleted stale custom domains", "deleted", deleted, "cutoff", cutoff)

	return deleted, nil
}

func (s *Service) getOwnedDomain(ctx context.Context, domainID string) (*CustomDomain, *site.Site, error) {
	if domainID == "" {
		return nil, nil, errMissingDomainID
	}

	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, nil, errAccessDenied
	}

	domain, err := s.repo.Get(ctx, domainID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, errNotFound
		}

		return nil, nil, err
	}

	if domain.User != user.ID {
		return nil, nil, errAccessDenied
	}

	siteRecord, err := s.siteRepo.Get(ctx, domain.Site)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, errNotFound
		}

		return nil, nil, err
	}

	return domain, siteRecord, nil
}

func (s *Service) deleteDomain(ctx context.Context, domain *CustomDomain) error {
	if domain.Status == DomainStatusActive {
		if err := s.hostSync.Remove(ctx, domain.Hostname); err != nil {
			return err
		}
	}

	return s.repo.Delete(ctx, domain.Id)
}

func (s *Service) verifyOwnership(ctx context.Context, domain *CustomDomain) error {
	txtName := verificationTXTName(s.dnsCfg.VerificationPrefix, domain.Hostname)
	values, err := s.dns.LookupTXT(ctx, txtName)
	if err != nil {
		return fmt.Errorf("Не удалось подтвердить TXT запись %s со значением %s", txtName, domain.VerificationToken)
	}

	for _, value := range values {
		if strings.TrimSpace(value) == domain.VerificationToken {
			return nil
		}
	}

	return fmt.Errorf("TXT запись %s должна содержать значение %s", txtName, domain.VerificationToken)
}

func (s *Service) verifyRouting(ctx context.Context, domain *CustomDomain, siteRecord *site.Site) error {
	isApex, err := isApexHostname(domain.Hostname)
	if err != nil {
		return err
	}

	if isApex {
		return errApexHostnameRejected
	}

	cname, err := s.dns.LookupCNAME(ctx, domain.Hostname)
	cnameTarget := s.cnameTarget()
	if err != nil {
		return fmt.Errorf("Укажите CNAME запись %s -> %s", domain.Hostname, cnameTarget)
	}

	if normalizePublicHost(cname) != cnameTarget {
		return fmt.Errorf("CNAME запись %s должна указывать на %s", domain.Hostname, cnameTarget)
	}

	return nil
}

func (s *Service) recordFromJoinedDomain(domain *DomainWithSite) DomainRecord {
	return s.recordFromDomain(&domain.CustomDomain, &site.Site{
		Id:     domain.Site,
		Name:   domain.SiteName,
		Domain: domain.SiteDomain,
	})
}

func (s *Service) recordFromDomain(domain *CustomDomain, siteRecord *site.Site) DomainRecord {
	record := DomainRecord{
		Id:                    domain.Id,
		SiteId:                siteRecord.Id,
		SiteName:              siteRecord.Name,
		SiteDomain:            siteRecord.Domain,
		Hostname:              domain.Hostname,
		Type:                  domain.Type,
		Status:                domain.Status,
		VerificationToken:     domain.VerificationToken,
		LastVerificationError: domain.LastVerificationError,
		CreatedAt:             domain.CreatedAt,
		UpdatedAt:             domain.UpdatedAt,
		Instructions:          s.instructionsForDomain(domain.Hostname, domain.VerificationToken),
	}

	if domain.LastCheckedAt.Valid {
		lastCheckedAt := domain.LastCheckedAt.Time
		record.LastCheckedAt = &lastCheckedAt
	}

	if domain.ActivatedAt.Valid {
		activatedAt := domain.ActivatedAt.Time
		record.ActivatedAt = &activatedAt
	}

	return record
}

func (s *Service) instructionsForDomain(hostname, token string) DNSInstructions {
	ownership := DNSRecord{
		Type:  "TXT",
		Name:  verificationTXTName(s.dnsCfg.VerificationPrefix, hostname),
		Value: token,
	}

	return DNSInstructions{
		Ownership: ownership,
		Routing: []DNSRecord{{
			Type:  "CNAME",
			Name:  hostname,
			Value: s.cnameTarget(),
		}},
	}
}

func (s *Service) cnameTarget() string {
	target := normalizePublicHost(s.dnsCfg.CNAMETarget)
	if target == "" {
		return "dns.hostamat.ru"
	}

	return target
}

func normalizeHostname(hostname string) (string, error) {
	hostname = strings.TrimSpace(strings.ToLower(hostname))
	hostname = strings.TrimSuffix(hostname, ".")
	if hostname == "" {
		return "", errMissingHostname
	}

	if strings.Contains(hostname, "://") || strings.Contains(hostname, "/") || strings.Contains(hostname, "@") || strings.Contains(hostname, ":") {
		return "", errInvalidHostname
	}

	labels := strings.Split(hostname, ".")
	if len(labels) < 2 {
		return "", errInvalidHostname
	}

	for _, label := range labels {
		if label == "" || len(label) > 63 {
			return "", errInvalidHostname
		}

		if strings.HasPrefix(label, "-") || strings.HasSuffix(label, "-") {
			return "", errInvalidHostname
		}

		for _, char := range label {
			if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-' {
				continue
			}

			return "", errInvalidHostname
		}
	}

	if _, err := publicsuffix.EffectiveTLDPlusOne(hostname); err != nil {
		return "", errInvalidHostname
	}

	return hostname, nil
}

func normalizePublicHost(host string) string {
	host = strings.TrimSpace(host)
	if value, _, err := net.SplitHostPort(host); err == nil {
		host = value
	}

	host = strings.TrimSuffix(host, ".")
	host = strings.ToLower(host)

	return host
}

func verificationTXTName(prefix, hostname string) string {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		prefix = "_hostamat"
	}

	return prefix + "." + hostname
}

func isApexHostname(hostname string) (bool, error) {
	rootDomain, err := publicsuffix.EffectiveTLDPlusOne(hostname)
	if err != nil {
		return false, errInvalidHostname
	}

	return hostname == rootDomain, nil
}

type netResolver struct{}

func (netResolver) LookupTXT(ctx context.Context, name string) ([]string, error) {
	return net.DefaultResolver.LookupTXT(ctx, name)
}

func (netResolver) LookupCNAME(ctx context.Context, host string) (string, error) {
	return net.DefaultResolver.LookupCNAME(ctx, host)
}
