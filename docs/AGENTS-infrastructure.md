# Infrastructure

## Stack

- **Database**: PostgreSQL (production), SQLite (tests)
- **Storage**: MinIO (S3-compatible)
- **Reverse proxy**: Caddy
- **Go version**: 1.25.0
- **Runtime**: Node/Bun for frontend

## Deployment

- Docker multi-platform builds target `linux/amd64`
- Container registry: `hostamat-container-registry.registry.twcstorage.ru`
- Monorepo structure with separate backend and frontend directories
