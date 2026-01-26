# Go Backend Style

## Error Handling

- Return errors, don't panic
- Wrap errors with context when appropriate
- Use custom error types/variables for domain-specific errors (e.g., `errUploadSize`, `ErrInvalidFileType`)
- Check `errors.Is()` for error comparison

## Logging

- Use structured logging with `log.InfoContext(ctx, message, key, value, ...)`
- Log significant business events (creation, deletion, errors)
- Don't test logging code

## Context Usage

- Always pass `context.Context` as first parameter
- Use `ctx` as variable name
- Extract user info with `auth.UserFromContext(ctx)`

## Types and Interfaces

- Define interfaces at the consumer level (in the file that uses them)
- Use constructor functions: `NewService(...)` pattern
- Struct fields for dependencies should use interface types when possible

## Common Patterns

**Service methods:**
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
