package yandexoauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	httpClient   httpClient
	authorizeURL string
	tokenURL     string
	userInfoURL  string
	clientID     string
	clientSecret string
	redirectURI  string
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

type UserInfo struct {
	ID           string `json:"id"`
	Login        string `json:"login"`
	DefaultEmail string `json:"default_email"`
}

func New(httpClient httpClient, authorizeURL, tokenURL, userInfoURL, clientID, clientSecret, redirectURI string) *Client {
	return &Client{
		httpClient:   httpClient,
		authorizeURL: authorizeURL,
		tokenURL:     tokenURL,
		userInfoURL:  userInfoURL,
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
	}
}

func (c *Client) AuthorizeURL(state string) string {
	values := url.Values{}
	values.Set("response_type", "code")
	values.Set("client_id", c.clientID)
	values.Set("redirect_uri", c.redirectURI)
	values.Set("state", state)

	return c.authorizeURL + "?" + values.Encode()
}

func (c *Client) ExchangeCode(ctx context.Context, code string) (string, error) {
	values := url.Values{}
	values.Set("grant_type", "authorization_code")
	values.Set("code", code)
	values.Set("client_id", c.clientID)
	values.Set("client_secret", c.clientSecret)
	values.Set("redirect_uri", c.redirectURI)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.tokenURL, strings.NewReader(values.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to exchange oauth code: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("yandex token request failed: %s", strings.TrimSpace(string(body)))
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	if tokenResp.AccessToken == "" {
		return "", fmt.Errorf("yandex token response is missing access token")
	}

	return tokenResp.AccessToken, nil
}

func (c *Client) GetUserInfo(ctx context.Context, accessToken string) (*UserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.userInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create userinfo request: %w", err)
	}

	req.Header.Set("Authorization", "OAuth "+accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch yandex user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("yandex user info request failed: %s", strings.TrimSpace(string(body)))
	}

	var userInfo UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info response: %w", err)
	}

	if userInfo.ID == "" {
		return nil, fmt.Errorf("yandex user info response is missing id")
	}

	return &userInfo, nil
}
