# Contribution guidelines

## Testing guidelines

- Always run tests before commiting them to the repository.
- Always run tests from `platform` directory
- Do not use any additional testing frameworks . Obly use `testing` package from standard library.
- When test object that depends on `*sqlx.DB` use inmemory sqlite database
