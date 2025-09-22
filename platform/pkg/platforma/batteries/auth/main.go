package auth

import (
	"context"
	"hostamat/pkg/platforma/batteries/user"
)

type contextKey string

const (
	UserContextKey contextKey = "user"
)

func UserFromContext(ctx context.Context) *user.User {
	user, ok := ctx.Value(UserContextKey).(*user.User)
	if !ok {
		return nil
	}
	return user
}
