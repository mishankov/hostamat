# HOSTAMAT - PROJECT KNOWLEDGE BASE

**Generated:** 2026-01-27
**Commit:** 1fe7a99
**Branch:** master

## OVERVIEW

Static site hosting platform with automatic subdomain provisioning. Go backend + SvelteKit frontend + Caddy reverse proxy + S3 storage + PostgreSQL.

**Runtime**: Bun (frontend), Go 1.25 (backend)
**Task Runner**: [Task](https://taskfile.dev/) (not make/npm scripts)
**Code Style**: See [docs/conventions/code-style.md](docs/conventions/code-style.md)

## STRUCTURE

```
hostamat/
├── platform/          # Go backend (domain-driven) → see platform/AGENTS.md
│   ├── cmd/           # Entry points: hostamat, dashboard
│   └── internal/      # Domains: site, file, subdomain, uploadsession
├── webapp/            # SvelteKit frontend (Bun runtime) → see webapp/AGENTS.md
│   └── src/
│       ├── lib/       # Components, API client, utilities → see webapp/src/lib/AGENTS.md
│       └── routes/    # Pages: dashboard, auth, profile
├── infra/             # Docker infrastructure → see infra/AGENTS.md
│   ├── test/          # Local development compose
│   └── prod/          # Docker Swarm production
├── docs/              # Project documentation
│   └── conventions/   # Code style, testing, patterns
├── bruno/             # API test collections
└── Taskfile.dist.yml  # Task runner for all operations
```

## DOMAIN-SPECIFIC GUIDES

| Component | Documentation | Focus |
|-----------|---------------|-------|
| **Backend** | [platform/AGENTS.md](platform/AGENTS.md) | Go conventions, domain-driven design, platforma framework |
| **Frontend** | [webapp/AGENTS.md](webapp/AGENTS.md) | SvelteKit, Result pattern, components |
| **Shared UI** | [webapp/src/lib/AGENTS.md](webapp/src/lib/AGENTS.md) | Reusable components, API client |
| **Site Domain** | [platform/internal/site/AGENTS.md](platform/internal/site/AGENTS.md) | Site lifecycle, handlers, upload flow |
| **Infrastructure** | [infra/AGENTS.md](infra/AGENTS.md) | Deployment, services, Docker Swarm |

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
