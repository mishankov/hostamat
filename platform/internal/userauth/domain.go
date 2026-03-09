package userauth

import (
	"hostamat/internal/adapters/yandexoauth"

	platformauth "github.com/platforma-dev/platforma/auth"
	"github.com/platforma-dev/platforma/httpserver"
)

type Domain struct {
	Repository  *Repository
	Service     *Service
	HandleGroup *httpserver.HandlerGroup
	Middleware  httpserver.Middleware
}

func New(db db, authStorage authStorage, yandexClient *yandexoauth.Client, sessionCookieName, redirectURI string, cleanupEnqueuer cleanupEnqueuer) *Domain {
	repository := NewRepository(db)
	service := NewService(repository, authStorage, yandexClient, sessionCookieName, redirectURI, cleanupEnqueuer)

	authMiddleware := platformauth.NewAuthenticationMiddleware(service)
	authAPI := httpserver.NewHandlerGroup()
	authAPI.Handle("GET /yandex/start", NewYandexStartHandler(service))
	authAPI.Handle("GET /yandex/callback", NewYandexCallbackHandler(service))
	authAPI.Handle("POST /logout", NewLogoutHandler(service))
	authAPI.Handle("GET /me", NewGetHandler(service))
	authAPI.Handle("DELETE /me", authMiddleware.Wrap(NewDeleteHandler(service)))

	return &Domain{
		Repository:  repository,
		Service:     service,
		HandleGroup: authAPI,
		Middleware:  authMiddleware,
	}
}

func (d *Domain) GetRepository() any {
	return d.Repository
}
