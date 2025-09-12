package auth

import (
	"context"
	"hostamat/pkg/platforma/batteries/session"
	"hostamat/pkg/platforma/batteries/user"
	"hostamat/pkg/platforma/log"
	"net/http"
)

type AuthenticationMiddleware struct {
	sessionService *session.Service
	userService    *user.Service
}

func NewAuthenticationMiddleware(sessionService *session.Service, userService *user.Service) *AuthenticationMiddleware {
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
