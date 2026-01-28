# PLATFORM (GO BACKEND) - AGENTS KNOWLEDGE BASE

**Go backend for static site hosting. Domain-driven, platforma framework, sqlx + MinIO.**

## STRUCTURE

```
platform/
├── cmd/
│   ├── hostamat/       # Main API server (port 8000)
│   └── dashboard/      # Admin dashboard backend
├── internal/
│   ├── site/           # Core domain: sites CRUD, file upload, serving
│   ├── file/           # File storage (S3/MinIO operations)
│   ├── subdomain/      # Subdomain generation + Caddy integration
│   ├── uploadsession/  # Session-based chunked file uploads
│   ├── adapters/       # External integrations
│   │   ├── caddy/      # Dynamic route config via Caddy API
│   │   ├── config/     # Env var parsing (HOSTAMAT_*)
│   │   └── s3/         # MinIO client wrapper
│   └── middlewares/    # HTTP middleware (logging, etc.)
└── go.mod              # Go 1.25, platforma, sqlx, minio-go
```

## WHERE TO LOOK

| Task | Location | Notes |
|------|----------|-------|
| HTTP handlers | `internal/*/handler_*.go` | One file per operation (create, delete, etc.) |
| Business logic | `internal/*/service.go` | Core domain logic, orchestration |
| Database ops | `internal/*/repository.go` | sqlx-based, PostgreSQL |
| Models/structs | `internal/*/model.go` | Domain entities, constants |
| Domain wiring | `internal/*/domain.go` | Assembles repository → service → handlers |
| Migrations | `internal/*/internal/migrations/` | Database schema changes |
| Config parsing | `internal/adapters/config/config.go` | All `HOSTAMAT_*` env vars |
| App entry point | `cmd/hostamat/hostamat.go` | Dependency injection, startup logic |
| Caddy integration | `internal/adapters/caddy/` | Dynamic subdomain routing |

## CONVENTIONS

**Domain Structure** (each domain has):
- `domain.go` - Wires repository → service → handlers, exposes Domain struct
- `service.go` - Business logic, interfaces defined inline (e.g., `type repository interface`)
- `repository.go` - Database operations via sqlx (NamedExecContext, GetContext, SelectContext)
- `model.go` - Structs with `db:"col"` and `json:"field"` tags
- `handler_*.go` - HTTP handlers, one per operation (e.g., `handler_create.go`, `handler_delete.go`)

**platforma Framework**:
- Uses `platforma/application` for lifecycle management (OnStart hooks, graceful shutdown)
- `platforma/httpserver` for routing (`HandlerGroup`, middleware chaining)
- `platforma/auth` for session-based auth (`UserFromContext`)
- Domains register via `app.RegisterDomain("name", "db_connection_id", domain)`
- Docs for agents - https://platforma-dev.github.io/platforma/llms.txt

**Handlers**:
- Interface-based service dependencies (e.g., `type createSiteService interface { Create(...) }`)
- Manual JSON encoding/decoding (no magic frameworks)
- Custom error codes via `handlerError{ErrorCode, ErrorMessage}`

**Repository**:
- sqlx named queries (`:param` placeholders)
- `db` interface defined per domain (not shared)
- Migrations via `Migrations() []database.Migration`

## ANTI-PATTERNS

- DO NOT create separate interface files - define in service.go
- DO NOT use ORM - use sqlx with raw SQL
- DO NOT test logging - only business/db/http logic
- DO NOT use `*_test.go` in same package - use `package foo_test`
- DO NOT add dependencies beyond: platforma, sqlx, minio-go, uuid, prometheus
