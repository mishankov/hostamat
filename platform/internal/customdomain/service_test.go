package customdomain_test

import (
	"context"
	"database/sql"
	"errors"
	"net"
	"testing"
	"time"

	"hostamat/internal/customdomain"
	"hostamat/internal/site"

	"github.com/platforma-dev/platforma/auth"
	"github.com/platforma-dev/platforma/log"
)

func TestCreateNormalizesHostnameAndRejectsDuplicates(t *testing.T) {
	t.Parallel()

	siteRepo := &siteRepositoryMock{
		sites: map[string]*site.Site{
			"site-1": {
				Id:     "site-1",
				User:   "user-1",
				Name:   "Docs",
				Domain: "docs.hostamat.test",
			},
		},
	}
	repo := &repositoryMock{}

	service := customdomain.NewService(
		repo,
		siteRepo,
		&hostSyncMock{},
		&dnsLookupMock{},
		customdomain.DNSConfig{VerificationPrefix: "_hostamat"},
	)

	record, err := service.Create(testContext("user-1"), "site-1", "WWW.Example.COM.")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if record.Hostname != "www.example.com" {
		t.Fatalf("Create() hostname = %q, want %q", record.Hostname, "www.example.com")
	}

	siteRepo.siteByDomain = map[string]*site.Site{
		"www.example.com": {
			Id:     "site-2",
			User:   "user-2",
			Domain: "www.example.com",
		},
	}

	_, err = service.Create(testContext("user-1"), "site-1", "www.example.com")
	if err == nil || err.Error() != "Этот домен уже подключен к другому сайту" {
		t.Fatalf("Create() duplicate error = %v", err)
	}
}

func TestVerifyActivatesDomainAfterSuccessfulChecks(t *testing.T) {
	t.Parallel()

	now := time.Now()
	repo := &repositoryMock{
		domains: map[string]*customdomain.CustomDomain{
			"domain-1": {
				Id:                "domain-1",
				Site:              "site-1",
				User:              "user-1",
				Hostname:          "www.example.com",
				Type:              customdomain.DomainTypeCustom,
				Status:            customdomain.DomainStatusPending,
				VerificationToken: "token",
				CreatedAt:         now,
				UpdatedAt:         now,
			},
		},
	}
	siteRepo := &siteRepositoryMock{
		sites: map[string]*site.Site{
			"site-1": {
				Id:     "site-1",
				User:   "user-1",
				Name:   "Docs",
				Domain: "docs.hostamat.test",
			},
		},
	}
	hostSync := &hostSyncMock{}
	dns := &dnsLookupMock{
		txt: map[string][]string{
			"_hostamat.www.example.com": {"token"},
		},
		cname: map[string]string{
			"www.example.com": "docs.hostamat.test.",
		},
	}

	service := customdomain.NewService(
		repo,
		siteRepo,
		hostSync,
		dns,
		customdomain.DNSConfig{VerificationPrefix: "_hostamat"},
	)

	record, err := service.Verify(testContext("user-1"), "domain-1")
	if err != nil {
		t.Fatalf("Verify() error = %v", err)
	}

	if record.Status != customdomain.DomainStatusActive {
		t.Fatalf("Verify() status = %q, want %q", record.Status, customdomain.DomainStatusActive)
	}

	if len(hostSync.added) != 1 || hostSync.added[0] != "www.example.com" {
		t.Fatalf("host sync added = %v, want [www.example.com]", hostSync.added)
	}
}

func TestVerifyStoresActionableFailure(t *testing.T) {
	t.Parallel()

	now := time.Now()
	repo := &repositoryMock{
		domains: map[string]*customdomain.CustomDomain{
			"domain-1": {
				Id:                "domain-1",
				Site:              "site-1",
				User:              "user-1",
				Hostname:          "example.com",
				Type:              customdomain.DomainTypeCustom,
				Status:            customdomain.DomainStatusPending,
				VerificationToken: "token",
				CreatedAt:         now,
				UpdatedAt:         now,
			},
		},
	}
	siteRepo := &siteRepositoryMock{
		sites: map[string]*site.Site{
			"site-1": {
				Id:     "site-1",
				User:   "user-1",
				Name:   "Docs",
				Domain: "docs.hostamat.test",
			},
		},
	}

	service := customdomain.NewService(
		repo,
		siteRepo,
		&hostSyncMock{},
		&dnsLookupMock{},
		customdomain.DNSConfig{
			VerificationPrefix: "_hostamat",
			ApexIPv4Targets:    []string{"192.0.2.10"},
		},
	)

	record, err := service.Verify(testContext("user-1"), "domain-1")
	if err != nil {
		t.Fatalf("Verify() error = %v", err)
	}

	if record.Status != customdomain.DomainStatusVerificationFailed {
		t.Fatalf("Verify() status = %q, want %q", record.Status, customdomain.DomainStatusVerificationFailed)
	}

	if record.LastVerificationError == "" {
		t.Fatal("Verify() last verification error is empty")
	}
}

func TestResolveSiteByHostSupportsGeneratedAndCustomHosts(t *testing.T) {
	t.Parallel()

	siteRepo := &siteRepositoryMock{
		siteByDomain: map[string]*site.Site{
			"docs.hostamat.test": {
				Id:     "site-1",
				User:   "user-1",
				Name:   "Docs",
				Domain: "docs.hostamat.test",
			},
		},
	}
	repo := &repositoryMock{
		resolvedSites: map[string]*site.Site{
			"www.example.com": {
				Id:     "site-1",
				User:   "user-1",
				Name:   "Docs",
				Domain: "docs.hostamat.test",
			},
		},
	}

	service := customdomain.NewService(
		repo,
		siteRepo,
		&hostSyncMock{},
		&dnsLookupMock{},
		customdomain.DNSConfig{VerificationPrefix: "_hostamat"},
	)

	if _, err := service.ResolveSiteByHost(context.Background(), "docs.hostamat.test:443"); err != nil {
		t.Fatalf("ResolveSiteByHost() generated error = %v", err)
	}

	if _, err := service.ResolveSiteByHost(context.Background(), "www.example.com"); err != nil {
		t.Fatalf("ResolveSiteByHost() custom error = %v", err)
	}
}

func TestDeleteBySiteRemovesActiveDomainsFromProxy(t *testing.T) {
	t.Parallel()

	now := time.Now()
	repo := &repositoryMock{
		bySite: map[string][]*customdomain.CustomDomain{
			"site-1": {
				{
					Id:                "domain-1",
					Site:              "site-1",
					User:              "user-1",
					Hostname:          "www.example.com",
					Type:              customdomain.DomainTypeCustom,
					Status:            customdomain.DomainStatusActive,
					VerificationToken: "token",
					CreatedAt:         now,
					UpdatedAt:         now,
				},
			},
		},
	}
	hostSync := &hostSyncMock{}

	service := customdomain.NewService(
		repo,
		&siteRepositoryMock{},
		hostSync,
		&dnsLookupMock{},
		customdomain.DNSConfig{VerificationPrefix: "_hostamat"},
	)

	if err := service.DeleteBySite(context.Background(), "site-1"); err != nil {
		t.Fatalf("DeleteBySite() error = %v", err)
	}

	if len(hostSync.removed) != 1 || hostSync.removed[0] != "www.example.com" {
		t.Fatalf("host sync removed = %v, want [www.example.com]", hostSync.removed)
	}

	if len(repo.deleted) != 1 || repo.deleted[0] != "domain-1" {
		t.Fatalf("deleted domains = %v, want [domain-1]", repo.deleted)
	}
}

type repositoryMock struct {
	domains       map[string]*customdomain.CustomDomain
	bySite        map[string][]*customdomain.CustomDomain
	resolvedSites map[string]*site.Site
	activeHosts   []string
	deleted       []string
}

func (m *repositoryMock) Create(_ context.Context, domain *customdomain.CustomDomain) error {
	if m.domains == nil {
		m.domains = make(map[string]*customdomain.CustomDomain)
	}
	copyDomain := *domain
	m.domains[domain.Id] = &copyDomain
	return nil
}

func (m *repositoryMock) Get(_ context.Context, id string) (*customdomain.CustomDomain, error) {
	domain, ok := m.domains[id]
	if !ok {
		return nil, sql.ErrNoRows
	}
	copyDomain := *domain
	return &copyDomain, nil
}

func (m *repositoryMock) GetByHostname(_ context.Context, hostname string) (*customdomain.CustomDomain, error) {
	for _, domain := range m.domains {
		if domain.Hostname == hostname {
			copyDomain := *domain
			return &copyDomain, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (m *repositoryMock) GetBySite(_ context.Context, siteID string) ([]*customdomain.CustomDomain, error) {
	domains := m.bySite[siteID]
	result := make([]*customdomain.CustomDomain, 0, len(domains))
	for _, domain := range domains {
		copyDomain := *domain
		result = append(result, &copyDomain)
	}
	return result, nil
}

func (m *repositoryMock) GetListByUser(_ context.Context, _ string) ([]*customdomain.DomainWithSite, error) {
	return nil, nil
}

func (m *repositoryMock) Update(_ context.Context, domain *customdomain.CustomDomain) error {
	if m.domains == nil {
		m.domains = make(map[string]*customdomain.CustomDomain)
	}
	copyDomain := *domain
	m.domains[domain.Id] = &copyDomain
	return nil
}

func (m *repositoryMock) Delete(_ context.Context, id string) error {
	m.deleted = append(m.deleted, id)
	delete(m.domains, id)
	return nil
}

func (m *repositoryMock) ResolveSiteByHostname(_ context.Context, hostname string) (*site.Site, error) {
	if m.resolvedSites == nil {
		return nil, sql.ErrNoRows
	}
	siteRecord, ok := m.resolvedSites[hostname]
	if !ok {
		return nil, sql.ErrNoRows
	}
	return siteRecord, nil
}

func (m *repositoryMock) ListActiveHosts(_ context.Context) ([]string, error) {
	return append([]string(nil), m.activeHosts...), nil
}

type siteRepositoryMock struct {
	sites        map[string]*site.Site
	siteByDomain map[string]*site.Site
}

func (m *siteRepositoryMock) Get(_ context.Context, id string) (*site.Site, error) {
	if m.sites == nil {
		return nil, sql.ErrNoRows
	}
	siteRecord, ok := m.sites[id]
	if !ok {
		return nil, sql.ErrNoRows
	}
	return siteRecord, nil
}

func (m *siteRepositoryMock) GetByDomain(_ context.Context, domain string) (*site.Site, error) {
	if m.siteByDomain == nil {
		return nil, sql.ErrNoRows
	}
	siteRecord, ok := m.siteByDomain[domain]
	if !ok {
		return nil, sql.ErrNoRows
	}
	return siteRecord, nil
}

func (m *siteRepositoryMock) GetList(_ context.Context, _ string) ([]*site.Site, error) {
	result := make([]*site.Site, 0, len(m.sites))
	for _, siteRecord := range m.sites {
		result = append(result, siteRecord)
	}
	return result, nil
}

type hostSyncMock struct {
	added   []string
	removed []string
}

func (m *hostSyncMock) Add(_ context.Context, host string) error {
	m.added = append(m.added, host)
	return nil
}

func (m *hostSyncMock) Remove(_ context.Context, host string) error {
	m.removed = append(m.removed, host)
	return nil
}

type dnsLookupMock struct {
	txt   map[string][]string
	cname map[string]string
	ip    map[string][]net.IPAddr
}

func (m *dnsLookupMock) LookupTXT(_ context.Context, name string) ([]string, error) {
	if m.txt == nil {
		return nil, errors.New("not found")
	}
	values, ok := m.txt[name]
	if !ok {
		return nil, errors.New("not found")
	}
	return values, nil
}

func (m *dnsLookupMock) LookupCNAME(_ context.Context, host string) (string, error) {
	if m.cname == nil {
		return "", errors.New("not found")
	}
	value, ok := m.cname[host]
	if !ok {
		return "", errors.New("not found")
	}
	return value, nil
}

func (m *dnsLookupMock) LookupIPAddr(_ context.Context, host string) ([]net.IPAddr, error) {
	if m.ip == nil {
		return nil, errors.New("not found")
	}
	value, ok := m.ip[host]
	if !ok {
		return nil, errors.New("not found")
	}
	return value, nil
}

func testContext(userID string) context.Context {
	ctx := context.WithValue(context.Background(), auth.UserContextKey, &auth.User{ID: userID})
	return context.WithValue(ctx, log.WideEventKey, log.NewEvent("test"))
}
