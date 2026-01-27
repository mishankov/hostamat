# SITE DOMAIN

Core domain managing site lifecycle: CRUD, file uploads, subdomain provisioning, public file serving.

## FILES

| File | Purpose |
|------|---------|
| `domain.go` | Domain wiring: 3 handler groups (site, file, upload session) |
| `service.go` | Business logic: auth checks, validation, orchestration |
| `repository.go` | DB layer: CRUD operations on `sites` table via sqlx |
| `model.go` | Core entity + limits: `Site`, max 5 sites/user, 20MB/site |
| `errors.go` | Validation errors + handler error codes |

## HANDLERS

Organized into 3 handler groups in `domain.go`:

### SiteHandlerGroup
| Handler | Method | Purpose |
|---------|--------|---------|
| `handler_create.go` | POST / | Create site with name, auto-provision subdomain |
| `handler_list.go` | GET / | List user's sites (max 5) |
| `handler_update.go` | PUT / | Update site name only |
| `handler_delete.go` | DELETE / | Delete site + subdomain + files (async) |

### FileHandlerGroup
| Handler | Method | Purpose |
|---------|--------|---------|
| `handler_uploadfile.go` | POST / | Upload file to upload session, validate ext/size |

### UploadSessionHandlerGroup
| Handler | Method | Purpose |
|---------|--------|---------|
| `handler_us_create.go` | POST / | Create upload session, wipe old files |
| `handler_us_close.go` | DELETE / | Close session, return uploaded file list |

### FileServeHandler
| Handler | Purpose |
|---------|---------|
| `handler_servefile.go` | Public file serving: resolve domain → site, serve `/path` or `/index.html` |

## DEPENDENCIES

**Service Interfaces** (defined in `service.go`):
- `subdomainService`: Create/delete subdomains (from `subdomain` domain)
- `*file.Service`: Upload/delete/retrieve files from S3
- `*uploadsession.Service`: Manage chunked upload sessions

**Repository Interface**:
- `repository`: CRUD on `sites` table

## KEY PATTERNS

**Access Control**: All mutations check `auth.UserFromContext` ownership via `checkAccess()`
**Upload Flow**: CreateUploadSession → UploadFile (N times) → CloseUploadSession
**File Resolution**: Host header → GetByDomainNoAuth → serve from S3 with fallback to `index.html`
**Limits**: Max 5 sites, 25 char name, 20MB total files per site
**Error Codes**: Structured `handlerError` with typed error codes for frontend consumption
