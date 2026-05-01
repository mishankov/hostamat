package apikey

import "errors"

var (
	ErrMissingName       = errors.New("name is required")
	ErrInvalidBody       = errors.New("invalid request body")
	ErrNameTooLong       = errors.New("name must be at most 80 characters")
	ErrMissingScopes     = errors.New("at least one scope is required")
	ErrInvalidScope      = errors.New("invalid scope")
	ErrInvalidExpiration = errors.New("expiration must be in the future")
	ErrActiveKeyLimit    = errors.New("active API key limit reached")
	ErrMissingBearer     = errors.New("missing bearer token")
	ErrMalformedBearer   = errors.New("malformed bearer token")
	ErrUnknownKey        = errors.New("API key not found")
	ErrExpiredKey        = errors.New("API key expired")
	ErrRevokedKey        = errors.New("API key revoked")
	ErrOwnerInactive     = errors.New("API key owner is inactive")
	ErrInsufficientScope = errors.New("insufficient API key scope")
)

type ErrorCode string

const (
	ErrorCodeInvalidBody       ErrorCode = "invalid_body"
	ErrorCodeMissingName       ErrorCode = "missing_name"
	ErrorCodeNameTooLong       ErrorCode = "name_too_long"
	ErrorCodeMissingScopes     ErrorCode = "missing_scopes"
	ErrorCodeInvalidScope      ErrorCode = "invalid_scope"
	ErrorCodeInvalidExpiration ErrorCode = "invalid_expiration"
	ErrorCodeActiveKeyLimit    ErrorCode = "active_key_limit_reached"
	ErrorCodeMissingBearer     ErrorCode = "missing_bearer_token"
	ErrorCodeMalformedBearer   ErrorCode = "malformed_bearer_token"
	ErrorCodeUnknownKey        ErrorCode = "unknown_api_key"
	ErrorCodeExpiredKey        ErrorCode = "expired_api_key"
	ErrorCodeRevokedKey        ErrorCode = "revoked_api_key"
	ErrorCodeOwnerInactive     ErrorCode = "api_key_owner_inactive"
	ErrorCodeInsufficientScope ErrorCode = "insufficient_scope"
	ErrorCodeNotFound          ErrorCode = "not_found"
	ErrorCodeInternal          ErrorCode = "internal_error"
)

func CodeForError(err error) ErrorCode {
	switch {
	case errors.Is(err, ErrMissingName):
		return ErrorCodeMissingName
	case errors.Is(err, ErrInvalidBody):
		return ErrorCodeInvalidBody
	case errors.Is(err, ErrNameTooLong):
		return ErrorCodeNameTooLong
	case errors.Is(err, ErrMissingScopes):
		return ErrorCodeMissingScopes
	case errors.Is(err, ErrInvalidScope):
		return ErrorCodeInvalidScope
	case errors.Is(err, ErrInvalidExpiration):
		return ErrorCodeInvalidExpiration
	case errors.Is(err, ErrActiveKeyLimit):
		return ErrorCodeActiveKeyLimit
	case errors.Is(err, ErrMissingBearer):
		return ErrorCodeMissingBearer
	case errors.Is(err, ErrMalformedBearer):
		return ErrorCodeMalformedBearer
	case errors.Is(err, ErrUnknownKey):
		return ErrorCodeUnknownKey
	case errors.Is(err, ErrExpiredKey):
		return ErrorCodeExpiredKey
	case errors.Is(err, ErrRevokedKey):
		return ErrorCodeRevokedKey
	case errors.Is(err, ErrOwnerInactive):
		return ErrorCodeOwnerInactive
	case errors.Is(err, ErrInsufficientScope):
		return ErrorCodeInsufficientScope
	default:
		return ErrorCodeInternal
	}
}
