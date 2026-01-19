# Agent Guidelines for Hostamat

This document provides essential information for AI coding agents working in the Hostamat repository. Hostamat is a static site hosting platform with a Go backend (`platform/`) and SvelteKit frontend (`webapp/`).

## Project Structure

```
hostamat/
├── platform/          # Go backend service
│   ├── cmd/          # Application entry points (hostamat, dashboard)
│   └── internal/     # Internal packages (site, file, subdomain, uploadsession, etc.)
├── webapp/           # SvelteKit frontend
├── infra/            # Infrastructure configurations
└── bruno/            # API testing collections
```

## Build, Test, and Run Commands

All commands use [Task](https://taskfile.dev/) as the build tool. View `Taskfile.dist.yml` for all available tasks.

### Backend (Go)

**Run tests:**
```bash
cd platform && go test ./... -v -cover
```

**Run a single test:**
```bash
cd platform && go test ./internal/<package> -v -run TestName
# Example: cd platform && go test ./internal/file -v -run TestService_GetFile
```

**Run backend service:**
```bash
task backend
# Or run directly in test environment
task test-infra  # Start dependencies first
cd infra/test && docker compose up hostamat
```

**Run dashboard:**
```bash
task dash
# Or manually:
cd platform && TOKEN=secret_token DB_CONN=postgresql://hostamat:hostamat@localhost:5432/hostamat?sslmode=disable go run ./cmd/dashboard
```

### Frontend (SvelteKit)

**Install dependencies:**
```bash
cd webapp && bun install
```

**Development server:**
```bash
cd webapp && bun run --bun dev
# Or: task webapp
```

**Build:**
```bash
cd webapp && bun run build
```

**Run tests:**
```bash
cd webapp && npm run test        # Run all tests once
cd webapp && npm run test:unit   # Run tests in watch mode
```

**Linting and formatting:**
```bash
cd webapp && npm run lint        # Check code with ESLint and Prettier
cd webapp && npm run format      # Format code with Prettier
cd webapp && npm run check       # Type-check with svelte-check
```

## Code Style Guidelines

### Go Backend

**Imports:**
- Standard library imports first
- External dependencies second
- Internal packages third
- Group with blank lines between groups

Example:
```go
import (
    "context"
    "errors"
    "fmt"
    
    "github.com/google/uuid"
    "github.com/jmoiron/sqlx"
    
    "hostamat/internal/file"
    "hostamat/internal/subdomain"
)
```

**Naming conventions:**
- Use camelCase for unexported identifiers
- Use PascalCase for exported identifiers
- Interface names: descriptive nouns (e.g., `repository`, `subdomainService`)
- Keep names concise but clear

**Types and interfaces:**
- Define interfaces at the consumer level (in the file that uses them)
- Use constructor functions: `NewService(...)` pattern
- Struct fields for dependencies should use interface types when possible

**Error handling:**
- Return errors, don't panic
- Wrap errors with context when appropriate
- Use custom error types/variables for domain-specific errors (e.g., `errUploadSize`, `ErrInvalidFileType`)
- Check `errors.Is()` for error comparison

**Logging:**
- Use structured logging with `log.InfoContext(ctx, message, key, value, ...)`
- Log significant business events (creation, deletion, errors)
- Don't test logging code

**Context:**
- Always pass `context.Context` as first parameter
- Use `ctx` as variable name
- Extract user info with `auth.UserFromContext(ctx)`

### TypeScript/SvelteKit Frontend

**Formatting:**
- Use tabs for indentation
- Single quotes for strings
- No trailing commas
- Print width: 100 characters
- Configured in `.prettierrc`

**Imports:**
- Use `$lib/` path alias for library imports
- Example: `import { Badge, Button, Card } from '$lib/components';`

**TypeScript:**
- Strict mode enabled
- Always add type annotations for function parameters
- Use `allowJs` and `checkJs` for JavaScript files

**Svelte components:**
- Use `<script>` for component logic
- Use `<svelte:head>` for page metadata
- Scoped styles in `<style>` blocks
- Follow SvelteKit's file-based routing conventions

**ESLint:**
- `no-undef` rule disabled (TypeScript handles this)
- Prettier integration configured
- TypeScript ESLint recommended rules enabled

## Testing Guidelines

**From CONTRIBUTION.md:**

- Always run tests before committing to the repository
- Always use separate package for tests with `_test` suffix (e.g., `package subdomain_test`)
- Use only the standard `testing` package - no additional testing frameworks
- Do not test internal implementation details
- When testing objects that depend on `*sqlx.DB`, use in-memory SQLite database
- Repository tests: only test database interaction logic
- Service tests: only test business logic
- Handler tests: only test HTTP/protocol interaction logic
- Do not test logging

**Test file naming:**
- Go: `*_test.go` (e.g., `service_test.go`, `repository_test.go`)
- Frontend: `*.test.ts` or `*.spec.ts`

## Architecture Patterns

**Backend layers:**
1. **Handlers**: HTTP request/response handling (`handler_*.go`)
2. **Services**: Business logic, orchestration (`service.go`)
3. **Repositories**: Database operations (`repository.go`)
4. **Domain**: Business models and validation (`domain.go`, `model.go`)

**Dependency injection:**
- Services receive dependencies through constructors
- Use interfaces for loose coupling
- Example: `NewService(repository repository, subdomainService subdomainService, ...)`

**Access control:**
- Check user access in service layer
- Use `checkAccess()` helper methods
- Extract user from context: `auth.UserFromContext(ctx)`

## Common Patterns

**Go service methods:**
```go
func (s *Service) MethodName(ctx context.Context, params...) (*ReturnType, error) {
    // 1. Validate access
    // 2. Validate input
    // 3. Execute business logic
    // 4. Log important events
    // 5. Return result
}
```

**Error creation:**
```go
var (
    errLongName  = errors.New("name too long")
    errShortName = errors.New("name too short")
)
```

## Infrastructure

- **Database**: PostgreSQL (production), SQLite (tests)
- **Storage**: MinIO (S3-compatible)
- **Reverse proxy**: Caddy
- **Package manager**: Go modules (backend), Bun (frontend)
- **Container registry**: `hostamat-container-registry.registry.twcstorage.ru`

## Additional Notes

- Go version: 1.25.0
- Node/Bun runtime for frontend
- Uses Task for build automation instead of Make
- Docker multi-platform builds target `linux/amd64`
- Monorepo structure with separate backend and frontend directories
