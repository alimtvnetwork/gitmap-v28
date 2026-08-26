# 15 — Tunable Constants (Single-Value Pins)

**Spec:** `19-main-worker-service`
**Version:** 1.4.0
**Created:** 2026-05-04
**Status:** Authoritative
**Resolves:** audit findings F-A-15, F-A-16, F-B-12 (top-10 fix #7). Closes AC-7, partially AC-6.
**Authority:** Single source of truth for every numeric tunable referenced anywhere in `spec/19-main-worker-service/` or `spec/14-update/28-worker-push-instruction.md`. On any conflict, **this file wins**. All values seedable via `spec/06-seedable-config-architecture/` (key/value `Categories` block, not `Tables`).

---

## 1. Why this file exists

Earlier drafts of spec/19 mentioned "3 retries" in one place and "5 retries" in another (audit F-A-15). Idempotency-key TTL was implied but never named. Heartbeat interval appeared as `30s` in `10-worker-bootstrap-protocol.md` §3.2 but was unmentioned in the routing doc — leaving the dumb-AI implementer free to invent. This file pins every such constant.

Each row below is **the** value. Implementations MAY override via Seedable-Config at install time but MUST start from this default and MUST NOT hard-code different defaults in source.

---

## 2. Canonical constants

### 2.1 Retry / backoff

| Key | Default | Unit | Used by | Notes |
|---|---:|---|---|---|
| `MainWorker.Retry.MaxAttempts` | **3** | count | All Main↔Worker HTTP calls | Resolves F-A-15 (3-vs-5 contradiction). Includes initial attempt. |
| `MainWorker.Retry.BackoffSeconds` | `[2, 8, 30]` | seconds[] | Same | Length MUST equal `MaxAttempts - 1`. Exponential-ish, ceiling 30s. |
| `MainWorker.Retry.JitterPct` | **20** | percent | Same | ±20% applied per attempt. |
| `WorkerPushUpdate.MaxRetries` | **3** | count | `spec/14-update/28` §3.1 `OnFailure.MaxRetries` | Mirrors above; `RetryBackoffSeconds` defaults to `[30, 120, 300]` per §3 of that file. |

### 2.2 Idempotency

| Key | Default | Unit | Used by | Notes |
|---|---:|---|---|---|
| `MainWorker.Idempotency.KeyTtlSeconds` | **86400** (24h) | seconds | Worker + Main idempotency stores | Single canonical TTL. Replays beyond this MAY re-execute. |
| `MainWorker.Idempotency.KeyMaxLength` | **64** | chars | Header validation | ULID = 26; reservation buffer for future formats. |
| `MainWorker.Idempotency.StoreCleanupSeconds` | **3600** (1h) | seconds | Background sweeper | Deletes rows older than `KeyTtlSeconds`. |

### 2.3 Heartbeat & quarantine

| Key | Default | Unit | Used by | Notes |
|---|---:|---|---|---|
| `MainWorker.Heartbeat.IntervalSeconds` | **30** | seconds | `10-worker-bootstrap-protocol.md` §3.2 / §7 | Worker→Main ping cadence. |
| `MainWorker.Heartbeat.MissedThreshold` | **3** | count | `10` §7, `04-worker-routing.md` | Consecutive misses → quarantine. |
| `MainWorker.Heartbeat.QuarantineCooldownSeconds` | **300** (5m) | seconds | `04-worker-routing.md` | Quarantined worker eligible for re-eval after this. |
| `MainWorker.Heartbeat.GraceWindowSeconds` | **5** | seconds | Main-side scheduler | Tolerated jitter before counting a miss. |

### 2.4 Auth & JWT

| Key | Default | Unit | Used by | Notes |
|---|---:|---|---|---|
| `MainWorker.Auth.WorkerJwtTtlSeconds` | **900** (15m) | seconds | `12-jwt-delivery-contract.md` §6 | Already pinned there; mirrored here for the single-table view. |
| `MainWorker.Auth.JwtRefreshLeadSeconds` | **60** | seconds | `12` §6 | React refreshes when within this window of `exp`. |
| `MainWorker.Auth.MainSessionTtlSeconds` | **28800** (8h) | seconds | Main session cookie | Sliding window: each qualifying request extends. See §7.2. |
| `MainWorker.Auth.MainSessionAbsoluteMaxSeconds` | **86400** (24h) | seconds | Main session cookie | Hard ceiling from initial login regardless of activity. Forces `Reauthenticate`. MUST be ≥ `MainSessionTtlSeconds` (T4 linter invariant). Decided in §7.2. |
| `MainWorker.Auth.SessionSlidingExtendOnReadOnly` | **true** | bool | Main session cookie | If `false`, only state-changing requests extend the sliding window. Decided in §7.2. |
| `MainWorker.Auth.ClockSkewToleranceSeconds` | **60** | seconds | `12` §7, `10` §3 | Same value across both contexts. |

### 2.5 Routing & pool

| Key | Default | Unit | Used by | Notes |
|---|---:|---|---|---|
| `MainWorker.Routing.DefaultStrategy` | **`LeastLoaded`** | enum | `04-worker-routing.md` §1 + §1.5 | Resolves OQ-2 (Phase 12.3). Allowed: `RoundRobin` \| `LeastLoaded` \| `Manual`. Default rationale + override guidance in `04-worker-routing.md` §1.5. Main MUST refuse to start if the configured value is not in the allow-list. |
| `MainWorker.Routing.HttpTimeoutSeconds` | **15** | seconds | Main → Worker request timeout | Distinct from retry attempt budget. |
| `MainWorker.Routing.HttpHandshakeTimeoutSeconds` | **30** | seconds | Push-update POST handshake (per `spec/14-update/28` §2) | Already pinned there. |
| `MainWorker.Routing.MaxConcurrentPerWorker` | **64** | count | Connection pool cap | Prevents stampede. |

### 2.6 Rate limiting (mirrors `06-core-api-endpoints.md` §6)

| Key | Default | Unit | Notes |
|---|---:|---|---|
| `MainWorker.RateLimit.AuthEndpointsPerMinutePerIp` | **10** | rpm | Already pinned in `06` §6. |
| `MainWorker.RateLimit.WorkerEndpointsPerMinutePerToken` | **60** | rpm | Same. |
| `MainWorker.RateLimit.OtherAuthenticatedPerMinutePerUser` | **600** | rpm | Same. |

### 2.7 Self-update execution window (mirrors `spec/14-update/28` §3)

| Key | Default | Unit | Notes |
|---|---:|---|---|
| `WorkerPushUpdate.MaxRunDurationSeconds` | **600** (10m) | seconds | Hard cap before self-abort + rollback. |
| `WorkerPushUpdate.HandoffTimeoutSeconds` | **60** | seconds | Per `spec/14-update/28` §5 step 8. |
| `WorkerPushUpdate.InstructionRetentionDays` | **14** | days | Per `spec/14-update/28` §7. |
| `WorkerPushUpdate.IssuedSkewSeconds` | **300** (5m) | seconds | Per `spec/14-update/28` §3 — Worker rejects an instruction whose `IssuedAtUtc` is older than this. Distinct from JWT clock-skew (§2.4). Bumped to v1.1.0 to retire the §28 line-82 TUNABLE-WAIVER. |

### 2.8 Self-update pointer (consumer: `09-self-update-pointer.md`)

| Key | Default | Unit | Used by | Notes |
|---|---:|---|---|---|
| `MainWorker.SelfUpdate.RedirectStaleHours` | **36** | hours | `09-self-update-pointer.md` §1 step 5 | If the cached redirect URL is older than this OR unreachable, re-resolve via the original endpoint. Bumped to v1.1.0 to retire the §09 line-41 TUNABLE-WAIVER. |

### 2.9 Worker bootstrap (consumer: `10-worker-bootstrap-protocol.md`)

| Key | Default | Unit | Used by | Notes |
|---|---:|---|---|---|
| `MainWorker.Bootstrap.RetryBackoffSeconds` | `[10, 30, 90, 300]` | seconds[] | `10-worker-bootstrap-protocol.md` §6 (`WORKER-100-01 OAUTH_HANDSHAKE_FAIL`) | Cold-bootstrap retry ladder for Worker→Main OAuth handshake; distinct from the steady-state `MainWorker.Retry.BackoffSeconds` of §2.1. After exhaustion the Worker exits and is restarted by its supervisor. Bumped to v1.1.0 to retire the §10 line-137 TUNABLE-WAIVER. |
| `MainWorker.Bootstrap.RetryMaxAttempts` | **4** | count | Same | MUST equal `len(RetryBackoffSeconds)`. |

### 2.10 Role-access cache bin (consumer: `17-cascading-roles-and-cache-bin.md`)

| Key | Default | Unit | Used by | Notes |
|---|---:|---|---|---|
| `MainWorker.RoleCache.TtlSeconds` | **600** (10m) | seconds | `17-…` §4 — TTL on `RoleAccessCache` rows in the Worker's in-memory Cache tier. | Safety net bound on staleness when the §5 invalidation broadcast fails. Must be ≤ `MainWorker.Auth.WorkerJwtTtlSeconds × 2` so a stale cache cannot outlive two JWT generations. |
| `MainWorker.RoleCache.RequireReauthOnCatalogBump` | **false** | bool | `17-…` §5.3 — when `true`, a JWT carrying an outdated `CatalogVersion` triggers `401 ReauthRequired` instead of transparent recompute. | Set `true` for high-security tenants where role demotions must take effect immediately even at the cost of a forced sign-in. |

### 2.11 Backup nodes (consumer: `18-backup-nodes.md`)

| Key | Default | Unit | Used by | Notes |
|---|---:|---|---|---|
| `MainWorker.Backup.Enabled` | **false** | bool | All `MainWorker.Backup.*` consumers | Master feature flag (added v2.0.0). When `false`, backup endpoints reject with `MAIN-800-01 BackupSubsystemDisabled`, the sync cron does not run, and every other `MainWorker.Backup.*` tunable is inert. Default `false` until the operator opts in (per D9 — no auto-failover). |
| `MainWorker.Backup.MaxBackupsPerPrimary` | **3** | count | `18-…` §2 R2 / §3.2 `AssertCapacityForBackup` | Hard ceiling enforced at backup registration. |
| `MainWorker.Backup.LagWarningSeconds` | **900** (15m) | seconds | `18-…` §4.3 — heartbeat watcher flips `BackupAttached` → `BackupLagging`. | No auto-failover; Power Admin acts on the warning. |
| `MainWorker.Backup.HeartbeatIntervalSeconds` | **60** | seconds | Backup → Main heartbeat cadence | Slower than primary heartbeat (`MainWorker.Heartbeat.IntervalSeconds = 30`) because backups never serve traffic. |
| `MainWorker.Backup.SyncIntervalSeconds` | **60** | seconds | `19-…` §3.1 cron driver on the primary | Per-backup pass cadence. Lower = lower lag, higher = lower CPU/IO. |
| `MainWorker.Backup.MaxRowsPerEnvelope` | **5000** | count | `19-…` §3.2 read query | Hard ceiling per envelope; prevents oversized blobs on bulk writes. |
| `MainWorker.Backup.TombstoneRetentionSeconds` | **604800** (7 d) | seconds | `19-…` §4 (Shape A compaction) | Minimum age before a delivered tombstone row is reclaimed. |
| `MainWorker.Backup.LogRetentionSeconds` | **604800** (7 d) | seconds | `19-…` §4 (Shape B compaction) | Minimum age before a delivered `BackupSyncLog` row is truncated. |
| `MainWorker.Backup.QuarantineCompactionOverrideSeconds` | **86400** (24 h) | seconds | `19-…` §4 stall guard | When a single lagging backup blocks compaction past this, Main raises `MAIN-810-01 BackupCompactionStalled`. Operator decides next step; never auto-detaches. |

### 2.12 Backup encryption + Pair-RSA rotation (consumer: `20-backup-encryption-and-keys.md`)

| Key | Default | Unit | Used by | Notes |
|---|---:|---|---|---|
| `MainWorker.Backup.MaxKeyAgeSeconds` | **7776000** (90 d) | seconds | `20-…` §7.1 trigger #1 | Scheduled rotation cadence per primary-backup couple. |
| `MainWorker.Backup.RotationAckTimeoutSeconds` | **120** | seconds | `20-…` §7.2 S2/S3/S6 | Per-step ACK budget; miss → `MAIN-820-01 RotationStepTimeout`. |
| `MainWorker.Backup.RotationActivationDelaySeconds` | **60** | seconds | `20-…` §7.2 S5→S6 | Settle window before flipping `Pending` → `Active`. |
| `MainWorker.Backup.RetiredKeyGraceSeconds` | **86400** (24 h) | seconds | `20-…` §6 background sweeper | How long `Retired` private material survives in-flight envelopes; after this, `WORKER-920-02 KeyEpochDiscarded`. Override = 0 when `Reason="Compromise"`. |
| `MainWorker.Backup.RsaKeySizeBits` | **4096** | bits | `20-…` §3 K1 | RSA-OAEP wrap + RSA-PSS sign modulus. Bumping requires v2.0 review. |

### 2.13 Backup endpoints (consumer: `21-backup-endpoints.md`)

| Key | Default | Unit | Used by | Notes |
|---|---:|---|---|---|
| `MainWorker.Backup.Endpoint.IncrementalDiffTimeoutSeconds` | **120** | seconds | BE-1 (`21-…` §4) | Larger than routing default — envelope upload may be MB-scale. |
| `MainWorker.Backup.Endpoint.RotateKeysTimeoutSeconds` | **30** | seconds | BE-2 (`21-…` §5) | Below `RotationAckTimeoutSeconds=120` so step retries fit inside the budget. |
| `MainWorker.Backup.Endpoint.RestoreByDateTimeoutSeconds` | **60** | seconds | BE-3 (`21-…` §6) | Accepts the job; restore itself runs as a long-running worker job. |
| `MainWorker.Backup.Endpoint.SnapshotsTimeoutSeconds` | **15** | seconds | BE-4 (`21-…` §7) | Catalogue read. |
| `MainWorker.Backup.Endpoint.HealthTimeoutSeconds` | **5** | seconds | BE-5 (`21-…` §8) | Tight on purpose so dashboard polling surfaces real outages. |

### 2.14 Backup apply pipeline (consumer: `22-backup-apply-logic.md`)

| Key | Default | Unit | Used by | Notes |
|---|---:|---|---|---|
| `MainWorker.Backup.Apply.MaxRetriesPerEnvelope` | **5** | count | `22-…` §6.2 | Exceeding = `MAIN-840-01 BackupApplyExhausted` surfaced via BE-5. |
| `MainWorker.Backup.Apply.TransactionTimeoutSeconds` | **30** | seconds | `22-…` §4 Stage-4 TX | Exceeding = `WORKER-930-04 ApplyTransactionTimeout`. |
| `MainWorker.Backup.Apply.DeadLetterRetentionDays` | **30** | days | DLQ sweeper (`22-…` §6) | Symmetric with snapshot retention default; finalised in Phase 11 with OQ-A4. |
| `MainWorker.Backup.Apply.IdempotencyRowRetentionDays` | **14** | days | DLQ sweeper (`22-…` §7) | `BackupApplyIdempotency.Status='Applied'` rows reaped after this; replay-protection window. |

### 2.15 Backup snapshot + restore (consumer: `23-snapshot-storage-and-restore.md`)

| Key | Default | Unit | Used by | Notes |
|---|---:|---|---|---|
| `MainWorker.Backup.SnapshotRetentionDays` | **30** | days | `23-…` §6 sweep | **Resolves OQ-A4.** Linter `BACKUP-SNAP-002` enforces ≥ 7 (compliance floor). |
| `MainWorker.Backup.Snapshot.BuildHourUtc` | **3** | hour-of-day (0-23) | `23-…` §3 cron | Off-peak default; overridable per node. |
| `MainWorker.Backup.Snapshot.QuiesceTimeoutSeconds` | **120** | seconds | `23-…` §3 B2 | Below `MaxRetriesPerEnvelope × TransactionTimeoutSeconds` (5×30=150). |
| `MainWorker.Backup.Snapshot.MaxBuildSeconds` | **1800** (30 m) | seconds | `23-…` §3 B3 | Hard ceiling on `sqlite3_backup_step` total elapsed; exceeded → `WORKER-940-03`. |
| `MainWorker.Backup.Restore.PrimaryAckTimeoutSeconds` | **600** (10 m) | seconds | `23-…` §7 R8 | Backup waits this long for primary's BE-6 200 before flipping job to `Failed`. |

All Backup-tier tunables now allocated. Phase 12 closes the backup work with diagrams + acceptance criteria.

### 2.16 In-process Main caches (consumer: `01-architecture.md` §5)

Promoted from §4.2 to a first-class catalogue entry so T3 (seed parity) covers them. Defaults match §4.2 verbatim; §4.2 retains the prose discussion.

| Key | Default | Unit | Used by | Notes |
|---|---:|---|---|---|
| `MainWorker.Cache.CompanyToWorkerTtlSeconds` | **900** (15 m) | seconds | `01-architecture.md` §5 — `CompanyId → WorkerNodeId` cache | Invalidate on worker reassignment. |
| `MainWorker.Cache.WorkerRegistryTtlSeconds` | **60** | seconds | `01-architecture.md` §5 — Worker registry | Invalidate on worker register/deregister. |
| `MainWorker.Cache.RecentCompanyPerUserTtlSeconds` | **`= MainWorker.Auth.MainSessionTtlSeconds`** | seconds | `01-architecture.md` §5 — Per-user recent-company | Bound to session TTL by definition; T3 omits parity check (non-numeric default). |

---

## 3. Single-value rule (for the dumb AI)

| Rule | Why |
|---|---|
| Every numeric tunable referenced in spec/19 or spec/14-update/28 MUST appear in §2. | Eliminates "guess what value the author meant." |
| If a value appears in two prose locations, both MUST cite this file by anchor (e.g. `15-tunable-constants.md` §2.3). | Forces single source of truth. |
| Updating a default = bump `config.seed.json` `Categories.MainWorker` `Version`, edit ONE row in §2, run linter. | No multi-file scavenger hunt. |
| Override at runtime = Seedable-Config Categories block (NOT `Tables`). | Per `spec/06-seedable-config-architecture/01-fundamentals.md`. |

---

## 4. `config.seed.json` Categories binding (paste-ready)

Add (or merge with) the following category at SemVer `2.0.0` of `config.seed.json` (bumped from `1.5.0` to materialize the 27 backup-tier tunables (`MainWorker.Backup.*`) defined in §2.11–2.15, closing the audit-09 §2.1 deferral; the major bump reflects the new `Backup` sub-namespace and the addition of the `Enabled` feature flag, not a breaking change to existing v1.5.0 keys — all v1.5.0 defaults are preserved verbatim):

```jsonc
"MainWorker": {
  "DisplayName": "Main / Worker tunables",
  "Description": "Single-value tunables for spec/19. See 15-tunable-constants.md.",
  "Version": "2.0.0",
  "AddedIn":  "1.3.0",
  "Settings": {

    "RetryMaxAttempts":            { "Type": "number",  "Default": 3,         "Min": 1,    "Max": 10 },
    "RetryBackoffSeconds":         { "Type": "string",  "Default": "2,8,30",  "Description": "Comma-separated; length = RetryMaxAttempts - 1." },
    "RetryJitterPct":              { "Type": "number",  "Default": 20,        "Min": 0,    "Max": 100 },

    "IdempotencyKeyTtlSeconds":    { "Type": "number",  "Default": 86400,     "Min": 60 },
    "IdempotencyKeyMaxLength":     { "Type": "number",  "Default": 64,        "Min": 26,   "Max": 256 },
    "IdempotencyStoreCleanupSec":  { "Type": "number",  "Default": 3600,      "Min": 60 },

    "HeartbeatIntervalSeconds":    { "Type": "number",  "Default": 30,        "Min": 5,    "Max": 600 },
    "HeartbeatMissedThreshold":    { "Type": "number",  "Default": 3,         "Min": 1,    "Max": 20 },
    "HeartbeatQuarantineCooldown": { "Type": "number",  "Default": 300,       "Min": 30 },
    "HeartbeatGraceWindowSeconds": { "Type": "number",  "Default": 5,         "Min": 0 },

    "WorkerJwtTtlSeconds":         { "Type": "number",  "Default": 900,       "Min": 60,   "Max": 3600 },
    "JwtRefreshLeadSeconds":       { "Type": "number",  "Default": 60,        "Min": 10 },
    "MainSessionTtlSeconds":         { "Type": "number",  "Default": 28800,     "Min": 300 },
    "MainSessionAbsoluteMaxSeconds": { "Type": "number",  "Default": 86400,     "Min": 300 },
    "SessionSlidingExtendOnRead":    { "Type": "boolean", "Default": true },
    "ClockSkewToleranceSeconds":   { "Type": "number",  "Default": 60,        "Min": 0,    "Max": 300 },

    "RoutingHttpTimeoutSeconds":   { "Type": "number",  "Default": 15,        "Min": 1 },
    "RoutingHandshakeTimeoutSec":  { "Type": "number",  "Default": 30,        "Min": 1 },
    "RoutingMaxConcurrentPerNode": { "Type": "number",  "Default": 64,        "Min": 1 },

    "RateAuthPerMinutePerIp":      { "Type": "number",  "Default": 10 },
    "RateWorkerPerMinutePerToken": { "Type": "number",  "Default": 60 },
    "RateOtherPerMinutePerUser":   { "Type": "number",  "Default": 600 },

    "PushUpdateMaxRunSeconds":     { "Type": "number",  "Default": 600,       "Min": 30 },
    "PushUpdateHandoffTimeoutSec": { "Type": "number",  "Default": 60,        "Min": 5 },
    "PushUpdateRetentionDays":     { "Type": "number",  "Default": 14,        "Min": 1 },
    "PushUpdateIssuedSkewSec":     { "Type": "number",  "Default": 300,       "Min": 30,   "Max": 3600 },

    "SelfUpdateRedirectStaleHours":{ "Type": "number",  "Default": 36,        "Min": 1,    "Max": 720 },

    "BootstrapRetryBackoffSec":    { "Type": "string",  "Default": "10,30,90,300", "Description": "Comma-separated; length = BootstrapRetryMaxAttempts." },
    "BootstrapRetryMaxAttempts":   { "Type": "number",  "Default": 4,         "Min": 1,    "Max": 10 },

    "CacheCompanyToWorkerTtlSeconds":     { "Type": "number", "Default": 900,   "Min": 30 },
    "CacheWorkerRegistryTtlSeconds":      { "Type": "number", "Default": 60,    "Min": 5  },
    "CacheRecentCompanyPerUserTtlSeconds":{ "Type": "number", "Default": 28800, "Min": 300, "Description": "Default mirrors MainSessionTtlSeconds; runtime resolver MAY substitute the live MainSessionTtlSeconds value to honour §4.2 binding." },

    // ─── Backup tier (added v2.0.0; gated by Backup.Enabled) ─────────────────
    // Source: §2.11–2.15. All defaults match the prose verbatim.

    "Backup.Enabled":                              { "Type": "boolean", "Default": false, "Description": "Master feature flag. When false, the entire backup subsystem is dormant: backup endpoints reject with MAIN-800-01 BackupSubsystemDisabled, sync cron does not run, and Backup.* tunables below have no effect. Default false until operator explicitly opts in (per D9 — no auto-failover)." },

    // §2.11 Backup nodes
    "Backup.MaxBackupsPerPrimary":                 { "Type": "number",  "Default": 3,       "Min": 1, "Max": 10 },
    "Backup.LagWarningSeconds":                    { "Type": "number",  "Default": 900,     "Min": 60 },
    "Backup.HeartbeatIntervalSeconds":             { "Type": "number",  "Default": 60,      "Min": 5,  "Max": 600 },
    "Backup.SyncIntervalSeconds":                  { "Type": "number",  "Default": 60,      "Min": 5,  "Max": 3600 },
    "Backup.MaxRowsPerEnvelope":                   { "Type": "number",  "Default": 5000,    "Min": 1,  "Max": 100000 },
    "Backup.TombstoneRetentionSeconds":            { "Type": "number",  "Default": 604800,  "Min": 86400 },
    "Backup.LogRetentionSeconds":                  { "Type": "number",  "Default": 604800,  "Min": 86400 },
    "Backup.QuarantineCompactionOverrideSeconds":  { "Type": "number",  "Default": 86400,   "Min": 3600 },

    // §2.12 Backup encryption + Pair-RSA rotation
    "Backup.MaxKeyAgeSeconds":                     { "Type": "number",  "Default": 7776000, "Min": 86400 },
    "Backup.RotationAckTimeoutSeconds":            { "Type": "number",  "Default": 120,     "Min": 10 },
    "Backup.RotationActivationDelaySeconds":       { "Type": "number",  "Default": 60,      "Min": 0 },
    "Backup.RetiredKeyGraceSeconds":               { "Type": "number",  "Default": 86400,   "Min": 0,  "Description": "Grace window for in-flight envelopes signed with retired keys. MUST be 0 when rotation Reason='Compromise' (override at request time, not via seed)." },
    "Backup.RsaKeySizeBits":                       { "Type": "number",  "Default": 4096,    "Min": 4096, "Max": 4096, "Description": "Locked at 4096 for v2.0; bumping requires a v3.0 seed review." },

    // §2.13 Backup endpoints
    "Backup.Endpoint.IncrementalDiffTimeoutSeconds": { "Type": "number", "Default": 120, "Min": 1 },
    "Backup.Endpoint.RotateKeysTimeoutSeconds":      { "Type": "number", "Default": 30,  "Min": 1 },
    "Backup.Endpoint.RestoreByDateTimeoutSeconds":   { "Type": "number", "Default": 60,  "Min": 1 },
    "Backup.Endpoint.SnapshotsTimeoutSeconds":       { "Type": "number", "Default": 15,  "Min": 1 },
    "Backup.Endpoint.HealthTimeoutSeconds":          { "Type": "number", "Default": 5,   "Min": 1 },

    // §2.14 Backup apply pipeline
    "Backup.Apply.MaxRetriesPerEnvelope":           { "Type": "number", "Default": 5,  "Min": 1, "Max": 20 },
    "Backup.Apply.TransactionTimeoutSeconds":       { "Type": "number", "Default": 30, "Min": 1 },
    "Backup.Apply.DeadLetterRetentionDays":         { "Type": "number", "Default": 30, "Min": 1 },
    "Backup.Apply.IdempotencyRowRetentionDays":     { "Type": "number", "Default": 14, "Min": 1 },

    // §2.15 Backup snapshot + restore
    "Backup.SnapshotRetentionDays":                 { "Type": "number", "Default": 30,   "Min": 7, "Description": "Linter BACKUP-SNAP-002 enforces ≥7 (compliance floor)." },
    "Backup.Snapshot.BuildHourUtc":                 { "Type": "number", "Default": 3,    "Min": 0, "Max": 23 },
    "Backup.Snapshot.QuiesceTimeoutSeconds":        { "Type": "number", "Default": 120,  "Min": 1 },
    "Backup.Snapshot.MaxBuildSeconds":              { "Type": "number", "Default": 1800, "Min": 60 },
    "Backup.Restore.PrimaryAckTimeoutSeconds":      { "Type": "number", "Default": 600,  "Min": 30 }

  }
}
```

> **Backup-tier seed (v2.0.0 — closes audit-09 §2.1):** All 27 `MainWorker.Backup.*` tunables from §2.11–2.15 are now materialized in the JSON payload above, gated by `MainWorker.Backup.Enabled` (default `false`). The T3 (seed parity) linter waiver previously held against the `MainWorker.Backup.*` namespace is **lifted as of v2.0.0**; T3 now applies to every key in §2 including the backup tier. Implementers reading defaults MUST consume the seed; §2.11–2.15 prose remains the SoT for *meaning* and unit conventions only.

---

## 4.1 Prose↔seed key alias map (Phase 13.2 — closes audit C-4)

The §2 prose uses full dotted keys (e.g. `MainWorker.RateLimit.AuthEndpointsPerMinutePerIp`); the §4 JSON seed flattens them under `Categories.MainWorker.Settings` (e.g. `RateAuthPerMinutePerIp`). Both are canonical; the runtime resolver MUST treat them as the same setting via this table. A literal AI MUST NOT invent its own short forms — only the aliases below are valid.

| Prose key (§2)                                            | Seed key (§4 `Settings.<key>`)   |
|-----------------------------------------------------------|----------------------------------|
| `MainWorker.Retry.MaxAttempts`                            | `RetryMaxAttempts`               |
| `MainWorker.Retry.BackoffSeconds`                         | `RetryBackoffSeconds`            |
| `MainWorker.Retry.JitterPct`                              | `RetryJitterPct`                 |
| `MainWorker.Idempotency.KeyTtlSeconds`                    | `IdempotencyKeyTtlSeconds`       |
| `MainWorker.Idempotency.KeyMaxLength`                     | `IdempotencyKeyMaxLength`        |
| `MainWorker.Idempotency.StoreCleanupSeconds`              | `IdempotencyStoreCleanupSec`     |
| `MainWorker.Heartbeat.IntervalSeconds`                    | `HeartbeatIntervalSeconds`       |
| `MainWorker.Heartbeat.MissedThreshold`                    | `HeartbeatMissedThreshold`       |
| `MainWorker.Heartbeat.QuarantineCooldownSeconds`          | `HeartbeatQuarantineCooldown`    |
| `MainWorker.Heartbeat.GraceWindowSeconds`                 | `HeartbeatGraceWindowSeconds`    |
| `MainWorker.Auth.WorkerJwtTtlSeconds`                     | `WorkerJwtTtlSeconds`            |
| `MainWorker.Auth.JwtRefreshLeadSeconds`                   | `JwtRefreshLeadSeconds`          |
| `MainWorker.Auth.MainSessionTtlSeconds`                   | `MainSessionTtlSeconds`          |
| `MainWorker.Auth.MainSessionAbsoluteMaxSeconds`           | `MainSessionAbsoluteMaxSeconds`  |
| `MainWorker.Auth.SessionSlidingExtendOnReadOnly`          | `SessionSlidingExtendOnRead`     |
| `MainWorker.Auth.ClockSkewToleranceSeconds`               | `ClockSkewToleranceSeconds`      |
| `MainWorker.Routing.HttpTimeoutSeconds`                   | `RoutingHttpTimeoutSeconds`      |
| `MainWorker.Routing.HttpHandshakeTimeoutSeconds`          | `RoutingHandshakeTimeoutSec`     |
| `MainWorker.Routing.MaxConcurrentPerWorker`               | `RoutingMaxConcurrentPerNode`    |
| `MainWorker.RateLimit.AuthEndpointsPerMinutePerIp`        | `RateAuthPerMinutePerIp`         |
| `MainWorker.RateLimit.WorkerEndpointsPerMinutePerToken`   | `RateWorkerPerMinutePerToken`    |
| `MainWorker.RateLimit.OtherAuthenticatedPerMinutePerUser` | `RateOtherPerMinutePerUser`      |
| `WorkerPushUpdate.MaxRunDurationSeconds`                  | `PushUpdateMaxRunSeconds`        |
| `WorkerPushUpdate.HandoffTimeoutSeconds`                  | `PushUpdateHandoffTimeoutSec`    |
| `WorkerPushUpdate.InstructionRetentionDays`               | `PushUpdateRetentionDays`        |
| `WorkerPushUpdate.IssuedSkewSeconds`                      | `PushUpdateIssuedSkewSec`        |
| `MainWorker.SelfUpdate.RedirectStaleHours`                | `SelfUpdateRedirectStaleHours`   |
| `MainWorker.Bootstrap.RetryBackoffSeconds`                | `BootstrapRetryBackoffSec`       |
| `MainWorker.Bootstrap.RetryMaxAttempts`                   | `BootstrapRetryMaxAttempts`      |
| `MainWorker.Cache.CompanyToWorkerTtlSeconds`              | `CacheCompanyToWorkerTtlSeconds` |
| `MainWorker.Cache.WorkerRegistryTtlSeconds`               | `CacheWorkerRegistryTtlSeconds`  |
| `MainWorker.Cache.RecentCompanyPerUserTtlSeconds`         | `CacheRecentCompanyPerUserTtlSeconds` |
| `MainWorker.Backup.Enabled`                               | `Backup.Enabled`                              |
| `MainWorker.Backup.MaxBackupsPerPrimary`                  | `Backup.MaxBackupsPerPrimary`                 |
| `MainWorker.Backup.LagWarningSeconds`                     | `Backup.LagWarningSeconds`                    |
| `MainWorker.Backup.HeartbeatIntervalSeconds`              | `Backup.HeartbeatIntervalSeconds`             |
| `MainWorker.Backup.SyncIntervalSeconds`                   | `Backup.SyncIntervalSeconds`                  |
| `MainWorker.Backup.MaxRowsPerEnvelope`                    | `Backup.MaxRowsPerEnvelope`                   |
| `MainWorker.Backup.TombstoneRetentionSeconds`             | `Backup.TombstoneRetentionSeconds`            |
| `MainWorker.Backup.LogRetentionSeconds`                   | `Backup.LogRetentionSeconds`                  |
| `MainWorker.Backup.QuarantineCompactionOverrideSeconds`   | `Backup.QuarantineCompactionOverrideSeconds`  |
| `MainWorker.Backup.MaxKeyAgeSeconds`                      | `Backup.MaxKeyAgeSeconds`                     |
| `MainWorker.Backup.RotationAckTimeoutSeconds`             | `Backup.RotationAckTimeoutSeconds`            |
| `MainWorker.Backup.RotationActivationDelaySeconds`        | `Backup.RotationActivationDelaySeconds`       |
| `MainWorker.Backup.RetiredKeyGraceSeconds`                | `Backup.RetiredKeyGraceSeconds`               |
| `MainWorker.Backup.RsaKeySizeBits`                        | `Backup.RsaKeySizeBits`                       |
| `MainWorker.Backup.Endpoint.IncrementalDiffTimeoutSeconds`| `Backup.Endpoint.IncrementalDiffTimeoutSeconds` |
| `MainWorker.Backup.Endpoint.RotateKeysTimeoutSeconds`     | `Backup.Endpoint.RotateKeysTimeoutSeconds`    |
| `MainWorker.Backup.Endpoint.RestoreByDateTimeoutSeconds`  | `Backup.Endpoint.RestoreByDateTimeoutSeconds` |
| `MainWorker.Backup.Endpoint.SnapshotsTimeoutSeconds`      | `Backup.Endpoint.SnapshotsTimeoutSeconds`     |
| `MainWorker.Backup.Endpoint.HealthTimeoutSeconds`         | `Backup.Endpoint.HealthTimeoutSeconds`        |
| `MainWorker.Backup.Apply.MaxRetriesPerEnvelope`           | `Backup.Apply.MaxRetriesPerEnvelope`          |
| `MainWorker.Backup.Apply.TransactionTimeoutSeconds`       | `Backup.Apply.TransactionTimeoutSeconds`      |
| `MainWorker.Backup.Apply.DeadLetterRetentionDays`         | `Backup.Apply.DeadLetterRetentionDays`        |
| `MainWorker.Backup.Apply.IdempotencyRowRetentionDays`     | `Backup.Apply.IdempotencyRowRetentionDays`    |
| `MainWorker.Backup.SnapshotRetentionDays`                 | `Backup.SnapshotRetentionDays`                |
| `MainWorker.Backup.Snapshot.BuildHourUtc`                 | `Backup.Snapshot.BuildHourUtc`                |
| `MainWorker.Backup.Snapshot.QuiesceTimeoutSeconds`        | `Backup.Snapshot.QuiesceTimeoutSeconds`       |
| `MainWorker.Backup.Snapshot.MaxBuildSeconds`              | `Backup.Snapshot.MaxBuildSeconds`             |
| `MainWorker.Backup.Restore.PrimaryAckTimeoutSeconds`      | `Backup.Restore.PrimaryAckTimeoutSeconds`     |

Linter T3 (§6) MUST validate this alias table is exhaustive — every §2 prose row appears here and every §4 settings key appears here. Adding a tunable means updating §2, §4, and §4.1 in the same commit.

## 4.2 Caching tunables (Phase 13.2 — closes audit M-2 / M-3)

The Main-side caches that previously hard-coded TTLs in `01-architecture.md` §5 are pinned here so a literal AI does not invent values. These govern in-process caches (distinct from the cross-tier protocol budgets above) but live under the same `MainWorker` category for single-source-of-truth.

| Key | Default | Unit | Used by | Notes |
|---|---:|---|---|---|
| `MainWorker.Cache.CompanyToWorkerTtlSeconds`      | **900** (15 m) | seconds | `01-architecture.md` §5 — `CompanyId → WorkerNodeId` cache | Invalidate on worker reassignment. |
| `MainWorker.Cache.WorkerRegistryTtlSeconds`       | **60**         | seconds | `01-architecture.md` §5 — Worker registry | Invalidate on worker register/deregister. |
| `MainWorker.Cache.RecentCompanyPerUserTtlSeconds` | **= `MainWorker.Auth.MainSessionTtlSeconds`** | seconds | `01-architecture.md` §5 — Per-user recent-company | Invalidate on logout. Tied to session TTL by definition. |

---

## 5. Stale-prose corrections (follow-ups)

The following docs cite tunables inline. Each MUST be edited to cite this file instead. Tracked as FU-14:

| File | Stale phrase | Replace with |
|---|---|---|
| `04-worker-routing.md` | "retry up to 5 times" / "3 retries" | "Retries per `15-tunable-constants.md` §2.1." |
| `08-error-contract.md` | implicit retry budget | Add cite to `15` §2.1. |
| `10-worker-bootstrap-protocol.md` §3.2 | "30 s" hard-coded literal | Replace with: "default `30 s` per `15-tunable-constants.md` §2.3." |
| `spec/14-update/28-worker-push-instruction.md` §3.1 `MaxRetries: 3` | already correct | Add cross-cite to `15` §2.1. |
| `spec/14-update/28-worker-push-instruction.md` §7 `WorkerUpdateInstructionRetentionDays` | scattered | Add cross-cite to `15` §2.7. |

---

## 6. Linter assertion (FU-15 + FU-16)

`linter-scripts/check-tunable-constants.py` MUST verify:

1. **T1 (presence)** — every numeric literal in spec/19 prose that ends with `s`, `sec`, `seconds`, `min`, `minutes`, `h`, `hours`, `attempts`, `times`, or `retries` is either
   (a) named in §2 of this file, or
   (b) explicitly waivered via `<!-- TUNABLE-WAIVER: rationale -->` comment.
2. **T2 (unique keys)** — no two §2 rows share the same Key.
3. **T3 (seed parity)** — `config.seed.json` `Categories.MainWorker.Settings.*.Default` matches §4 verbatim. **As of seed v2.0.0**, this includes the full `MainWorker.Backup.*` namespace (the previous Phase-13.4 waiver is lifted). The only suppression remaining: `Backup.RetiredKeyGraceSeconds=0` may be substituted at request time when rotation `Reason="Compromise"` — this is a runtime override, not a seed default change, and T3 ignores it.
4. **T4 (session-TTL invariant — FU-16)** — `MainWorker.Auth.MainSessionAbsoluteMaxSeconds` >= `MainWorker.Auth.MainSessionTtlSeconds` in BOTH §2 catalogue defaults AND §4 seed defaults. Sliding TTL must never exceed the absolute cap (otherwise the cap is unreachable). Resolves FU-16; cited in §7.2.

Failure = build break.

---

## 7. Resolved Decisions (formerly Open Questions)

### 7.1 OQ-15-1 — Retry backoff shape: **explicit array** ✅ RESOLVED 2026-05-04

**Decision:** Keep explicit `[2, 8, 30]` array (`MainWorker.Retry.BackoffSeconds`). **Reject** `base^n` exponential.

**Why:**
- **Dumb-AI friendly** — every value is visible; no implementer needs to compute `2^n` mentally and re-derive a ceiling clamp.
- **Bounded worst case** — explicit ceiling at element `[N-1]`; exponential needs an extra `Min(base^n, Cap)` rule which is a second tunable AND a second source of disagreement (audit F-A-15 root cause was exactly this kind of derived value).
- **Operationally tweakable** — ops can paste `[5, 30, 120]` into Seedable-Config without re-reading any formula doc. Exponential requires editing `Base` AND `Cap` AND mentally reconciling.
- **Length contract** — `len(BackoffSeconds) == MaxAttempts - 1` is a single-line linter check (already in `check-tunable-constants.py`); `base^n` would need range validation per attempt index.

**Trade-off accepted:** Cannot smoothly extend to N=10+ attempts without ugly arrays. Tolerated — `MainWorker.Retry.MaxAttempts=3` is pinned and any move to N≥6 requires a v2.0 design review anyway.

**No values change.** This decision codifies the existing `[2, 8, 30]` / `[30, 120, 300]` defaults already shipped in §2.1 and `spec/14-update/28-…md` §3.1.

### 7.2 OQ-15-2 — Session TTL semantics: **sliding with absolute cap** ✅ RESOLVED 2026-05-04

**Decision:** Adopt **sliding** TTL by default (matches Laravel Sanctum / Express-Session / ASP.NET Core), bounded by an **absolute** maximum-lifetime ceiling.

**Why:**
- **Sliding** matches every default-stack target (Laravel/.NET/Express) — zero surprise for implementers.
- **Pure sliding is unbounded** — a user keeping a tab open for 90 days never re-authenticates. Compliance frameworks (SOC 2, ISO 27001 §9.4.2) require periodic re-authentication.
- **Sliding + absolute cap** is the industry compromise (used by Auth0, Okta, AWS Cognito) — refresh on activity, but force re-login at the absolute boundary regardless of activity.

**Existing tunable retained, two new ones added in §2.4 above + §4** (not re-tabled here to avoid catalogue duplication; this section explains the contract):

- `MainWorker.Auth.MainSessionTtlSeconds` (existing, **28800s/8h**) — sliding window; reset on every qualifying authenticated request.
- `MainWorker.Auth.MainSessionAbsoluteMaxSeconds` (NEW, **86400s/24h**) — hard ceiling from initial login regardless of activity; forces `Reauthenticate`. MUST be ≥ sliding TTL.
- `MainWorker.Auth.SessionSlidingExtendOnReadOnly` (NEW, **true**) — if `false`, only state-changing requests (POST/PUT/PATCH/DELETE) extend the window — mitigates background-poll abuse.

**Implementation contract (consumed by `05-auth-and-2fa.md` §6):**
1. On each authenticated request: if `(now - SessionStartedAt) >= MainSessionAbsoluteMaxSeconds` → `401 + X-Auth-Action: Reauthenticate`.
2. Else if request qualifies (write request, OR sliding-extend-on-read flag true): `SessionLastSeenAt = now`; cookie `Max-Age` reset to sliding TTL.
3. Else: leave `SessionLastSeenAt` and cookie `Max-Age` untouched.

**Linter follow-up (FU-16):** `check-tunable-constants.py` to assert sliding TTL ≤ absolute max (T4 invariant).


---

## 8. Cross-references

- `spec/06-seedable-config-architecture/01-fundamentals.md` — `Categories` block schema (§4 conforms).
- `spec/19-main-worker-service/04-worker-routing.md` — consumer (FU-14).
- `spec/19-main-worker-service/10-worker-bootstrap-protocol.md` §3.2/§7 — consumer (FU-14).
- `spec/19-main-worker-service/12-jwt-delivery-contract.md` §6 — consumer (already cites tunables but pinned here).
- `spec/14-update/28-worker-push-instruction.md` §3/§7 — consumer (FU-14).
- `spec/19-main-worker-service/13-error-codes.md` — `WORKER-300-01` triggered when `IdempotencyKeyTtlSeconds` exceeded.

---

*Tunable constants v1.10.0 — 2026-05-06 (Phase 11: §2.15 added — SnapshotRetentionDays=30 resolving OQ-A4, Snapshot.BuildHourUtc, Snapshot.QuiesceTimeoutSeconds, Snapshot.MaxBuildSeconds, Restore.PrimaryAckTimeoutSeconds)*
