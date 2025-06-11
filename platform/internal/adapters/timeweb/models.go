package timeweb

type Domain struct {
	Fqdn       string
	Subdomains []Domain
}
