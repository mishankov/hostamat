# HOSTAMAT - PROJECT KNOWLEDGE BASE

**Generated:** 2026-01-27
**Commit:** 1fe7a99
**Branch:** master

## OVERVIEW

Static site hosting platform with automatic subdomain provisioning. Go backend + SvelteKit frontend + Caddy reverse proxy + S3 storage + PostgreSQL.

## STRUCTURE

```
hostamat/
├── platform/          # Go backend (domain-driven)
│   ├── cmd/           # Entry points: hostamat, dashboard
│   └── internal/      # Domains: site, file, subdomain, uploadsession
├── webapp/            # SvelteKit frontend (Bun runtime)
│   └── src/
│       ├── lib/       # Components, API client, utilities
│       └── routes/    # Pages: dashboard, auth, profile
├── infra/             # Docker infrastructure
│   ├── test/          # Local development compose
│   └── prod/          # Docker Swarm production
├── bruno/             # API test collections
└── Taskfile.dist.yml  # Task runner for all operations
```

## WHERE TO LOOK

| Task | Location | Notes |
|------|----------|-------|
| Backend API handlers | `platform/internal/site/handler_*.go` | HTTP handlers per operation |
| Domain business logic | `platform/internal/*/service.go` | Core service layer |
| Database operations | `platform/internal/*/repository.go` | sqlx-based |
| Frontend pages | `webapp/src/routes/**/*.svelte` | SvelteKit file routing |
| UI components | `webapp/src/lib/components/*.svelte` | Reusable components |
| API client | `webapp/src/lib/api-client/` | Typed fetch wrappers |
| Config/env vars | `platform/internal/adapters/config/config.go` | All `HOSTAMAT_*` vars |
| Caddy integration | `platform/internal/adapters/caddy/` | Dynamic route config |

## CONVENTIONS

### Go Backend
- Domain-driven: each domain has `domain.go`, `service.go`, `repository.go`, `model.go`
- Uses `github.com/platforma-dev/platforma` framework
- Interfaces defined in service files, not separate files
- Tests use standard `testing` package only, no frameworks
- Test packages: `*_test` (separate package)
- Repository tests: in-memory SQLite via sqlx
- DO NOT test logging, only test business/db/http logic

### TypeScript Frontend
- Result pattern for error handling: `{ data: T, error: null } | { data: null, error: E }`
- `safeFetch()` wrapper - never use bare fetch
- Prettier: tabs, single quotes, no trailing commas, 100 char width
- ESLint + TypeScript strict mode
- Component props: `interface Props` with destructuring

## ANTI-PATTERNS

- DO NOT use additional Go testing frameworks
- DO NOT test internal implementation details
- DO NOT use bare `fetch()` - use `safeFetch()` or `tryCatch()`
- DO NOT throw exceptions in TypeScript - return Result types
- DO NOT use console.log for error handling

## COMMANDS

```bash
# Development (requires Task: go install github.com/go-task/task/v3/cmd/task@latest)
task test-infra        # Start PostgreSQL, MinIO, Caddy, proxy
task backend           # Run Go backend with hot reload
task webapp            # Run SvelteKit dev server (bun)
task dash              # Run dashboard

# Testing
cd platform && go test ./... -v -cover
cd webapp && bun run test

# Build & Deploy
task push-backend      # Build + push Docker image
task push-webapp       # Build + push frontend image
task deploy-stack      # Full deployment
```

## ENV VARS (Backend)

| Variable | Purpose |
|----------|---------|
| `HOSTAMAT_PORT` | HTTP server port |
| `HOSTAMAT_DOMAIN` | Base domain for sites (e.g., hostamat.ru) |
| `HOSTAMAT_FILESERVER_HOST` | Internal file server host |
| `HOSTAMAT_CADDY_URL` | Caddy admin API URL |
| `HOSTAMAT_CADDY_SHORTHAND_ID` | Caddy config identifier |
| `HOSTAMAT_S3_*` | S3/MinIO connection |
| `HOSTAMAT_POSTGRES_*` | Database connection |

## NOTES

- Uses Caddy for automatic SSL via Let's Encrypt
- Dynamic subdomain creation: `{random}.hostamat.ru`
- File uploads use session-based chunked upload flow
- Production: Docker Swarm with Portainer, Grafana, Prometheus monitoring
