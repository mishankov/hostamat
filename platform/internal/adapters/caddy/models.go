package caddy

type Route struct {
	Handle   []Handle `json:"handle,omitempty"`
	Match    []Match  `json:"match,omitempty"`
	Terminal bool     `json:"terminal,omitempty"`
	Group    string   `json:"group,omitempty"`
}

type Handle struct {
	Handler   string     `json:"handler,omitempty"`
	Uri       string     `json:"uri,omitempty"`
	Routes    []Route    `json:"routes,omitempty"`
	Upstreams []Upstream `json:"upstreams,omitempty"`
}

type Upstream struct {
	Dial string `json:"dial,omitempty"`
}

type Match struct {
	Host []string `json:"host,omitempty"`
}
