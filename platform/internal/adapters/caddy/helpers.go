package caddy

import "errors"

var (
	ErrConfigNotFound     = errors.New("caddy config not found")
	ErrInvalidRequest     = errors.New("invalid request to caddy")
	ErrServerError        = errors.New("caddy server error")
	ErrUnexpectedResponse = errors.New("unexpected response from caddy")
)

func statusOk(status int) bool {
	return status >= 200 && status < 300
}

func statusBadRequest(status int) bool {
	return status >= 400 && status < 500
}

func statusServerError(status int) bool {
	return status >= 500
}
