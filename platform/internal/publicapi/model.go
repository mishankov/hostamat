package publicapi

import (
	"time"

	"hostamat/internal/apikey"
	"hostamat/internal/customdomain"
	"hostamat/internal/site"
)

type meResponse struct {
	User   apiUser `json:"user"`
	APIKey apiKey  `json:"apiKey"`
}

type apiUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type apiKey struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	KeyPrefix string         `json:"keyPrefix"`
	Scopes    []apikey.Scope `json:"scopes"`
	ExpiresAt *time.Time     `json:"expiresAt"`
}

type siteResponse struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	HostamatDomain string    `json:"hostamatDomain"`
	URL            string    `json:"url"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type sitesResponse struct {
	Sites []siteResponse `json:"sites"`
}

type siteEnvelopeResponse struct {
	Site siteResponse `json:"site"`
}

type createSiteRequest struct {
	Name string `json:"name"`
}

type updateSiteRequest struct {
	Name string `json:"name"`
}

type uploadResponse struct {
	Summary uploadSummary        `json:"summary"`
	Files   []uploadFileResponse `json:"files"`
}

type uploadSummary struct {
	UploadedCount int   `json:"uploadedCount"`
	SkippedCount  int   `json:"skippedCount"`
	TotalBytes    int64 `json:"totalBytes"`
}

type uploadFileResponse struct {
	Path   string `json:"path"`
	Size   int64  `json:"size"`
	Status string `json:"status"`
}

type domainsResponse struct {
	Domains []domainResponse `json:"domains"`
}

type domainEnvelopeResponse struct {
	Domain domainResponse `json:"domain"`
}

type domainResponse struct {
	ID                    string                    `json:"id"`
	Hostname              string                    `json:"hostname"`
	Status                customdomain.DomainStatus `json:"status"`
	Site                  linkedSiteResponse        `json:"site"`
	DNSInstructions       []dnsInstruction          `json:"dnsInstructions"`
	LastVerificationError *string                   `json:"lastVerificationError"`
	LastCheckedAt         *time.Time                `json:"lastCheckedAt"`
	ActivatedAt           *time.Time                `json:"activatedAt"`
	CreatedAt             time.Time                 `json:"createdAt"`
	UpdatedAt             time.Time                 `json:"updatedAt"`
}

type linkedSiteResponse struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	HostamatDomain string `json:"hostamatDomain"`
}

type dnsInstruction struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type createDomainRequest struct {
	SiteID   string `json:"siteId"`
	Hostname string `json:"hostname"`
}

type updateDomainRequest struct {
	SiteID string `json:"siteId"`
}

func siteDTO(value *site.Site) siteResponse {
	return siteResponse{
		ID:             value.Id,
		Name:           value.Name,
		HostamatDomain: value.Domain,
		URL:            "https://" + value.Domain,
		CreatedAt:      value.CreatedAt,
		UpdatedAt:      value.UpdatedAt,
	}
}

func domainDTO(value customdomain.DomainRecord) domainResponse {
	instructions := make([]dnsInstruction, 0, 1+len(value.Instructions.Routing))
	if value.Instructions.Ownership.Type != "" {
		instructions = append(instructions, dnsRecordDTO(value.Instructions.Ownership))
	}
	for _, record := range value.Instructions.Routing {
		instructions = append(instructions, dnsRecordDTO(record))
	}

	var lastVerificationError *string
	if value.LastVerificationError != "" {
		lastVerificationError = &value.LastVerificationError
	}

	return domainResponse{
		ID:       value.Id,
		Hostname: value.Hostname,
		Status:   value.Status,
		Site: linkedSiteResponse{
			ID:             value.SiteId,
			Name:           value.SiteName,
			HostamatDomain: value.SiteDomain,
		},
		DNSInstructions:       instructions,
		LastVerificationError: lastVerificationError,
		LastCheckedAt:         value.LastCheckedAt,
		ActivatedAt:           value.ActivatedAt,
		CreatedAt:             value.CreatedAt,
		UpdatedAt:             value.UpdatedAt,
	}
}

func dnsRecordDTO(value customdomain.DNSRecord) dnsInstruction {
	return dnsInstruction{
		Type:  value.Type,
		Name:  value.Name,
		Value: value.Value,
	}
}
