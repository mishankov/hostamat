# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Hostamat is a static site hosting platform with a Go backend (`platform/`) and SvelteKit frontend (`webapp/`). Users can create sites with subdomains, upload files, and manage web content.

## Build, Test, and Run Commands

All commands use [Task](https://taskfile.dev/) as the build tool. See `Taskfile.dist.yml` for all available tasks.

### Backend (Go)

```bash
# Run tests
cd platform && go test ./... -v -cover

# Run a single test
cd platform && go test ./internal/<package> -v -run TestName
# Example: cd platform && go test ./internal/file -v -run TestService_GetFile

# Start backend (starts infra dependencies first)
task backend

# Start dashboard admin interface
task dash
```

### Frontend (SvelteKit)

```bash
# Install dependencies
cd webapp && bun install

# Development server
task webapp
# Or: cd webapp && bun run --bun dev

# Build
cd webapp && bun run build

# Run tests
cd webapp && npm run test        # Run all tests once
cd webapp && npm run test:unit   # Run tests in watch mode

# Linting and formatting
cd webapp && npm run lint        # Check code with ESLint and Prettier
cd webapp && npm run format      # Format code with Prettier
cd webapp && npm run check       # Type-check with svelte-check
```

### Infrastructure

```bash
# Start test infrastructure (Caddy, MinIO, PostgreSQL)
task test-infra

# Deploy (builds and pushes all images, then deploys)
task deploy-stack
```

## Architecture

### Backend Layers

The Go backend follows a layered architecture in `platform/internal/`:

1. **Handlers** (`handler_*.go`): HTTP request/response handling
2. **Services** (`service.go`): Business logic and orchestration
3. **Repositories** (`repository.go`): Database operations
4. **Domain** (`domain.go`, `model.go`): Business models and validation

Domain modules: `site/` (main CRUD), `file/` (S3 storage), `subdomain/` (Caddy integration), `uploadsession/` (upload tracking), `adapters/` (external integrations).

### Frontend Structure

SvelteKit with file-based routing in `webapp/src/routes/`:
- Layout groups: `(auth)/` for auth pages, `(legal)/` for legal pages
- Shared components in `src/lib/components/`
- API client in `src/lib/api-client/`

### Infrastructure

- **Database**: PostgreSQL (production), SQLite (in-memory for tests)
- **Storage**: MinIO (S3-compatible)
- **Reverse proxy**: Caddy
- **Container registry**: `hostamat-container-registry.registry.twcstorage.ru`

## Code Style

### Go

**Imports order** (with blank lines between groups):
1. Standard library
2. External dependencies
3. Internal packages

**Patterns**:
- Define interfaces at the consumer level
- Use constructor functions: `NewService(...)`
- Return errors, don't panic; wrap errors with context
- Use custom error variables: `var errLongName = errors.New("name too long")`
- Structured logging: `log.InfoContext(ctx, message, key, value, ...)`
- Always pass `context.Context` as first parameter

### TypeScript/Svelte

- Tabs for indentation, single quotes, no trailing commas
- Print width: 100 characters (configured in `.prettierrc`)
- Use `$lib/` path alias for imports
- Strict TypeScript mode enabled

## Testing Guidelines

- Use separate package for tests with `_test` suffix (e.g., `package subdomain_test`)
- Use only the standard `testing` package - no additional testing frameworks
- When testing objects that depend on `*sqlx.DB`, use in-memory SQLite database
- Repository tests: only test database interaction logic
- Service tests: only test business logic
- Handler tests: only test HTTP/protocol interaction logic
- Do not test logging or internal implementation details
