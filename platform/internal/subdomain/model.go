package subdomain

type Subdomain struct {
	Id         string `db:"id" json:"id"`
	User       string `db:"user" json:"user"`
	Name       string `db:"name" json:"name"`
	FullDomain string `db:"full_domain" json:"fullDomain"`
}
