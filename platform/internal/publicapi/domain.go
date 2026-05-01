package publicapi

import (
	"net/http"

	"hostamat/internal/apikey"

	"github.com/platforma-dev/platforma/httpserver"
)

type Domain struct {
	HandlerGroup *httpserver.HandlerGroup
}

func New(apiKeyService *apikey.Service, sites siteService, domains domainService) *Domain {
	errorWriter := ErrorWriter{}
	authMiddleware := apikey.NewBearerAuthMiddleware(apiKeyService, errorWriter)
	handler := NewHandler(sites, domains)

	group := httpserver.NewHandlerGroup()
	group.Use(authMiddleware)

	group.Handle("GET /me", http.HandlerFunc(handler.me))
	group.Handle("GET /sites", withScope(apikey.ScopeSitesRead, errorWriter, handler.listSites))
	group.Handle("POST /sites", withScope(apikey.ScopeSitesWrite, errorWriter, handler.createSite))
	group.Handle("GET /sites/{siteId}", withScope(apikey.ScopeSitesRead, errorWriter, handler.getSite))
	group.Handle("PATCH /sites/{siteId}", withScope(apikey.ScopeSitesWrite, errorWriter, handler.updateSite))
	group.Handle("DELETE /sites/{siteId}", withScope(apikey.ScopeSitesWrite, errorWriter, handler.deleteSite))
	group.Handle("POST /sites/{siteId}/upload", withScope(apikey.ScopeSitesUpload, errorWriter, handler.uploadSite))
	group.Handle("GET /domains", withScope(apikey.ScopeDomainsRead, errorWriter, handler.listDomains))
	group.Handle("POST /domains", withScope(apikey.ScopeDomainsWrite, errorWriter, handler.createDomain))
	group.Handle("GET /domains/{domainId}", withScope(apikey.ScopeDomainsRead, errorWriter, handler.getDomain))
	group.Handle("PATCH /domains/{domainId}", withScope(apikey.ScopeDomainsWrite, errorWriter, handler.updateDomain))
	group.Handle("POST /domains/{domainId}/verify", withScope(apikey.ScopeDomainsWrite, errorWriter, handler.verifyDomain))
	group.Handle("DELETE /domains/{domainId}", withScope(apikey.ScopeDomainsWrite, errorWriter, handler.deleteDomain))

	return &Domain{HandlerGroup: group}
}

func withScope(scope apikey.Scope, errorWriter ErrorWriter, handler func(http.ResponseWriter, *http.Request)) http.Handler {
	return apikey.NewScopeMiddleware(scope, errorWriter).Wrap(http.HandlerFunc(handler))
}
