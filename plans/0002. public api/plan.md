# Plan: Public API for Sites, Domains, and Uploads

> Source PRD: `plans/0002. public api/prd.md`

## Architectural decisions

Durable decisions that apply across all phases:

- **Routes**:
  Public API routes are served from the API subdomain and versioned from the start under `https://api.hostamat.ru/v1/...`.
  The backend does not enforce the request host; infrastructure and documentation define the public API host.
  HTTPS for public API traffic is an infrastructure concern; docs advertise HTTPS and backend-level HTTPS enforcement is not required.
  DNS/proxy routing for `api.hostamat.ru` is manual deployment work and is not included in this implementation plan.
  Public API includes `GET /v1/me` as a stable authentication test endpoint returning minimal non-sensitive user identity plus key metadata; it requires a valid API key but no resource scope.
  Public API uses `PATCH` for partial resource updates.
  Public site routes are `GET /v1/sites`, `POST /v1/sites`, `GET /v1/sites/{siteId}`, `PATCH /v1/sites/{siteId}`, `DELETE /v1/sites/{siteId}`, and `POST /v1/sites/{siteId}/upload`.
  Site zip upload is modeled as `POST /v1/sites/{siteId}/upload`.
  Public custom-domain routes are `GET /v1/domains`, `POST /v1/domains`, `GET /v1/domains/{domainId}`, `PATCH /v1/domains/{domainId}`, `POST /v1/domains/{domainId}/verify`, and `DELETE /v1/domains/{domainId}`.
  Custom domain verification is modeled as an action endpoint: `POST /v1/domains/{domainId}/verify`.
  Dashboard session routes continue to exist for the webapp and are not treated as the public contract.
  Existing dashboard/private API routes do not need to be renamed as part of this feature.
  API key management is a dashboard-authenticated feature, while public API resource routes authenticate with bearer API keys.
  API key management UI lives in the dashboard navigation as its own page, for example `/dashboard/api-keys`.
  Hosted documentation is available publicly inside the main webapp, with an interactive Scalar OpenAPI page and written tutorial pages linked from the product, public navigation, and API keys dashboard page.
- **Authentication**:
  Public API clients authenticate with `Authorization: Bearer <api-key>`.
  API key secrets use the format `hmat_<id>_<secret>`, where `<id>` is a UUID that supports lookup/display and the secret portion is verified securely.
  API key secret verification uses a high-entropy random secret stored as a peppered HMAC-SHA-256-style digest rather than plaintext.
  The API key digest pepper comes from a stable backend environment variable such as `HOSTAMAT_API_KEY_PEPPER`.
  API key authentication resolves to the owning Hostamat user so existing service-level ownership and quota rules remain authoritative.
  API key authentication fails if the owning user cannot be resolved as an active user.
  API key secrets are shown once at creation time, never stored in plaintext, and never returned by list/detail endpoints.
  API key UI shows a non-secret key identifier/prefix so users can match dashboard records to configured integrations.
  Expired, revoked, deleted, malformed, and missing keys produce distinct public API error codes.
- **Authorization and scopes**:
  Keys carry explicit scopes: `sites:read`, `sites:write`, `sites:upload`, `domains:read`, and `domains:write`.
  Scopes are literal and non-hierarchical; no scope implies another scope.
  `sites:write` covers creating, updating the site name, and deleting sites.
  `sites:upload` allows uploading a zip archive to any site owned by the key owner, but does not allow creating, updating, or deleting sites.
  Public zip upload preserves the current archive path handling and uploaded/skipped-file behavior used by the existing file service, and its response contains only upload summary/file-list data rather than the site object.
  Public zip upload is synchronous in v1 and does not expose upload status/history endpoints.
  `domains:write` covers creating, updating, verifying, moving, and deleting custom domains.
  Scope enforcement happens at the public API boundary before resource operations run.
  User-level limits remain attached to the key owner, not to the key.
- **Schema**:
  Add persistent API key records containing owner, display name, key prefix or public identifier, secret digest, scopes, expiration timestamp, last-used timestamp, creation timestamp, update timestamp, and revocation/deletion state.
  Public API resource responses use documented external DTOs rather than accidental internal handler response shapes.
  Public API handlers use a thin facade over existing product services to preserve current behavior while mapping external DTOs and errors deliberately.
  A dedicated public API module owns the `/v1` route group and depends on existing product services plus API key middleware.
  API key lifecycle, verification, scopes, bearer authentication middleware, and scope middleware live in their own backend domain/module rather than being folded into browser auth.
  Product behavior, validation, limits, and default ordering mirror the current dashboard/service behavior unless this plan explicitly says otherwise.
- **Key models**:
  API key is a credential owned by one user.
  API scope is an explicit permission string checked by public API handlers.
  API key names are required labels, do not need to be unique per user, and have an 80-character maximum length.
  API keys may have no expiration date.
  API key expiration is stored as a timestamp and evaluated with backend server time; date-only UI choices are converted to an end-of-day timestamp before storage.
  API key metadata represents no-expiration keys with `expiresAt: null`.
  API key deletion is implemented as soft revocation; revoked keys are hidden from the UI and rejected by the public API.
  Only the API key name can be edited after creation.
  API key creation offers expiration presets including a visible "never expires" option.
  API key creation shows the generated secret in a one-time modal with a copy action and an explicit acknowledgement that it cannot be viewed again.
  Users can have at most 20 active API keys; expired and revoked keys do not count toward the active-key limit.
  Public API error is a stable JSON envelope with a machine-readable code and English human-readable message: `{ "error": { "code": "...", "message": "..." } }`.
  Request IDs are exposed through the existing `Hostamat-Trace-Id` response header only, including error responses.
  Public API JSON fields use `camelCase`.
  Public API timestamps are RFC 3339 strings.
  Public list endpoints do not add pagination, filtering, or sorting in v1.
  OpenAPI spec is a hand-maintained committed contract artifact at `docs/api/openapi.yaml` that reflects public API behavior.
  Each public API implementation phase updates `docs/api/openapi.yaml` for the endpoints it introduces; the hosted documentation phase publishes and validates the completed reference.
  Written tutorial documentation is Russian-only in v1.
  OpenAPI title, summaries, and descriptions are Russian in v1, while schema fields and error codes remain stable code-oriented identifiers.
  OpenAPI documents only the public bearer-token API, not dashboard/session API key management endpoints.
  Tutorial examples use `HOSTAMAT_API_BASE_URL`, documented with default value `https://api.hostamat.ru/v1`.
- **Security and observability**:
  Authorization headers and generated secrets stay redacted in logs.
  Last-used tracking records only `lastUsedAt` for successful authentication, without storing request payloads, IP addresses, or user agents, and updates are throttled to avoid a database write on every API request.
  Public API does not enable CORS in v1.
  Public API does not add app-level rate limiting in v1.
  Public API and API key creation are enabled immediately for authenticated users; no feature flag is required for v1.
  Public API behavior is tested at service and HTTP boundaries, with emphasis on externally visible contracts.

---

## Phase 1: API Key Lifecycle and Authenticated Public Entry Point

**User stories**: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 30, 32, 37, 38, 39

### What to build

Create the first complete path for public API authentication. A dashboard user can create a named, scoped, expiring API key, copy the secret once, see non-secret key metadata in the dashboard, delete or revoke the key, and use the bearer token against `GET /v1/me` to prove authentication, scope parsing, expiration checks, revocation checks, and last-used tracking work end to end.

This phase establishes the API key domain, persistence, dashboard management UI, public API middleware, shared public error envelope, and baseline automated tests before exposing real site or domain mutations.

### Acceptance criteria

- [ ] A logged-in dashboard user can create an API key with name, expiration, and one or more scopes.
- [ ] A logged-in dashboard user can create an API key with no expiration date.
- [ ] The full API key secret is displayed once immediately after creation and is not returned by later API responses.
- [ ] The one-time secret view includes a copy action and requires acknowledgement that the secret cannot be viewed again.
- [ ] API key management is available as its own dashboard navigation page.
- [ ] API key creation includes expiration presets with a visible "never expires" option.
- [ ] API key creation uses checkboxes plus presets for the five supported scopes, and requires at least one scope.
- [ ] The dashboard lists API keys with name, scopes, creation date, expiration date, last-used date, and active/expired/revoked state.
- [ ] The dashboard shows each active/expired key's non-secret identifier or prefix.
- [ ] Expired keys remain visible until deleted; revoked keys are hidden from the UI.
- [ ] Only the API key name can be edited after creation.
- [ ] A user can delete or revoke an API key from the dashboard and the key stops authenticating immediately.
- [ ] A user cannot create more than 20 active API keys.
- [ ] Expired and revoked keys do not count toward the active-key limit.
- [ ] Public API authentication accepts valid bearer keys and rejects missing, malformed, expired, revoked, and unknown keys with stable JSON errors.
- [ ] Public API authentication rejects keys whose owning user is missing or inactive.
- [ ] API key verification uses a stable configured pepper so keys survive backend restarts and deploys.
- [ ] `GET /v1/me` returns minimal non-sensitive user identity plus key metadata, including key ID, key name, scopes, and expiration, and does not require a resource scope.
- [ ] Successful public API authentication updates last-used tracking without logging the key secret, with throttling to avoid a database write on every request.
- [ ] `docs/api/openapi.yaml` documents the public authentication model and `GET /v1/me`.
- [ ] Backend tests cover key creation, secret verification, expiration, revocation, missing scopes, and last-used updates.

---

## Phase 2: Public Sites API

**User stories**: 12, 13, 14, 15, 16, 17, 30, 32, 37, 40

### What to build

Expose site listing and site lifecycle operations through the public API using scoped bearer-key authentication. The slice should let external clients list sites, create a site, rename a site, and delete a site while preserving the same ownership, validation, generated-domain, and account-limit behavior as the dashboard.

This phase proves the public API can wrap an existing product domain without bypassing service rules or leaking dashboard-specific response behavior.

### Acceptance criteria

- [ ] A key with site read scope can list only the owner user's sites.
- [ ] A key with site read scope can fetch a single owned site by ID.
- [ ] A key with site write scope can create, rename, and delete only the owner user's sites.
- [ ] Public site update is limited to the site name in v1.
- [ ] Because scopes are non-hierarchical, keys need `sites:read` for site reads and `sites:write` for site mutations.
- [ ] Site upload is not covered by `sites:write` and requires `sites:upload`.
- [ ] Read-only site keys cannot create, rename, delete, upload files, or manage domains.
- [ ] Site creation through the public API respects the existing per-user site limit.
- [ ] Site validation errors, not-found errors, ownership failures, and scope failures use the public API error envelope.
- [ ] Site responses include stable identifiers, site names, generated Hostamat domains, and timestamps suitable for external clients.
- [ ] `docs/api/openapi.yaml` documents public site endpoints.
- [ ] Backend HTTP tests cover successful site operations, insufficient scopes, cross-user access prevention, and limit enforcement.

---

## Phase 3: Public Site Zip Upload API

**User stories**: 25, 26, 27, 28, 29, 31, 32, 37, 40

### What to build

Expose a public upload path for static files. A scoped key can upload one raw zip archive to a site it owns, and the upload replaces the site's existing files. The behavior should reuse the existing archive extraction, file validation, file-size limits, and site ownership checks where possible so API uploads match dashboard constraints without exposing dashboard upload sessions publicly.

This phase makes the public API useful for CI/CD while keeping upload permission separate from broader site administration.

### Acceptance criteria

- [ ] A key with `sites:upload` can upload a raw `application/zip` request body to a site owned by the key owner.
- [ ] The upload endpoint accepts zip archives only.
- [ ] Zip archive paths and partial uploaded/skipped file results are interpreted using the current file-service archive logic.
- [ ] Successful upload replaces all existing files for the target site.
- [ ] The response includes only a summary plus uploaded file list, not the site object.
- [ ] Upload completes synchronously in the request and v1 does not include upload status/history endpoints.
- [ ] Upload-only keys cannot read sites, create sites, delete sites, or manage domains unless they also have those explicit scopes.
- [ ] Upload validation, allowed file behavior, total size limits, missing site handling, and ownership checks match dashboard behavior.
- [ ] Upload failures return stable public API error codes and useful messages for CI logs.
- [ ] `docs/api/openapi.yaml` documents public zip upload.
- [ ] Backend tests cover raw zip upload, replacement behavior, validation failures, limit failures, insufficient scopes, and cross-user access prevention.

---

## Phase 4: Public Domains API

**User stories**: 18, 19, 20, 21, 22, 23, 24, 30, 32, 37, 40

### What to build

Expose custom domain management through the public API. A scoped key can list custom domains, fetch a custom domain by ID, create a custom domain for an owned site, verify a custom domain, move a custom domain to another owned site, and delete a custom domain. Responses include domain status, connected site information, DNS instructions, verification timestamps, and actionable verification errors.

This phase completes the resource-management surface promised by the PRD while preserving existing custom-domain lifecycle rules.

### Acceptance criteria

- [ ] A key with domain read scope can list only custom domains owned by the key owner.
- [ ] A key with domain read scope can fetch a single owned custom domain by ID.
- [ ] Generated Hostamat domains are exposed on site responses, not through the domains API.
- [ ] Domain listing is global for the user; v1 does not include per-site domain routes.
- [ ] A key with domain write scope can create custom domains only for sites owned by the key owner, and custom domain creation requires `siteId`.
- [ ] A key with domain write scope can verify, update, and delete owned custom domains.
- [ ] Domain operations preserve existing hostname normalization, apex rejection, uniqueness, DNS verification, and activation behavior.
- [ ] Domain responses include DNS instructions and last verification error details when applicable.
- [ ] Because scopes are non-hierarchical, keys need `domains:read` for domain reads and `domains:write` for domain mutations.
- [ ] Read-only domain keys cannot create, verify, update, or delete domains.
- [ ] `docs/api/openapi.yaml` documents public custom domain endpoints.
- [ ] Backend tests cover successful domain operations, duplicate hostnames, ownership failures, DNS verification outcomes, and insufficient scopes.

---

## Phase 5: Hosted OpenAPI Reference and Written Tutorials

**User stories**: 33, 34, 35, 36, 40

### What to build

Publish the public API contract as an OpenAPI specification and host interactive documentation inside the app. Add written tutorial documentation that walks a user through creating an API key in the dashboard, making an authenticated request, creating or selecting a site, uploading a zip archive, and inspecting the resulting URL or domain state.

This phase makes the API self-serve and reviewable as a product contract rather than an undocumented set of endpoints.

### Acceptance criteria

- [ ] A hand-maintained committed OpenAPI spec at `docs/api/openapi.yaml` describes public authentication, scopes, errors, sites, zip upload flow, and custom domains.
- [ ] OpenAPI excludes dashboard/session API key management endpoints.
- [ ] The webapp hosts a public Scalar OpenAPI documentation page reachable from the product.
- [ ] Public API documentation is linked from public navigation and from the dashboard API keys page.
- [ ] If Scalar provides interactive request functionality, it remains available.
- [ ] Russian written tutorial documentation explains key creation, bearer authentication, common curl usage, site creation, zip upload, and domain inspection.
- [ ] Documentation examples use `HOSTAMAT_API_BASE_URL`, with default value `https://api.hostamat.ru/v1`, and the public error model.
- [ ] The OpenAPI spec is validated by an automated check or documented local command.
- [ ] Manual browser verification confirms the hosted docs page loads and references the current spec.

---

## Phase 6: Contract Hardening and Release Readiness

**User stories**: 32, 37, 38, 39, 40

### What to build

Tighten the public API release by auditing contract consistency, security-sensitive behavior, logging, migration safety, and failure modes across all resource groups. This phase focuses on making the first version dependable before users build automation against it.

The outcome should be a stable v1 surface with clear errors, secure credential handling, adequate tests, and documented operational expectations.

### Acceptance criteria

- [ ] All public API endpoints use the same authentication, scope, and error-response conventions.
- [ ] Public API routes do not depend on browser cookies or CSRF assumptions.
- [ ] Public API routes are documented for `https://api.hostamat.ru/v1/...`; backend host enforcement is not required.
- [ ] Public API does not enable CORS in v1.
- [ ] Authorization headers, API key secrets, and sensitive upload payload details are not logged.
- [ ] Expired and revoked key behavior is consistent across all public API endpoints.
- [ ] Migration and rollback considerations for API key persistence are documented.
- [ ] Backend test coverage includes representative success and failure paths for every public API resource group.
- [ ] Final verification includes backend tests, OpenAPI validation, and manual dashboard/docs checks.
