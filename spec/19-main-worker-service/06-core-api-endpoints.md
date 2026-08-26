# 06 — Core API Endpoints

**Spec:** `19-main-worker-service`
**Version:** 1.3.0

Authoritative REST surface for both tiers. All paths are `/API/V1/...`. JSON request/response with PascalCase keys.

---

## 1. Conventions

- **Versioning:** path-based (`/API/V1`).
- **Content-Type:** `application/json` request and response.
- **JSON keys:** PascalCase.
- **Errors:** envelope per `08-error-contract.md`.
- **Headers:** `X-Correlation-Id` mandatory inbound (generated server-side if missing); `X-Idempotency-Key` mandatory on POST/PUT/PATCH.
- **Auth column meanings:** `Session` = cookie; `JWT` = worker JWT (RS256); `OAuth` = client-credentials; `None` = public.

---

## 2. Endpoint Catalog

### 2.1 Auth (Main)

| Method | Path | Auth | Purpose |
|--------|------|------|---------|
| POST | `/API/V1/Auth/SignUp` | None | Create user (and optionally company) |
| POST | `/API/V1/Auth/SignIn` | None | Verify password, start session |
| POST | `/API/V1/Auth/Verify2FA` | None (challenge-bound) | Submit TOTP for active challenge |
| POST | `/API/V1/Auth/SignOut` | Session | End current session |
| POST | `/API/V1/Auth/SignOutAll` | Session | End all sessions for user |
| POST | `/API/V1/Auth/PasswordResetRequest` | None | Email reset link |
| POST | `/API/V1/Auth/PasswordResetConfirm` | None (token-bound) | Set new password |
| POST | `/API/V1/Auth/Enroll2FA` | Session | Begin TOTP enrollment |
| POST | `/API/V1/Auth/Confirm2FA` | Session | Confirm TOTP code, enable 2FA |
| POST | `/API/V1/Auth/Disable2FA` | Session + TOTP | Disable 2FA |
| POST | `/API/V1/Auth/RefreshWorkerToken` | Session | Mint a fresh worker JWT |

### 2.2 Company (Main → routes to Worker)

| Method | Path | Auth | Purpose |
|--------|------|------|---------|
| POST | `/API/V1/Company` | Session | Create company (Main routes to Worker) |
| GET | `/API/V1/Company/{CompanySlug}/Resolve` | Session | Returns `WorkerEndpoint` + worker JWT |
| GET | `/API/V1/Company/{CompanySlug}` | JWT (on Worker) | Read company (after resolve) |
| PATCH | `/API/V1/Company/{CompanySlug}` | JWT (on Worker) | Update company |
| DELETE | `/API/V1/Company/{CompanySlug}` | JWT (on Worker) + PowerAdmin | Hard-delete |

### 2.3 User (Main → routes to Worker)

| Method | Path | Auth | Purpose |
|--------|------|------|---------|
| POST | `/API/V1/Users` | Session + access | Create user under current company |
| GET | `/API/V1/Users/{UserId}` | JWT (on Worker) | Read user |
| PATCH | `/API/V1/Users/{UserId}` | JWT (on Worker) + access | Update user |

### 2.4 Status / Version (both tiers, intentionally minimal)

| Method | Path | Auth | Purpose |
|--------|------|------|---------|
| GET | `/API/V1/Status` | None (default; configurable) | Liveness + readiness |
| GET | `/API/V1/Version` | None (default; configurable) | App name, title, current version, update available flag |

### 2.5 Workers (Main only — Power Admin)

| Method | Path | Auth | Purpose |
|--------|------|------|---------|
| GET | `/API/V1/Workers` | Session + PowerAdmin | List worker registry |
| POST | `/API/V1/Workers/Register` | OAuth | Worker registers itself with Main |
| POST | `/API/V1/Workers/{WorkerNodeId}/Heartbeat` | OAuth | Worker liveness ping |
| POST | `/API/V1/Workers/{WorkerNodeId}/Update` | Session + PowerAdmin | Push update to one worker |
| POST | `/API/V1/Workers/All/Update` | Session + PowerAdmin | Push update to all workers |
| POST | `/API/V1/Workers/PublishZip` | Session + PowerAdmin (multipart) | Upload deployment zip via PowerShell |

### 2.6 Self-Update (both tiers — pointer only)

| Method | Path | Auth | Purpose |
|--------|------|------|---------|
| POST | `/API/V1/SelfUpdate` | OAuth | Trigger self-update workflow (see `09-self-update-pointer.md`) |

### 2.7 Settings (Main — Power Admin)

| Method | Path | Auth | Purpose |
|--------|------|------|---------|
| GET | `/API/V1/Settings/EndpointAuth` | Session + PowerAdmin | List per-endpoint auth-mechanism overrides (see §5) |
| PATCH | `/API/V1/Settings/EndpointAuth` | Session + PowerAdmin | Replace one row's mechanism set + IsEnabled (see §5) |
| GET | `/API/V1/Settings/UpdateSchedule` | Session + PowerAdmin | Read schedule |
| PATCH | `/API/V1/Settings/UpdateSchedule` | Session + PowerAdmin | Update schedule (see §4) |

---

## 3. Reference Payloads

> **Wire-format authority (Phase 13.3):** Canonical golden-file fixtures for every endpoint listed in §2 live in [`fixtures/endpoints/`](./fixtures/endpoints/). When this prose and a fixture disagree on **field name, ordering, casing, null-vs-omitted, or timestamp encoding**, the **fixture wins**. Encoding rules are defined once in [`fixtures/conventions.md`](./fixtures/conventions.md). Behavior (algorithm, validation rule, retry policy) remains owned by this file.

### 3.1 `POST /API/V1/Company` request

```json
{
  "CompanyName": "Riseup Asia LLC",
  "CompanyWebsite": "https://riseup-asia.com",
  "CompanySlug": "riseup-asia",
  "Address": "Kuala Lumpur, Malaysia",
  "PhoneNumber": "+60-...",
  "NumberOfPeople": 25,
  "Calendar": "https://cal.com/...",
  "WhatsApp": "+60-...",
  "Facebook": "https://facebook.com/...",
  "LinkedIn": "https://linkedin.com/company/...",
  "PreferredWorkerNodeId": null
}
```

Field nullability (resolves F-A-01 — replaces prior "Most fields Non-Nullable"):

| Field | Nullable | Notes |
|-------|----------|-------|
| `CompanyName` | NO | |
| `CompanyWebsite` | YES | |
| `CompanySlug` | NO | Unique. |
| `Address` | YES | |
| `PhoneNumber` | NO | E.164 preferred. |
| `NumberOfPeople` | YES | Integer ≥ 1 when set. |
| `Calendar` | YES | |
| `WhatsApp` | YES | |
| `Facebook` | YES | |
| `LinkedIn` | YES | |
| `PreferredWorkerNodeId` | YES | Only honored when worker-selection strategy is `Manual` AND caller has `PowerAdmin` access; ignored otherwise (no error). |

### 3.2 `GET /API/V1/Company/{CompanySlug}/Resolve` response

```json
{
  "CompanySlug": "riseup-asia",
  "CompanyId": 42,
  "WorkerNodeId": 3,
  "WorkerEndpoint": "https://w3.example.com",
  "WorkerJwt": "<RS256 token>",
  "JwtExpiresAt": "2026-05-04T12:15:00Z"
}
```

### 3.3 `GET /API/V1/Version` response

```json
{
  "ApplicationName": "MainServer",
  "ApplicationTitle": "Coordinator",
  "CurrentVersion": "1.0.0",
  "UpdateAvailable": false,
  "LatestVersion": null
}
```

---

## 4. Update Schedule Settings (per verbatim §Update Schedule Settings)

`Settings.UpdateSchedule` shape:

```json
{
  "Cadence": "Weekly",
  "EveryNHours": null,
  "SpecificTimeOfDay": "04:00",
  "TimeZone": "Asia/Kuala_Lumpur",
  "Enabled": true
}
```

`Cadence` enum: `Hourly`, `EveryNHours`, `Daily`, `Weekly`, `Monthly`, `Yearly`. When `EveryNHours`, populate `EveryNHours` (allowed values 5, 6, 12, 24). `SpecificTimeOfDay` is null unless `Cadence` ∈ {Daily, Weekly, Monthly, Yearly}. Default: `Weekly @ 04:00 Asia/Kuala_Lumpur`.

---

## 5. Per-Endpoint Auth-Mechanism Overrides (OQ-1 — RESOLVED 2026-05-04)

Resolves OQ-1 (`05-auth-and-2fa.md` §8) and audit finding F-M-10. The Settings.EndpointAuth surface lets a Power Admin tighten or loosen the default auth mechanism on a per-endpoint-pattern basis, within a fixed allow-list. Defaults from §2 always apply when no override row matches.

### 5.1 Schema (canonical — Main DB, Settings tier per `11-split-db-tier-reconciliation.md`)

```sql
CREATE TABLE AuthMechanism (
    AuthMechanismId    INTEGER PRIMARY KEY AUTOINCREMENT,
    AuthMechanismCode  TEXT NOT NULL UNIQUE,   -- Session | Jwt | OAuth | None
    AuthMechanismLabel TEXT NOT NULL,
    Description        TEXT NULL
);

CREATE TABLE EndpointAuthSetting (
    EndpointAuthSettingId       INTEGER PRIMARY KEY AUTOINCREMENT,
    EndpointPathPattern         TEXT NOT NULL UNIQUE,           -- e.g. "/API/V1/Status", "/API/V1/Company/*"
    HttpMethodMask              TEXT NOT NULL,                  -- CSV of GET|POST|PATCH|PUT|DELETE or "*"
    IsEnabled                   INTEGER NOT NULL,               -- 1 = override active, 0 = ignore (fallback to default)
    UpdatedByUserDirectoryId    INTEGER NOT NULL REFERENCES UserDirectory(UserDirectoryId),  -- v2.1.0: legacy `User` table removed; routing index is the canonical actor reference (see `03-main-db-schema.md` §2.4)
    UpdatedByUserEmail          TEXT NOT NULL,                  -- snapshotted at write time so audit/listing survives `UserDirectory` deletion
    UpdatedAt                   INTEGER NOT NULL,               -- Epoch seconds, UTC (Rule 7.1 v2)
    Notes                       TEXT NULL,
    Comments                    TEXT NULL,
    Description                 TEXT NULL
);

CREATE TABLE EndpointAuthSettingMechanism (
    EndpointAuthSettingId INTEGER NOT NULL REFERENCES EndpointAuthSetting(EndpointAuthSettingId) ON DELETE CASCADE,
    AuthMechanismId       INTEGER NOT NULL REFERENCES AuthMechanism(AuthMechanismId),
    PRIMARY KEY (EndpointAuthSettingId, AuthMechanismId)
);
```

Seed rows for `AuthMechanism` come from `14-rbac-and-status-seed.md` §AuthMechanism (4 rows: `Session`, `Jwt`, `OAuth`, `None`).

### 5.2 Resolution algorithm (deterministic)

For each inbound request `(Method, Path)` the auth middleware MUST execute, in order:

1. Look up the FIRST `EndpointAuthSetting` row where `IsEnabled=1` AND `EndpointPathPattern` matches `Path` (longest-prefix wins; exact `/API/V1/Status` beats `/API/V1/*`) AND `HttpMethodMask` contains `Method` (or `*`).
2. If found, build `AcceptedMechanisms` = JOIN through `EndpointAuthSettingMechanism`. Request is authorized when its presented credential matches ANY entry. Empty set is REJECTED at PATCH time (see §5.4 R-3).
3. If no row matches, apply the default from `05-auth-and-2fa.md` §8 table.

A request that fails resolution returns the `AccessDenied` envelope per `08-error-contract.md` §3.4 with `X-Auth-Action: Reauthenticate`.

### 5.3 PATCH semantics

`PATCH /API/V1/Settings/EndpointAuth` is a **whole-row replace, single-row scope** (NOT JSON-Patch RFC 6902, NOT bulk). Each request mutates exactly ONE `EndpointAuthSetting` row, identified by `EndpointPathPattern` (its natural key). This avoids the partial-merge ambiguity flagged in F-M-10.

Request body:

```json
{
  "EndpointPathPattern": "/API/V1/Company/*",
  "HttpMethodMask": "GET,PATCH",
  "IsEnabled": true,
  "AcceptedMechanisms": ["Session", "Jwt"],
  "Notes": "Allow direct React→Worker calls during migration window.",
  "Description": "Tightened from default (Session-only on Main side)."
}
```

Semantics:

| Field | Required | Effect |
|-------|----------|--------|
| `EndpointPathPattern` | YES | Upsert key. Created if absent, replaced if present. |
| `HttpMethodMask` | YES | Replaces stored value verbatim. CSV of `GET|POST|PATCH|PUT|DELETE` or literal `"*"`. |
| `IsEnabled` | YES | `false` retains the row but reverts the endpoint to its §8 default (audit-friendly soft-disable). |
| `AcceptedMechanisms` | YES | Array of `AuthMechanismCode` strings. The set REPLACES the prior `EndpointAuthSettingMechanism` rows transactionally. |
| `Notes` | optional | Free-text. Nullable. |
| `Description` | optional | Free-text. Nullable. |

Response: `200 OK` with the resolved row (post-write) including `EndpointAuthSettingId`, `UpdatedByUserDirectoryId`, `UpdatedByUserEmail`, `UpdatedAt` (server-stamped, epoch seconds UTC). 

**Idempotency:** mandatory `X-Idempotency-Key` per §1. Two PATCHes with the same key + identical body within the TTL window (`MainWorker.Idempotency.KeyTtlSeconds`, see `15-tunable-constants.md`) return the original response without re-writing.

**Atomicity:** the `EndpointAuthSetting` upsert and the `EndpointAuthSettingMechanism` set-replacement MUST execute in a single SQLite transaction. Partial application is forbidden.

### 5.4 Validation rules (server-side, fail-closed)

The server MUST reject the PATCH with `400 ValidationFailed` envelope (per `08-error-contract.md` §3.2) and `Error.FieldErrors[]` populated when ANY of:

| ID | Rule | Field |
|----|------|-------|
| R-1 | `EndpointPathPattern` MUST start with `/API/V1/` and contain only `[A-Za-z0-9/_*{}-]`. | `EndpointPathPattern` |
| R-2 | `HttpMethodMask` MUST be `*` OR a CSV subset of `{GET,POST,PATCH,PUT,DELETE}` with no duplicates. | `HttpMethodMask` |
| R-3 | `AcceptedMechanisms` MUST be non-empty when `IsEnabled=true`. (Empty + enabled would lock the endpoint out entirely — refused.) | `AcceptedMechanisms` |
| R-4 | Every `AcceptedMechanisms` entry MUST resolve to a row in `AuthMechanism`. | `AcceptedMechanisms` |
| R-5 | Patterns matching `/API/V1/Workers/*` or `/API/V1/SelfUpdate` are LOCKED — PATCH returns `403 EndpointAuthLocked` (`MAIN-400-10` / flat `21170`, catalogued in `13-error-codes.md` §3.4 — FU-18 RESOLVED 2026-05-05). These endpoints are always-protected per §8. | `EndpointPathPattern` |
| R-6 | Combining `None` with any other mechanism in the same `AcceptedMechanisms` set is forbidden — `None` is mutually exclusive (otherwise it silently bypasses the others). | `AcceptedMechanisms` |
| R-7 | Caller MUST hold `EnumPage.PowerAdminPage` access; otherwise `403 AccessDenied`. | session |

### 5.5 GET shape

`GET /API/V1/Settings/EndpointAuth` returns:

```json
{
  "Defaults": [
    { "EndpointPathPattern": "/API/V1/Status", "AcceptedMechanisms": ["None"], "Source": "Default" }
  ],
  "Overrides": [
    {
      "EndpointAuthSettingId": 7,
      "EndpointPathPattern": "/API/V1/Company/*",
      "HttpMethodMask": "GET,PATCH",
      "IsEnabled": true,
      "AcceptedMechanisms": ["Session", "Jwt"],
      "UpdatedByUserDirectoryId": 1,
      "UpdatedByUserEmail": "admin@example.com",
      "UpdatedAt": "2026-05-04T08:15:00Z",
      "Notes": null,
      "Description": "Tightened from default."
    }
  ]
}
```

Defaults are derived from `05-auth-and-2fa.md` §8 at runtime (NOT stored), so spec edits there propagate without a migration.

### 5.6 Audit trail (FU-17 — RESOLVED)

Every successful PATCH MUST emit exactly one `EndpointAuthAuditEvent` row (table canonicalised in `03-main-db-schema.md` §2.6.4) inside the same SQLite transaction as the `EndpointAuthSetting` upsert + `EndpointAuthSettingMechanism` set-replacement (per §5.3 *Atomicity*). The transaction is the audit boundary — partial commits are forbidden, and a failed audit insert MUST roll back the entire PATCH.

Field stamping rules:

| Audit column | Source |
|--------------|--------|
| `EndpointAuthSettingId` | The post-upsert PK of the parent row. |
| `EndpointPathPattern` | Snapshotted from the request body (NOT a join) so audit survives parent-row deletion. |
| `HttpMethodMaskOld` / `IsEnabledOld` / `OldMechanismsJson` | Read from the prior row inside the same transaction (`SELECT … FOR UPDATE` semantics; SQLite uses `BEGIN IMMEDIATE`). NULL when `ChangeKind=Create`. |
| `HttpMethodMaskNew` / `IsEnabledNew` / `NewMechanismsJson` | Post-write values. `NewMechanismsJson` MUST be a JSON array sorted ascending by `AuthMechanismCode` so diffs are stable. |
| `ChangeKindId` | Resolved per the rules in `03-main-db-schema.md` §2.6.5 (`Create` / `Replace` / `SoftDisable` / `Reenable`). |
| `UpdatedByUserDirectoryId` | The authenticated Power Admin (same `UserDirectoryId` stamped on the parent row's `UpdatedByUserDirectoryId`). |
| `UpdatedByUserEmail` | Snapshot of the actor's email (same value stamped on the parent row). Survives later `UserDirectory` deletion. |
| `CorrelationId` | The inbound `X-Correlation-Id` header. The middleware MUST reject the PATCH if absent (per `spec/04-database-conventions/06-rest-api-format.md`). |
| `IdempotencyKey` | The inbound `X-Idempotency-Key` header. Unique-indexed (per §2.6.4) — guarantees idempotent replays do not double-write the audit row. |
| `OccurredAt` | Server clock, equals the parent row's `UpdatedAt`. |

Idempotent replay: when the idempotency-cache hit returns the prior response (per §5.3 *Idempotency*), the audit insert MUST be skipped — the existing audit row already records the original write. The unique index on `IdempotencyKey` is the belt-and-braces enforcement.

Observability hook: implementations MUST also log the diff at `INFO` level via the `apperror` package with `OperationId=EndpointAuthChange`, the resolved `ChangeKindCode`, and the `X-Correlation-Id` — for log-pipeline consumers that do not query the DB. The log is a sibling of the audit row, not a substitute.

### 5.7 Cross-references

- Default mechanisms: `05-auth-and-2fa.md` §8
- Lock-list rationale: `05-auth-and-2fa.md` §8 (rows marked "No / always protected")
- ErrorCodes: `13-error-codes.md` §3.4 — `MAIN-400-10 EndpointAuthLocked` (flat `21170`, HTTP 403; FU-18 RESOLVED)
- Header conventions: `spec/04-database-conventions/06-rest-api-format.md` §X-Auth-Action
- Idempotency TTL: `15-tunable-constants.md` §2.3

---

## 6. Rate Limiting (defaults — MUST apply unless overridden via Seedable-Config)

Defaults below are MANDATORY out-of-the-box. Implementations MUST apply them unless explicitly overridden via Seedable-Config keys per `15-tunable-constants.md` §2.6 (resolves F-A-02 — replaces prior "recommended defaults" softening).

| Endpoint group | Default | Seedable-Config key |
|----------------|---------|---------------------|
| `/API/V1/Auth/*` | 10 / minute / IP | `MainWorker.RateAuthPerMinutePerIp` |
| `/API/V1/Workers/*` | 60 / minute / token | `MainWorker.RateWorkerPerMinutePerToken` |
| Other authenticated | 600 / minute / user | `MainWorker.RateOtherPerMinutePerUser` |

Implementer uses framework-native middleware (e.g. Laravel `throttle`). On limit-exceeded, return `WorkerOverloaded` envelope per `08-error-contract.md` §3.6 with HTTP 429 + `Retry-After` header.

---

*Core API endpoints v1.2.0 — 2026-05-05 (FU-17: §5.6 Audit trail wired to `EndpointAuthAuditEvent` (`03-main-db-schema.md` §2.6.4) — same-transaction insert, idempotent-replay skip, ChangeKind resolution, `apperror` sibling log)*

---

## §6 — Backup-Tier Endpoint Catalogue (Phase 12 stub)

**Added:** Phase 12 (Backup-tier consolidation). **Authority for full contract:** `21-backup-endpoints.md`. This section is a **directory pointer** — request/response shapes, error codes, and idempotency rules live in the authoritative file.

All Backup-tier endpoints are **S2S only** (`aud="Backup"`, mandatory `PairingId` claim per `12-jwt-delivery-contract.md` §13). They MUST NOT be exposed via the public proxy chain that serves `/API/V1/Auth/*` and `/API/V1/Company/*`.

### 6.1 Catalogue

| ID | Method | Path | Direction | Scope | Auth | Purpose |
|----|--------|------|-----------|-------|------|---------|
| **BE-1** | POST | `/API/V1/Backup/SyncEnvelope` | Primary → Backup | `Backup.Sync.Write` | S2S `aud=Backup` | Deliver sealed CDC envelope (apply pipeline per `22-backup-apply-logic.md`). |
| **BE-2** | GET | `/API/V1/Backup/SyncWatermark` | Backup → Primary (poll) | `Backup.Ack.Read` | S2S `aud=Backup` | Returns `LastAcceptedSyncOpSeq` for catch-up replay. |
| **BE-3** | POST | `/API/V1/Backup/RestoreByDate` | Operator → Backup | `Backup.Restore.Write` | S2S `aud=Backup` | Operator-initiated restore (per `23-snapshot-storage-and-restore.md` §6 + `seq-backup-restore.mmd`). |
| **BE-4** | GET | `/API/V1/Backup/RestoreJob/{RestoreJobId}` | Operator → Backup | `Backup.Restore.Write` | S2S `aud=Backup` | Poll restore job status. |
| **BE-5** | POST | `/API/V1/Backup/Pairing/Confirm` | Either side | `Backup.Sync.Write` | S2S `aud=Backup` | Two-phase pairing handshake (per `19-incremental-backup-sync.md` §5). |
| **BE-6** | POST | `/API/V1/Backup/RestoreInbox` | Backup → Primary | `Backup.Restore.Apply` | S2S `aud=Backup` | Receive re-sealed snapshot, perform offline App-tier import, realign watermark. |

### 6.2 Cross-cutting rules

1. **Audience isolation:** §5.2 default-deny middleware MUST treat `aud=Backup` as a **separate** auth surface. A token valid for `aud=worker` or `aud=main-orchestration` MUST NOT satisfy any BE-* endpoint.
2. **Misroute detection:** All BE-* endpoints reject with **HTTP 421 Misdirected Request** + `MAIN-800-04` when `PairingId` does not match the receiving node's `BackupPairing` table. This is a **CODE RED** path — never silently 404 or 401.
3. **Rate limits:** BE-1 and BE-2 are exempted from §5.4 user/IP throttles and instead governed by `MainWorker.Backup.PerPairingEnvelopesPerMinute` (per `15-tunable-constants.md` §2.15).
4. **Audit trail:** BE-3 and BE-6 MUST write to `EndpointAuthAuditEvent` (per §5.6) with `ChangeKind=BackupRestore`.
5. **Endpoint not yet enumerated in §2:** Intentional. The §2.x catalogue is for UI/orchestration surfaces; the BE-* surface is segregated to make audience-isolation review trivial.

*See `21-backup-endpoints.md` for full payloads, error envelopes, and idempotency contracts.*

---

*Core API endpoints v1.3.0 — 2026-05-06 (Phase 12: §6 Backup-tier endpoint catalogue stub added)*
