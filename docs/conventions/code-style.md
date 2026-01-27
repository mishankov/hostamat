# CODE STYLE CONVENTIONS

Code formatting, linting, and style guidelines for Hostamat project.

## TypeScript/JavaScript (Frontend)

### Prettier Configuration

- **Indentation**: Tabs (not spaces)
- **Quotes**: Single quotes (`'`)
- **Trailing commas**: None
- **Line width**: 100 characters

### ESLint

- TypeScript strict mode enabled
- No bare `fetch()` - use `safeFetch()` wrapper
- No `console.log` for error handling - use Result types

### Component Props Pattern

```typescript
// Define props interface
interface Props {
	kind?: 'primary' | 'confirm' | 'danger' | 'action';
	disabled?: boolean;
	filled?: boolean;
	onclick?: (e: Event) => void;
}

// Destructure with $props()
let { kind = 'primary', disabled = false }: Props = $props();
```

### Error Handling (CRITICAL)

**Result Pattern** - NEVER throw exceptions:

```typescript
type Result<T, E = Error> = { data: T; error: null } | { data: null; error: E };
type AsyncResult<T, E = Error> = Promise<Result<T, E>>;

// CORRECT
const { data, error } = await safeFetch('/api/endpoint');
if (error) return { data: null, error: error.message };

// WRONG
throw new Error('Something failed');
```

**API Client Pattern**:

```typescript
export const createSite = async (name: string): AsyncResult<boolean, string> => {
	const { data: resp, error } = await safeFetch('/api/proxy', {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ name })
	});
	if (error) return { data: null, error: error.message };
	if (!resp.ok) {
		const errorData = await resp.json();
		return { data: null, error: errorData.errorMessage };
	}
	return { data: true, error: null };
};
```

## Go (Backend)

### General Style

- Follow standard Go conventions (`gofmt`, `golint`)
- Interfaces defined inline in service files (not separate interface files)
- One handler per file: `handler_create.go`, `handler_delete.go`, etc.

### Testing

- **Package**: Use `*_test` suffix for test packages (e.g., `package site_test`)
- **Framework**: Standard `testing` package ONLY - no additional frameworks
- **Focus**: Test business logic, database operations, HTTP handlers
- **DO NOT TEST**: Logging, internal implementation details

### Domain-Driven Structure

Each domain follows this pattern:

```
domain/
├── domain.go      # Wiring: repository → service → handlers
├── service.go     # Business logic, inline interfaces
├── repository.go  # Database operations via sqlx
├── model.go       # Structs with db:"col" and json:"field" tags
└── handler_*.go   # One file per HTTP operation
```

### Repository Pattern

```go
// Inline interface definition in service.go
type repository interface {
	Create(ctx context.Context, site Site) error
	GetByID(ctx context.Context, id string) (Site, error)
}

// sqlx named queries in repository.go
const createQuery = `
	INSERT INTO sites (id, name, user_id) 
	VALUES (:id, :name, :user_id)
`
```

### Handler Pattern

```go
// Interface-based service dependencies
type createSiteService interface {
	Create(ctx context.Context, userID, name string) (Site, error)
}

// Manual JSON encoding/decoding
func (h *CreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Handle error
	}
	// ...
}
```

## File Naming

| Type | Pattern | Example |
|------|---------|---------|
| Go handlers | `handler_*.go` | `handler_create.go`, `handler_delete.go` |
| Go tests | `*_test.go` | `service_test.go`, `repository_test.go` |
| Svelte components | `PascalCase.svelte` | `Button.svelte`, `Modal.svelte` |
| TypeScript utilities | `kebab-case.ts` | `safe-fetch.ts`, `api-client.ts` |

## Import Organization

### TypeScript

```typescript
// 1. External dependencies
import { goto } from '$app/navigation';

// 2. Internal lib imports
import Button from '$lib/components/Button.svelte';
import { safeFetch } from '$lib/safe';

// 3. Relative imports
import type { Props } from './types';
```

### Go

```go
import (
	// 1. Standard library
	"context"
	"encoding/json"
	
	// 2. External dependencies
	"github.com/platforma-dev/platforma/httpserver"
	
	// 3. Internal packages
	"hostamat/internal/site"
)
```

## Comments

### TypeScript

- Use JSDoc for public APIs
- Inline comments for complex logic only
- No obvious comments (e.g., `// Create user` above `createUser()`)

### Go

- Exported functions/types MUST have doc comments
- Start with the name: `// CreateSite creates a new site...`
- No obvious comments
