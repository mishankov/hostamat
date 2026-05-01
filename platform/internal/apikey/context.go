package apikey

import "context"

type contextKey string

const authenticatedKeyContextKey contextKey = "authenticated_api_key"

func WithAuthenticatedKey(ctx context.Context, key *AuthenticatedKey) context.Context {
	return context.WithValue(ctx, authenticatedKeyContextKey, key)
}

func AuthenticatedKeyFromContext(ctx context.Context) *AuthenticatedKey {
	key, ok := ctx.Value(authenticatedKeyContextKey).(*AuthenticatedKey)
	if !ok {
		return nil
	}

	return key
}
