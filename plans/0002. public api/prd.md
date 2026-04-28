# PRD: Public API for Sites, Domains, and Uploads

## Problem Statement

Hostamat users can manage sites, custom domains, and uploads through the authenticated dashboard, but they cannot safely automate the same workflows from external tools. Users who want CI/CD uploads, custom scripts, internal admin panels, or integrations must rely on browser sessions and dashboard-only flows.

This limits Hostamat's usefulness for technical users. A static hosting platform should let users create sites, upload files, manage domains, and inspect state through a documented, stable API. At the same time, exposing these capabilities publicly introduces security, authorization, quota, documentation, and support risks that need to be designed deliberately rather than by exposing the current dashboard API as-is.

## Solution

Introduce a first-class public API that lets authenticated users create scoped API keys from the dashboard and use those keys to interact with their own Hostamat resources.

The public API will cover site management, custom domain management, and raw zip upload workflows. API keys will support explicit scopes, human-readable names, optional expiration dates, revocation, and last-used tracking. All existing ownership checks, validation rules, and account-level limits continue to apply to API requests exactly as they do for dashboard sessions.

The dashboard will gain an API key management surface where users can create, inspect, edit names for, and delete keys. The app will also host public Scalar-based interactive OpenAPI documentation and Russian written tutorial-style documentation so users can discover endpoints, authenticate requests, and build an upload workflow without reading source code.

## User Stories

1. As a Hostamat user, I want to create an API key from the dashboard, so that I can automate site and domain workflows outside the browser.
2. As a Hostamat user, I want to give each API key a readable name, so that I can recognize which script, project, or integration uses it.
3. As a Hostamat user, I want API keys to have expiration dates, so that stale credentials do not remain valid forever.
4. As a Hostamat user, I want to choose scopes for an API key, so that an integration only gets the permissions it needs.
5. As a Hostamat user, I want to delete or revoke an API key, so that compromised or obsolete credentials immediately stop working.
6. As a Hostamat user, I want to see when an API key was last used, so that I can detect unused, stale, or suspicious keys.
7. As a Hostamat user, I want to see when each API key was created and when it expires, so that I can rotate keys before automation breaks.
8. As a Hostamat user, I want the full API key secret shown only once after creation, so that Hostamat does not expose reusable credentials later.
9. As a Hostamat user, I want copied API keys to be easy to use in command-line tools, so that setup is not error-prone.
10. As a Hostamat user, I want expired keys to be visibly marked in the dashboard, so that I understand why an integration may have stopped working.
11. As a Hostamat user, I want revoked keys to be rejected by the public API, so that deletion has immediate security effect.
12. As a Hostamat user, I want to list my sites through the public API, so that automation can discover current resources.
13. As a Hostamat user, I want to create a site through the public API, so that CI or setup scripts can provision projects automatically.
14. As a Hostamat user, I want to rename a site through the public API, so that automation can keep site metadata in sync.
15. As a Hostamat user, I want to delete a site through the public API, so that cleanup scripts can remove unused projects.
16. As a Hostamat user, I want site creation through the public API to respect my account's site limit, so that public API access cannot bypass product constraints.
17. As a Hostamat user, I want API responses to include stable site identifiers and generated Hostamat domains, so that clients can chain calls reliably.
18. As a Hostamat user, I want to create a custom domain through the public API, so that automation can attach hostnames to sites.
19. As a Hostamat user, I want to list custom domains through the public API, so that integrations can inspect active and pending custom hostnames.
20. As a Hostamat user, I want to verify a custom domain through the public API, so that scripts can retry DNS verification after records change.
21. As a Hostamat user, I want to update which site a custom domain points to through the public API, so that automation can move a hostname between my sites.
22. As a Hostamat user, I want to delete a custom domain through the public API, so that cleanup scripts can detach hostnames.
23. As a Hostamat user, I want domain API responses to include DNS instructions and verification errors, so that tools can display actionable setup steps.
24. As a Hostamat user, I want domain management through the public API to respect ownership and uniqueness rules, so that another user cannot claim or mutate my hostnames.
25. As a Hostamat user, I want to upload a zip archive through the public API, so that a tool can publish site files.
26. As a Hostamat user, I want zip uploads to replace the current site files, so that CI/CD publishes a complete build artifact.
27. As a Hostamat user, I want upload responses to include a summary and file list, so that automation logs show what changed.
28. As a Hostamat user, I want upload limits and file validation to match the dashboard behavior, so that public API uploads cannot bypass account or site limits.
29. As a Hostamat user, I want upload errors to identify invalid file types, size problems, missing sites, and authorization failures, so that upload scripts can fail clearly.
30. As a Hostamat user, I want read-only API keys to be unable to create, update, delete, or upload, so that analytics or inventory scripts are low risk.
31. As a Hostamat user, I want upload-only or narrowly scoped keys where possible, so that CI does not need domain-management access.
32. As a Hostamat user, I want public API errors to use a consistent JSON shape, so that client code can handle failures predictably.
33. As a Hostamat user, I want public API authentication documented with examples, so that I can make my first request quickly.
34. As a Hostamat user, I want interactive OpenAPI documentation hosted inside the app, so that I can inspect endpoints and request schemas from the product.
35. As a Hostamat user, I want Russian written tutorial documentation, so that I can follow a complete example from key creation to upload.
36. As a Hostamat user, I want documentation examples using common tools like curl, so that I can test the API without writing an SDK first.
37. As a Hostamat operator, I want public API authentication to be separate from browser session authentication, so that key lifecycle and audit behavior are explicit.
38. As a Hostamat operator, I want API key secrets stored securely, so that database access does not expose plaintext active keys.
39. As a Hostamat operator, I want last-used tracking to update without making every request fragile, so that observability does not reduce API reliability.
40. As a Hostamat operator, I want public API routes to be documented and versionable, so that future changes can avoid breaking existing clients unexpectedly.

## Implementation Decisions

- Public API routes will be treated as a first-class external contract, separate from dashboard-oriented session routes.
- The initial public API version will be served from `https://api.hostamat.ru/v1/...`.
- Backend host enforcement is not required; infrastructure and documentation define the public API host.
- DNS/proxy routing for `api.hostamat.ru` is manual deployment work and is not included in this implementation scope.
- Existing dashboard/private API routes do not need to be renamed as part of this feature.
- Public API includes `GET /v1/me` as a stable authentication test endpoint returning minimal non-sensitive user identity plus key metadata.
- Public API uses `PATCH` for partial resource updates.
- Public site routes are `GET /v1/sites`, `POST /v1/sites`, `GET /v1/sites/{siteId}`, `PATCH /v1/sites/{siteId}`, `DELETE /v1/sites/{siteId}`, and `POST /v1/sites/{siteId}/upload`.
- Site zip upload is modeled as `POST /v1/sites/{siteId}/upload`.
- Public custom-domain routes are `GET /v1/domains`, `POST /v1/domains`, `GET /v1/domains/{domainId}`, `PATCH /v1/domains/{domainId}`, `POST /v1/domains/{domainId}/verify`, and `DELETE /v1/domains/{domainId}`.
- Custom domain verification is modeled as an action endpoint: `POST /v1/domains/{domainId}/verify`.
- API key authentication will use bearer tokens in the Authorization header.
- API key secrets will use the format `hmat_<id>_<secret>`, where `<id>` is a UUID and the secret portion is high entropy.
- API key records will store metadata needed for management: owner, name, key identifier, scopes, optional expiration timestamp, creation timestamp, last-used timestamp, and revocation/deletion state.
- API key secrets must not be stored in plaintext. The system will store a peppered HMAC-SHA-256-style digest and only show the full secret once immediately after creation.
- The API key digest pepper will come from a stable backend environment variable such as `HOSTAMAT_API_KEY_PEPPER`.
- API key display in the dashboard will use a non-secret identifier or prefix so users can distinguish keys without revealing the key value.
- API keys will support explicit non-hierarchical scopes: `sites:read`, `sites:write`, `sites:upload`, `domains:read`, and `domains:write`.
- No scope implies another scope.
- `sites:write` covers site create, site name update, and site delete.
- `sites:upload` covers raw zip uploads to any site owned by the key owner and does not allow site create, update, or delete.
- `domains:write` covers custom domain create, update, verify, move between owned sites, and delete.
- Scope checks will happen after successful API key authentication and before invoking resource mutations.
- Existing user ownership and account-limit rules remain authoritative. API keys inherit the owning user's limits; they do not define separate quotas.
- Existing site lifecycle behavior remains the source of truth for create, list, update, and delete semantics.
- Public site update is limited to site name in v1.
- Public site reads include list and single-site detail endpoints.
- Domains API exposes custom domains only. Generated Hostamat domains are exposed on site responses, not through the domains API.
- Domain list is global for the user in v1, and custom domain creation requires `siteId`.
- Existing custom-domain lifecycle behavior remains the source of truth for custom domain create, list, detail, update, verify, and delete semantics.
- Public upload accepts a raw `application/zip` request body only, replaces all existing files for the target site, and returns only a summary plus file list.
- Public zip upload preserves the current archive path handling and uploaded/skipped-file behavior used by the existing file service.
- Public zip upload is synchronous in v1 and does not expose upload status/history endpoints.
- Public API response models should be stable, documented, and intentionally shaped for external consumers rather than leaking accidental internal structures.
- Public API JSON fields use `camelCase`.
- Public API timestamps are RFC 3339 strings.
- Error responses should use a consistent JSON envelope with a machine-readable code and English human-readable message.
- Request IDs are exposed through the existing `Hostamat-Trace-Id` response header only, including error responses.
- Authentication failures, missing scopes, missing resources, validation failures, expired keys, and revoked keys should be distinguishable.
- API key authentication fails if the owning user cannot be resolved as an active user.
- Last-used tracking should record only `lastUsedAt` for successful key authentication, should not store request bodies or sensitive payloads, and should be throttled to avoid a database write on every request.
- API keys may have no expiration date.
- API key expiration is stored as a timestamp and evaluated with backend server time. Date-only UI choices are converted to an end-of-day timestamp before storage.
- API key metadata represents no-expiration keys with `expiresAt: null`.
- API key deletion is implemented as soft revocation. Revoked keys are rejected by the public API and hidden from the UI.
- Only the API key name can be edited after creation.
- API key names are required labels, do not need to be unique per user, and have an 80-character maximum length.
- Users can have at most 20 active API keys; expired and revoked keys do not count toward that limit.
- The dashboard will add an API keys page in dashboard navigation where users can list keys, create keys with scope checkboxes and presets, choose expiration presets including "never expires", copy the newly created secret once, edit key names, and delete keys.
- The one-time secret view will include a copy action and require acknowledgement that the secret cannot be viewed again.
- The dashboard should prevent creating unusable keys, such as a key without scopes.
- Hosted interactive documentation will be available publicly inside the main app and backed by Scalar over an OpenAPI specification.
- Written tutorial documentation will be Russian-only in v1 and reachable from public navigation and the API keys dashboard page.
- OpenAPI title, summaries, and descriptions will be Russian in v1, while schema fields and error codes remain stable code-oriented identifiers.
- The OpenAPI specification will be a hand-maintained committed artifact at `docs/api/openapi.yaml`.
- Each public API implementation phase should update `docs/api/openapi.yaml` for the endpoints it introduces.
- OpenAPI documents only the public bearer-token API, not dashboard/session API key management endpoints.
- Tutorial examples will use `HOSTAMAT_API_BASE_URL`, documented with default value `https://api.hostamat.ru/v1`.
- The public API should avoid depending on browser cookies or CSRF behavior because external clients will authenticate with bearer keys.
- Public API does not enable CORS in v1.
- Public API does not add app-level rate limiting in v1.
- Public API and API key creation are enabled immediately for authenticated users; no feature flag is required for v1.
- Request logging must continue to redact Authorization headers and must not log API key secrets.
- API key lifecycle, verification, scopes, bearer authentication middleware, and scope middleware should live in their own backend domain/module rather than being folded into browser auth.
- A dedicated public API module should own the `/v1` route group and depend on existing product services plus API key middleware.
- Public API handlers should use a thin facade over existing product services to preserve current behavior while mapping external DTOs and errors deliberately.
- The implementation should leave room for future SDKs, rate limits, usage analytics, and organization/team keys without requiring them for the first release.

## Testing Decisions

- Good tests should verify externally visible behavior: accepted requests, rejected requests, persisted state, ownership enforcement, scope enforcement, expiration behavior, and response/error contracts.
- Tests should not assert private implementation details such as hashing algorithm internals beyond verifying that plaintext API keys are not retrievable after creation.
- Backend service tests should cover API key creation, key verification, key expiration, key deletion/revocation, scope checks, owner-user status checks, and throttled last-used updates.
- Backend HTTP tests should cover authentication success and failure, missing bearer tokens, invalid bearer tokens, expired keys, revoked keys, insufficient scopes, and successful scoped access.
- Site API tests should cover list, create, update, and delete through public API authentication, including preservation of user-level site limits.
- Site upload API tests should cover raw zip upload, replacement behavior, current uploaded/skipped-file behavior, file validation, size limits, missing site access, and insufficient `sites:upload` scope.
- Custom domain API tests should cover list, detail, create, update, verify, and delete through public API authentication, including ownership, duplicate hostname, and validation failures.
- Documentation tests or CI checks should validate that the OpenAPI spec is syntactically valid and remains available to the hosted documentation UI.
- Frontend work should follow the existing policy of not adding webapp unit tests. Manual and browser-level verification should cover the key management UI and hosted docs page.
- Prior art for backend tests exists in service-level tests around custom domains and file behavior; new tests should follow the same package and behavior-driven style.

## Out of Scope

- Organization or team-owned API keys.
- OAuth application flows for third-party apps.
- Personal access token import/export.
- Fine-grained per-site or per-domain key restrictions beyond scopes.
- Rate limiting, billing metering, and usage dashboards.
- Browser CORS support for direct public API calls.
- Webhooks.
- Official SDKs or CLI tooling.
- Automated key rotation.
- API key management through public API keys.
- Recovery of a full API key secret after creation.
- Public upload sessions or multipart per-file upload endpoints.
- Public upload status or upload history endpoints.
- Public API pagination, filtering, or sorting for v1 list endpoints.
- Bruno/API collection updates.
- Public API access to account deletion, billing, or user profile management.
- Changing existing account-level limits for sites, domains, file size, or upload validation.
- Replacing the dashboard's existing session-authenticated API surface.

## Further Notes

- This feature introduces a new external contract, so the first version should be conservative and stable rather than exposing every internal operation.
- File upload documentation should describe the raw zip upload sequence clearly: create or identify site, upload the zip archive, inspect uploaded/skipped files, then verify the public URL.
- The docs should include at least one complete curl-based tutorial and one reference example for scoped key creation from the UI.
- The OpenAPI documentation should make authentication, scopes, and error responses prominent because these are the most likely integration stumbling blocks.
- Security-sensitive UX details matter: key creation should make it obvious that the secret is shown once, and deletion should clearly communicate that existing integrations will stop working.
