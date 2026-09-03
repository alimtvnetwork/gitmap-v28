# 26 — Inherited Rules (Inline Reference for Blind-AI Implementers)

**Spec:** `19-main-worker-service`
**Version:** 1.0.0
**Created:** 2026-05-06 (Phase 13.3 — closes audit-08 §2.4 cross-spec inlining gap)
**Status:** Authoritative inline reference. Authoritative *source* for each rule remains the linked external spec — this file is a **paste-ready cheat sheet**, not a fork. On any divergence, the linked source spec wins.

---

## 1. Why this file exists

The dumb-AI audit (`audit/07-…v2.md`, `audit/08-…v3.md`) flagged that several rules `02-spec/19-…/` *inherits by reference* from sibling specs are not resolvable when an implementer is handed only the `02-spec/19-…` folder. This file inlines the **minimum operative subset** of those rules so a literal AI never needs to leave the folder during implementation.

**Five upstream specs are inherited:**

1. `02-spec/03-error-manage/` — generic error handling.
2. `02-spec/04-database-conventions/` — naming, REST headers, timestamp format, PK rules.
3. `02-spec/05-split-db-architecture/` — 6-tier model (Worker uses 4, Main uses 3 per `11-§FU-1`).
4. `02-spec/06-seedable-config-architecture/` — config seeding semantics.
5. `mem://architecture/error-handling`, `mem://architecture/caching-policy`, `mem://style/naming-conventions` — repo-level rules already in Core memory.

---

## 2. Inherited from `02-spec/03-error-manage/`

Operative subset for spec/19 implementers (full text in [`../03-error-manage/`](../03-error-manage/)):

1. **Catch → log → rethrow or handle.** Every `catch` block MUST take one of these three actions. Silent catch is forbidden.
2. **Log level matches severity:** `debug` (verbose internals) | `info` (normal lifecycle) | `warn` (recoverable anomaly) | `error` (failed operation, surfaced) | `fatal` (process-ending).
3. **Log context MUST include explicit file path + operation name** (per `mem://architecture/error-handling`). Format: `{file: "path/to/file.go", op: "Company.Create", err: <raw>}`.
4. **Use the `apperror` package equivalent** in the implementer's stack. Wrap every error crossing a layer boundary; preserve the cause chain.
5. **Never `String(caught)` — pass the raw error object through.** Stringification at the boundary is the loggers' job.
6. **Error codes are catalogued in [`14-error-codes.md`](./14-error-codes.md);** the cross-tier envelope is defined in [`09-error-contract.md`](./09-error-contract.md) §2. The envelope is the sole wire shape.

**Conflict rule:** if any of the above conflicts with `09-error-contract.md`, **`02-spec/03-error-manage/` wins** for the rule itself; the envelope shape is owned by `08-§2`.

---

## 3. Inherited from `02-spec/04-database-conventions/`

### 3.1 Naming (from `02-naming-conventions.md`)

- **PascalCase** for every table, column, index, view, JSON key, type name, and internal ID. No `snake_case`, no `camelCase`. Acronyms fully uppercased (`UserId`, `JwtExpiresAt`, `WorkerJwt`).
- **Primary keys:** `INTEGER PRIMARY KEY AUTOINCREMENT`, named `{TableName}Id` (e.g. `WorkerNodeId`, `CompanyId`). **No UUIDs as PKs.** UUIDs may appear as opaque request-scoped strings (`CorrelationId`, `IdempotencyKey`) but never as a PK.
- **Foreign keys** match the referenced PK name verbatim. If a row holds two FKs to the same target, prefix with role: `UpdatedByUserDirectoryId`, `CreatedByUserDirectoryId`.
- **Timestamps on the wire:** RFC 3339 with `Z` (`"2026-05-06T10:22:14Z"`).
  **Timestamps in DB columns:** epoch seconds, UTC, INTEGER. (Per Rule 7.1 v2; see `04-main-db-schema.md` §banner.)
- **Booleans in DB columns:** INTEGER 0/1 (SQLite has no native bool).
- **Zero-Underscore policy** for identifiers (per `mem://style/naming-conventions`). Underscores allowed only in SQL reserved patterns (`_lib`, `_temp`).

### 3.2 REST API headers (from `07-rest-api-format.md`)

Every Main↔Worker call MUST honour:

| Header | Direction | Required when | Notes |
| --- | --- | --- | --- |
| `X-Correlation-Id` | inbound + echoed | always | ULID (Crockford base32, 26 chars). Server generates if missing on inbound and stamps response. |
| `X-Idempotency-Key` | inbound | POST / PUT / PATCH | UUID v4. Replays inside `MainWorker.Idempotency.KeyTtlSeconds` MUST return the original response (`08-§3.7`). |
| `X-Auth-Action` | response | when re-auth required | Frontend treats this as the **sole** "force re-login" signal. Do NOT infer from status codes. Values: `Reauthenticate`, `RegenerateBackupCodes`. (`08-§3.4`) |
| `Content-Type` | bodies | always on bodies | `application/json; charset=utf-8`. |
| `Authorization: Bearer <jwt>` | inbound | when endpoint row in `06-§2` lists `JWT (on Worker)` | RS256, verified per `12-§7`. |

### 3.3 Schema rules (from `03-schema-design.md` + Rules 10/11/12)

- **Entity / ref tables** MUST have `Description TEXT NULL` (no DEFAULT).
- **Transactional tables** MUST have `Notes TEXT NULL` and `Comments TEXT NULL` (no DEFAULT).
- **Join tables** are exempt.
- All free-text columns are nullable. No DEFAULT empty-string.

### 3.4 Views (from `04-orm-and-views.md`)

- View names prefixed `Vw` (e.g. `VwWorkerNodeCurrentVersion`).
- Read-only by contract; mutations go through tables + ORM.

---

## 4. Inherited from `02-spec/05-split-db-architecture/`

### 4.1 The 6-tier model (operative subset for spec/19)

Full model has 6 tiers: Root, Settings, App, Session, Cache, Document. Spec/19 uses **4 of 6** in v1.0:

| Tier | Used by Main? | Used by Worker? | Holds | Notes |
| --- | :-: | :-: | --- | --- |
| **Root** | ✅ | ✅ | Top-level identity registry | Main: `WorkerNode`, `Company`, `UserDirectory`. Worker: `RootCompany` shard pointers. |
| **Settings** | ✅ | ✅ | Per-tier tunables + bootstrap state | Includes `EndpointAuthSetting` on Main; `WorkerBootstrapState` on Worker. |
| **App** | ❌ | ✅ | Per-company business data | Authoritative `AppUser` (with `PasswordHash`, `TotpSecret`, `BackupCodesHash`) lives here. Main MUST NOT host this tier. |
| **Session** | ✅ | ✅ | Per-session ephemeral state | Includes JWT denylist, in-flight idempotency rows. |
| Cache | (reserved) | (reserved) | — | Unused in v1.0. |
| Document | (reserved) | (reserved) | — | Unused in v1.0. |

Authoritative tier mapping for spec/19: [`12-split-db-tier-reconciliation.md`](./12-split-db-tier-reconciliation.md) §5.

### 4.2 Operative SQLite settings (per tier)

- **WAL mode** is set **once per session** at the App tier by the Worker's bootstrap (`PRAGMA journal_mode=WAL`). Per-envelope WAL toggling is **prohibited** (see `22-§6.4`).
- **Foreign keys ON** (`PRAGMA foreign_keys = ON;`) at every connection open.
- **Busy timeout** ≥ 5000 ms.

### 4.3 Casbin (RBAC store)

The RBAC adapter lives in the App tier per `mem://architecture/split-database`. Spec/19 references this only via JWT claim `roles[]`; the policy resolver itself is out of scope for spec/19.

---

## 5. Inherited from `02-spec/06-seedable-config-architecture/`

### 5.1 Operative semantics

- Config seed file: `config.seed.json`. Categories block holds tunables (e.g. `Categories.MainWorker.Settings.WorkerJwtTtlSeconds`); Tables block holds reference rows.
- **SemVer-merged via GORM at boot.** A category's `Version` controls migration order: bumping = forward migration on next boot.
- **Defaults are authoritative** in the seed; runtime DB rows override only via explicit admin write. A literal AI MUST NOT hard-code defaults that already live in the seed.
- Spec/19's seed payload is in [`16-tunable-constants.md`](./16-tunable-constants.md) §4 (current `MainWorker` category SemVer: **1.5.0**).

### 5.2 Bump rules

- Adding a new key = **minor** bump (`1.x.0 → 1.(x+1).0`).
- Renaming or removing a key = **major** bump (`1.x.y → 2.0.0`) + add migration entry in §5 of `15-`.
- Changing a default value = **patch** bump (`1.x.y → 1.x.(y+1)`); document old → new in §5.

---

## 6. Inheritance precedence (when in doubt)

If two inherited rules disagree, resolve in this order (highest authority first):

1. **Locked decisions D1–D10** (listed in `01-index.md` §3) — never override.
2. The **source spec** named in §2–§5 above.
3. The local file in `02-spec/19-…/` that re-states the rule.
4. This `26-inherited-rules.md` summary.
5. Diagrams in `diagrams/` (non-authoritative projections per their banners).

A literal AI that finds a contradiction MUST stop and surface it via `MAIN-900-01 SpecContradiction` rather than pick. (`13-§5`.)

---

## 7. Cross-references

- `01-index.md` §Inherits — points to this file from the spec entry-point.
- `09-error-contract.md` §1 — names `02-spec/03-error-manage/` as the source of catch/log/rethrow.
- `16-tunable-constants.md` §4 — names `02-spec/06-seedable-config-architecture/` as the source of seed semantics.
- `12-split-db-tier-reconciliation.md` §5 — names `02-spec/05-split-db-architecture/` as the source of the tier model.
- `mem://architecture/error-handling`, `mem://architecture/caching-policy`, `mem://style/naming-conventions` — repo-level Core rules already inlined into project memory.

— end of 25 —
