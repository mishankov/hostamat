# Hostamat Development Guidelines

## Build/Lint/Test Commands

**Webapp (Svelte/TypeScript):**
- `npm run dev` - Start dev server
- `npm run build` - Build for production
- `npm run lint` - Run Prettier + ESLint
- `npm run format` - Format code with Prettier
- `npm test` - Run unit tests
- `npm run test:unit` - Run tests in watch mode

**Platform (Go):**
- `go test ./...` - Run all tests
- `go test -v ./internal/site/...` - Run specific package tests
- `go build ./cmd/hostamat` - Build binary

**Infrastructure:**
- `task test-infra` - Start test infrastructure
- `task backend` - Start backend with dependencies

## Code Style Guidelines

**Go:**
- Use standard Go formatting (gofmt)
- Imports: stdlib, third-party, local (grouped)
- Error handling: return errors, use context
- Naming: camelCase, interfaces end with "er"
- Tests: `func TestXxx(t *testing.T)` pattern

**TypeScript/Svelte:**
- Prettier: tabs, single quotes, no trailing commas
- ESLint: TypeScript + Svelte rules
- Components: PascalCase, props typing
- State management: Svelte 5 runes (`$state`, `$props`)

**General:**
- Error handling: proper logging and user feedback
- Testing: unit tests for business logic
- File structure: follow existing patterns