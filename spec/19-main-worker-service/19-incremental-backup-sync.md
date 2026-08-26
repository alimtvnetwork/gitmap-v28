# 19 — Incremental Backup Sync (CDC)

**Spec:** `19-main-worker-service`
**Version:** 1.0.0
**Created:** 2026-05-06
**Status:** Authoritative (Phase 7)
**Resolves:** locked decision D10 (`SyncOp` flag on synced rows). Builds on Phase 6 (`18-backup-nodes.md`) which defined the backup-node concept and `KnownBackupNode` mirror.

> **Phase 7 scope.** Define **how** the primary Worker captures changes to its App-tier rows and ships a deterministic, replayable diff to each attached backup. Encryption / zip-password / endpoints / apply-side semantics are deferred to Phases 8–10.
>
> **Mind-frame.** This is **change-data-capture (CDC)** by sequence number, not log shipping. Every mutating statement on an App-tier table writes a `SyncOp` row to the change log. The diff generator drains that log between two watermarks, packages it, and hands it to Phase 8 for encryption. Idempotency is keyed by `(SourceTable, SourceRowId, SyncOpSeq)` so a backup can safely re-receive the same envelope.

---

## 1. The `SyncOp` Convention (D10)

Every App-tier table whose rows must mirror to backups carries a logical `SyncOp` lifecycle. Two implementation shapes are defined below; Workers MUST pick one **per database**, not per table.

### 1.0 Default for v1.0 (binding)

**Default for v1.0 is Shape A (inline column).** Workers MUST adopt Shape A unless an existing production table cannot accept three new columns without unacceptable downtime; in that exceptional case the Worker MAY adopt Shape B and MUST record the choice in its Seedable-Config under `MainWorker.Backup.SyncOpShape` (allowed values: `Inline` | `SideLog`). A blind-AI implementer MUST treat this as an unconditional MUST and choose Shape A.

### 1.1 Shape A — Inline column (recommended for small tables)

Add three columns to the table itself:

| Column | Type | Null | Notes |
|---|---|---|---|
| `SyncOpCode` | TEXT | NO | One of `Insert`, `Update`, `Delete`. `Delete` is a **soft-tombstone**: the row is retained until the next compaction window so backups can replay the deletion. |
| `SyncOpSeq` | INTEGER | NO | Monotonic per-database sequence (see §2). Set on every write. |
| `SyncOpAt` | INTEGER | NO | Epoch seconds, UTC. |

Pros: zero join cost when reading recent changes. Cons: schema noise; tombstones occupy table space.

### 1.2 Shape B — Side log (recommended for hot/large tables)

Leave the source table unchanged; emit one row per mutation into a sibling log:

```sql
-- Worker, App tier
CREATE TABLE BackupSyncLog (
    BackupSyncLogId   INTEGER PRIMARY KEY AUTOINCREMENT,
    SourceTable       TEXT    NOT NULL,           -- e.g. 'AppUser'
    SourceRowId       INTEGER NOT NULL,           -- the source PK value
    SyncOpCode        TEXT    NOT NULL,           -- Insert | Update | Delete
    SyncOpSeq         INTEGER NOT NULL UNIQUE,    -- DB-wide monotonic (§2)
    RowSnapshotJson   TEXT    NULL,               -- full post-mutation row as JSON; NULL on Delete
    SyncOpAt          INTEGER NOT NULL,           -- epoch seconds, UTC
    Description       TEXT    NULL                -- Rule 11
);
CREATE INDEX IX_BackupSyncLog_SourceTable_RowId ON BackupSyncLog(SourceTable, SourceRowId);
CREATE INDEX IX_BackupSyncLog_SyncOpSeq ON BackupSyncLog(SyncOpSeq);
```

**Conventions applied:** PK `BackupSyncLogId` per Rule 1. `SyncOpCode` stays a TEXT enum (not an FK ref) because the value set is closed at three and a join on every change-log read would dominate diff-build cost — explicit waiver, justified inline.

Pros: zero impact on the source table; supports large rows. Cons: one extra write per mutation.

### 1.3 Required: `SyncOp` ref table (catalog only)

Whichever shape is used, the catalog of permissible codes is seeded centrally so linters can validate values:

```sql
-- Worker, App tier
CREATE TABLE SyncOp (
    SyncOpId    INTEGER PRIMARY KEY AUTOINCREMENT,
    Code        TEXT NOT NULL UNIQUE,    -- Insert | Update | Delete
    Label       TEXT NOT NULL,
    Description TEXT NULL
);
```

Seeded via Seedable-Config (3 rows). The `BackupSyncLog.SyncOpCode` and inline `SyncOpCode` columns are **not** FK-bound to this table to avoid a join on every diff read; the linter validates code values against this catalog at boot.

---

## 2. The `BackupSyncWatermark` Table (per-primary)

The primary Worker tracks, **per attached backup**, how far that backup has acknowledged. This is the "where did we last leave off" pointer; it is independent of `KnownBackupNode.LastSuccessfulSyncAt` (Phase 6) which is a wall-clock observation.

```sql
-- Worker, App tier (Phase 7 addition)
CREATE TABLE BackupSyncWatermark (
    BackupSyncWatermarkId  INTEGER PRIMARY KEY AUTOINCREMENT,
    BackupWorkerNodeId     INTEGER NOT NULL UNIQUE,    -- one row per attached backup
    LastShippedSyncOpSeq   INTEGER NOT NULL DEFAULT 0, -- highest seq put on the wire
    LastAckedSyncOpSeq     INTEGER NOT NULL DEFAULT 0, -- highest seq the backup ACKed (Phase 9)
    LastAttemptedAt        INTEGER NULL,
    LastSuccessfulAt       INTEGER NULL,
    ConsecutiveFailureCount INTEGER NOT NULL DEFAULT 0,
    Description            TEXT NULL                    -- Rule 11
);
CREATE INDEX IX_BackupSyncWatermark_LastAckedSeq ON BackupSyncWatermark(LastAckedSyncOpSeq);
```

**Sequence allocation.** `SyncOpSeq` is a **single per-database monotonic counter** allocated by a tiny helper table:

```sql
CREATE TABLE BackupSyncSequence (
    BackupSyncSequenceId INTEGER PRIMARY KEY AUTOINCREMENT,
    NextSeq              INTEGER NOT NULL DEFAULT 1,
    Description          TEXT NULL
);
INSERT INTO BackupSyncSequence (NextSeq) VALUES (1);
```

The allocator MUST be wrapped in the same transaction as the source mutation so a row never escapes without a sequence number. `INTEGER` in SQLite is 64-bit signed — at one mutation per millisecond, exhaustion is ≈292 million years away; no rollover handling required.

---

## 3. Diff-Generation Algorithm

Run on the **primary** Worker on a cron (default `MainWorker.Backup.SyncIntervalSeconds = 60`). One pass per attached backup; passes for different backups MAY run in parallel but MUST NOT share envelopes.

### 3.1 Top-level driver (CODE RED ≤15 lines, positive guards)

```
function ShipDiffsForAllBackups():
    backups := KnownBackupNode.ListAttached()
    foreach b in backups:
        ShipDiffForBackup(b)

function ShipDiffForBackup(b):
    watermark := BackupSyncWatermark.FindOrCreate(b.BackupWorkerNodeId)
    AssertWatermarkConsistent(watermark)        -- LastAcked <= LastShipped
    rows := ReadChangeLogAfter(watermark.LastAckedSyncOpSeq)
    if IsEmptyDiff(rows) -> return
    envelope := PackageEnvelope(rows)
    Phase8_EncryptAndZip(envelope, b)            -- handed off to Phase 8
    Phase9_PostIncrementalDiff(envelope, b)      -- handed off to Phase 9
    BumpLastShipped(watermark, MaxSeq(rows))
```

**Resume from `LastAcked`, not `LastShipped`.** A previous envelope that was sent but never ACKed (Worker crash mid-flight, network drop) MUST be re-sent. Idempotency on the backup side (Phase 10) deduplicates by `(SourceTable, SourceRowId, SyncOpSeq)`.

### 3.2 Read-after-watermark (Shape B example)

```sql
SELECT BackupSyncLogId, SourceTable, SourceRowId, SyncOpCode,
       SyncOpSeq, RowSnapshotJson, SyncOpAt
  FROM BackupSyncLog
 WHERE SyncOpSeq > :LastAckedSyncOpSeq
 ORDER BY SyncOpSeq ASC
 LIMIT  :MainWorker_Backup_MaxRowsPerEnvelope;     -- new tunable, §6
```

For Shape A, the equivalent query unions over each tracked source table:

```sql
SELECT 'AppUser' AS SourceTable, AppUserId AS SourceRowId,
       SyncOpCode, SyncOpSeq, SyncOpAt
  FROM AppUser
 WHERE SyncOpSeq > :LastAckedSyncOpSeq
UNION ALL
SELECT ...
ORDER BY SyncOpSeq ASC
LIMIT :MainWorker_Backup_MaxRowsPerEnvelope;
```

The set of **tracked source tables** (Shape A only) is declared once in the seed file (Phase 12 follow-up) so the diff generator can build the UNION ALL deterministically.

### 3.3 Envelope structure (handed to Phase 8)

A diff envelope is a small **SQLite database file** named `diff-{LastAckedSyncOpSeq+1}-{MaxSyncOpSeq}.sqlite` containing exactly two tables:

```sql
CREATE TABLE Envelope (
    EnvelopeId        INTEGER PRIMARY KEY AUTOINCREMENT,
    PrimaryWorkerNodeId INTEGER NOT NULL,
    BackupWorkerNodeId  INTEGER NOT NULL,
    LowSyncOpSeq      INTEGER NOT NULL,     -- inclusive
    HighSyncOpSeq     INTEGER NOT NULL,     -- inclusive
    RowCount          INTEGER NOT NULL,
    GeneratedAt       INTEGER NOT NULL,     -- epoch seconds, UTC
    SchemaVersion     TEXT NOT NULL,        -- e.g. '1.0.0'
    Description       TEXT NULL
);

CREATE TABLE EnvelopeRow (
    EnvelopeRowId    INTEGER PRIMARY KEY AUTOINCREMENT,
    SourceTable      TEXT NOT NULL,
    SourceRowId      INTEGER NOT NULL,
    SyncOpCode       TEXT NOT NULL,
    SyncOpSeq        INTEGER NOT NULL UNIQUE,
    RowSnapshotJson  TEXT NULL,
    SyncOpAt         INTEGER NOT NULL,
    Description      TEXT NULL
);
CREATE INDEX IX_EnvelopeRow_SourceTable_RowId ON EnvelopeRow(SourceTable, SourceRowId);
```

Why a SQLite file inside the envelope? It gives backups a free integrity check (the file opens or it doesn't), supports queries during apply, and ships through any blob transport (zipped + encrypted by Phase 8) without bespoke parsing.

---

## 4. Compaction (tombstone & log retention)

Without compaction the change log grows unbounded. Two policies, both Worker-local:

| Policy | Trigger | Action |
|---|---|---|
| **Tombstone reclaim** (Shape A only) | Hourly cron. | `DELETE FROM <table> WHERE SyncOpCode = 'Delete' AND SyncOpAt < (Now - MainWorker.Backup.TombstoneRetentionSeconds) AND SyncOpSeq <= MIN(BackupSyncWatermark.LastAckedSyncOpSeq);` |
| **Log truncation** (Shape B only) | Hourly cron. | `DELETE FROM BackupSyncLog WHERE SyncOpSeq <= MIN(BackupSyncWatermark.LastAckedSyncOpSeq) AND SyncOpAt < (Now - MainWorker.Backup.LogRetentionSeconds);` |

The `MIN(LastAckedSyncOpSeq)` clause is the critical safety net: a row is reclaimed only when **every** attached backup has acknowledged past it. A single quarantined backup cannot stall compaction beyond `MainWorker.Backup.QuarantineCompactionOverrideSeconds` (default 24 h) — after which Main MUST be notified via `MAIN-810-01 BackupCompactionStalled` (new code, §5) and the operator decides whether to detach the lagging backup.

---

## 5. New Error Codes (additions to `13-error-codes.md`)

A new sub-range `MAIN-810` for sync-pipeline errors observed at Main, and a new `WORKER-910` sub-range for diff-generation errors observed at the primary.

| Code | Flat | Name | Message | HTTP | Source |
|---|---|---|---|---|---|
| `WORKER-910-01` | `21092` | `BackupSyncWatermarkInconsistent` | "`BackupSyncWatermark.LastAckedSyncOpSeq > LastShippedSyncOpSeq` (impossible state)." | 500 | `19-incremental-backup-sync.md` §3.1 |
| `WORKER-910-02` | `21093` | `BackupEnvelopeBuildFailed` | "Diff envelope SQLite file could not be created." | 500 | `19-incremental-backup-sync.md` §3.3 |
| `WORKER-910-03` | `21094` | `BackupChangeLogQueryFailed` | "Read-after-watermark query failed." | 500 | `19-incremental-backup-sync.md` §3.2 |
| `MAIN-810-01` | `21185` | `BackupCompactionStalled` | "Backup acknowledgement lag exceeds `QuarantineCompactionOverrideSeconds`; compaction blocked." | n/a | `19-incremental-backup-sync.md` §4 |

Reserved-range table (§4 of `13-…`) gains:
- `WORKER-21092-21094` consumed by `WORKER-910-01..03` (Phase 7).
- `21185` consumed by `MAIN-810-01` (Phase 7).
- `MAIN-21186-21199` remains for future expansion.

---

## 6. New Tunables (additions to `15-tunable-constants.md` §2.11)

Appended to the §2.11 "Backup nodes" block opened in Phase 6:

| Key | Default | Unit | Used by | Notes |
|---|---:|---|---|---|
| `MainWorker.Backup.SyncIntervalSeconds` | **60** | seconds | `19-…` §3.1 cron driver | Per-backup pass cadence on the primary. Lower = lower lag, higher = lower CPU/IO. |
| `MainWorker.Backup.MaxRowsPerEnvelope` | **5000** | count | `19-…` §3.2 read query | Hard ceiling per envelope; prevents oversized blobs on bulk writes. |
| `MainWorker.Backup.TombstoneRetentionSeconds` | **604800** (7 d) | seconds | `19-…` §4 (Shape A compaction) | Minimum age before a delivered tombstone is reclaimed. |
| `MainWorker.Backup.LogRetentionSeconds` | **604800** (7 d) | seconds | `19-…` §4 (Shape B compaction) | Minimum age before a delivered `BackupSyncLog` row is truncated. |
| `MainWorker.Backup.QuarantineCompactionOverrideSeconds` | **86400** (24 h) | seconds | `19-…` §4 stall guard | Operator escalation point; never auto-deletes the backup. |

Phase 8 will add encryption tunables; Phase 9 adds endpoint timeouts; Phase 11 adds `MainWorker.Backup.SnapshotRetentionDays` (resolves OQ-A4).

---

## 7. Linter Hooks (Phase-12 punch list)

Two new linter rules to be added in Phase 12 alongside the diagram regen:

1. **DB-SYNCOP-001** — every App-tier table declared in the "tracked tables" seed MUST either (a) carry the inline `SyncOpCode` / `SyncOpSeq` / `SyncOpAt` triple, **or** (b) have at least one trigger or application-layer write path that emits a `BackupSyncLog` row. Linter scans schema + a hand-maintained allowlist.
2. **DB-SYNCOP-002** — `SyncOpCode` values observed in `BackupSyncLog` (or the inline column) MUST be a subset of `SyncOp.Code` seed values. Catches code drift.

---

## 8. What Phase 7 Does NOT Define

- **Encryption / zip-password / RSA key handling** → Phase 8 (`20-backup-encryption-and-keys.md`).
- **HTTP envelope, idempotency-key derivation, ACK semantics on the wire** → Phase 9 (`21-backup-endpoints.md`). The §3.1 driver hands off; the wire contract is Phase 9's job.
- **Apply-side dispatch (decrypt → unzip → SyncOp dispatch)** → Phase 10.
- **Date-named full snapshots and restore-by-date** → Phase 11.

---

## 9. Cross-References

- `03-main-db-schema.md` §2.1 — `WorkerNode.IsBackup` / `BackupOfWorkerNodeId`.
- `11-split-db-tier-reconciliation.md` §5 — Worker App tier (where `BackupSyncLog`, `BackupSyncWatermark`, `BackupSyncSequence`, `SyncOp` live).
- `13-error-codes.md` §2.10 (`WORKER-910-*`) and §3.8 (`MAIN-810-*`) — added in this phase.
- `15-tunable-constants.md` §2.11 — extended in this phase.
- `18-backup-nodes.md` §5 — `KnownBackupNode` mirror (the `LastSyncWatermark` field is an aggregate view of `BackupSyncWatermark.LastAckedSyncOpSeq` for human dashboards).

---

*Incremental Backup Sync (CDC) v1.0.0 — 2026-05-06 (Phase 7)*
