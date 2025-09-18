# Contribution guidelines

## Testing guidelines

- Always run tests before commiting them to the repository.
- Always run tests with command `cd platform && go test <test files>`
- Do not use any additional testing frameworks . Only use `testing` package from standard library.
- When test object that depends on `*sqlx.DB` use inmemory sqlite database
- Repository tests should only test database interaction logic.
- Service tests should only test business logic.
- Handlers tests should only test http (or other protocol) interaction logic.
- Do not test logging
