# 17 — Cascading Roles & Role-Access Cache Bin

**Spec:** `19-main-worker-service`
**Version:** 1.1.0
**Created:** 2026-05-06
**Updated:** 2026-07-19 (v1.1.0: OQ-A1 promoted to D15, OQ-A2 promoted to D16; §7 rewritten as "Resolved decisions" and both items removed from `.lovable/plan.md` Open Questions.)
**Status:** Authoritative (Phase 5, all open questions in scope now closed)
**Resolves:** locked decisions D11 (cascading = union), D12 (cache-bin in ER diagram), **D15** (simple union is final, no role hierarchy), **D16** (per-process SQLite `:memory:` is final for the cache bin storage tier; Redis and in-process map remain configurable alternatives per §4).

> **Phase 5 scope.** Define how a user that holds **multiple roles simultaneously** resolves to a single effective `AccessItem` set, where that resolution happens, and how the result is cached safely under a credential-blind Main + authoritative Worker split (Phase 3).
>
> **Phase 5 invariant.** The catalog (`Role`, `AccessItem`, `RoleAccessItem`) lives on **Main** (Settings tier). Per-user role assignments (`AppUserRole`) live on the **Worker** (App tier, `11-split-db-tier-reconciliation.md` §5). Effective access is computed on the Worker and cached on the Worker. Main never sees a per-user effective set; Main only **broadcasts invalidation** when the catalog mutates.

---

## 1. The Cascading Rule (D11)

A user MAY hold N roles (N ≥ 1). The user's **effective AccessItem set** is the **union** of every grant attached to any of those roles:

```
EffectiveAccess(User)
  = ⋃ RoleAccessItem(R) for each R in AppUserRole(User)
```

Per `(AccessItemId)`, the effective `(CanRead, CanWrite)` pair is the **bitwise OR** across roles:

```
EffectiveCanRead (User, AccessItem)  = ANY(role.CanRead  = 1)
EffectiveCanWrite(User, AccessItem)  = ANY(role.CanWrite = 1)
```

**Examples** (using the seed in `14-rbac-and-status-seed.md`):

| User holds | Effective set |
|---|---|
| `PowerAdmin` only | All 9 AccessItems, full read+write where granted to PowerAdmin. |
| `Member` only | Whatever Member is granted in the seed (typically Dashboard + own profile). |
| `AdminUser` + `Member` | Union of the two — `Member` adds nothing new since AdminUser ⊇ Member, but the union is computed unconditionally. |
| `AdminUser` + custom `BillingViewer` | AdminUser's set ∪ BillingViewer's set. If both grant `BillingPage` and one is read-only and the other read-write, the OR yields read-write. |

**No role hierarchy.** PowerAdmin does NOT implicitly imply AdminUser via inheritance. PowerAdmin happens to be granted everything **because the seed grants it everything**, not because of an inheritance edge. This keeps the model trivially auditable: every effective grant traces to exactly one row in `RoleAccessItem`. (Hierarchy is the rejected alternative — see §7 OQ-A1.)

**Empty role set.** A user with zero `AppUserRole` rows has an empty effective set and MUST be treated as `AccessDenied` for every gated route. Seed migration MUST refuse to delete a role if it would leave any user with zero roles, unless that user is also being deleted in the same transaction.

---

## 2. Where Resolution Happens

| Step | Node | Action |
|---|---|---|
| 1 | Worker (JWT mint) | On `/Auth/InternalSignIn` success, the Worker reads `AppUserRole` for the user, computes the union against its **mirrored** copy of `Role` / `AccessItem` / `RoleAccessItem` (Settings-tier mirror per `11-…` §4), and embeds the resolved `AccessItem.Code[]` (read set) and a separate `AccessItemWrite.Code[]` (write set) into the JWT body. |
| 2 | Worker (per request) | Route guard intersects the request's matched `AccessItem.Code` against the JWT's read/write list. Zero match → `AccessDenied` per `08-error-contract.md` §3.5 → audit row written via `AccessDenialEvent` on Main. |
| 3 | Main (catalog mutation) | When a Power Admin edits `Role`, `AccessItem`, or `RoleAccessItem`, Main fans out an **invalidation broadcast** (§5). |

The cache bin (§4) sits between steps 1 and 2 to avoid hitting SQLite for every guard check.

---

## 3. Cache Bin — Conceptual Model

The cache bin holds the **compiled per-role access set** keyed by `RoleId`, NOT per user. Per-user resolution is then a cheap N-way union of small bitmaps (N = roles per user, typically 1–3).

| Entry | Key | Value | Lifetime |
|---|---|---|---|
| Compiled role set | `RoleId` | `{ ReadAccessItemIds: int[], WriteAccessItemIds: int[], CatalogVersion, CompiledAt }` | TTL = `MainWorker.RoleCache.TtlSeconds` (default 600 s) **or** until invalidation broadcast received, whichever comes first. |
| Catalog version | `_CatalogVersion` | INTEGER, monotonic | Bumped by Main on any catalog write (Phase 5). Worker compares on fetch; mismatch ⇒ recompile. |

> **Why per-role, not per-user.** Roles are O(10) for the lifetime of the system; users are O(10⁶). Caching per role makes invalidation cheap (drop ≤10 entries on any catalog change) and keeps memory bounded.

---

## 4. Cache Bin — Storage Tier (OQ-A2 default: SQLite `:memory:`)

The cache bin is a **per-process SQLite `:memory:` database** opened on Worker boot, with the following schema. This sits in the Worker's **Cache tier** (`05-split-db-architecture/01-fundamentals.md` Cache tier) and is **never** persisted.

```sql
-- Worker, Cache tier, in-memory only
CREATE TABLE RoleAccessCache (
    RoleAccessCacheId INTEGER PRIMARY KEY AUTOINCREMENT,
    RoleId            INTEGER NOT NULL UNIQUE,
    ReadAccessItemIdsJson  TEXT NOT NULL,   -- JSON array of INTEGER
    WriteAccessItemIdsJson TEXT NOT NULL,   -- JSON array of INTEGER
    CatalogVersion    INTEGER NOT NULL,
    CompiledAt        INTEGER NOT NULL,     -- epoch seconds, UTC (Rule 7.1 v2)
    ExpiresAt         INTEGER NOT NULL,     -- CompiledAt + TtlSeconds
    Description       TEXT NULL              -- Rule 11
);
CREATE INDEX IX_RoleAccessCache_ExpiresAt ON RoleAccessCache(ExpiresAt);

CREATE TABLE RoleCacheCatalogVersion (
    RoleCacheCatalogVersionId INTEGER PRIMARY KEY AUTOINCREMENT,
    CatalogVersion INTEGER NOT NULL,
    UpdatedAt      INTEGER NOT NULL,
    Description    TEXT NULL
);
```

**PK convention.** `{TableName}Id` per Rule 1 (universal). `Code` / `Label` simplification (Rule 13) is N/A — these are not ref tables.

**Lookup procedure** (CODE RED compliant, ≤15 lines, positive guards):

```
function LoadEffectiveAccess(UserId):
    roleIds := ReadAppUserRole(UserId)
    return UnionRoleSets(roleIds)

function UnionRoleSets(roleIds):
    sets := map(roleIds, FetchRoleSet)
    return MergeUnionBitwiseOr(sets)

function FetchRoleSet(roleId):
    cached := RoleAccessCache.Get(roleId)
    if IsCacheFresh(cached) -> return cached
    fresh := RecompileRoleSet(roleId)
    RoleAccessCache.Upsert(fresh)
    return fresh

function IsCacheFresh(entry):
    return entry != null
       AND entry.ExpiresAt > Now()
       AND entry.CatalogVersion == CurrentCatalogVersion()
```

> **Alternative tiers (rejected as default, kept as configurable).** Redis (cross-process invalidation built in, but adds an external dependency); plain in-process map (zero dep but no eviction tooling). Implementers MAY swap the storage by re-implementing the four functions above against the same contract. The contract — not the SQLite — is the spec.

---

## 5. Invalidation Protocol (Main → Worker broadcast)

When a Power Admin mutates the catalog on Main, Main MUST notify every registered Worker so each Worker's cache bin can drop stale entries before serving the next request.

### 5.1 Trigger events on Main

| Mutation | Affected RoleIds |
|---|---|
| `INSERT/UPDATE/DELETE Role` | The single `RoleId` |
| `INSERT/UPDATE/DELETE AccessItem` | **All** roles (broad invalidation — AccessItem changes ripple through every role's compiled set) |
| `INSERT/UPDATE/DELETE RoleAccessItem` | The single `RoleId` referenced |

After any of the above commits, Main:

1. Increments a global `CatalogVersion` counter (stored in Main's `Settings` tier).
2. Calls each Worker's `POST /API/V1/Cache/InvalidateRoleAccess` over the existing S2S channel (`05-auth-and-2fa.md` Worker-internal auth).
3. Audits the broadcast in a new `RoleAccessInvalidationEvent` table on Main (Phase 12 ER update).

### 5.2 Endpoint contract

**Endpoint:** `POST /API/V1/Cache/InvalidateRoleAccess` (Worker-side)

**Auth:** S2S Worker-internal (per `05-auth-and-2fa.md` §7).

**Idempotent:** yes (key = `CatalogVersion`).

**Body:**

```json
{
  "CatalogVersion": 42,
  "InvalidatedRoleIds": [1, 3],
  "InvalidateAll": false,
  "OccurredAt": 1730000000,
  "CorrelationId": "01J9Z..."
}
```

**Semantics:**
- `InvalidateAll = true` ⇒ Worker MUST `DELETE FROM RoleAccessCache;` AND bump its local `CatalogVersion`.
- `InvalidateAll = false` ⇒ Worker MUST `DELETE FROM RoleAccessCache WHERE RoleId IN (...);` AND bump its local `CatalogVersion`.
- A request whose `CatalogVersion` is **older than** the Worker's current value is a **no-op** (out-of-order delivery is safe; the higher version already invalidated).

**Response:** `200 OK` → `{ "AppliedAt": 1730000000, "EvictedCount": 2 }`.

**Failure modes:**
- Worker unreachable → Main retries per `15-tunable-constants.md` §2.1; final failure logs `MAIN-700-01 CacheInvalidationDeliveryFailed` (new error code, see §6) but **does not roll back the Main-side mutation**. The TTL safety net (default 600 s) bounds staleness. <!-- TUNABLE-WAIVER: 600 s = MainWorker.RoleCache.TtlSeconds, see §2.10 -->
- Worker reachable but returns 5xx → same retry policy; same fallback.

### 5.3 JWT staleness

A JWT minted before the invalidation still embeds the old AccessItem list. Two complementary mitigations:

1. **Short JWT TTL** — recommended `MainWorker.Jwt.TtlSeconds` ≤ 900 s.
2. **Catalog-version stamp in JWT** — the JWT body MUST include `CatalogVersion` (the value at mint time). Route guards compare against the Worker's current `CatalogVersion`; on mismatch the guard MUST re-resolve the user's roles and either issue a refresh hint (`401 + ReauthRequired`) or transparently recompute against the live cache. Default behaviour = transparent recompute; strict deployments MAY enforce reauth via `MainWorker.RoleCache.RequireReauthOnCatalogBump = true`.

---

## 6. New Error Codes (additions to `13-error-codes.md`)

| Code | Flat | Name | Message | HTTP | Source |
|---|---|---|---|---|---|
| `MAIN-700-01` | `21171` | `CacheInvalidationDeliveryFailed` | "Worker did not ACK `InvalidateRoleAccess` within retry budget." | 502 | `17-cascading-roles-and-cache-bin.md` §5.2 |
| `WORKER-900-01` | `21090` | `RoleCacheRecompileFailed` | "Worker failed to recompile `RoleAccessCache` after invalidation." | 500 | `17` §4 |
| `WORKER-900-02` | `21091` | `EmptyEffectiveAccessSet` | "User has zero `AppUserRole` rows; access denied to gated route." | 403 | `17` §1 |

These rows MUST be added to `13-error-codes.md` and `error-codes.json` in the Phase 5 commit.

---

## 7. Resolved Decisions (was "Open Questions")

Both OQ-A1 and OQ-A2 were promoted to locked decisions on 2026-07-19 after the defaults below shipped in v1.0.0 and drove Phases 6, 8, 11, and 12 without contradiction.

- **OQ-A1 → D15 (locked).** Cascading semantics is **simple union** per §1. Role hierarchy is rejected: every effective grant must trace to exactly one `RoleAccessItem` row for auditability. Reopening this would require dropping `WORKER-900-02 EmptyEffectiveAccessSet` and every guard that currently reads `AppUserRole` as a flat set, and would invalidate the JWT `AccessItem.Code[]` embedding in §2 step 1.
- **OQ-A2 → D16 (locked).** Cache-bin storage tier is **per-process SQLite `:memory:`** with the `RoleAccessCache` / `RoleCacheCatalogVersion` schema in §4 and the invalidation contract in §5. Redis and plain in-process map remain **configurable alternatives** (implementer swaps the four functions in §4 against the same contract); they are not the default and MUST NOT be assumed by Phase 6+ code. `MainWorker.RoleCache.TtlSeconds` (default 600 s) and `MainWorker.RoleCache.RequireReauthOnCatalogBump` (default `false`) in `15-tunable-constants.md` are part of this decision.
- **OQ-A3: Backup zip password derivation.** ✅ Resolved in Phase 8 (`20-backup-encryption-and-keys.md` §2.12).
- **OQ-A4: Snapshot retention.** ✅ Resolved in Phase 11 (`MainWorker.Backup.SnapshotRetentionDays = 30`).

---

## 8. Cross-References

- `03-main-db-schema.md` §2.6 — `Role`, `AccessItem`, `RoleAccessItem` (catalog).
- `07-role-based-dashboards.md` §5 — access check implementation flow.
- `11-split-db-tier-reconciliation.md` §4 (Main Settings tier) and §5 (Worker App tier).
- `13-error-codes.md` — new MAIN-700 / WORKER-900 series.
- `14-rbac-and-status-seed.md` §3 — concrete row sets the cache compiles from.
- `15-tunable-constants.md` — `MainWorker.RoleCache.TtlSeconds` (new tunable, default 600 s), `MainWorker.RoleCache.RequireReauthOnCatalogBump` (new tunable, default `false`).

---

*Cascading Roles & Cache Bin v1.0.0 — 2026-05-06 (Phase 5)*
