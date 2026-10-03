# 14 — Error Codes (Main/Worker Service)

**Spec:** `19-main-worker-service`
**Version:** 1.1.0
**Created:** 2026-05-04
**Status:** Authoritative
**Project Prefix:** `MWS`
**Range:** Worker tier `21000-21099` (primary) + `21200-21299` (overflow expansion per §4 Slot-overflow rule), Main tier `21100-21199`
**Resolves:** audit findings F-X-08, F-A-21, F-B-08 (top-10 fix #4). Unblocks AC-6, AC-1.
**Registered in:** `02-spec/03-error-manage/03-error-code-registry/01-registry.md` (line 61–62 entries).

> **Envelope authority (Phase 13.2 — closes audit C-5):** the JSON error envelope shape is owned exclusively by **`09-error-contract.md` §2**. This file (`13-`) is a **code catalogue only** — its rows define `ErrorCode`, `Flat`, `Name`, `Message`, `HTTP`, and `Source` for use *inside* the envelope, not envelope fields themselves. Implementations MUST emit the envelope from `08-` §2 verbatim and populate `Error.ErrorCode` with the prefixed string (e.g. `WORKER-100-01`) from §2/§3 below. The **flat integer is for Go internals and DB columns only** and MUST NOT appear as a JSON envelope field. There is no `ErrorCodeFlat` or `ErrorName` JSON field on the wire; if older drafts mention them, ignore.

---

## 1. Overview

Catalogues every error code referenced by `02-spec/19-main-worker-service/` and by `02-spec/14-update/28-worker-push-instruction.md`. Codes are split into two sub-ranges by issuing tier so dumb-AI implementers can tell at a glance which side throws what.

Both formats are valid (per `02-spec/03-error-manage/03-error-code-registry/01-registry.md` §Format):

- **Prefixed (3-segment):** `WORKER-100-01` — used in spec prose, JSON error envelopes, and PHP/TS code.
- **Flat integer:** `21001` — used in Go internals and DB columns.

The mapping is mechanical: `WORKER-{XYY}-{ZZ}` ↔ `21{XYY}` for worker, `MAIN-{XYY}-{ZZ}` ↔ `211{YY}` for main. Both columns appear below.

> **Slot-overflow rule (added v1.1.0).** When a sub-range's strict `211{YY}` slot is exhausted, the next code is assigned the lowest free flat from the same tier's reserved expansion range (§4) and the §3 sub-range table is annotated. This preserves the prefixed↔flat bijection (Rule R3) without renumbering existing codes. First instance: `MAIN-400-10 EndpointAuthLocked` → `21170` (4xx routing flats `21140-21149` were exhausted by tasks #32 + #39; allocated from `MAIN-21170-21199` reserved range per FU-18).

---

## 2. Worker tier — `WORKER-*` (21000-21099)

### 2.1 Bootstrap (000-099 → 21000-21009)

| Code (prefixed) | Flat | Name | Message | HTTP | Source |
|---|---|---|---|---|---|
| `WORKER-000-01` | `21001` | `BootstrapConfigMissing` | "Required bootstrap config key missing." | 500 | `11-worker-bootstrap-protocol.md` §6 |
| `WORKER-000-02` | `21002` | `SplitDbTierMissing` | "Required split-DB tier failed self-test." | 500 | `10` §6 + `11` §6 |

### 2.2 Authentication (100-199 → 21010-21019)

| Code | Flat | Name | Message | HTTP | Source |
|---|---|---|---|---|---|
| `WORKER-100-01` | `21010` | `OAuthHandshakeFail` | "OAuth client-credentials handshake failed." | 401 | `10` §6 |
| `WORKER-100-02` | `21011` | `KidUnknown` | "JWT signing key id not in trust store." | 401 | `13-jwt-delivery-contract.md` §7 |
| `WORKER-100-03` | `21012` | `WrongWorker` | "JWT `wnk` claim does not match this worker." | 403 | `12` §7 |
| `WORKER-100-04` | `21013` | `JwtSignatureInvalid` | "JWT RS256 signature failed verification." | 401 | `12` §7 |
| `WORKER-100-05` | `21014` | `JwtExpired` | "JWT `exp` is in the past (allowing `ClockSkewToleranceSeconds` per `16-tunable-constants.md` §2.4 — default 60 s)." | 401 | `12` §7 |
| `WORKER-100-06` | `21015` | `JwtIssuerMismatch` | "JWT `iss` claim does not match configured Main host." | 401 | `12` §7 |

### 2.3 Authorization (200-299 → 21020-21029)

| Code | Flat | Name | Message | HTTP | Source |
|---|---|---|---|---|---|
| `WORKER-200-01` | `21020` | `InstructionAlreadyApplied` | "Idempotent replay; instruction already applied." | 200 | `02-spec/14-update/28-worker-push-instruction.md` §6 |
| `WORKER-200-02` | `21021` | `RoleMissingForPage` | "Caller role lacks `RolePageAccess` for requested page." | 403 | `08-role-based-dashboards.md` |

### 2.4 Validation (300-399 → 21030-21039)

| Code | Flat | Name | Message | HTTP | Source |
|---|---|---|---|---|---|
| `WORKER-300-01` | `21030` | `IdempotencyKeyMissing` | "Required `X-Idempotency-Key` header absent on POST/PUT/PATCH." | 400 | `07-core-api-endpoints.md` §1 |
| `WORKER-300-02` | `21031` | `CorrelationIdMissing` | "Required `X-Correlation-Id` header absent." | 400 | `06` §1 |
| `WORKER-300-03` | `21032` | `RequestBodyInvalid` | "Request body fails JSON schema validation." | 400 | `06` |
| `WORKER-300-04` | `21033` | `BackupNotRoutable` | "Selected `WorkerNode` has `IsBackup = 1`; backup nodes never serve traffic." | 409 | `05-worker-routing.md` §1.4 |

### 2.5 Business Logic / Versioning (400-499 → 21040-21049)

| Code | Flat | Name | Message | HTTP | Source |
|---|---|---|---|---|---|
| `WORKER-400-01` | `21040` | `WorkerVersionTooOld` | "Worker `WorkerVersionPin` below `AcceptedVersionRange` lower bound." | 409 | `10` §6 |
| `WORKER-400-02` | `21041` | `VersionMismatch` | "Worker version not in Main's accepted range." | 409 | `10` §6 |
| `WORKER-400-03` | `21042` | `InstructionExpired` | "Push-update instruction past `LatestStartUtc`." | 409 | `02-spec/14-update/28` §6 |
| `WORKER-400-04` | `21043` | `InstructionKindUnsupported` | "PayloadKind not implemented in this worker version." | 409 | `02-spec/14-update/28` §4 |
| `WORKER-401-01` | `21046` | `InvalidPeerCert` | "Invalid peer certificate." | 401 | `27-trust-boundaries-and-isolation.md` |
| `WORKER-403-01` | `21044` | `PushDisabledInProduction` | "Push endpoints REFUSE with 403 when `Env=Production` at Main." | 403 | `17-update-channels.md` §A (Push) |
| `WORKER-403-02` | `21045` | `PayloadHostNotAllowed` | "`PayloadUrl` host not on `AllowedHostsAllowlist` (Channel C defence-in-depth)." | 403 | `17-update-channels.md` §C |

### 2.6 Database / Persistence (500-599 → 21050-21059)

| Code | Flat | Name | Message | HTTP | Source |
|---|---|---|---|---|---|
| `WORKER-500-01` | `21050` | `ClockSkewTooLarge` | "Local clock differs from Main `ServerTimeUtc` beyond `MainWorker.Auth.ClockSkewToleranceSeconds` (per `16-tunable-constants.md` §2.4, default 60 s)." | 500 | `10` §6 |
| `WORKER-500-02` | `21051` | `HandoffFailed` | "Post-update handoff did not confirm within `WorkerPushUpdate.HandoffTimeoutSeconds` (per `16-tunable-constants.md` §2.7, default 60 s); rolled back." | 500 | `02-spec/14-update/28` §6 |
| `WORKER-500-03` | `21052` | `SplitDbWriteFail` | "Write to App or Session tier DB failed." | 500 | `12-split-db-tier-reconciliation.md` |

### 2.7 External Services / Payload (600-699 → 21060-21069)

| Code | Flat | Name | Message | HTTP | Source |
|---|---|---|---|---|---|
| `WORKER-600-01` | `21060` | `PayloadVerificationFail` | "Payload RS256 signature did not verify." | 502 | `02-spec/14-update/28` §6 |
| `WORKER-600-02` | `21061` | `PayloadChecksumFail` | "Downloaded payload SHA256 mismatch." | 502 | `02-spec/14-update/28` §6 |
| `WORKER-600-03` | `21062` | `PayloadSizeFail` | "Downloaded payload size mismatch." | 502 | `02-spec/14-update/28` §6 |
| `WORKER-600-04` | `21063` | `PayloadDownloadFail` | "Could not download payload from `PayloadUrl`." | 502 | `02-spec/14-update/28` §3 |

### 2.8 File System / Deploy (700-799 → 21070-21079)

| Code | Flat | Name | Message | HTTP | Source |
|---|---|---|---|---|---|
| `WORKER-700-01` | `21070` | `SettingsPersistFail` | "Failed to persist `WorkerBootstrapState` to Settings tier." | 500 | `10` §6 |
| `WORKER-700-02` | `21071` | `SchemaMigrateFail` | "Tier schema migration failed." | 500 | `10` §6 |
| `WORKER-700-03` | `21072` | `DeployIoFail` | "Disk-IO failure during rename-first deploy." | 500 | `02-spec/14-update/28` §6 |

### 2.9 Network (800-899 → 21080-21089)

| Code | Flat | Name | Message | HTTP | Source |
|---|---|---|---|---|---|
| `WORKER-800-01` | `21080` | `ListenerBindFail` | "Failed to bind public listener." | 500 | `10` §6 |
| `WORKER-800-02` | `21081` | `WorkerUnreachable` | "Main could not reach worker on `WorkerEndpointPublic`." | 502 | `05-worker-routing.md` |
| `WORKER-503-01` | `21082` | `MainUnreachable` | "Channel B pull: Main host unreachable; exponential backoff, fall through to Channel C if enabled." | 503 | `17-update-channels.md` §B |
| `WORKER-503-02` | `21083` | `ManifestUnreachable` | "Channel C pull: remote manifest unreachable; retry per `PollIntervalSeconds`." | 503 | `17-update-channels.md` §C |
| `WORKER-870-01` | `21084` | `GitBackupRepoCreateFailed` | "Git backup repo creation failed." | 500 | `28-git-backup-targets.md` |
| `WORKER-870-02` | `21085` | `GitBackupAuthFailed` | "SSH key rejected." | 401 | `28-git-backup-targets.md` |
| `WORKER-870-03` | `21086` | `GitBackupRepoMissing` | "Repo missing and auto-create disabled." | 404 | `28-git-backup-targets.md` |
| `WORKER-870-04` | `21087` | `GitBackupPushRejected` | "Push rejected (non-fast-forward)." | 409 | `28-git-backup-targets.md` |
| `WORKER-870-05` | `21088` | `GitBackupReadmeTampered` | "README tamper detected." | 409 | `28-git-backup-targets.md` |
| `WORKER-870-06` | `21089` | `GitBackupNetworkUnreachable` | "Network unreachable." | 503 | `28-git-backup-targets.md` |

### 2.10 Cache Coherence (900-999 → 21090-21099)

| Code | Flat | Name | Message | HTTP | Source |
|---|---|---|---|---|---|
| `WORKER-900-01` | `21090` | `RoleCacheRecompileFailed` | "Worker failed to recompile `RoleAccessCache` after invalidation." | 500 | `18-cascading-roles-and-cache-bin.md` §4 |
| `WORKER-900-02` | `21091` | `EmptyEffectiveAccessSet` | "User has zero `AppUserRole` rows; access denied to gated route." | 403 | `18-cascading-roles-and-cache-bin.md` §1 |
| `WORKER-910-01` | `21092` | `BackupSyncWatermarkInconsistent` | "`BackupSyncWatermark.LastAckedSyncOpSeq > LastShippedSyncOpSeq` (impossible state)." | 500 | `20-incremental-backup-sync.md` §3.1 |
| `WORKER-910-02` | `21093` | `BackupEnvelopeBuildFailed` | "Diff envelope SQLite file could not be created." | 500 | `20-incremental-backup-sync.md` §3.3 |
| `WORKER-910-03` | `21094` | `BackupChangeLogQueryFailed` | "Read-after-watermark query failed." | 500 | `20-incremental-backup-sync.md` §3.2 |
| `WORKER-920-01` | `21095` | `KeyEpochNotYetActive` | "Envelope received under a `Pending` `KeyEpoch`; cannot decrypt." | 409 | `21-backup-encryption-and-keys.md` §9 |
| `WORKER-920-02` | `21096` | `KeyEpochDiscarded` | "`KeyEpoch` past `RetiredKeyGraceSeconds`; private material wiped." | 410 | `21-backup-encryption-and-keys.md` §6, §9 |
| `WORKER-920-03` | `21097` | `UnknownKeyEpoch` | "No `BackupKeyEpoch` row for the supplied epoch." | 404 | `21-backup-encryption-and-keys.md` §9 |
| `WORKER-920-04` | `21098` | `ZipCipherTooWeak` | "Outer zip used ZipCrypto / non-AES; refused." | 415 | `21-backup-encryption-and-keys.md` §5 |
| `WORKER-920-05` | `21099` | `EnvelopeSignatureInvalid` | "RSA-PSS signature failed verification." | 422 | `21-backup-encryption-and-keys.md` §9 |
| `WORKER-930-01` | `21200` | `EnvelopeRowCountMismatch` | "`Envelope.RowCount` ≠ `COUNT(EnvelopeRow)` (V4)." | 422 | `23-backup-apply-logic.md` §3 |
| `WORKER-930-02` | `21201` | `UntrackedSourceTable` | "`SourceTable` not in `AppBackupTrackedTable` allowlist (V5)." | 422 | `23-backup-apply-logic.md` §3 |
| `WORKER-930-03` | `21202` | `UnknownSyncOpCode` | "`SyncOpCode` not in {`Insert`,`Update`,`Delete`} (Stage-4 guard)." | 422 | `23-backup-apply-logic.md` §4 |
| `WORKER-930-04` | `21203` | `ApplyTransactionTimeout` | "Stage-4 TX exceeded `MainWorker.Backup.Apply.TransactionTimeoutSeconds`." | 504 | `23-backup-apply-logic.md` §4 |
| `WORKER-940-01` | `21204` | `SnapshotQuiesceTimeout` | "B2 quiesce wait exceeded `Snapshot.QuiesceTimeoutSeconds`." | 504 | `24-snapshot-storage-and-restore.md` §3 |
| `WORKER-940-02` | `21205` | `RestoreImportFailed` | "BE-6 R7 offline import write failed." | 500 | `24-snapshot-storage-and-restore.md` §8 |
| `WORKER-940-03` | `21206` | `SnapshotBuildTimeout` | "B3 `sqlite3_backup_step` exceeded `Snapshot.MaxBuildSeconds`." | 504 | `24-snapshot-storage-and-restore.md` §3 |
| `WORKER-940-04` | `21207` | `SnapshotSealFailed` | "B4 AES-256-zip seal raised IO/crypto error." | 500 | `24-snapshot-storage-and-restore.md` §3 |

---

## 3. Main tier — `MAIN-*` (21100-21199)

### 3.1 Authentication (100-199 → 21110-21119)

| Code | Flat | Name | Message | HTTP | Source |
|---|---|---|---|---|---|
| `MAIN-100-01` | `21110` | `AuthHandshakeFail` | "Sign-in credentials invalid." | 401 | `06-auth-and-2fa.md` |
| `MAIN-100-02` | `21111` | `SessionExpired` | "Main session cookie expired." | 401 | `05` |
| `MAIN-100-03` | `21112` | `TwoFactorRequired` | "User has 2FA enabled; supply TOTP." | 401 | `05` |
| `MAIN-100-04` | `21113` | `TwoFactorInvalid` | "TOTP code invalid or replayed." | 401 | `05` |

### 3.2 Authorization (200-299 → 21120-21129)

| Code | Flat | Name | Message | HTTP | Source |
|---|---|---|---|---|---|
| `MAIN-200-01` | `21120` | `PowerAdminRequired` | "Endpoint requires PowerAdmin role." | 403 | `06` §2.5/§2.7 |
| `MAIN-200-02` | `21121` | `RoleMissingForPage` | "Caller role lacks `RolePageAccess` for requested page." | 403 | `07` |

### 3.3 Validation (300-399 → 21130-21139)

| Code | Flat | Name | Message | HTTP | Source |
|---|---|---|---|---|---|
| `MAIN-300-01` | `21130` | `CorrelationIdMissing` | "Required `X-Correlation-Id` header absent." | 400 | `05-worker-routing.md` §7.4 + `02-spec/04-database-conventions/07-rest-api-format.md` |
| `MAIN-300-04` | `21131` | `IdempotencyBodyMismatch` | "Replay with different request body for same `X-Idempotency-Key`." | 409 | `05-worker-routing.md` §7.3 |
| `MAIN-300-05` | `21132` | `RoutingResolveFail` | "Could not resolve `(CompanyId) -> WorkerNode` mapping." | 404 | `05-worker-routing.md` §7.2 |

### 3.4 Routing (400-499 → 21140-21149)

| Code | Flat | Name | Message | HTTP | Source |
|---|---|---|---|---|---|
| `MAIN-400-01` | `21140` | `TenantNotFound` | "No `Tenant` row for given `CompanySlug`." | 404 | `05-worker-routing.md` |
| `MAIN-400-02` | `21141` | `WorkerQuarantined` | "Resolved worker is quarantined; routing refused." | 503 | `04` + `10` §7 |
| `MAIN-400-03` | `21142` | `NoEligibleWorker` | "No worker matches placement strategy." | 503 | `04` |
| `MAIN-400-05` | `21143` | `TwoFactorChallengeUnknown` | "TOTP submitted for unknown / expired challenge id." | 401 | `05-worker-routing.md` §7.2 (`/Auth/TwoFactor/Verify`) |
| `MAIN-400-08` | `21144` | `RefreshNotEligible` | "JWT not within refresh window or already rotated." | 401 | `05-worker-routing.md` §7.2 (`/Auth/Refresh`) |
| `MAIN-400-09` | `21145` | `RefreshReplay` | "Single-use refresh JWT replayed after rotation." | 401 | `05-worker-routing.md` §7.2 + diagrams/seq-login-routing |
| `MAIN-400-04` | `21147` | `WorkerRegisterRejected` | "Main refused Worker registration (version pin mismatch, IP not in allow-list, or duplicate `WorkerNodeName`)." | 409 | `09-error-contract.md` §9 + `11-worker-bootstrap-protocol.md` |
| `MAIN-400-06` | `21148` | `WorkerHeartbeatRejected` | "Worker is `Quarantined` or `Offline`; Worker MUST stop sending heartbeats until restart." | 410 | `09-error-contract.md` §9 + `10` §7 |
| `MAIN-400-07` | `21149` | `WorkerPushAckUnknownJid` | "PushAck received for unknown / expired Job-Id; Worker logs and discards." | 404 | `09-error-contract.md` §9 + `02-spec/14-update/28-worker-push-instruction.md` |
| `MAIN-400-11` | `21146` | `AuthActionMissing` | "Required `X-Auth-Action` header absent on multi-step auth flow." | 400 | `05-worker-routing.md` §7.4 |
| `MAIN-400-10` | `21170` | `EndpointAuthLocked` | "Endpoint pattern matches the lock-list (`/API/V1/Workers/*` or `/API/V1/SelfUpdate`) and cannot be reconfigured via `PATCH /API/V1/Settings/EndpointAuth`." | 403 | `07-core-api-endpoints.md` §5.4 R-5 + `06-auth-and-2fa.md` §8 (always-protected rows) |

### 3.5 Database (500-599 → 21150-21159)

| Code | Flat | Name | Message | HTTP | Source |
|---|---|---|---|---|---|
| `MAIN-500-01` | `21150` | `MainDbWriteFail` | "Write to Main Root tier DB failed." | 500 | `04-main-db-schema.md` |

### 3.6 External Services (600-699 → 21160-21169)

| Code | Flat | Name | Message | HTTP | Source |
|---|---|---|---|---|---|
| `MAIN-600-01` | `21160` | `WorkerUnreachable` | "Main could not reach worker on `WorkerEndpointPublic`." | 502 | `04` |
| `MAIN-600-02` | `21161` | `WorkerHeartbeatStale` | "Worker missed ≥3 heartbeats; quarantining." | n/a | `10` §7 |

### 3.7 Cache Coherence (700-799 → 21170-21179)

| Code | Flat | Name | Message | HTTP | Source |
|---|---|---|---|---|---|
| `MAIN-700-01` | `21171` | `CacheInvalidationDeliveryFailed` | "Worker did not ACK `InvalidateRoleAccess` within retry budget." | 502 | `18-cascading-roles-and-cache-bin.md` §5.2 |

### 3.8 Backup Lifecycle (800-899 → 21180-21189)

| Code | Flat | Name | Message | HTTP | Source |
|---|---|---|---|---|---|
| `MAIN-800-01` | `21181` | `BackupChainNotAllowed` | "Cannot register a backup whose primary is itself a backup." | 422 | `19-backup-nodes.md` §1, §3.2 |
| `MAIN-800-02` | `21182` | `PrimaryNotFound` | "`BackupOfWorkerIdentity` does not resolve to any registered Worker." | 404 | `19-backup-nodes.md` §3.2 |
| `MAIN-800-03` | `21183` | `BackupCapacityExceeded` | "Primary already has `MainWorker.Backup.MaxBackupsPerPrimary` backups." | 409 | `19-backup-nodes.md` §2 R2 |
| `MAIN-800-04` | `21184` | `TrafficOnBackupRejected` | "Backup node rejected an inbound user-facing request (D9 invariant)." | 421 | `19-backup-nodes.md` §1, §8 |
| `MAIN-810-01` | `21185` | `BackupCompactionStalled` | "Backup acknowledgement lag exceeds `QuarantineCompactionOverrideSeconds`; compaction blocked." | n/a | `20-incremental-backup-sync.md` §4 |

### 3.9 Backup Encryption / Rotation (820-899 → 21186-21188)

| Code | Flat | Name | Message | HTTP | Source |
|---|---|---|---|---|---|
| `MAIN-820-01` | `21186` | `RotationStepTimeout` | "Pair-RSA rotation step (S2/S3/S6) ACK missed within `RotationAckTimeoutSeconds`." | 504 | `21-backup-encryption-and-keys.md` §7.2 |
| `MAIN-820-02` | `21187` | `RotationAlreadyInProgress` | "A `Pending` `BackupKeyEpoch` already exists for this pairing." | 409 | `21-backup-encryption-and-keys.md` §6 |
| `MAIN-820-03` | `21188` | `RotationActivationSplitBrain` | "One side activated, counterpart did not; manual recovery required." | 500 | `21-backup-encryption-and-keys.md` §7.2 |

### 3.10 Backup Endpoints / Restore Wire (830-899 → 21189-21190)

| Code | Flat | Name | Message | HTTP | Source |
|---|---|---|---|---|---|
| `MAIN-830-01` | `21189` | `SnapshotNotFound` | "Requested `SnapshotDate` not present on backup filesystem." | 404 | `22-backup-endpoints.md` §6.3, §11 (Phase 11 owns storage contract) |
| `MAIN-830-02` | `21190` | `RestoreAlreadyInProgress` | "A `RestoreJobId` for the same `(BackupWorkerNodeId, TargetPrimaryWorkerNodeId)` is still running." | 409 | `22-backup-endpoints.md` §6.3 |

### 3.11 Backup Apply Pipeline (840-899 → 21191)

| Code | Flat | Name | Message | HTTP | Source |
|---|---|---|---|---|---|
| `MAIN-840-01` | `21191` | `BackupApplyExhausted` | "`MaxRetriesPerEnvelope` exceeded for one `EnvelopeId`; surfaced via BE-5 Health." | n/a | `23-backup-apply-logic.md` §6.2 |
| `MAIN-840-02` | `21192` | `SnapshotCorrupt` | "Restore R3 SHA-256 mismatch against `BackupSnapshotCatalog.Sha256Hex`; surfaced via BE-5." | n/a | `24-snapshot-storage-and-restore.md` §7 |

### 3.12 Spec Integrity / AuthZ (900-999 → 21193-21196)

| Code | Flat | Name | Message | HTTP | Source |
|---|---|---|---|---|---|
| `MAIN-900-01` | `21193` | `SpecContradiction` | "Literal AI/implementation found a precedence-rule contradiction; halt rather than guess." | 500 | `26-inherited-rules.md` §6 |
| `MAIN-960-01` | `21194` | `WorkerWriteAttempt` | "Worker JWT used against any non-allow-listed endpoint." | 403 | `27-trust-boundaries-and-isolation.md` |
| `MAIN-960-02` | `21195` | `WorkerCrossTenantQuery` | "Worker queries data for a tenant not assigned to it." | 403 | `27-trust-boundaries-and-isolation.md` |
| `MAIN-960-03` | `21196` | `WorkerImpersonatesMain` | "Worker presents a token claiming role `Power` / `MainAdmin`." | 403 | `27-trust-boundaries-and-isolation.md` |

---

## 4. Reserved sub-ranges

| Sub-range | Reserved for |
|---|---|
| _(consumed)_ | `WORKER-21090-21091` consumed by `WORKER-900-01 RoleCacheRecompileFailed` and `WORKER-900-02 EmptyEffectiveAccessSet` per Phase 5 (`18-cascading-roles-and-cache-bin.md`) |
| _(consumed)_ | `WORKER-21092-21094` consumed by `WORKER-910-01..03` Backup Sync per Phase 7 (`20-incremental-backup-sync.md`) |
| _(consumed)_ | `WORKER-21095-21099` consumed by `WORKER-920-01..05` Backup Encryption per Phase 8 (`21-backup-encryption-and-keys.md`) |
| _(consumed)_ | `WORKER-21200-21203` consumed by `WORKER-930-01..04` Backup Apply per Phase 10 (`23-backup-apply-logic.md`) |
| _(consumed)_ | `WORKER-21204-21207` consumed by `WORKER-940-01..04` Backup Snapshot/Restore per Phase 11 (`24-snapshot-storage-and-restore.md`) |
| `WORKER-21208-21299` | Worker future expansion (snapshot/restore overflow + future overflow windows per §1 Slot-overflow rule) |
| `MAIN-21133-21139` | Main validation future expansion |
| _(consumed)_ | 21147-21149 consumed by `WorkerRegisterRejected` / `WorkerHeartbeatRejected` / `WorkerPushAckUnknownJid` per task #32 (was: Main routing future expansion) |
| _(consumed)_ | 21170 consumed by `MAIN-400-10 EndpointAuthLocked` per FU-18 (overflow from exhausted 21140-21149 4xx-routing range; see §1 *Slot-overflow rule*) |
| _(consumed)_ | 21171 consumed by `MAIN-700-01 CacheInvalidationDeliveryFailed` per Phase 5 (`18-cascading-roles-and-cache-bin.md`) |
| _(consumed)_ | 21181-21184 consumed by `MAIN-800-01..04` Backup Lifecycle per Phase 6 (`19-backup-nodes.md`) |
| _(consumed)_ | 21185 consumed by `MAIN-810-01 BackupCompactionStalled` per Phase 7 (`20-incremental-backup-sync.md`) |
| `MAIN-21162-21169` | Main external-services future expansion |
| `MAIN-21172-21180` | Main future expansion (file-system, network, additional cache-coherence overflow) |
| _(consumed)_ | 21186-21188 consumed by `MAIN-820-01..03` Backup Encryption rotation per Phase 8 (`21-backup-encryption-and-keys.md`) |
| _(consumed)_ | 21189-21190 consumed by `MAIN-830-01..02` Backup Endpoints wire per Phase 9 (`22-backup-endpoints.md`); rows 21191-21199 reserved for Phase 11 snapshot/restore |
| _(consumed)_ | 21191 consumed by `MAIN-840-01 BackupApplyExhausted` per Phase 10 (`23-backup-apply-logic.md`) |
| _(consumed)_ | 21192 consumed by `MAIN-840-02 SnapshotCorrupt` per Phase 11 (`24-snapshot-storage-and-restore.md`) |
| `MAIN-21193-21199` | Main future expansion (additional snapshot/restore + wire overflow) |

---

## 5. Usage examples

```go
// Go (worker)
return errors.New(ErrMws21041, "worker version 1.3.9 below 1.4.0")
```

```typescript
// TypeScript (React)
throw new AppError(ErrorCodes.MAIN_400_01, `No tenant for ${companySlug}`);
```

```php
// PHP (Laravel main)
throw new AppError('MAIN-100-03', 'TOTP required for user '.$userId);
```

All three formats route to the same JSON error envelope per `09-error-contract.md`.

---

## 6. JSON error envelope (per `09-error-contract.md`)

```jsonc
{
  "ErrorCode":      "WORKER-400-02",
  "ErrorCodeFlat":  21041,
  "ErrorName":      "VersionMismatch",
  "ErrorMessage":   "Worker version 1.3.9 not in accepted range >=1.4.0 <2.0.0",
  "CorrelationId":  "01J...ULID",
  "TimestampUtc":   "2026-05-04T12:00:00Z",
  "Tier":           "Worker",
  "WorkerNodeId":   3
}
```

`Tier` is `Main` or `Worker`. `WorkerNodeId` is omitted on Main-tier errors.

---

## 7. Linter assertions (CI) — implemented as `linter-scripts/check-mws-error-codes.py`

Implemented 2026-05-04 (FU-9 closed). Rules enforced:

1. **R1 Presence** — every `WORKER-XYY-ZZ` / `MAIN-XYY-ZZ` literal in `02-spec/19/`, `02-spec/14-update/`, `src/`, and `linter-scripts/tests/` MUST appear in this file's tables.
2. **R2 No orphans** — every code catalogued here MUST be referenced from ≥1 source location outside the catalogue files (`14-error-codes.md`, `error-codes.json`, `error-codes-master.json`). Codes documented only by range-notation cross-reference (e.g. "`WORKER-100-01..05`") may be waived in §7.1.
3. **R3 Bijection** — prefixed ↔ flat mapping MUST be one-to-one.
4. **R4 Range** — `WORKER-*` flats in 21000-21099, `MAIN-*` flats in 21100-21199.

### 7.1 Orphan waivers (Rule R2 exemptions)

Codes referenced only via prose range-notation. Each waiver names the file + range expression.

| Code | Waiver source | Notation |
|---|---|---|
| `WORKER-100-04`, `WORKER-100-05`, `WORKER-100-06` | `13-jwt-delivery-contract.md` §7 | "verification failure → registered error code" |
| `WORKER-200-02` | `08-role-based-dashboards.md` | "RolePageAccess denial" |
| `WORKER-300-02` | `04-database-conventions/07-rest-api-format.md` §Validation | "Missing X-Correlation-Id" |
| `WORKER-400-04` | `02-spec/14-update/28-worker-push-instruction.md` §4 | "PayloadKind not supported" |
| `WORKER-500-02`, `WORKER-500-03` | `02-spec/14-update/28` §6 + `12-split-db-tier-reconciliation.md` | range mention |
| `WORKER-600-04` | `02-spec/14-update/28` §3 | "PayloadDownloadFail" referenced as table heading only |
| `WORKER-700-03` | `02-spec/14-update/28` §6 | "Disk-IO failure" |
| `WORKER-800-02` | `05-worker-routing.md` §3.1 | "WorkerUnreachable" referenced by name |
| `MAIN-100-01..04` | `06-auth-and-2fa.md` | catalogue-only sub-range |
| `MAIN-200-01`, `MAIN-200-02` | `07-core-api-endpoints.md` §2.5/§2.7, `08-role-based-dashboards.md` | catalogue-only |
| `MAIN-400-02`, `MAIN-400-03` | `05-worker-routing.md` §3.1 | range "WorkerUnreachable / quarantine" |
| `MAIN-500-01` | `04-main-db-schema.md` | catalogue-only |
| `MAIN-600-02` | `11-worker-bootstrap-protocol.md` §7 | "missed ≥3 heartbeats" |

Waivers are loaded from `linter-scripts/check-mws-error-codes.waivers.txt` (one prefixed code per line, `#` for comments). Adding a waiver requires a row in the table above + a same-PR change to the waiver file. Removing a waiver after its referencing code is added to source MUST be done in the same PR that adds the reference.

Failure = build break.

---

## 8. Cross-references

- `02-spec/03-error-manage/03-error-code-registry/01-registry.md` — registry entries (added in this task).
- `02-spec/03-error-manage/03-error-code-registry/error-codes-master.json` — must be regenerated to include MWS range (follow-up FU-10).
- `02-spec/19-main-worker-service/09-error-contract.md` — JSON envelope shape.
- `02-spec/19-main-worker-service/11-worker-bootstrap-protocol.md` §6 — bootstrap codes referenced.
- `02-spec/19-main-worker-service/13-jwt-delivery-contract.md` §7 — JWT verification codes referenced.
- `02-spec/14-update/28-worker-push-instruction.md` §6 — push-update codes referenced.

---

## 9. Open Questions (logged, non-blocking)

- **OQ-13-1** Should `MAIN-200-02` and `WORKER-200-02` collapse to a single shared code? Inferred: keep separate so logs identify the throwing tier without an extra field.
- **OQ-13-2** Should `WORKER-600-04 PayloadDownloadFail` carry the HTTP status of the failed GET? Inferred: yes — add as `Details.UpstreamHttpStatus` in the envelope; not breaking.

---

*Error codes (Main/Worker Service) v1.5.0 — 2026-05-06 (Phase 11: +`WORKER-940-01..04` Backup Snapshot/Restore consuming `WORKER-21204-21207`; +`MAIN-840-02 SnapshotCorrupt` consuming `MAIN-21192`; §4 reserved-range refresh — `MAIN-21193-21199` reserved for future overflow)*
