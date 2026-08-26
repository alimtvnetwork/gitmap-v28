# 22 — Backup Apply Logic

**Spec:** `19-main-worker-service`
**Version:** 1.0.0
**Created:** 2026-05-06
**Status:** Authoritative (spec-only, per plan §Mode)
**Resolves:** Phase 10 — what BE-1 (`21-backup-endpoints.md` §4) actually does **after** the sealed envelope is received.
**Depends on:** [`19-incremental-backup-sync.md`](./19-incremental-backup-sync.md) §6 (envelope SQLite shape), [`20-backup-encryption-and-keys.md`](./20-backup-encryption-and-keys.md) §9 (decrypt path), [`21-backup-endpoints.md`](./21-backup-endpoints.md) §4, [`13-error-codes.md`](./13-error-codes.md), [`15-tunable-constants.md`](./15-tunable-constants.md).

---

## Keywords

`backup-apply` · `decrypt-pipeline` · `idempotency-key` · `dispatch-by-syncop` · `dead-letter-queue` · `transaction-boundary`

---

## 1. Purpose

Pin the **server-side processing pipeline** that runs on the backup node once BE-1 has accepted a sealed envelope. Wire format is owned by Phase 9; encryption is owned by Phase 8; CDC source-side is owned by Phase 7. This file is the **consumer** that turns one sealed zip into idempotent SQLite mutations on the backup's mirrored App tier.

---

## 2. Pipeline overview

Five stages, strictly sequential, no nesting (CODE RED).

```
Stage 1  Decrypt        — per 20-…  §9 steps 1..7
Stage 2  Open envelope  — sqlite_open(envelope.sqlite, READONLY)
Stage 3  Validate       — schema sanity, AAD↔meta cross-check, idempotency lookup
Stage 4  Dispatch       — per-row apply by SyncOpCode inside one TX per envelope
Stage 5  Persist ACK    — update BackupSyncWatermark, mark idempotency row, return BE-1 200
```

Failure on any stage short-circuits to **§6 Failure handling**. Partial state MUST NOT leak — Stage 4 runs in a single SQLite transaction; Stage 5 is the only commit point.

---

## 3. Stage 3 — Validation rules

Run **before** opening the apply transaction. All checks return early with a named error; none allow fallback.

| ID | Rule | Failure code |
|---|---|---|
| V1 | `Envelope.PrimaryWorkerNodeId` matches BE-1 header `X-Primary-Worker-Identity`. | `MAIN-800-04 TrafficOnBackupRejected` |
| V2 | `Envelope.BackupWorkerNodeId` equals self. | `WORKER-300-04 BackupNotRoutable` |
| V3 | `Envelope.KeyEpoch` matches the `BackupKeyEpoch` row used to decrypt (no epoch mixing). | `WORKER-920-03 UnknownKeyEpoch` |
| V4 | `Envelope.RowCount` equals `SELECT COUNT(*) FROM EnvelopeRow`. | `WORKER-930-01 EnvelopeRowCountMismatch` |
| V5 | Every `EnvelopeRow.SourceTable` appears in the **tracked-table allowlist** (`AppBackupTrackedTable` ref, Phase 12 seed). | `WORKER-930-02 UntrackedSourceTable` |
| V6 | `Envelope.MinSyncOpSeq` strictly greater than `BackupSyncWatermark.LastAcceptedSyncOpSeq` for this primary. | `WORKER-910-01 BackupSyncWatermarkInconsistent` |
| V7 | `EnvelopeId` not already present in `BackupApplyIdempotency` with `Status='Applied'`. | Treat as replay → return original 200 (per §7). |

V1–V3 protect against operator misrouting; V4–V5 protect against truncated or schema-drift envelopes; V6 protects against out-of-order delivery; V7 is the idempotency short-circuit.

---

## 4. Stage 4 — Dispatch

One SQLite TX per envelope (`BEGIN IMMEDIATE` on the App tier). Iterate `EnvelopeRow ORDER BY SyncOpSeq ASC`. Row dispatch:

| `SyncOpCode` | Action on backup mirror |
|---|---|
| `Insert` | `INSERT … ON CONFLICT(SourceTable PK) DO UPDATE SET …` (upsert; tolerates redelivery within a single envelope retry). |
| `Update` | Same upsert path — primary's view is authoritative; backup never branches. |
| `Delete` | `DELETE FROM <SourceTable> WHERE <PkColumn> = <SourceRowId>`; missing row is **not** an error (idempotent). |

CODE-RED-compliant per-row pseudocode (≤15 lines, zero nesting, ≤2 operands per guard):

```
ApplyOneRow(row):
    AssertKnownSyncOp(row.SyncOpCode)            // guard, returns or raises
    AssertKnownTable(row.SourceTable)            // guard
    AssertNonEmptyPk(row.SourceRowId)            // guard
    DispatchInsertOrUpdate(row)  if row.SyncOpCode != 'Delete'
    DispatchDelete(row)          if row.SyncOpCode == 'Delete'
    RecordRowAppliedSeq(row.SyncOpSeq)
```

Constraints:

- **No nested `if`** — guards are positively named functions (`AssertKnownSyncOp`, `AssertKnownTable`, `AssertNonEmptyPk`).
- **Max 2 operands per boolean** — split compound conditions into named guards.
- **No swallowed errors** — each guard either passes or raises a typed error; no silent skips.

---

## 5. New tables (App tier on backup; mirrored conceptually on Main only as views)

### 5.1 `BackupApplyIdempotency`

```
CREATE TABLE BackupApplyIdempotency (
    BackupApplyIdempotencyId  INTEGER PRIMARY KEY AUTOINCREMENT,
    EnvelopeId                TEXT    NOT NULL,
    PrimaryWorkerNodeId       INTEGER NOT NULL,
    KeyEpoch                  INTEGER NOT NULL,
    MinSyncOpSeq              INTEGER NOT NULL,
    MaxSyncOpSeq              INTEGER NOT NULL,
    RowCount                  INTEGER NOT NULL,
    Status                    TEXT    NOT NULL,        -- Received | Applied | DeadLettered
    ReceivedAtEpoch           INTEGER NOT NULL,
    AppliedAtEpoch            INTEGER NULL,
    OriginalResponseJson      TEXT    NULL,            -- replay payload for V7 short-circuit
    Notes                     TEXT    NULL,            -- transactional table → Rule 11
    Comments                  TEXT    NULL,            -- transactional → Rule 11
    UNIQUE (EnvelopeId, PrimaryWorkerNodeId)
);
```

### 5.2 `BackupApplyDeadLetter`

```
CREATE TABLE BackupApplyDeadLetter (
    BackupApplyDeadLetterId   INTEGER PRIMARY KEY AUTOINCREMENT,
    EnvelopeId                TEXT    NOT NULL,
    PrimaryWorkerNodeId       INTEGER NOT NULL,
    FailedAtSyncOpSeq         INTEGER NULL,            -- NULL = pre-Stage-4 failure
    FailureCode               TEXT    NOT NULL,        -- WORKER-9xx-yy / MAIN-8xx-yy
    FailureMessage            TEXT    NOT NULL,
    EnvelopeBlobPath          TEXT    NOT NULL,        -- relative path under var/dlq/
    FirstSeenAtEpoch          INTEGER NOT NULL,
    LastSeenAtEpoch           INTEGER NOT NULL,
    SeenCount                 INTEGER NOT NULL,
    Notes                     TEXT    NULL,
    Comments                  TEXT    NULL,
    UNIQUE (EnvelopeId, PrimaryWorkerNodeId)
);
```

Memory compliance:

- PascalCase, `{TableName}Id` PK (Memory: DB Schema).
- Transactional tables → `Notes` + `Comments TEXT NULL`, no DEFAULT (Memory: Rule 11).
- INTEGER epoch-seconds for every `*At` (D2).

---

## 6. Failure handling

> **Blind-AI MUST read §6.4 first.** Per-envelope `PRAGMA journal_mode=WAL` is **forbidden** — the App tier sets WAL once per session and that mode is authoritative. A literal AI that flips WAL inside Stage 4 will defeat the App tier's shared cache and regress to Phase-13.2-era contradictions. (Phase 13.3 callout — closes audit-08 §2.5 trap #2.)

### 6.1 Pre-Stage-4 failures (V1–V6, decrypt errors)

- Insert a `BackupApplyDeadLetter` row with `FailedAtSyncOpSeq = NULL`.
- Move the raw zip to `var/dlq/<EnvelopeId>.zip`.
- Return the named error to BE-1 caller. Primary will retry per `MainWorker.Backup.SyncIntervalSeconds`.

### 6.2 Stage-4 mid-row failures

- `ROLLBACK` the apply transaction.
- Insert a `BackupApplyDeadLetter` row with `FailedAtSyncOpSeq = <row.SyncOpSeq>`.
- **Do not** advance `BackupSyncWatermark`.
- Increment `SeenCount` and update `LastSeenAtEpoch` on subsequent retries of the same `EnvelopeId`.
- After `MainWorker.Backup.Apply.MaxRetriesPerEnvelope` retries (default **5**), raise `MAIN-840-01 BackupApplyExhausted` to Main on the next BE-5 Health poll. Operator decides next step; the system MUST NOT auto-skip.

### 6.3 No silent skips

CODE-RED rule: never swallow errors. A row that cannot be applied because the backup mirror lacks a target column = **dead-letter the envelope**, not skip the row. Schema drift is an operational issue, not an apply-time recoverable.

### 6.4 Journal mode (binding — RESOLVED 2026-05-06, was OQ-22-1)

Stage 4 MUST NOT issue per-envelope `PRAGMA journal_mode=WAL`. The App tier's session-wide WAL mode (pinned in `spec/05-split-db-architecture/`) is **authoritative**. A per-envelope override would defeat the App tier's shared cache and is forbidden. A blind-AI implementer MUST treat this as an unconditional MUST.

---

## 7. Idempotency (V7 short-circuit)

On envelope receipt, **before** decryption (cheap path):

```
1. Look up BackupApplyIdempotency by (EnvelopeId, PrimaryWorkerNodeId).
2. If Status='Applied'      → return OriginalResponseJson verbatim. No decrypt, no apply.
3. If Status='Received'     → another worker process is mid-apply. Return 409 + WORKER-300-01.
4. If Status='DeadLettered' → return the original failure code; do not retry from BE-1.
5. Else (no row)            → INSERT row with Status='Received', proceed to decrypt.
```

Step 5's INSERT uses the table's `UNIQUE` constraint as the lock — second concurrent BE-1 call lands on step 3 automatically. No advisory locks, no separate mutex.

---

## 8. Tunables introduced (mirrored verbatim into `15-tunable-constants.md` §2.14)

| Key | Default | Unit | Used by | Notes |
|---|---:|---|---|---|
| `MainWorker.Backup.Apply.MaxRetriesPerEnvelope` | **5** | count | §6.2 | Primary's BE-1 retries inherit this; exceeding = `MAIN-840-01`. |
| `MainWorker.Backup.Apply.TransactionTimeoutSeconds` | **30** | seconds | Stage 4 `BEGIN IMMEDIATE` | Wraps the per-envelope TX; exceeding = `WORKER-930-04 ApplyTransactionTimeout`. |
| `MainWorker.Backup.Apply.DeadLetterRetentionDays` | **30** | days | DLQ sweeper | Symmetric with snapshot retention default; finalised in Phase 11 with OQ-A4. |
| `MainWorker.Backup.Apply.IdempotencyRowRetentionDays` | **14** | days | DLQ sweeper | `Applied` rows reaped after this; replay protection window. |

---

## 9. Errors introduced (mirrored verbatim into `13-error-codes.md`)

Worker tier (apply-pipeline failures, range `WORKER-930-*` — first three slots of the freshly-allocated overflow range, cited below):

| Prefixed | Flat | HTTP | Name | Meaning |
|---|---:|---:|---|---|
| `WORKER-930-01` | _(see §9.1)_ | 422 | `EnvelopeRowCountMismatch` | V4: `Envelope.RowCount` ≠ `COUNT(EnvelopeRow)`. |
| `WORKER-930-02` | _(see §9.1)_ | 422 | `UntrackedSourceTable` | V5: `SourceTable` not in tracked-table allowlist. |
| `WORKER-930-03` | _(see §9.1)_ | 422 | `UnknownSyncOpCode` | Stage-4 guard `AssertKnownSyncOp` failure (not `Insert`/`Update`/`Delete`). |
| `WORKER-930-04` | _(see §9.1)_ | 504 | `ApplyTransactionTimeout` | Stage-4 TX exceeded `TransactionTimeoutSeconds`. |

Main tier:

| Prefixed | Flat | HTTP | Name | Meaning |
|---|---:|---:|---|---|
| `MAIN-840-01` | _(see §9.1)_ | n/a | `BackupApplyExhausted` | §6.2 — `MaxRetriesPerEnvelope` exceeded for one `EnvelopeId`; surfaced via BE-5 Health. |

### 9.1 Range-allocation note

Worker future-expansion range `WORKER-21095-21099` was **fully consumed** by Phase 8 (`WORKER-920-01..05`). Phase 10 therefore opens a new overflow window per `13-error-codes.md` §1 *Slot-overflow rule*: **`WORKER-21200-21299`** (next free hundred). Concrete slot assignments:

| Prefixed | Flat |
|---|---:|
| `WORKER-930-01` | `21200` |
| `WORKER-930-02` | `21201` |
| `WORKER-930-03` | `21202` |
| `WORKER-930-04` | `21203` |
| `MAIN-840-01` | `21191` (still inside the `MAIN-21191-21199` Phase-11-reserved window; see §9.2) |

### 9.2 Coexistence with Phase 11

Phase 11 reserves `MAIN-21191-21199` for snapshot/restore. `MAIN-840-01` consumes the **first** slot (`21191`). Phase 11 will open at `21192`. The reserved-range table in `13-…` §4 is updated accordingly.

---

## 10. Linter hooks queued for Phase 12

| ID | Rule |
|---|---|
| `BACKUP-APPLY-001` | Every `BackupApplyIdempotency.Status='Applied'` row MUST have `AppliedAtEpoch IS NOT NULL` and `OriginalResponseJson IS NOT NULL`. |
| `BACKUP-APPLY-002` | Every `BackupApplyDeadLetter` row MUST have `EnvelopeBlobPath` resolving to an existing file under `var/dlq/`. |
| `BACKUP-APPLY-003` | `EnvelopeRow.SourceTable` values referenced anywhere in spec/19 MUST appear in the `AppBackupTrackedTable` seed (Phase 12). |
| `BACKUP-APPLY-004` | Apply pipeline functions MUST be ≤15 lines, zero nested `if`, ≤2 operands per boolean (CODE RED). Violations fail CI. |

---

## 11. Cross-references

- `19-incremental-backup-sync.md` §6 — envelope SQLite shape (input to Stage 2).
- `20-backup-encryption-and-keys.md` §9 — decrypt path (Stage 1).
- `21-backup-endpoints.md` §4 — BE-1 wire (caller).
- `13-error-codes.md` §2.10 / §3.10 / §4 — code allocations + range table.
- `15-tunable-constants.md` §2.14 — mirrored tunables.
- `17-consolidated-guidelines/` (CODE RED) — guard-extraction style enforced by `BACKUP-APPLY-004`.

---

## 12. Open Questions (logged, non-blocking)

- ✅ **OQ-22-1 RESOLVED 2026-05-06.** Per-envelope `PRAGMA journal_mode=WAL` is forbidden — see §6.4 (binding MUST). Rely on the App tier's session-wide WAL.
- **OQ-22-2** Should `BackupApplyDeadLetter` rows older than `DeadLetterRetentionDays` be reaped automatically, or require operator action? Inferred: automatic sweep, but emit a `MAIN-840-01` summary metric on every sweep so dashboards never go silent.
- **OQ-22-3** Should the tracked-table allowlist (`AppBackupTrackedTable`) be hard-coded in the seed or derived dynamically from rows that carry Shape A `SyncOpCode` columns? Inferred: hard-coded seed wins (dumb-AI friendly + audit-trail friendly); Phase 12 owns the seed.

---

*Backup apply logic v1.0.0 — 2026-05-06 (Phase 10). CODE-RED-compliant pipeline; idempotent dispatch; explicit DLQ.*
