package httpclient

import (
	"context"
	"net/http"
	"strings"
)

type logger interface {
	DebugContext(ctx context.Context, msg string, args ...any)
}

type HttpClient struct {
	client *http.Client
	logger logger
}

func New(logger logger) *HttpClient {
	// TODO: probably need better http client
	return &HttpClient{http.DefaultClient, logger}
}

func (c *HttpClient) Do(req *http.Request) (*http.Response, error) {
	// TODO: mask or do not log Authorization header
	c.logger.DebugContext(req.Context(), "request started", "url", req.URL, "headers", maskedHeaders(req.Header))
	resp, err := c.client.Do(req)
	c.logger.DebugContext(req.Context(), "request made", "status", resp.Status, "headers", maskedHeaders(resp.Header))

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
