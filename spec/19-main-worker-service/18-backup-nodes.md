# 18 — Backup Nodes

**Spec:** `19-main-worker-service`
**Version:** 1.0.0
**Created:** 2026-05-06
**Status:** Authoritative (Phase 6)
**Resolves:** locked decisions D8 (`IsBackup` flag), D9 (backups never serve traffic), D10 (`SyncOp` CDC flag — referenced here, defined in detail in `19-incremental-backup-sync.md`).

> **Phase 6 scope.** Define **what a backup node is**, how it **registers**, how the Main↔Worker↔Backup relationship is **propagated** to all three tiers, and the **"never serve traffic"** invariant. Wire-format and CDC mechanics are deferred to Phases 7–11.
>
> **Kubernetes mind-frame (per user verbatim).** A backup node is to its primary Worker what a Kubernetes pod replica is to its leader: same image, same data shape, **only the leader serves**. The backup is a hot mirror — not a read replica, not a load-shedding peer. Failover is an explicit Power-Admin decision, not an automatic election.

---

## 1. Definition

A **Backup Node** is a `WorkerNode` row with `IsBackup = 1` and `BackupOfWorkerNodeId = <primary>.WorkerNodeId` (per `03-main-db-schema.md` §2.1, Phase 4).

| Property | Primary Worker | Backup Node |
|---|---|---|
| `IsBackup` | `0` | `1` |
| `BackupOfWorkerNodeId` | `NULL` | NOT NULL → primary's `WorkerNodeId` |
| Serves user-facing traffic | YES | **NO** (D9 invariant) |
| Receives `Company` assignments | YES | NO |
| Holds split-DB tiers (Root, App, Session, Cache, Document) | YES | YES — but **App tier is the only one synced from primary**; Root/Settings remain backup-local; Session is **never** copied. |
| Heartbeats to Main | YES | YES (different status code, see §4) |
| Listed in `LeastLoaded` / `RoundRobin` candidate set | YES | NO (excluded by `IsPrimary` guard, `04-worker-routing.md` §1.4) |
| Accepts `POST /API/V1/Backup/*` endpoints | NO | YES (defined in `21-backup-endpoints.md`) |

A **backup chain** (backup-of-a-backup) is **forbidden** by the schema CHECK in `03-…` §2.1 and re-asserted here: any attempt to register a backup whose `BackupOfWorkerNodeId` points at another `IsBackup = 1` row MUST be rejected with `MAIN-800-01 BackupChainNotAllowed` (new code, see §6).

---

## 2. The Three-Tier Relationship Model

```
                ┌──────────────────────┐
                │        MAIN          │
                │  (catalog + routing) │
                └─────────┬────────────┘
                          │  registers, broadcasts
                          ▼
        ┌────────────── WorkerNode (IsBackup=0) ───────────────┐
        │                                                       │
        │  serves traffic ◄──── tenants                         │
        │  owns App tier (authoritative)                        │
        │                                                       │
        └────────────────────┬──────────────────────────────────┘
                             │  encrypted incremental diff (Phase 7+8)
                             ▼
                    WorkerNode (IsBackup=1)
                    BackupOfWorkerNodeId → primary
                    receives + stores ONLY
```

Three relationship facts that the spec MUST guarantee:

| # | Fact | Enforced by |
|---|---|---|
| R1 | Every backup has exactly one primary. | Schema CHECK in `03-…` §2.1. |
| R2 | A primary MAY have **0..N** backups. N is bounded only by `MainWorker.Backup.MaxBackupsPerPrimary` (new tunable, default `3`). | `15-…` §2.11 (added in this phase). |
| R3 | The primary is the **only** writer of its App tier. Backups MUST refuse any inbound write that did not arrive via the `/API/V1/Backup/*` channel. | `21-backup-endpoints.md` (Phase 9) auth contract; reinforced by §5 below. |

---

## 3. Registration Flow

Backup registration extends the existing Worker bootstrap protocol (`10-worker-bootstrap-protocol.md`) with two extra fields. **No new endpoint** — backups boot through `POST /API/V1/Workers/Register` like any other Worker.

### 3.1 Request additions

The base bootstrap request body is **canonically defined in `10-worker-bootstrap-protocol.md` §3.1** (fields `WorkerNodeDisplayName`, `WorkerEndpointPublic`, `WorkerVersionPin`, `BootInstanceUlid`, `Description`). A backup node MUST send the same body and append exactly two extra fields:

```jsonc
{
  // ... all base fields per `10-worker-bootstrap-protocol.md` §3.1 ...
  "IsBackup":                true,
  "BackupOfWorkerIdentity":  "<WorkerNodeIdentity of the primary>"
}
```

`BackupOfWorkerIdentity` is the **`WorkerNodeIdentity`** of the primary (per `03-main-db-schema.md` §2.1), not the numeric ID — the booting backup does not yet know Main's PKs. Main resolves it server-side. Field names from earlier drafts (`WorkerNodeCode`, `WorkerNodeName`, `AdvertisedHost`, `AdvertisedPort`, `WorkerVersionCode`, `Capabilities[]`) are obsolete; use the `10-` §3.1 names verbatim.

### 3.2 Main-side acceptance procedure (CODE RED, ≤15 lines, positive guards)

```
function RegisterBackup(request):
    primary := WorkerNode.FindByIdentity(request.BackupOfWorkerIdentity)
    AssertPrimaryExists(primary)              -- else MAIN-800-02 PrimaryNotFound
    AssertPrimaryIsLeader(primary)            -- IsBackup = 0; else MAIN-800-01 BackupChainNotAllowed
    AssertCapacityForBackup(primary)          -- count(IsBackup=1, Of=primary.Id) < Tunable; else MAIN-800-03 BackupCapacityExceeded
    backup := WorkerNode.Insert({
        ...request,
        IsBackup: 1,
        BackupOfWorkerNodeId: primary.WorkerNodeId,
        Sequence: primary.Sequence,           -- backups inherit primary Sequence for grouping
        WorkerNodeStatusId: Status('Provisioning').Id
    })
    BroadcastBackupRegistered(primary, backup)  -- §4
    return MintBootstrapResponse(backup)
```

### 3.3 Response

Standard `10-…` `201 Created` body, plus:

```json
{
  "WorkerNodeId": 42,
  "JwtPublicKeyPem": "...",
  "AssignedSubdomain": null,
  "HeartbeatIntervalSeconds": 30,
  "BackupOfWorkerNodeId": 7,
  "BackupSyncSourceUrl": "https://w7.example.com/API/V1/Backup/Outbound"
}
```

`AssignedSubdomain` is `null` for backups — they have no public surface.

---

## 4. Propagation: Main → Primary, Main → Backup

After §3.2 completes, Main MUST notify both ends of the new pairing so the primary knows where to send diffs and the backup knows what to expect.

### 4.1 Main → Primary

**Endpoint:** `POST /API/V1/Worker/Inbound/BackupAttached` (already-defined Worker S2S surface; new instruction kind).

```json
{
  "InstructionKind": "BackupAttached",
  "BackupWorkerNodeId": 42,
  "BackupEndpointPublic": "https://w7-backup-1b.example.com",
  "BackupPublicKeyPem": "...",
  "AttachedAt": 1730000000,
  "CorrelationId": "01J9..."
}
```

The primary appends a row to its **`KnownBackupNode`** mirror table (§5).

### 4.2 Main → Backup

Sent as part of the §3.3 bootstrap response **and** re-sent on any rotation event (Phase 8). The backup persists its mirror of the primary's identity and public key.

### 4.3 Heartbeat status codes

`WorkerNodeStatus` (`14-rbac-and-status-seed.md` §3) is extended with two backup-specific codes (seedable, additive):

| `Code` | `Label` | Applies to | Meaning |
|---|---|---|---|
| `BackupAttached` | "Backup Attached" | `IsBackup = 1` rows | Healthy, last diff applied successfully. |
| `BackupLagging` | "Backup Lagging" | `IsBackup = 1` rows | Diff lag exceeds `MainWorker.Backup.LagWarningSeconds` (new tunable, default 900 s). Still attached, no auto-failover. |

Existing `Active` / `Draining` / `Offline` / `Quarantined` continue to apply to **primaries only**. A migration script MUST refuse to assign `Active` to a row where `IsBackup = 1`.

---

## 5. Worker-Side Mirror: `KnownBackupNode`

Each primary Worker keeps a small App-tier table listing the backups it must push diffs to. This is **not** authoritative — Main is — but it lets the primary operate during a temporary Main outage.

```sql
-- Worker, App tier (per 11-split-db-tier-reconciliation.md §5)
CREATE TABLE KnownBackupNode (
    KnownBackupNodeId       INTEGER PRIMARY KEY AUTOINCREMENT,
    BackupWorkerNodeId      INTEGER NOT NULL UNIQUE,    -- Main's PK, snapshotted
    BackupEndpointPublic    TEXT    NOT NULL,
    BackupPublicKeyPem      TEXT    NOT NULL,           -- RSA public key (Phase 8 detail)
    AttachedAt              INTEGER NOT NULL,           -- epoch seconds, UTC
    LastSuccessfulSyncAt    INTEGER NULL,               -- epoch seconds, UTC; NULL until first diff
    LastSyncWatermark       INTEGER NULL,               -- last `SyncOpSeq` shipped (Phase 7)
    StatusCode              TEXT    NOT NULL,           -- mirror of `WorkerNodeStatus.Code`
    Description             TEXT    NULL                -- Rule 11
);
CREATE INDEX IX_KnownBackupNode_StatusCode ON KnownBackupNode(StatusCode);
```

**Conventions applied:**
- PK `KnownBackupNodeId` per Rule 1 (universal).
- `BackupEndpointPublic`, `BackupPublicKeyPem`, `StatusCode` are NOT enum-FK refs — they are **snapshots** kept locally so the Worker remains operational without Main. The authoritative copy lives on Main; reconciliation occurs on the next Main-initiated `BackupAttached` / `BackupRotated` instruction.
- `LastSyncWatermark` is the bridge into Phase 7 (`SyncOp` CDC).

**Mutation rules:**
- INSERT only on `BackupAttached` instruction.
- UPDATE on `BackupRotated` (Phase 8) or after each successful diff push.
- DELETE on `BackupDetached` (instruction added in Phase 9 endpoint set).

---

## 6. New Error Codes (additions to `13-error-codes.md`)

A new sub-range **Backup Lifecycle** is opened on the Main tier. Worker-side codes for actually applying diffs are reserved for Phase 10.

| Code | Flat | Name | Message | HTTP | Source |
|---|---|---|---|---|---|
| `MAIN-800-01` | `21181` | `BackupChainNotAllowed` | "Cannot register a backup whose primary is itself a backup." | 422 | `18-backup-nodes.md` §1, §3.2 |
| `MAIN-800-02` | `21182` | `PrimaryNotFound` | "`BackupOfWorkerIdentity` does not resolve to any registered Worker." | 404 | `18-…` §3.2 |
| `MAIN-800-03` | `21183` | `BackupCapacityExceeded` | "Primary already has `MainWorker.Backup.MaxBackupsPerPrimary` backups." | 409 | `18-…` §2 R2 |
| `MAIN-800-04` | `21184` | `TrafficOnBackupRejected` | "Backup node rejected an inbound user-facing request (D9 invariant)." | 421 | `18-…` §1 |

Reserved-range table in `13-…` §4 must mark `21181-21184` as consumed; `MAIN-21172-21180` and `MAIN-21185-21199` remain available for future expansion. The flat numbers `21181-21184` were chosen to live above the previously-used `21171` (Phase 5) so there is no collision with the `MAIN-400-10` overflow at `21170`.

---

## 7. New Tunables (additions to `15-tunable-constants.md`)

New §2.11 "Backup nodes":

| Key | Default | Unit | Used by | Notes |
|---|---:|---|---|---|
| `MainWorker.Backup.MaxBackupsPerPrimary` | **3** | count | `18-…` §2 R2, §3.2 `AssertCapacityForBackup` | Hard ceiling enforced at registration time. |
| `MainWorker.Backup.LagWarningSeconds` | **900** (15m) | seconds | `18-…` §4.3 (`BackupLagging` status threshold) | Heartbeat watcher flips `BackupAttached` → `BackupLagging` when `now - LastSuccessfulSyncAt > value`. |
| `MainWorker.Backup.HeartbeatIntervalSeconds` | **60** | seconds | Backup → Main heartbeat cadence | Slower than primary heartbeat (`MainWorker.Heartbeat.IntervalSeconds = 30`) because backups never serve traffic. |

Phase 7 will add `MainWorker.Backup.SyncIntervalSeconds`; Phase 8 adds key-rotation tunables; Phase 11 adds `MainWorker.Backup.SnapshotRetentionDays` (resolves OQ-A4).

---

## 8. The "Never Serve Traffic" Invariant (D9)

Three independent enforcement points, defence in depth:

| Layer | Enforcement |
|---|---|
| **Routing (Main)** | `04-worker-routing.md` §1.4 `IsPrimary(node)` guard excludes backups from every selection strategy, including `Manual`. Manual attempts return `WORKER-300-04 BackupNotRoutable` (Phase 4). |
| **Reverse-proxy (Main edge)** | Main MUST NOT proxy any user-facing path (everything except `/API/V1/Backup/*`) to a `WorkerNode` row where `IsBackup = 1`. Violation logs `MAIN-800-04 TrafficOnBackupRejected` and returns HTTP `421 Misdirected Request`. |
| **Application (Backup)** | Backup nodes boot with their public router **disabled**; only the `/API/V1/Backup/*` surface is mounted. A grep over the backup binary's route table for any non-`/Backup/` mount is a release-blocker. |

`421 Misdirected Request` is the correct HTTP status: the request reached a node that cannot service it, and the client (or upstream proxy) should retry against a different node — exactly the semantics RFC 7540 §9.1.2 defines.

---

## 9. What Phase 6 Does NOT Define

Listed here so reviewers do not look for it in this file:

- **The diff wire format / `SyncOp` mechanics** → `19-incremental-backup-sync.md` (Phase 7).
- **Encryption, key rotation, zip password** → `20-backup-encryption-and-keys.md` (Phase 8).
- **Endpoint URLs, request/response bodies, idempotency keys** → `21-backup-endpoints.md` (Phase 9).
- **Apply algorithm (decrypt → unzip → SyncOp dispatch)** → Phase 10.
- **Snapshot storage, retention, restore-by-date** → `23-snapshot-storage-and-restore.md` (Phase 11).

---

## 10. Cross-References

- `03-main-db-schema.md` §2.1 — `WorkerNode.IsBackup` / `BackupOfWorkerNodeId` / `Sequence` columns.
- `04-worker-routing.md` §1.4 — `IsPrimary(node)` eligibility guard.
- `10-worker-bootstrap-protocol.md` — host bootstrap flow this phase extends.
- `11-split-db-tier-reconciliation.md` §5 — Worker App tier (where `KnownBackupNode` lives).
- `13-error-codes.md` §3.8 (added in this phase) — `MAIN-800-*` series.
- `14-rbac-and-status-seed.md` §3 — `WorkerNodeStatus` seed extended with `BackupAttached` and `BackupLagging`.
- `15-tunable-constants.md` §2.11 — backup tunables.

---

*Backup Nodes v1.0.0 — 2026-05-06 (Phase 6)*
