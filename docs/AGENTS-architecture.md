# Architecture Patterns

## Backend Layers

1. **Handlers**: HTTP request/response handling (`handler_*.go`)
2. **Services**: Business logic, orchestration (`service.go`)
3. **Repositories**: Database operations (`repository.go`)
4. **Domain**: Business models and validation (`domain.go`, `model.go`)

## Dependency Injection

- Services receive dependencies through constructors
- Use interfaces for loose coupling
- Example: `NewService(repository repository, subdomainService subdomainService, ...)`

## Access Control

- Check user access in service layer
- Use `checkAccess()` helper methods
- Extract user from context: `auth.UserFromContext(ctx)`
