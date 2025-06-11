package timeweb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type TimewebAPI struct {
	httpClient  httpClient
	apiBasePath string
	apiToken    string
}

func New(httpClient httpClient, apiBasePath string, apiToken string) *TimewebAPI {
	return &TimewebAPI{httpClient: httpClient, apiBasePath: apiBasePath, apiToken: apiToken}
}

func (t *TimewebAPI) AddSubdomain(ctx context.Context, domain, subdomain string) error {
	req, err := http.NewRequestWithContext(ctx, "POST",
		fmt.Sprintf(t.apiBasePath+"/domains/%s/subdomains/%s", domain, subdomain), nil,
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+t.apiToken)
	resp, err := t.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (t *TimewebAPI) DeleteSubdomain(ctx context.Context, domain, subdomain string) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE",
		fmt.Sprintf(t.apiBasePath+"/domains/%s/subdomains/%s", domain, subdomain), nil,
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+t.apiToken)
	resp, err := t.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (t *TimewebAPI) GetDomain(ctx context.Context, domain string) (Domain, error) {
	req, err := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf(t.apiBasePath+"/domains/%s", domain), nil,
	)
	if err != nil {
		return Domain{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+t.apiToken)
	resp, err := t.httpClient.Do(req)
	if err != nil {
		return Domain{}, err
	}
	defer resp.Body.Close()

	var respParsed struct{ Domain Domain }
	err = json.NewDecoder(resp.Body).Decode(&respParsed)

	return respParsed.Domain, err
}
