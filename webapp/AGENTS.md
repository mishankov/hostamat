# WEBAPP - FRONTEND KNOWLEDGE BASE

**Generated:** 2026-01-27

## OVERVIEW

SvelteKit 2.x frontend with Bun adapter. Result pattern for error handling, component-driven UI, API proxy to backend.
Unit tests are not used in this webapp.

## STRUCTURE

```
webapp/
├── src/
│   ├── lib/
│   │   ├── components/        # Reusable UI: Button, Modal, Card, Input, Error, etc.
│   │   ├── icons/             # SVG icon components
│   │   ├── api-client/        # Backend API wrappers (sites, uploads, auth)
│   │   ├── safe/              # Result pattern utilities (safeFetch, tryCatch)
│   │   └── types.ts           # Shared TypeScript types
│   └── routes/
│       ├── (auth)/            # Login, register
│       ├── (legal)/           # Privacy, terms, cookies, rights
│       ├── dashboard/         # Main site management UI
│       ├── profile/           # User profile page
│       ├── api/               # SvelteKit API routes (+server.ts)
│       ├── +layout.svelte     # Root layout
│       └── +page.svelte       # Landing page
├── svelte.config.js           # Bun adapter, experimental tracing
└── package.json               # Scripts, Prettier, ESLint
```

## WHERE TO LOOK

| Task             | Location                        | Notes                   |
| ---------------- | ------------------------------- | ----------------------- |
| Add UI component | `src/lib/components/*.svelte`   | 13 components available |
| Add page         | `src/routes/**/+page.svelte`    | File-based routing      |
| Server-side data | `src/routes/**/+page.server.ts` | Load functions          |
| API endpoint     | `src/routes/api/**/+server.ts`  | GET/POST handlers       |
| Backend API call | `src/lib/api-client/*.ts`       | Returns `Result<T, E>`  |
| Error handling   | `src/lib/safe/index.ts`         | `safeFetch`, `tryCatch` |
| Types            | `src/lib/types.ts`              | Shared interfaces       |

## CONVENTIONS

### Result Pattern (CRITICAL)

- NEVER use bare `fetch()` - always `safeFetch()`
- API functions return `AsyncResult<T, E> = Promise<Result<T, E>>`
- `Result<T, E> = { data: T, error: null } | { data: null, error: E }`
- Check errors: `if (error) { return { data: null, error: ... } }`

### Components (Svelte 5)

- Props: `interface Props` → destructure with `$props()`
- Children: `children?: any` → render with `{@render children?.()}`
- Bindable refs: `ref = $bindable(null)`
- CSS variables for theming: `--color-primary-90`, `--border-radius-regular`

### API Client Pattern

```typescript
export const createSite = async (name: string): AsyncResult<boolean, string> => {
  const { data: resp, error } = await safeFetch('/api/proxy', { ... });
  if (error) return { data: null, error: error.message };
  if (!resp.ok) {
    const errorData = await resp.json();
    return { data: null, error: errorData.errorMessage };
  }
  return { data: true, error: null };
};
```

### Server Routes

- Pass SvelteKit `fetch` + `request` to API client for SSR
- Example: `getSites(fetch, request)` in `+page.server.ts`
- API proxy: `/api/proxy` with `X-Method` header routes to backend

### Testing Policy

- Do not add or maintain unit tests in this webapp.
- Do not create `*.test.*` or `*.spec.*` files.

## ANTI-PATTERNS

- ❌ Bare `fetch()` → Use `safeFetch()`
- ❌ Throwing exceptions → Return `Result<T, E>`
- ❌ `console.log` for errors → Return error in Result
- ❌ Ignoring `.error` in Result → Always check both paths
- ❌ Direct backend calls from client → Use `/api/proxy` routes
- ❌ Adding unit test files (`*.test.*`, `*.spec.*`)
