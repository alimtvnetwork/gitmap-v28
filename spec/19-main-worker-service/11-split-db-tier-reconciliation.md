# 11 — Split-DB Tier Reconciliation (Main + Worker)

**Spec:** `19-main-worker-service`
**Version:** 2.0.0
**Created:** 2026-05-04 · **Last applied:** 2026-05-06 (Phase 13.2)
**Status:** ✅ Applied — was a follow-up tracker, now collapsed into a 5-line stub.
**Authority:** This file is the canonical **mapping** between Main/Worker spec and `spec/05-split-db-architecture/`. On any tier-count or tier-name conflict, **`spec/05` wins** and this file translates.

---

## 1. Authoritative tier table (binding for spec/19)

`spec/05-split-db-architecture/` defines **6 tiers** (Root, Settings, App, Session, Cache, Document). Spec/19 uses subsets:

- **Main** uses **3 tiers** — Root, Settings, Session. Main has no App tier (it owns no business data).
- **Worker** uses **4 tiers** — Root, Settings, App, Session.
- **Cache** and **Document** tiers are reserved by `spec/05` for future RAG and are **not used in v1.0**.

Any "3-tier (Root/App/Session)" wording in older drafts is stale; the table above wins.

---

## 2. Per-tier table allocation (Main)

| Table | Tier | Source spec |
|---|---|---|
| `Company` | Root | `19/03-main-db-schema.md` §2.3 |
| `WorkerNode` | Root | `19/03-main-db-schema.md` §2.1 (canonical) |
| `WorkerNodeStatus`, `WorkerNodeKind` | Root | `19/03-main-db-schema.md` §2.2 |
| `UserDirectory` (routing index, no secrets) | Root | `19/03-main-db-schema.md` §2.4 |
| `Role`, `AccessItem`, `RoleAccessItem` | Settings | `19/03-main-db-schema.md` §2.6/§2.6.1/§2.6.2 |
| `EndpointAuthSetting`, `AuthMechanism`, `EndpointAuthSettingMechanism` | Settings | `19/06-core-api-endpoints.md` §5 |
| `EndpointAuthChangeKind` | Settings | `19/03-main-db-schema.md` §2.6.5 |
| `UpdateSchedule` | Settings | `19/06-core-api-endpoints.md` §4 |
| `WorkerSelectionStrategy` | Settings | `19/03-main-db-schema.md` §2.9 |
| `AuthSession` | Session | `19/05-auth-and-2fa.md` |
| `TwoFactorChallenge` | Session | `19/05-auth-and-2fa.md` (relayed; ephemeral) |
| `WorkerVersion` | Root | `19/03-main-db-schema.md` §2.7 |
| `WorkerSelectionEvent` | Root | `19/03-main-db-schema.md` §2.8 |
| `AccessDenialEvent` | Root | `19/03-main-db-schema.md` §2.6.3 |
| `EndpointAuthAuditEvent` | Root | `19/03-main-db-schema.md` §2.6.4 |

> `User` and `UserRole` are **REMOVED from Main as of v2.1.0** — moved to Worker as `AppUser` / `AppUserRole`. See §3 below.

---

## 3. Per-tier table allocation (Worker)

| Table | Tier | Source spec |
|---|---|---|
| `RootCompany`, `RootCompanyStatus`, `RootCompanyContact` | Root | `19/diagrams/erd-worker-split-db.mmd` |
| `AppCompanyShard` | Root | (registry of App-tier DBs) |
| `WorkerBootstrapState` | Settings | `19/10-worker-bootstrap-protocol.md` §9 |
| `WorkerUpdateInstruction` | Settings | `spec/14-update/28-worker-push-instruction.md` §7 |
| `AppUser` (authoritative identity, v2.1.0) | App | `19/05-auth-and-2fa.md` §3, `19/diagrams/erd-worker-split-db.mmd` |
| `AppUserRole` (user→role assignments) | App | `19/14-rbac-and-status-seed.md` §6 |
| `AppBusinessEntity` | App | `19/diagrams/erd-worker-split-db.mmd` |
| `AppSession` | Session | same |

---

## 4. Provisioning order on worker boot

Per `10-worker-bootstrap-protocol.md` §3 step 2 ("Self-test split-DB tiers"), the Worker MUST verify in this exact order:

1. Root tier — open `data/root.db`, verify schema version.
2. Settings tier — verify `WorkerBootstrapState` table exists in Root DB.
3. App tier — for each row in `AppCompanyShard`, verify `data/{CompanySlug}/app.db` exists and migrates.
4. Session tier — directory `data/{CompanySlug}/sessions/` exists and is writable.

Failure at any step → exit `WORKER-000-02 SPLIT_DB_TIER_MISSING` per `10/§6`.

---

## 5. Diagram banner template

Every Mermaid file under `spec/19/diagrams/` that touches DB tiers MUST include:

```
%% Authority: spec/05-split-db-architecture/ defines tier semantics.
%% Tier mapping for this spec: spec/19-main-worker-service/11-split-db-tier-reconciliation.md
%% This diagram is a non-authoritative projection. On conflict, prose wins.
```

---

## 6. Follow-up status (all applied — closed in Phase 13.2)

| # | Task | Target | Status |
|---|---|---|---|
| FU-1 | "3-tier" prose replaced by tier-authority block | `19/01-architecture.md` top | ✅ Applied |
| FU-2 | Main `App`-tier refs moved to Root/Settings | `19/03-main-db-schema.md` §2 banner | ✅ Applied |
| FU-3 | Settings tier note added | `19/09-self-update-pointer.md` top | ✅ Applied |
| FU-4 | Settings box + banner on worker ERD | `19/diagrams/erd-worker-split-db.mmd` | ✅ Applied (banner present; Settings tables documented in §2/§3 above rather than redrawn) |
| FU-5 | `WorkerUpdateInstruction` Settings tier | `spec/14-update/28-…md` §7 | ⏭ Out of scope for spec/19 — tracked in spec/14 |
| FU-6 | Back-link from `spec/05` | `spec/05-split-db-architecture/01-fundamentals.md` | ⏭ Out of scope for spec/19 — tracked in spec/05 |

---

## 7. Cross-references

- `spec/05-split-db-architecture/01-fundamentals.md` — source of truth for all tier semantics.
- `spec/19-main-worker-service/03-main-db-schema.md` — Main-side schema, post-FU-2.
- `spec/19-main-worker-service/10-worker-bootstrap-protocol.md` — uses §4 provisioning order.
- `spec/06-seedable-config-architecture/` — seeds Settings-tier rows.

---

*Split-DB tier reconciliation v2.0.0 — 2026-05-06 (Phase 13.2: §3/§8 follow-ups applied; file collapsed from tracker → applied stub).*
