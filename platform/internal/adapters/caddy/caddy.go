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
}

func New(httpClient httpClient, apiBasePath string) *CaddyAPI {
	return &CaddyAPI{httpClient: httpClient, apiBasePath: apiBasePath}
}

func (capi *CaddyAPI) AddSubdomain(ctx context.Context, baseDomain, subdomain, target string, handlerPrefix string) error {
	// TODO: use shorthand syntax
	reqBody := Route{
		Handle: []Handle{{
			Handler: "subroute",
			Routes: []Route{
				{
					Group: "group2",
					Handle: []Handle{{
						Handler: "rewrite",
						Uri:     handlerPrefix + "{http.request.uri}",
					}},
				},
				{
					Handle: []Handle{{
						Handler: "reverse_proxy",
						Upstreams: []Upstream{{
							Dial: target,
						}},
					}},
				}},
		}},
		Match: []Match{{
			// TODO: only add hosts to matcher. routes can stay the same and be updated on start or when needed
			Host: []string{subdomain + "." + baseDomain},
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
