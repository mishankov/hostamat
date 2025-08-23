package subdomain

type Subdomain struct {
	Id         string `db:"id" json:"id"`
	Name       string `db:"name" json:"name"`
	FullDomain string `db:"full_domain" json:"fullDomain"`
}
