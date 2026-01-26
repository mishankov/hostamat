# TypeScript/SvelteKit Style

## Formatting

- Use tabs for indentation
- Single quotes for strings
- No trailing commas
- Print width: 100 characters
- Configured in `.prettierrc`

## TypeScript

- Strict mode enabled
- Always add type annotations for function parameters
- Use `allowJs` and `checkJs` for JavaScript files

## Imports

- Use `$lib/` path alias for library imports
- Example: `import { Badge, Button, Card } from '$lib/components';`

## ESLint

- `no-undef` rule disabled (TypeScript handles this)
- Prettier integration configured
- TypeScript ESLint recommended rules enabled
