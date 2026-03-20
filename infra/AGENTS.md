# INFRA - DEPLOYMENT KNOWLEDGE BASE

**Generated:** 2026-01-27  
**Focus:** Docker infrastructure, local dev and production deployment

## OVERVIEW

Docker-based infrastructure for local development (compose) and production (Docker Swarm).

## STRUCTURE

```
infra/
├── test/              # Local development environment
│   ├── compose.yml    # All services for local dev
│   ├── config/        # Service configs: proxymini, Caddyfile
│   └── .env           # Backend/webapp env vars
└── prod/              # Production Docker Swarm deployment
    ├── compose.yml    # Swarm stack definition
    ├── config/        # Configs: Caddyfile, Gatus, Prometheus, Dozzle
    └── .env           # Production env vars
```

## WHERE TO LOOK

| Task | Location | Notes |
|------|----------|-------|
| Local service ports | `test/compose.yml` | All exposed ports for dev |
| Dev proxy routing | `test/config/proxymini.conf.toml` | Routes `/backend`, `/caddy`, `/minio` |
| Production stack | `prod/compose.yml` | Node placement, replicas, external configs |
| Caddy reverse proxy | `prod/config/Caddyfile` | HTTPS, domains, routes |
| Health checks | `prod/config/gatus/config.yaml` | Gatus uptime monitoring |
| Metrics scraping | `prod/config/prometheus.yaml` | Prometheus targets |

## KEY SERVICES

### Local Development (test/)

| Service | Purpose | Port | Notes |
|---------|---------|------|-------|
| proxymini | Dev reverse proxy | 1000 | Routes to backend/caddy/minio |
| caddy | HTTPS proxy | 80, 443, 2019 | Admin API on 2019 |
| minio | S3-compatible storage | 9000, 9001 | Console on 9001 |
| db | PostgreSQL | 5432 | User/password: hostamat |
| hostamat | Backend app | 1001 | Hot reload via volume mount |
| dashboard | Ops dashboard | 1002 (→4444) | Hot reload via volume mount |

### Production Swarm (prod/)

| Service | Purpose | Port | Node Type |
|---------|---------|------|-----------|
| caddy | HTTPS proxy + Let's Encrypt | 80, 443 | webserver |
| hostamat | Backend app | - | app |
| webapp | SvelteKit frontend | - | app |
| portainer | Swarm UI management | 8080, 8000 | support |
| gatus | Uptime monitoring | 8082 | support |
| dashboard | Ops dashboard | 8083 | support |
| prometheus | Metrics collection | 9090 | support |
| grafana | Metrics visualization | 8084 | support |
| dozzle | Log viewer | 8081 | global (all nodes) |
| node_exporter | Host metrics | 9100 | global |
| cadvisor | Container metrics | 9080 | global |

## DEPLOYMENT NOTES

- **Node labels**: `kind=webserver`, `kind=app`, `kind=support` control service placement
- **External configs**: Caddyfile, Gatus, Prometheus configs created via `docker config create`
- **External volumes**: `caddy-data`, `caddy-conf`, `portainer_data`, `prometheus-data` persist data
- **Registry**: Custom registry at `hostamat-container-registry.registry.twcstorage.ru/hostamat/*`
- **Backend runtime**: dashboard images default to `run`; backend images default to `migrate && run` via Docker `ENTRYPOINT`
- **Migration flow**: backend container startup runs migrations before `run`
- **Deploy command**: `./scripts/deploy-stack.sh` refreshes configs and executes `docker stack deploy --compose-file swarm.yml --prune --with-registry-auth hostamat`
