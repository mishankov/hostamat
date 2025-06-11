package caddy

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type CaddyAPI struct {
	httpClient  httpClient
	apiBasePath string
	serverHost  string
	baseDomain  string
}

func New(httpClient httpClient, apiBasePath string, serverHost string, baseDomain string) *CaddyAPI {
	return &CaddyAPI{httpClient: httpClient, apiBasePath: apiBasePath, serverHost: serverHost, baseDomain: baseDomain}
}

func (capi *CaddyAPI) AddSubdomain(ctx context.Context, subdomain string) error {
	reqBody := Route{
		Handle: []Handle{{
			Handler: "subroute",
			Routes: []Route{{
				Handle: []Handle{{
					Handler: "reverse_proxy",
					Upstreams: []Upstream{{
						Dial: capi.serverHost,
					}},
				}},
			}},
		}},
		Match: []Match{{
			Host: []string{subdomain + "." + capi.baseDomain},
		}},
		Terminal: true,
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(reqBody)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST",
		capi.apiBasePath+"/apps/http/servers/srv0/routes", // TODO: remove srv0 harcode
		&buf,
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := capi.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
