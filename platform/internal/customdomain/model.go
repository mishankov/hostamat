package customdomain

import (
	"database/sql"
	"time"
)

type DomainType string

const (
	DomainTypeGenerated DomainType = "generated"
	DomainTypeCustom    DomainType = "custom"
)

type DomainStatus string

const (
	DomainStatusPending            DomainStatus = "pending"
	DomainStatusVerificationFailed DomainStatus = "verification_failed"
	DomainStatusVerified           DomainStatus = "verified"
	DomainStatusActive             DomainStatus = "active"
)

type CustomDomain struct {
	Id                    string       `db:"id"`
	Site                  string       `db:"site"`
	User                  string       `db:"user"`
	Hostname              string       `db:"hostname"`
	Type                  DomainType   `db:"domain_type"`
	Status                DomainStatus `db:"status"`
	VerificationToken     string       `db:"verification_token"`
	LastVerificationError string       `db:"last_verification_error"`
	LastCheckedAt         sql.NullTime `db:"last_checked_at"`
	ActivatedAt           sql.NullTime `db:"activated_at"`
	CreatedAt             time.Time    `db:"created_at"`
	UpdatedAt             time.Time    `db:"updated_at"`
}

type DNSRecord struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type DNSInstructions struct {
	Ownership DNSRecord   `json:"ownership"`
	Routing   []DNSRecord `json:"routing"`
}

type DomainRecord struct {
	Id                    string          `json:"id"`
	SiteId                string          `json:"siteId"`
	SiteName              string          `json:"siteName"`
	SiteDomain            string          `json:"siteDomain"`
	Hostname              string          `json:"hostname"`
	Type                  DomainType      `json:"type"`
	Status                DomainStatus    `json:"status"`
	IsApex                bool            `json:"isApex"`
	VerificationToken     string          `json:"verificationToken,omitempty"`
	LastVerificationError string          `json:"lastVerificationError,omitempty"`
	LastCheckedAt         *time.Time      `json:"lastCheckedAt,omitempty"`
	ActivatedAt           *time.Time      `json:"activatedAt,omitempty"`
	CreatedAt             time.Time       `json:"createdAt"`
	UpdatedAt             time.Time       `json:"updatedAt"`
	Instructions          DNSInstructions `json:"instructions,omitempty"`
}

type DomainWithSite struct {
	CustomDomain
	SiteName   string `db:"site_name"`
	SiteDomain string `db:"site_domain"`
}
