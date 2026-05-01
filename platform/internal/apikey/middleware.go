package apikey

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/platforma-dev/platforma/auth"
	"github.com/platforma-dev/platforma/httpserver"
	"github.com/platforma-dev/platforma/log"
)

type publicErrorWriter interface {
	WritePublicError(w http.ResponseWriter, err error)
}

type BearerAuthMiddleware struct {
	service     *Service
	errorWriter publicErrorWriter
}

func NewBearerAuthMiddleware(service *Service, errorWriter publicErrorWriter) *BearerAuthMiddleware {
	return &BearerAuthMiddleware{
		service:     service,
		errorWriter: errorWriter,
	}
}

func (m *BearerAuthMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := bearerToken(r.Header.Get("Authorization"))
		if err != nil {
			m.errorWriter.WritePublicError(w, err)
			return
		}

		authenticated, err := m.service.Authenticate(r.Context(), token)
		if err != nil {
			m.errorWriter.WritePublicError(w, err)
			return
		}

		user := &auth.User{
			ID:       authenticated.Owner.ID,
			Username: authenticated.Owner.Username,
			Created:  authenticated.Owner.Created,
			Updated:  authenticated.Owner.Updated,
			Status:   authenticated.Owner.Status,
		}

		ctx := context.WithValue(r.Context(), log.UserIDKey, user.ID)
		ctx = context.WithValue(ctx, auth.UserContextKey, user)
		ctx = WithAuthenticatedKey(ctx, authenticated)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type ScopeMiddleware struct {
	scope       Scope
	errorWriter publicErrorWriter
}

func NewScopeMiddleware(scope Scope, errorWriter publicErrorWriter) httpserver.Middleware {
	return &ScopeMiddleware{scope: scope, errorWriter: errorWriter}
}

func (m *ScopeMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authenticated := AuthenticatedKeyFromContext(r.Context())
		if authenticated == nil || authenticated.Key == nil {
			m.errorWriter.WritePublicError(w, ErrMissingBearer)
			return
		}

		if !authenticated.Key.HasScope(m.scope) {
			m.errorWriter.WritePublicError(w, ErrInsufficientScope)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func bearerToken(header string) (string, error) {
	if strings.TrimSpace(header) == "" {
		return "", ErrMissingBearer
	}

	parts := strings.Fields(header)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", ErrMalformedBearer
	}

	if parts[1] == "" {
		return "", ErrMalformedBearer
	}

	return parts[1], nil
}

func PublicStatusForError(err error) int {
	switch {
	case errors.Is(err, ErrMissingBearer), errors.Is(err, ErrMalformedBearer), errors.Is(err, ErrUnknownKey), errors.Is(err, ErrExpiredKey), errors.Is(err, ErrRevokedKey), errors.Is(err, ErrOwnerInactive):
		return http.StatusUnauthorized
	case errors.Is(err, ErrInsufficientScope):
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
