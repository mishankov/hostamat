package caddy

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
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

func (capi *CaddyAPI) CheckConfigWithId(ctx context.Context, shorthandId string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", capi.apiBasePath+"/id/"+shorthandId, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := capi.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return errors.New("config not found")
	}

	return nil
}

func (capi *CaddyAPI) CreateBaseConfig(ctx context.Context, target, handlerPrefix, shorthandId string) error {
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
			Id:   shorthandId,
			Host: []string{},
		}},
		Terminal: true,
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(reqBody)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST",
		capi.apiBasePath+"/config/apps/http/servers/srv0/routes", // TODO: remove srv0 harcode
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

func (capi *CaddyAPI) AddSubdomain(ctx context.Context, baseDomain, subdomain, shorthandId string) error {

	req, err := http.NewRequestWithContext(ctx, "POST",
		capi.apiBasePath+"/id/"+shorthandId+"/host",
		strings.NewReader(`"`+subdomain+"."+baseDomain+`"`),
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

func (capi *CaddyAPI) ClearSubdomains(ctx context.Context, shorthandId string) error {
	reqDelete, err := http.NewRequestWithContext(ctx, "DELETE",
		capi.apiBasePath+"/id/"+shorthandId+"/host",
		nil,
	)
	if err != nil {
		return err
	}

	reqDelete.Header.Set("Content-Type", "application/json")

	respDelete, err := capi.httpClient.Do(reqDelete)
	if err != nil {
		return err
	}
	defer respDelete.Body.Close()

	reqPost, err := http.NewRequestWithContext(ctx, "POST",
		capi.apiBasePath+"/id/"+shorthandId+"/host",
		strings.NewReader(`[]`),
	)
	if err != nil {
		return err
	}

	reqPost.Header.Set("Content-Type", "application/json")

	respPost, err := capi.httpClient.Do(reqPost)
	if err != nil {
		return err
	}
	defer respPost.Body.Close()

	return nil
}
