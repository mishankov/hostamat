# LIB - FRONTEND SHARED CODE

Reusable components, API client, utilities for webapp.

## DIRECTORIES

| Directory     | Purpose                                                               |
| ------------- | --------------------------------------------------------------------- |
| `components/` | UI primitives: Button, Card, Modal, Input, Badge, Link, etc.          |
| `icons/`      | SVG icon components: Upload, Trash, Edit, Spinner, Check, X, etc.     |
| `api-client/` | Typed API wrappers: `sites.ts`, `upload-session.ts`                   |
| `safe/`       | Result pattern utilities: `Result<T, E>`, `tryCatch()`, `safeFetch()` |
| `server/`     | Server-only config/utilities                                          |

## KEY PATTERNS

### Result Type (safe/)

```typescript
type Result<T, E = Error> = { data: T; error: null } | { data: null; error: E };

// Usage in API client
const { data, error } = await safeFetch('/api/endpoint');
if (error) return { data: null, error: error.message };
```

**Critical**: NEVER use bare `fetch()` or throw exceptions. Always return `Result<T, E>`.

### Component Props (components/)

```typescript
// Define props interface
type ButtonProps = {
	tone?: 'accent' | 'success' | 'danger';
	emphasis?: 'quiet' | 'solid';
	disabled?: boolean;
	onclick?: (e: Event) => void;
};

// Destructure with $props()
let { tone = 'accent', emphasis = 'quiet', disabled = false }: ButtonProps = $props();
```

### API Client (api-client/)

- All functions return `AsyncResult<T, string>`
- Use `safeFetch()` for requests
- Custom headers: `X-Method`, `H-Site-Id`
- Server-side: pass `fetch` and `request` from `load()` functions
- Client-side: omit `fetch` param (uses global)

## CONVENTIONS

- Components: single `<script>`, props via `$props()`, scoped `<style>`
- Icons: pure SVG, no logic, accept `class` prop
- No console.log for errors - use Result types
- Server-only code: `lib/server/` (excluded from client bundle)
