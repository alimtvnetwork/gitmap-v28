# 01 — Architecture

**Spec:** `19-main-worker-service`
**Version:** 1.1.0

> **Split-DB tier authority (FU-1):** Per [`11-split-db-tier-reconciliation.md`](./11-split-db-tier-reconciliation.md), Main uses **3 tiers** (Root / Settings / Session); Worker uses **4 tiers** (Root / Settings / App / Session). Both are subsets of `spec/05-split-db-architecture/`'s 6-tier model. Cache + Document tiers are unused in v1.0. Any "3-tier (Root/App/Session)" wording in older docs is stale — the reconciliation file wins.

---

## 1. Topology

```
                ┌────────────────────────────┐
                │      React Frontend        │
                │   (served by Main Server)  │
                └────────────┬───────────────┘
                             │ HTTPS
                ┌────────────▼───────────────┐
                │       Main Server          │
                │  • Serves UI + edge API    │
                │  • Routes business calls   │
                │  • Auth, 2FA, session      │
                │  • Tracks worker registry  │
                │  • SQLite (thin catalog)   │
                └─────┬────────────┬─────────┘
                      │            │
            JWT/OAuth │            │ JWT/OAuth
                      │            │
       ┌──────────────▼───┐  ┌─────▼────────────┐
       │   Worker Node 1  │  │   Worker Node N  │
       │ • Business logic │  │ • Business logic │
       │ • Split-DB       │  │ • Split-DB       │
       │ • No UI          │  │ • No UI          │
       └──────────────────┘  └──────────────────┘
```

ASCII above is a quick reference. Authoritative diagram: `diagrams/seq-login-routing.mmd` and `diagrams/erd-main-db.mmd`. The author's full mindmap is below; per-tenant subdomain layout is in image 03.

![EndPointService — full overview](./images/04-endpoint-service-full-overview.png)

![Per-tenant `wN.<domain>` subdomain routing](./images/03-worker-subdomain-routing.png)

---

## 2. Boundaries (what runs where)

| Concern | Main | Worker |
|---------|------|--------|
| React UI hosting | ✅ | ❌ |
| Edge REST endpoints (`/API/V1/Company` first hit) | ✅ | ❌ |
| Business logic execution | ❌ | ✅ |
| Tenant business data | ❌ (only minimal identity) | ✅ (full, in split-DB) |
| Worker registry | ✅ | ❌ |
| Tenant→Worker mapping | ✅ | ❌ |
| Auth, 2FA, session, sign-up, sign-in | ✅ (both) | ✅ (both) |
| JWT issuance to React | ✅ | ✅ (after Main resolves worker) |
| Push-update orchestration | ✅ | ❌ |
| Self-update receiver | ✅ (own) | ✅ (own) |
| Version reporting | ✅ (aggregate) | ✅ (own) |

---

## 3. Request Lifecycle

### 3.1 New Company creation

1. React → Main: `POST /API/V1/Company` (authenticated session).
2. Main validates payload, runs worker-selection strategy (`04-worker-routing.md`).
3. Main writes minimal identity row (`CompanyId`, `CompanySlug`, `CompanyName`, `WorkerNodeId`) into main DB.
4. Main forwards full payload to chosen Worker.
5. Worker creates the company in its split-DB (per `spec/05-split-db-architecture/`).
6. Worker returns success → Main returns 201 to React.

Authoritative: `diagrams/seq-company-creation.mmd`.

### 3.2 Subsequent business requests

1. React → Main: first request after login.
2. Main resolves `CompanyId → WorkerNodeId` (cache hit if recent).
3. Main returns the worker endpoint + a worker-scoped JWT to React.
4. React → Worker: all further dashboard/data requests go directly.
5. Worker validates JWT, executes business logic, returns data.

Authoritative: `diagrams/seq-login-routing.mmd`.

### 3.3 Push update

1. Power Admin → Main: `POST /API/V1/Workers/{WorkerNodeId}/Update` (or `/All/Update`).
2. Main calls Worker's `/SelfUpdate` endpoint with auth.
3. Worker performs self-update (mechanism is pointer-only — see `09-self-update-pointer.md`).
4. Worker reports new version → Main updates `WorkerVersion` row.

Authoritative: `diagrams/seq-push-update.mmd`.

---

## 4. Main↔Worker Communication Contract

| Aspect | Rule |
|--------|------|
| Transport | HTTPS only |
| Auth | JWT or OAuth (configurable). See `05-auth-and-2fa.md` |
| Correlation | Every cross-tier request carries `X-Correlation-Id` header (UUID v4 string used as opaque ID — NOT a DB key) |
| Idempotency | Mutating cross-tier calls carry `X-Idempotency-Key` |
| Timeout | Default per `15-tunable-constants.md` §2.5 (`MainWorker.Routing.HttpTimeoutSeconds`); per-endpoint overrides via Seedable-Config |
| Retry | Main retries Worker calls on 5xx per `15-tunable-constants.md` §2.1 (`MainWorker.Retry.MaxAttempts`, `RetryBackoffSeconds`, `RetryJitterPct`). NEVER retry on 4xx |
| Failure surfacing | See `08-error-contract.md` |

Note: correlation IDs and idempotency keys are opaque request-scoped strings, not database primary keys. The "no UUIDs" rule in `spec/04-database-conventions/` applies to PKs only.

---

## 5. Caching Boundaries

| Cache | Scope | TTL | Invalidation |
|-------|-------|-----|--------------|
| `CompanyId → WorkerNodeId` | Main process / session | Per `MainWorker.Cache.CompanyToWorkerTtlSeconds` (default 900 s / 15 m) — see `15-tunable-constants.md` §4.2 | On worker reassignment |
| Worker registry | Main process | Per `MainWorker.Cache.WorkerRegistryTtlSeconds` (default 60 s) — see `15-tunable-constants.md` §4.2 | On worker register/deregister |
| Per-user recent-company | User session | Per `MainWorker.Cache.RecentCompanyPerUserTtlSeconds` (= `MainWorker.Auth.MainSessionTtlSeconds`) — see `15-tunable-constants.md` §4.2 | On logout |

Caching policy (explicit TTL, deterministic keys, invalidate on mutation) is the local rule for this spec. The repo-level memory `mem://architecture/caching-policy` is informational; the binding values live in `15-tunable-constants.md` §4.2 (closes audit M-2 / M-3).

---

## 6. Configuration

Both Main and Worker consume **Seedable-Config** (per `spec/06-seedable-config-architecture/`). For terminology mapping (including replaced legacy terms), see `02-glossary.md` §Reserved.

Config keys consumed by this spec live under namespace `MainWorker.*`. Full key list in `06-core-api-endpoints.md` §Config.

---

## 7. Compliance Hooks

- All endpoint handlers obey `.lovable/coding-guidelines/coding-guidelines.md` (≤15 line functions, zero nested `if`, positive booleans, no `any`).
- All access checks use the `User has access to {EnumPage}` pattern (`07-role-based-dashboards.md`), never `if user.role == 'admin'`.
- Every `catch` logs per `spec/03-error-manage/` and the inline contract in `08-error-contract.md`.

---

*Architecture v1.0.0 — 2026-05-04*
