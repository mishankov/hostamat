package auth

import (
	"context"
	"net/http"

	"hostamat/pkg/platforma/batteries/session"
	"hostamat/pkg/platforma/batteries/user"
	"hostamat/pkg/platforma/log"
)

type sessionService interface {
	Get(ctx context.Context, id string) (*session.Session, error)
}

type userService interface {
	Get(ctx context.Context, id string) (*user.User, error)
}

type AuthenticationMiddleware struct {
	sessionService sessionService
	userService    userService
}

func NewAuthenticationMiddleware(sessionService sessionService, userService userService) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{sessionService: sessionService, userService: userService}
}

func (m *AuthenticationMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		session, err := m.sessionService.Get(r.Context(), cookie.Value)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		user, err := m.userService.Get(r.Context(), session.User)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctxWithUser := context.WithValue(r.Context(), log.UserIdKey, user.ID)
		requestWithUser := r.WithContext(ctxWithUser)

		next.ServeHTTP(w, requestWithUser)
	})
}
