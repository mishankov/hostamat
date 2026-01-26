# Testing Patterns

## Go Testing

- Always run tests before committing
- Use separate package for tests with `_test` suffix (e.g., `package subdomain_test`)
- Use only the standard `testing` package - no additional testing frameworks
- Do not test internal implementation details
- When testing objects that depend on `*sqlx.DB`, use in-memory SQLite database

### Test Organization

- **Repository tests**: Only test database interaction logic
- **Service tests**: Only test business logic
- **Handler tests**: Only test HTTP/protocol interaction logic
- Do not test logging

### Test File Naming

- Go: `*_test.go` (e.g., `service_test.go`, `repository_test.go`)
- Frontend: `*.test.ts` or `*.spec.ts`

## Running Tests

**Backend:**
```bash
cd platform && go test ./... -v -cover              # All tests
cd platform && go test ./internal/<package> -v -run TestName  # Single test
```

**Frontend:**
```bash
cd webapp && npm run test        # Run all tests once
cd webapp && npm run test:unit   # Run tests in watch mode
```
