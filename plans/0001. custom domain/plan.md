# Plan: Custom domains

> Source PRD: GitHub issue [#119](https://github.com/mishankov/hostamat/issues/119)

## Architectural decisions

Durable decisions that apply across all phases:

- **Routes**:
  The existing dashboard route stays centered at `/dashboard/domains` and becomes the primary management surface for custom domains.
  Backend domain-management endpoints should live under the authenticated app API beside the current site routes and expose list, create, verify, and delete flows.
  Public site serving keeps the current host-based entry path, but request resolution must move from `sites.domain` to a dedicated hostname-mapping lookup that supports both generated and custom hostnames.
- **Schema**:
  Keep the generated Hostamat hostname on the site record as the platform fallback URL.
  Introduce a dedicated domain table for custom hostname lifecycle instead of overloading `sites`.
  Domain records must carry normalized hostname, site reference, owner reference, domain type, status, verification token, last verification result, and timestamps suitable for retries and future background rechecks.
  Hostname uniqueness is global across the platform.
- **Key models**:
  `Site` remains the owner of content and the generated Hostamat hostname.
  A new domain entity represents attachable custom hostnames and their lifecycle.
  The MVP supports multiple custom domains per site, but only verified domains become active.
- **Verification**:
  Verification is DNS-based and requires both ownership proof and routing proof.
  Ownership proof comes from a per-domain TXT token.
  Custom domains are limited to subdomains such as `www.example.com` or `docs.example.com`; apex/root domains such as `example.com` are out of scope.
  Routing proof uses CNAME guidance to the site's generated Hostamat hostname.
  Verification logic should sit behind a service boundary so it can be reused later by schedulers or async jobs.
- **Proxy integration**:
  Caddy configuration is not the source of truth for domain state.
  Activation is a deliberate sync step that adds only verified hostnames to the active Caddy host list.
  Removal detaches the hostname from both persistent state and the Caddy host list without affecting the site, files, or generated Hostamat subdomain.
  Host-list rewrites stay serialized, consistent with the current subdomain flow.
- **Authorization and ownership**:
  Custom domains can only be attached to sites the caller owns.
  Duplicate hostnames and cross-account claims must be rejected before activation.
- **Configuration**:
  Public ingress targets used in DNS instructions and routing validation must come from environment configuration rather than hardcoded values.
- **Testing**:
  Backend automated coverage should carry the feature.
  Focus tests on normalization, uniqueness, authorization, status transitions, verification outcomes, host activation/removal, and public host resolution.
  Do not add frontend unit tests.

---

## Phase 1: Establish domain records and read model

**User stories**: 1, 6, 8, 9, 11, 12, 15, 19, 21, 22

### What to build

Create first-class custom-domain persistence and a service/API surface that lets an authenticated site owner add a hostname to an existing site and list all domains connected to that site. The first slice should support normalized hostname storage, global duplicate rejection, explicit pending status, site ownership checks, and a response shape the dashboard can render without yet activating public traffic on the custom hostname.

### Acceptance criteria

- [ ] A site owner can create a custom-domain record for one of their existing sites and receive a pending domain entry with an ownership-verification token.
- [ ] Hostname normalization is enforced before persistence so case and trailing-dot variants do not create duplicate records.
- [ ] The same hostname cannot be attached to two sites, and callers cannot attach domains to sites they do not own.
- [ ] Listing domains shows both the generated Hostamat hostname and any custom hostnames with explicit status and connected site information.
- [ ] The existing generated Hostamat subdomain continues to work unchanged for public site access.

---

## Phase 2: Expose DNS instructions and manual verification

**User stories**: 2, 3, 4, 5, 6, 7, 13, 14, 18, 21

### What to build

Add manual verification for pending domains. The slice should return DNS instructions for TXT ownership proof and CNAME routing, allow the user to trigger verification on demand, and persist actionable verification results. At the end of this phase the system knows whether a domain is ready for activation, but verified domains do not need to be publicly served yet.

### Acceptance criteria

- [ ] The dashboard and API expose DNS setup instructions for supported subdomain hostnames.
- [ ] Verification checks both ownership TXT data and CNAME routing toward the site's generated Hostamat hostname.
- [ ] A user can retry verification after changing DNS without recreating the domain.
- [ ] Verification failures are stored and returned as readable, actionable errors instead of generic failure states.
- [ ] Successful verification transitions the domain into a verified-ready state without breaking the generated Hostamat fallback domain.

---

## Phase 3: Activate verified domains in host resolution and proxy routing

**User stories**: 1, 5, 8, 16, 17, 18, 19, 20

### What to build

Wire verified custom domains into the public serving path. This slice moves host resolution off the single `sites.domain` lookup and into a domain-mapping layer that supports generated and verified custom hostnames. Verification success should trigger controlled Caddy host activation, and removal should safely detach the hostname from both request routing and proxy configuration.

### Acceptance criteria

- [ ] Public site serving resolves both generated Hostamat hostnames and verified custom hostnames to the same site content.
- [ ] Only verified domains are added to the active Caddy host list.
- [ ] Caddy sync failures do not silently mark a domain active in persistent state.
- [ ] Removing a custom domain stops serving the site from that hostname while leaving the site, files, and generated Hostamat hostname intact.
- [ ] Host-list updates remain serialized so one activation or removal cannot corrupt other active host mappings.

---

## Phase 4: Replace the placeholder dashboard with full domain management

**User stories**: 1, 3, 6, 7, 9, 10, 13, 14, 15, 22

### What to build

Turn `/dashboard/domains` from an informational placeholder into the full control surface for the feature. The page should let users list platform and custom domains, add a hostname to a selected site, inspect DNS instructions, see verification status and errors, retry verification, and remove a custom domain. The navigation should stop advertising the page as "Soon" once the feature is usable.

### Acceptance criteria

- [ ] `/dashboard/domains` provides add, list, verify, and remove flows for custom hostnames.
- [ ] Each domain entry clearly shows its status, the connected site, and any next action required from the user.
- [ ] Users can view the ownership token and DNS records needed for the selected hostname shape.
- [ ] Users can trigger verification retries and see updated status without leaving the domains area.
- [ ] The app navigation no longer labels the domains section as upcoming-only placeholder functionality.

---

## Phase 5: Hardening, recovery, and release validation

**User stories**: 6, 7, 11, 12, 17, 18, 19, 20, 21

### What to build

Finish the MVP by tightening edge cases, observability, and automated coverage around the new domain lifecycle. This phase is not about new end-user flows; it is about proving the system behaves predictably across duplicates, authorization failures, partial verification, activation errors, detach flows, and restart recovery.

### Acceptance criteria

- [ ] Automated backend tests cover domain creation rules, normalization, uniqueness, authorization, verification outcomes, activation/removal behavior, and host resolution for generated and custom domains.
- [ ] Repository migration tests cover the new domain persistence schema.
- [ ] Operator-visible and user-visible errors distinguish validation problems, DNS misconfiguration, ownership-proof failures, and proxy-sync failures.
- [ ] Persistent domain state is sufficient for the dashboard and backend to recover accurately after restart without inferring truth from Caddy alone.
- [ ] Environment-driven ingress targets are documented and validated for local and deployed environments.
