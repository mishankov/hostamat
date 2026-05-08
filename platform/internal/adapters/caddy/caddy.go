package caddy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type CaddyAPI struct {
	httpClient  httpClient
	apiBasePath string
	serverName  string
}

func New(httpClient httpClient, apiBasePath, serverName string) *CaddyAPI {
	return &CaddyAPI{
		httpClient:  httpClient,
		apiBasePath: apiBasePath,
		serverName:  serverName,
	}
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

	switch {
	case statusOk(resp.StatusCode):
		return nil
	case resp.StatusCode == 404:
		return ErrConfigNotFound
	case statusBadRequest(resp.StatusCode):
		return ErrInvalidRequest
	case statusServerError(resp.StatusCode):
		return ErrServerError
	default:
		return ErrUnexpectedResponse
	}
}

func (capi *CaddyAPI) CreateBaseConfig(ctx context.Context, target, handlerPrefix, shorthandId string) error {
	reqBody := Route{
		Handle: []Handle{{
			Handler: "subroute",
			Routes: []Route{
				{
					Match: []Match{{
						Path: []string{"/forms/submit/*"},
					}},
					Handle: []Handle{{
						Handler: "reverse_proxy",
						Upstreams: []Upstream{{
							Dial: target,
						}},
					}},
				},
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
		return fmt.Errorf("failed to encode request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST",
		capi.apiBasePath+"/config/apps/http/servers/"+capi.serverName+"/routes",
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

	switch {
	case statusOk(resp.StatusCode):
		return nil
	case statusBadRequest(resp.StatusCode):
		return ErrInvalidRequest
	case statusServerError(resp.StatusCode):
		return ErrServerError
	default:
		return ErrUnexpectedResponse
	}
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

	switch {
	case statusOk(resp.StatusCode):
		return nil
	case statusBadRequest(resp.StatusCode):
		return ErrInvalidRequest
	case statusServerError(resp.StatusCode):
		return ErrServerError
	default:
		return ErrUnexpectedResponse
	}
}

func (capi *CaddyAPI) GetHostsFromShorthandId(ctx context.Context, shorthandId string) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", capi.apiBasePath+"/id/"+shorthandId+"/host", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := capi.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch {
	case statusOk(resp.StatusCode):
		var hosts []string
		if err := json.NewDecoder(resp.Body).Decode(&hosts); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}
		return hosts, nil
	case resp.StatusCode == 404:
		return nil, ErrConfigNotFound
	case statusBadRequest(resp.StatusCode):
		return nil, ErrInvalidRequest
	case statusServerError(resp.StatusCode):
		return nil, ErrServerError
	default:
		return nil, ErrUnexpectedResponse
	}
}

func (capi *CaddyAPI) RewriteHostsInShorthandId(ctx context.Context, shorthandId string, hosts []string) error {
	if len(hosts) == 0 {
		return capi.ClearSubdomains(ctx, shorthandId)
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(hosts); err != nil {
		return fmt.Errorf("failed to encode request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH",
		capi.apiBasePath+"/id/"+shorthandId+"/host",
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

	switch {
	case statusOk(resp.StatusCode):
		return nil
	case resp.StatusCode == 404:
		return ErrConfigNotFound
	case statusBadRequest(resp.StatusCode):
		return ErrInvalidRequest
	case statusServerError(resp.StatusCode):
		return ErrServerError
	default:
		return ErrUnexpectedResponse
	}
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

	switch {
	case statusBadRequest(respDelete.StatusCode):
		return ErrInvalidRequest
	case statusServerError(respDelete.StatusCode):
		return ErrServerError
	}

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

	switch {
	case statusOk(respPost.StatusCode):
		return nil
	case statusBadRequest(respPost.StatusCode):
		return ErrInvalidRequest
	case statusServerError(respPost.StatusCode):
		return ErrServerError
	default:
		return ErrUnexpectedResponse
	}
}
