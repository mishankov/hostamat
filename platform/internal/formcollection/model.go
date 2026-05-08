package formcollection

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

const (
	MaxSubmissionBodyBytes   = 64 * 1024
	MaxSubmissionFields      = 100
	MaxSubmissionValueLength = 10000
)

type UploadedHTMLFile struct {
	FileName string
	PagePath string
	Body     []byte
}

type DiscoverySummary struct {
	DiscoveryID      string           `json:"discoveryId,omitempty"`
	TotalForms       int              `json:"totalForms"`
	SupportedForms   int              `json:"supportedForms"`
	UnsupportedForms int              `json:"unsupportedForms"`
	Forms            []DiscoveredForm `json:"forms"`
}

type DiscoveredForm struct {
	ID                string    `db:"id" json:"id"`
	DiscoveryID       string    `db:"discovery_id" json:"-"`
	SiteID            string    `db:"site" json:"-"`
	FileName          string    `db:"file_name" json:"fileName"`
	PagePath          string    `db:"page_path" json:"pagePath"`
	FormKey           string    `db:"form_key" json:"formKey"`
	Label             string    `db:"label" json:"label"`
	OriginalAction    string    `db:"original_action" json:"originalAction"`
	Method            string    `db:"method" json:"method"`
	Fields            FieldList `db:"fields" json:"fields"`
	Supported         bool      `db:"supported" json:"supported"`
	UnsupportedReason string    `db:"unsupported_reason" json:"unsupportedReason,omitempty"`
	Active            bool      `db:"active" json:"active"`
	PublicToken       string    `db:"public_token" json:"-"`
	CreatedAt         time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt         time.Time `db:"updated_at" json:"updatedAt"`
}

type Field struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type FieldList []Field

func (f FieldList) Value() (driver.Value, error) {
	return json.Marshal(f)
}

func (f *FieldList) Scan(value any) error {
	if value == nil {
		*f = FieldList{}
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, f)
	case string:
		return json.Unmarshal([]byte(v), f)
	default:
		*f = FieldList{}
		return nil
	}
}

type Discovery struct {
	ID        string    `db:"id"`
	SiteID    string    `db:"site"`
	CreatedAt time.Time `db:"created_at"`
}

type Submission struct {
	ID        string          `db:"id" json:"id"`
	SiteID    string          `db:"site" json:"siteId"`
	FormID    string          `db:"form" json:"formId"`
	PagePath  string          `db:"page_path" json:"pagePath"`
	Fields    json.RawMessage `db:"fields" json:"fields"`
	CreatedAt time.Time       `db:"created_at" json:"createdAt"`

	SiteName   string `db:"site_name" json:"siteName,omitempty"`
	FormLabel  string `db:"form_label" json:"formLabel,omitempty"`
	Active     bool   `db:"active" json:"active"`
	Historical bool   `db:"historical" json:"historical"`
}

type SubmissionFilters struct {
	SiteID string
	FormID string
}
