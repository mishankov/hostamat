package httpclient

import (
	"hostamat/pkg/platforma/log"
	"net/http"
	"strings"
)

type Client struct {
	client *http.Client
}

func New() *Client {
	// TODO: probably need better http client
	return &Client{http.DefaultClient}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	log.DebugContext(req.Context(), "request started", "url", req.URL, "headers", maskedHeaders(req.Header))
	resp, err := c.client.Do(req)
	log.DebugContext(req.Context(), "request made", "status", resp.Status, "headers", maskedHeaders(resp.Header))

	return resp, err
}

func maskedHeaders(headers http.Header) http.Header {
	newHeaders := http.Header{}
	for hn, hvs := range headers {
		if strings.ToLower(hn) == "authorization" {
			for _, hv := range hvs {
				if strings.HasPrefix(strings.ToLower(hv), "basic ") || strings.HasPrefix(strings.ToLower(hv), "bearer ") {
					newHeaders.Add(hn, strings.Split(hv, " ")[0]+" ***")
					continue
				}
				newHeaders.Add(hn, hv)
			}
		} else {
			for _, hv := range hvs {
				newHeaders.Add(hn, hv)
			}
		}
	}

	return newHeaders
}
