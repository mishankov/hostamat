# Agent Guidelines for Hostamat

Hostamat is a static site hosting platform with a Go backend (`platform/`) and SvelteKit frontend (`webapp/`).

## Quick Start

**Backend commands:**
```bash
cd platform && go test ./... -v -cover      # Run all tests
task backend                                   # Run backend service
task dash                                      # Run dashboard
```

**Frontend commands:**
```bash
cd webapp && bun install                      # Install dependencies
cd webapp && bun run build                    # Build production
cd webapp && npm run test                     # Run tests
cd webapp && npm run lint                     # Check code
cd webapp && npm run format                   # Format code
cd webapp && npm run check                    # Type-check
```

## Detailed Guides

- [Go Backend Style](docs/AGENTS-go-style.md)
- [TypeScript/SvelteKit Style](docs/AGENTS-frontend-style.md)
- [Testing Patterns](docs/AGENTS-testing.md)
- [Architecture Patterns](docs/AGENTS-architecture.md)
- [Infrastructure](docs/AGENTS-infrastructure.md)

## Package Managers

- **Build automation**: [Task](https://taskfile.dev/) (instead of Make)
- **Backend**: Go modules
- **Frontend**: Bun

## Plan Mode

- Make the plan extremely concise. Sacrifice grammar for the sake of concision.
- At the end of each plan, give me a list of unresolved questions to answer, if any.
