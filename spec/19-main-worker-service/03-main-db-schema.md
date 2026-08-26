# 03 — Main Server DB Schema

**Spec:** `19-main-worker-service`
**Version:** 2.2.0

> **v2.2.0 (Phase 4 — WorkerNode backup & ordering fields):**
> - `WorkerNode` gains three columns: `Sequence INTEGER NOT NULL` (registry display + RoundRobin order), `IsBackup INTEGER NOT NULL DEFAULT 0` (boolean flag), `BackupOfWorkerNodeId INTEGER NULL` (self-FK → primary worker this node mirrors).
> - **UI label rule:** the user-facing label for a `WorkerNode` row is **"Region"**. The internal code, table name, columns, and API contracts keep the name `Worker` / `WorkerNode` unchanged. Only the rendered string in dashboards, forms, and copy says "Region". See `07-role-based-dashboards.md` §UI labels.
> - **Backup invariant (locked decision D9):** when `IsBackup = 1`, `BackupOfWorkerNodeId` MUST be NOT NULL and MUST reference a row where `IsBackup = 0`. Backup nodes are excluded from every routing strategy (`04-worker-routing.md` §1.4). Backup chains (backup-of-a-backup) are forbidden.
>
> **v2.1.0 (Phase 3 — Move Users off Main):**
> - The legacy Main `User` table is **REMOVED**. Replaced by `UserDirectory` — a routing-only index containing `(UserEmail, CompanyId, WorkerNodeId)` and **no secrets**. Authoritative `AppUser` rows (with `PasswordHash`, `TotpSecret`, `BackupCodesHash`, etc.) now live exclusively in the assigned Worker's split-DB App tier per `11-split-db-tier-reconciliation.md` §5.
> - The legacy Main `UserRole` table is **REMOVED**. Role assignments live on the Worker as `AppUserRole`. Cascading-roles union is computed per-request on the Worker (Phase 5).
> - Audit tables `AccessDenialEvent` and `EndpointAuthAuditEvent` swap their `UserId` FKs for `UserDirectoryId` (nullable / non-nullable respectively) plus a snapshotted email so audit rows survive Worker-side `AppUser` deletion.
> - All TOTP columns (`UserTotpSecret`, `UserTotpEnrolledAt`, `UserTotpBackupCodesHash`) are deleted from Main.
> - Sign-in flow rewritten in `05-auth-and-2fa.md` §2.1, §6: Main forwards credentials to the resolved Worker; **Main never verifies passwords or TOTP codes**.
>
> **v2.0.0 (Phase 2 — DB convention overhaul):**
> - All `*At` timestamp columns are now `INTEGER` (epoch seconds, UTC) per `spec/04-database-conventions/01-naming-conventions.md` Rule 7.1 v2.
> - All ref / enum-like tables now use the canonical `(Id, Code, Label)` shape per Rule 13. The legacy `{Table}Code` / `{Table}Label` column names are **removed** in this spec; readers MUST migrate.
> - `Company` columns renamed: `CompanySlug` → `Slug`, `CompanyName` → `Name`.
>
> **v1.4.0 carryover (Phase 1):** `EnumPage` → `AccessItem`, `RolePageAccess` → `RoleAccessItem`. New column `AccessItem.PageUrlSuffix` is the route-matcher source of truth. Deprecation aliases for the old names remain accepted through v1.5.0 per `98-changelog.md`.
>
> **DB:** SQLite (default). Same schema portable to PostgreSQL/MySQL.

> **Split-DB tier authority (FU-2):** Main uses only **3 tiers** — Root, Settings, Session — per [`11-split-db-tier-reconciliation.md`](./11-split-db-tier-reconciliation.md) §4. **Main has no App tier** (it owns no business data). Any prior reference placing Main tables in an App tier is a bug; per the reconciliation file, such tables belong in Root or Settings. Tier assignments per table are listed in `11-…` §4.

---

## 1. Principles

- **Thin catalog only.** Main DB stores routing metadata, not business data.
- **PascalCase** for tables, columns, JSON keys (per `spec/04-database-conventions/`).
- **PKs:** `{TableName}Id INTEGER PRIMARY KEY AUTOINCREMENT`. **No UUIDs.**
- `Type` / `Status` / `Category` / `Kind` columns → join tables, never inline strings.
- Entity/ref tables include `Description TEXT NULL`. Transactional tables include `Notes TEXT NULL` and `Comments TEXT NULL`. All nullable, no DEFAULT (memory rules 10/11/12).
- Join tables exempt from the description/notes rule.

ERD: `diagrams/erd-main-db.mmd`.

---

## 2. Tables

### 2.1 `WorkerNode` (entity)

| Column | Type | Null | Notes |
|--------|------|------|-------|
| `WorkerNodeId` | INTEGER | NO | PK, AUTOINCREMENT |
| `WorkerNodeTitle` | TEXT | NO | Human label, may repeat across nodes. **UI renders this under the column header "Region"** (D7). |
| `WorkerNodeIdentity` | TEXT | NO | Unique stable identifier (e.g. machine fingerprint) |
| `WorkerNodeEndpoint` | TEXT | NO | Base URL, e.g. `https://w1.example.com` |
| `WorkerNodeStatusId` | INTEGER | NO | FK → `WorkerNodeStatus.WorkerNodeStatusId` |
| `WorkerNodeKindId` | INTEGER | NO | FK → `WorkerNodeKind.WorkerNodeKindId` |
| `Sequence` | INTEGER | NO | Monotonic display + RoundRobin order. Unique among non-backup peers. Backups inherit their primary's `Sequence` for grouping but never participate in RoundRobin (Phase 4, D6). |
| `IsBackup` | INTEGER | NO | Boolean (0/1). `1` ⇒ this row is a backup mirror of `BackupOfWorkerNodeId`. Backup rows MUST NOT serve traffic (D9). DEFAULT 0 (Phase 4, D8). |
| `BackupOfWorkerNodeId` | INTEGER | YES | Self-FK → `WorkerNode.WorkerNodeId`. NOT NULL when `IsBackup = 1`; MUST be NULL when `IsBackup = 0`. The referenced row MUST have `IsBackup = 0` (no chains). |
| `WorkerNodeRegisteredAt` | INTEGER | NO | Epoch seconds, UTC |
| `WorkerNodeLastSeenAt` | INTEGER | NO | Epoch seconds, UTC; updated on heartbeat |
| `Description` | TEXT | YES | Per Rule 11 |

Unique: `(WorkerNodeIdentity)`, `(Sequence) WHERE IsBackup = 0`.

CHECK constraints (enforced at migration time):
- `(IsBackup = 0 AND BackupOfWorkerNodeId IS NULL) OR (IsBackup = 1 AND BackupOfWorkerNodeId IS NOT NULL)`
- A row referenced by `BackupOfWorkerNodeId` MUST itself have `IsBackup = 0` (validated by trigger; SQLite cannot express this in pure CHECK).

### 2.2 `WorkerNodeStatus` (ref) and `WorkerNodeKind` (ref)

Both follow the canonical Rule 13 ref shape — PK keeps the full `{TableName}Id` form (Rule 1, universal); only `Code` / `Label` / `Description` drop the prefix:

**`WorkerNodeStatus`:**

| Column | Type | Null | Notes |
|--------|------|------|-------|
| `WorkerNodeStatusId` | INTEGER | NO | PK, AUTOINCREMENT (Rule 1) |
| `Code` | TEXT | NO | Unique, e.g. `Active`, `Draining`, `Offline` |
| `Label` | TEXT | NO | Human-readable |
| `Description` | TEXT | YES | Per Rule 10 |

**`WorkerNodeKind`:**

| Column | Type | Null | Notes |
|--------|------|------|-------|
| `WorkerNodeKindId` | INTEGER | NO | PK, AUTOINCREMENT (Rule 1) |
| `Code` | TEXT | NO | Unique, e.g. `Standard`, `HighMemory`, `Reserved` |
| `Label` | TEXT | NO | Human-readable |
| `Description` | TEXT | YES | Per Rule 10 |

Seed values via Seedable-Config. Statuses: `Active`, `Draining`, `Offline`, `Quarantined`. Kinds: `Standard`, `HighMemory`, `Reserved` (extensible).

### 2.3 `Company` (entity, MINIMAL identity only)

> ⚠ Full company data lives in the assigned Worker's split-DB. Main stores only what's needed to route.

| Column | Type | Null | Notes |
|--------|------|------|-------|
| `CompanyId` | INTEGER | NO | PK |
| `Slug` | TEXT | NO | Unique, URL-safe (renamed from `CompanySlug` in v2.0.0) |
| `Name` | TEXT | NO | Display name (renamed from `CompanyName` in v2.0.0) |
| `WorkerNodeId` | INTEGER | NO | FK → `WorkerNode` |
| `CompanyAssignedAt` | INTEGER | NO | Epoch seconds, UTC |
| `Description` | TEXT | YES |

Unique: `(Slug)`. Seedable-Config aliases `CompanySlug → Slug` / `CompanyName → Name` accepted through v2.1.0 then removed.

### 2.4 `UserDirectory` (lookup index — auth credentials live on Worker)

> 🔒 **Phase 3 — v2.1.0 (REMOVED from Main):** Per locked decision D5, the legacy Main `User`, `UserRole`, and all TOTP columns are **deleted from Main**. All identity, password, and 2FA state now live in the assigned **Worker's split-DB App tier** (table `AppUser` per `11-split-db-tier-reconciliation.md` §5; full schema in `spec/05-split-db-architecture/`). Main retains **only the routing pointer** below so an inbound `/Auth/SignIn` POST can be forwarded to the correct Worker without leaking credentials through Main.

`UserDirectory` is a **routing-only index**. It contains no secrets and no PII beyond email. The Worker — not Main — is the authoritative identity store.

| Column | Type | Null | Notes |
|--------|------|------|-------|
| `UserDirectoryId` | INTEGER | NO | PK, AUTOINCREMENT |
| `UserEmail` | TEXT | NO | Unique. Lower-cased before insert. |
| `CompanyId` | INTEGER | NO | FK → `Company`. Used to look up the assigned WorkerNode. |
| `WorkerNodeId` | INTEGER | NO | FK → `WorkerNode`. Denormalized from `Company.WorkerNodeId` so a sign-in lookup is a single-row read. Updated whenever `Company.WorkerNodeId` changes (cascade per §5). |
| `CreatedAt` | INTEGER | NO | Epoch seconds, UTC. |
| `LastSeenAt` | INTEGER | YES | Epoch seconds, UTC. Updated by Main on each successful forward. NULL before first sign-in. |
| `Description` | TEXT | YES | Per Rule 11. |

Unique: `(UserEmail)`. Indexed on `(CompanyId)` and `(WorkerNodeId)`.

> **What is NOT here (deliberately):** no `UserPasswordHash`, no `UserPasswordSalt`, no `UserTotpSecret`, no `UserTotpEnrolledAt`, no `UserTotpBackupCodesHash`, no role assignments. **Main MUST NOT verify passwords or TOTP codes.** It only resolves `email → WorkerNode` and proxies the credentialed body to the Worker. See `05-auth-and-2fa.md` §2.1 for the new flow.

### 2.5 `UserRole` — REMOVED in v2.1.0

Removed. Role assignments now live on the Worker as `AppUserRole` (Worker App tier). Cascading-roles union semantics defined in Phase 5 (`14-rbac-and-status-seed.md` §6) and remain a Worker-side computation. Main never reads role assignments.

### 2.6 `Role` (ref)

| Column | Type | Null |
|--------|------|------|
| `RoleId` | INTEGER | NO (PK) |
| `Code` | TEXT | NO (unique, e.g. `PowerAdmin`, `AdminUser`, `Member`) |
| `Label` | TEXT | NO |
| `Description` | TEXT | YES |

Seeded via Seedable-Config.

### 2.6.1 `AccessItem` (ref) — renamed from `EnumPage` in v1.4.0

Defines the closed set of access-controlled dashboard surfaces. Single source of truth for what RBAC governs. Replaces the older name `EnumPage`; old name is **deprecated** but accepted for one release as an alias in seed loaders.

| Column | Type | Null | Notes |
|--------|------|------|-------|
| `AccessItemId` | INTEGER | NO | PK, AUTOINCREMENT |
| `Code` | TEXT | NO | Unique PascalCase identifier (e.g. `WorkerRegistry`, `PushUpdate`, `UserManagement`). Used by code as a stable enum value. |
| `Label` | TEXT | NO | Human-readable UI string (e.g. `Worker Registry`). |
| `PageUrlSuffix` | TEXT | YES | Route matcher — the trailing path fragment used to associate a request with this AccessItem (e.g. `/admin/workers`, `/billing`). NULL for non-route capabilities (e.g. background actions). |
| `Description` | TEXT | YES | Per Rule 11. |

Unique: `(Code)`. Indexed on `(PageUrlSuffix)` for matcher lookups.

> **Matcher rule.** A request path matches an `AccessItem` when its normalized path **ends with** `PageUrlSuffix`. Suffix matching keeps the table portable across deployments mounted under different base paths. Multiple AccessItems sharing the same suffix is a seed error.

Seeded via Seedable-Config (9 rows enumerated in `14-rbac-and-status-seed.md`).

### 2.6.2 `RoleAccessItem` (join, exempt from Description rule) — renamed from `RolePageAccess` in v1.4.0

Per-role access grant for each `AccessItem`. (Resolves F-A-23 / F-B-10.)

| Column | Type | Null |
|--------|------|------|
| `RoleAccessItemId` | INTEGER | NO (PK) |
| `RoleId` | INTEGER | NO (FK → `Role`) |
| `AccessItemId` | INTEGER | NO (FK → `AccessItem`) |
| `CanRead` | INTEGER | NO (0/1) |
| `CanWrite` | INTEGER | NO (0/1) |

Unique: `(RoleId, AccessItemId)`. Seeded with 19 rows per `14-rbac-and-status-seed.md`.

### 2.6.3 `AccessDenialEvent` (transactional, audit)

Audit row written by Workers on every 403 returned for an `AccessDenied` envelope (per `08-error-contract.md` §3.5 and `07-§8`). (Resolves F-A-17.)

| Column | Type | Null | Notes |
|--------|------|------|------|
| `AccessDenialEventId` | INTEGER | NO (PK) |
| `UserDirectoryId` | INTEGER | YES (FK → `UserDirectory`; NULL when actor is unknown / pre-auth) |
| `ActorEmail` | TEXT | NO | Snapshotted email at time of denial. Survives `UserDirectory` deletion for audit. |
| `AccessItemId` | INTEGER | NO (FK → `AccessItem`) | The Main-side AccessItem catalog row that was denied. The catalog stays on Main per Phase 5; only user-to-role assignments move to Worker. |
| `WorkerNodeId` | INTEGER | YES (FK; the Worker that authoritatively denied. NULL when denied at Main edge) |
| `CorrelationId` | TEXT | NO |
| `OccurredAt` | INTEGER | NO | Epoch seconds, UTC (Rule 7.1 v2). |
| `Notes` | TEXT | YES |
| `Comments` | TEXT | YES |

Indexed on `(UserDirectoryId, OccurredAt)` and `(AccessItemId, OccurredAt)` for audit queries.

> **v2.1.0 change.** `UserId` FK was replaced by `UserDirectoryId` (nullable) plus a snapshotted `ActorEmail`, because authoritative `User` rows now live on the Worker. The Worker forwards the denial event to Main with the resolved `AccessItemId` and the actor's email; Main stores it for cross-Worker audit aggregation.

### 2.6.4 `EndpointAuthAuditEvent` (transactional, audit) — FU-17

Audit row written on every successful `PATCH /API/V1/Settings/EndpointAuth` (per `06-core-api-endpoints.md` §5.6). Sibling shape to `AccessDenialEvent` (§2.6.3). One row per accepted PATCH; idempotent replays (same `X-Idempotency-Key`, same body within `MainWorker.Idempotency.KeyTtlSeconds`) MUST NOT emit a duplicate row. Resolves audit follow-up FU-17.

| Column | Type | Null | Notes |
|--------|------|------|-------|
| `EndpointAuthAuditEventId` | INTEGER | NO | PK, AUTOINCREMENT |
| `EndpointAuthSettingId` | INTEGER | NO | FK → `EndpointAuthSetting.EndpointAuthSettingId` (the row that was written; created or replaced) |
| `EndpointPathPattern` | TEXT | NO | Snapshotted at write time so audit survives later row deletion |
| `HttpMethodMaskOld` | TEXT | YES | NULL when the row was newly created (no prior state) |
| `HttpMethodMaskNew` | TEXT | NO | Post-PATCH value |
| `IsEnabledOld` | INTEGER | YES | NULL on create; `0`/`1` on replace |
| `IsEnabledNew` | INTEGER | NO | Post-PATCH value |
| `OldMechanismsJson` | TEXT | YES | JSON array of prior `AuthMechanismCode[]`; NULL on create |
| `NewMechanismsJson` | TEXT | NO | JSON array of post-PATCH `AuthMechanismCode[]` (sorted ascending for diffability) |
| `ChangeKindId` | INTEGER | NO | FK → `EndpointAuthChangeKind.EndpointAuthChangeKindId` (`Create`, `Replace`, `SoftDisable`, `Reenable`) |
| `UpdatedByUserDirectoryId` | INTEGER | NO | FK → `UserDirectory.UserDirectoryId`. Identifies the Power Admin actor; full identity / role check happened on the Worker that minted the JWT carried by the PATCH. |
| `UpdatedByUserEmail` | TEXT | NO | Snapshotted at write time so the audit row survives later `UserDirectory` deletion. |
| `CorrelationId` | TEXT | NO | Echoes the inbound `X-Correlation-Id` header per `spec/04-database-conventions/06-rest-api-format.md` |
| `IdempotencyKey` | TEXT | NO | The `X-Idempotency-Key` that produced the write. Index supports replay-detection joins. |
| `OccurredAt` | INTEGER | NO | Epoch seconds, UTC; server-stamped, equals the parent row's `UpdatedAt` |
| `Notes` | TEXT | YES | Per Rule 12 |
| `Comments` | TEXT | YES | Per Rule 12 |

Unique: `(IdempotencyKey)` — guarantees the no-duplicate-on-replay invariant above. Indexed on `(EndpointAuthSettingId, OccurredAt DESC)` and `(UpdatedByUserDirectoryId, OccurredAt DESC)` for audit queries.

### 2.6.5 `EndpointAuthChangeKind` (ref)

| Column | Type | Null |
|--------|------|------|
| `EndpointAuthChangeKindId` | INTEGER | NO (PK) |
| `Code` | TEXT | NO (unique: `Create`, `Replace`, `SoftDisable`, `Reenable`) |
| `Label` | TEXT | NO |
| `Description` | TEXT | YES |

Seeded via Seedable-Config — 4 rows. Resolution rules:

- `Create` — no prior `EndpointAuthSetting` row for the pattern existed.
- `Replace` — prior row existed AND (`HttpMethodMask` changed OR `AcceptedMechanisms` set changed) AND `IsEnabled` did not transition.
- `SoftDisable` — `IsEnabled` transitioned `1 → 0`.
- `Reenable` — `IsEnabled` transitioned `0 → 1`.

### 2.7 `WorkerVersion` (transactional)

Tracks which version each Worker is currently running.

| Column | Type | Null | Notes |
|--------|------|------|-------|
| `WorkerVersionId` | INTEGER | NO | PK |
| `WorkerNodeId` | INTEGER | NO | FK |
| `WorkerVersionSemver` | TEXT | NO | e.g. `1.4.2` |
| `WorkerVersionRecordedAt` | INTEGER | NO | Epoch seconds, UTC |
| `Notes` | TEXT | YES | Per Rule 12 |
| `Comments` | TEXT | YES | Per Rule 12 |

### 2.8 `WorkerSelectionEvent` (transactional, audit)

Records every routing decision. Useful for debugging load distribution.

| Column | Type | Null | Notes |
|--------|------|------|------|
| `WorkerSelectionEventId` | INTEGER | NO (PK) |
| `CompanyId` | INTEGER | NO (FK) |
| `WorkerNodeId` | INTEGER | NO (FK) |
| `WorkerSelectionStrategyId` | INTEGER | NO (FK → `WorkerSelectionStrategy`) |
| `WorkerSelectionEventAt` | INTEGER | NO | Epoch seconds, UTC |
| `Notes` | TEXT | YES |
| `Comments` | TEXT | YES |

### 2.9 `WorkerSelectionStrategy` (ref)

| Column | Type | Null |
|--------|------|------|
| `WorkerSelectionStrategyId` | INTEGER | NO (PK) |
| `Code` | TEXT | NO (`RoundRobin`, `LeastLoaded`, `Manual`) |
| `Label` | TEXT | NO |
| `Description` | TEXT | YES |

---

## 3. Indexes

| Index | Columns |
|-------|---------|
| `IX_Company_WorkerNodeId` | `Company(WorkerNodeId)` |
| `IX_WorkerNode_BackupOf` | `WorkerNode(BackupOfWorkerNodeId)` WHERE `BackupOfWorkerNodeId IS NOT NULL` |
| `IX_WorkerNode_PrimaryEligible` | `WorkerNode(WorkerNodeStatusId, Sequence)` WHERE `IsBackup = 0` (covers routing eligibility scans) |
| `IX_UserDirectory_CompanyId` | `UserDirectory(CompanyId)` |
| `IX_UserDirectory_WorkerNodeId` | `UserDirectory(WorkerNodeId)` |
| `UX_UserDirectory_UserEmail` | `UserDirectory(UserEmail)` UNIQUE |
| `IX_WorkerVersion_WorkerNodeId_RecordedAt` | `WorkerVersion(WorkerNodeId, WorkerVersionRecordedAt DESC)` |
| `IX_WorkerSelectionEvent_At` | `WorkerSelectionEvent(WorkerSelectionEventAt DESC)` |
| `UX_EndpointAuthAuditEvent_IdempotencyKey` | `EndpointAuthAuditEvent(IdempotencyKey)` UNIQUE |
| `IX_EndpointAuthAuditEvent_Setting_At` | `EndpointAuthAuditEvent(EndpointAuthSettingId, OccurredAt DESC)` |
| `IX_EndpointAuthAuditEvent_Actor_At` | `EndpointAuthAuditEvent(UpdatedByUserDirectoryId, OccurredAt DESC)` |

---

## 4. What Main DB does NOT store

- Company business fields (address, social media, employee count, etc.)
- **User passwords, password salts, password hashes, or pepper** (Worker only).
- **TOTP secrets, TOTP enrollment timestamps, TOTP backup codes** (Worker only).
- **User → Role assignments** (Worker only as `AppUserRole`).
- User profile data beyond auth-routing (no name, no avatar, no preferences).
- Any per-tenant business state.
- Session bodies (kept in cache or session store, not the catalog).

All of the above belong in the Worker's split-DB per `spec/05-split-db-architecture/`.

> **Phase 3 invariant.** A grep over the Main DB for `password`, `totp`, `secret`, or `hash` MUST return zero column hits. Cache-bin tables (Phase 5) are likewise password-free.

---

## 5. Migration Notes

- Use the implementer's standard migration tool (Laravel migrations for the default stack).
- Migrations are idempotent and forward-only.
- **v2.1.0 migration (Phase 3):** drop tables `User`, `UserRole`; create `UserDirectory`; backfill `UserDirectory` from any existing Main `User` rows (`UserEmail`, `CompanyId`, `WorkerNodeId` resolved via `Company.WorkerNodeId`) **before** the drop; then forward `(UserPasswordHash, UserPasswordSalt, UserTotpSecret, UserTotpEnrolledAt, UserTotpBackupCodesHash, UserRole rows)` to each Worker via the bootstrap protocol (`10-worker-bootstrap-protocol.md`) using a one-shot `MigrateLegacyUsers` instruction. After Worker ACK, delete the rows from Main. The migration script MUST refuse to drop `User` if any `UserDirectory.WorkerNodeId` is unresolved.
- Seed data for `Role`, `WorkerNodeStatus`, `WorkerNodeKind`, `WorkerSelectionStrategy`, `AccessItem`, `RoleAccessItem` ships via Seedable-Config (`spec/06-seedable-config-architecture/`). `Role` and `RoleAccessItem` remain on Main as the **catalog** of available roles (Worker resolves `User → Role → AccessItem` against this catalog at JWT mint time).

---

*Main DB schema v2.1.0 — 2026-05-06 (Phase 3: Users moved to Worker; UserDirectory routing index added)*
