package auth

import (
	"hostamat/pkg/platforma/batteries/session"
	"hostamat/pkg/platforma/batteries/user"
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

		ctx := m.userService.SetCurrentUserToContext(r.Context(), user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
