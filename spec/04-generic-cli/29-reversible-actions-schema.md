# Reversible Actions Schema (v4.18.0+)

Companion to [`28-transaction-revert.md`](28-transaction-revert.md). This
spec defines the **typed forward + reverse payload contract** every
state-mutating gitmap-v27 operation must produce so that revert is
deterministic, auditable, and survives schema growth.

The transaction journal in spec 28 already stores file-level snapshots.
This spec adds a layer **above** files: a closed set of typed `Action`
records that carry the structured arguments needed to undo the
operation. File snapshots remain the byte-restore mechanism; actions
describe **what to do with those bytes**.

## 1. Goals

- **Deterministic replay**: given an action record, the reverse step
  must produce the exact pre-state regardless of host, time, or
  unrelated DB rows.
- **Closed action set**: a fixed enum of `ActionKind` values. New
  operations MUST map onto an existing kind or open a spec PR adding
  a new one — never invent ad-hoc payloads.
- **Self-describing**: every action carries enough JSON to be replayed
  without consulting the original command-line.
- **No partial reverts**: an action either reverts cleanly or fails
  with a typed error; never leaves a half-state.
- **Forward parity**: the same record could in principle drive the
  forward op (for replay/dry-run tooling). Today only the reverse half
  is consumed, but the schema is symmetric so the forward half is free.

## 2. Closed action set

| ActionKind        | Forward intent                              | Reverse intent                                                   |
|-------------------|---------------------------------------------|------------------------------------------------------------------|
| `create_path`     | Create a new file or directory at `Dst`     | `os.RemoveAll(Dst)` if its sha matches `PostSha`                 |
| `delete_path`     | Remove file/dir at `Src` (snapshot first)   | Restore bytes from `BackupRef` to `Src`                          |
| `edit_file`       | Overwrite `Src` (snapshot first)            | Restore bytes from `BackupRef` to `Src`                          |
| `rename_path`     | Rename `Src` → `Dst`                        | Rename `Dst` → `Src`                                             |
| `clone_repo`      | `git clone Url Dst` (+ DB row)              | `os.RemoveAll(Dst)` + delete `Repo` row by `RepoSlug`            |
| `db_upsert_repo`  | Insert/update a `Repo` row                  | Restore prior `RepoRowJson` (or delete row if it didn't exist)   |
| `db_delete_repo`  | Delete a `Repo` row (snapshot JSON first)   | Re-INSERT from `RepoRowJson`                                     |
| `tag_set`         | Set VS Code PM tags on `RepoSlug`           | Restore `PriorTagsJson`                                          |
| `workspace_write` | Write a `.code-workspace` file at `Dst`     | Restore `BackupRef` if existed; else delete                      |

No other kinds exist. A new operation either reduces to N of these or
the spec is amended.

## 3. Payload schema

Stored as one row per action in a new SQLite table sibling to
`TransactionFile`. JSON columns are validated against the schemas in
§3.2 on insert.

### 3.1 Table

```sql
CREATE TABLE IF NOT EXISTS TransactionAction (
  TransactionActionId INTEGER PRIMARY KEY AUTOINCREMENT,
  TransactionId       INTEGER NOT NULL REFERENCES "Transaction"(TransactionId) ON DELETE CASCADE,
  Seq                 INTEGER NOT NULL,            -- 1-based, dense, per-transaction
  Kind                TEXT    NOT NULL,            -- one of section 2
  ForwardJson         TEXT    NOT NULL,            -- section 3.2 forward payload
  ReverseJson         TEXT    NOT NULL,            -- section 3.2 reverse payload
  BackupRef           TEXT    NOT NULL DEFAULT '', -- TransactionFile.AbsPath when Kind needs bytes; '' otherwise
  AppliedAt           INTEGER,                     -- nullable; unix seconds when forward succeeded
  RevertedAt          INTEGER,                     -- nullable; unix seconds when reverse succeeded
  UNIQUE(TransactionId, Seq)
);
```

`Seq` is the canonical replay order. Forward = ascending; reverse =
descending. `Seq` MUST be dense and 1-based so a missing row is a
detectable corruption.

### 3.2 Per-kind payloads

All paths are **absolute** at write time. Slugs match the existing
`Repo.Slug` column. Sha values are SHA-256 hex (lowercase, 64 chars).

```text
create_path
  { "Dst": "/abs/path", "IsDir": false,
    "PostSha": "<sha256-of-created-bytes-or-'' for dir>" }

delete_path
  { "Src": "/abs/path", "IsDir": false,
    "PreSha": "<sha256-before>",
    "BackupRef": "<txn-file-abs-path>" }

edit_file
  { "Src": "/abs/path",
    "PreSha":  "<sha256-before>",
    "PostSha": "<sha256-after>",
    "BackupRef": "<txn-file-abs-path>" }

rename_path
  { "Src": "/abs/old", "Dst": "/abs/new", "IsDir": false }

clone_repo
  { "Url": "https://.../foo.git",
    "Dst": "/abs/foo",
    "RepoSlug": "foo",
    "GitShaAtClone": "<head-after-clone>" }

db_upsert_repo
  { "RepoSlug": "foo",
    "Existed": true,
    "RepoRowJson": "<full prior row as JSON>" }

db_delete_repo
  { "RepoSlug": "foo",
    "RepoRowJson": "<full prior row as JSON>" }

tag_set
  { "RepoSlug": "foo",
    "NewTagsJson":   "[\"go\",\"git\"]",
    "PriorTagsJson": "[\"go\"]" }

workspace_write
  { "Dst": "/abs/gitmap.code-workspace",
    "Existed": true,
    "BackupRef": "<txn-file-abs-path-or-''>" }
```

Forward and reverse halves share the same JSON shape per kind — the
engine just reads different fields when applying vs reverting.
`ForwardJson` and `ReverseJson` are stored separately (rather than
derived) so future spec drift cannot retroactively change how an old
transaction reverts.

## 4. Engine contract

### 4.1 Forward path

```
j := txn.Begin(kind)              // spec 28 section 4
j.Plan(action)                    // INSERT TransactionAction row (Seq, ForwardJson, ReverseJson, BackupRef)
j.Apply(action)                   // performs FS/DB op, stamps AppliedAt
// repeat for every step
j.Commit()
```

`Plan` is a pure DB write. `Apply` is the side-effecting half. The
split lets `--dry-run` surfaces print every planned action without
touching disk.

### 4.2 Reverse path

```
for action in db.ListActions(txnId) ORDER BY Seq DESC:
    if action.RevertedAt is not NULL: continue   # idempotent resume
    txn.ApplyReverse(action)                     # dispatch by Kind
    db.MarkActionReverted(action.Id)
db.MarkTransactionReverted(txnId)
```

`ApplyReverse` MUST be idempotent at the action level: re-running a
revert after a partial failure resumes from the first
`RevertedAt IS NULL` row.

### 4.3 Determinism guarantees

- Reverse order is `Seq DESC`. No other ordering is legal.
- Each kind's reverse handler is **pure** w.r.t. the payload + the
  bytes at `BackupRef`. It MUST NOT consult the live DB or live FS for
  anything other than the literal `Src/Dst` it owns.
- Sha checks gate every byte-restore (`delete_path`, `edit_file`,
  `workspace_write`). Mismatch → typed error `ErrBackupShaDrift`,
  bypassed only by `revert --force`.
- Every reverse handler treats "already in target state" as success
  (e.g. `delete_path` reverse where `Src` already exists with matching
  sha is a no-op, not a failure).

## 5. Operation → action mapping

| Op                   | Action sequence (forward order)                                             |
|----------------------|------------------------------------------------------------------------------|
| `clone` / `cn`       | `clone_repo` → `db_upsert_repo`                                              |
| `clone-fix-repo`     | `clone_repo` → `db_upsert_repo` → N×`edit_file`                              |
| `mv`                 | `rename_path` (one row, no byte snapshot)                                    |
| `merge-both`         | N×`edit_file` (one per touched file)                                         |
| `merge-left/right`   | N×`edit_file`                                                                |
| `fix-repo`           | N×`edit_file`                                                                |
| `fix-repo --replace` | 1×`edit_file` per file with `BackupRef=''` and `ReverseJson.NoBytes=true`; reverse handler errors with the spec 28 §6.1 message |
| `history-purge`      | 1×`edit_file` per restored repo file (folder snapshot still under `TransactionFile` for audit) |
| `history-pin`        | same as `history-purge`                                                      |
| `vscode-workspace`   | 1×`workspace_write`                                                          |
| `vpm sync`           | 1×`tag_set` per repo whose tags changed                                      |

`scan` / `rescan` remain out of scope (spec 28 §2). If scan ever needs
revert, it maps to N×`db_upsert_repo` + N×`db_delete_repo` and nothing
else — no new kinds required.

## 6. Error taxonomy

| Error                       | When                                                       | Exit |
|-----------------------------|------------------------------------------------------------|------|
| `ErrActionUnknownKind`      | DB row `Kind` not in section 2                             | 1    |
| `ErrActionPayloadMalformed` | JSON fails schema validation                               | 1    |
| `ErrActionSeqGap`           | `Seq` values not dense / 1-based for the txn               | 1    |
| `ErrBackupMissing`          | `BackupRef` path doesn't exist on disk                     | 1    |
| `ErrBackupShaDrift`         | Backup file sha != recorded `PreSha`/`PostSha`             | 1    |
| `ErrReverseUnsupported`     | `fix-repo --replace` action with `NoBytes:true`            | 1    |
| `ErrLiveStateConflict`      | Reverse target moved/edited since commit (no `--force`)    | 1    |

All errors are typed (Go `errors.Is`-friendly per project Code Red
policy) and surface via `fmt.Fprintf(os.Stderr, ...)` with the standard
`<command>: <message>\n` shape.

## 7. Retention & GC

Action rows inherit the parent transaction's retention (50-row cap per
spec 28 §7). When a `Transaction` row is pruned, `ON DELETE CASCADE`
drops every child `TransactionAction` and `TransactionFile`. The
filesystem `BackupRef` files are removed by the same prune sweep that
`os.RemoveAll`s the per-txn backup directory.

## 8. Versioning

- The schema lives at SQLite `SchemaVersion = 23` (one bump above the
  v22 introduced by spec 28).
- Adding a new `ActionKind` requires:
  1. A spec amendment to §2 + §3.2 (this file).
  2. A new constant in `constants/constants_transaction_actions.go`.
  3. Forward + reverse handlers in `gitmap-v27/txn/actions/`.
  4. A migration that backfills nothing — old rows keep their old
     kinds; only new rows can carry the new one.
- Renaming or removing a kind is **forbidden**. Old transactions must
  always be revertable by binaries built against any future schema.

## 9. Implementation hooks (forward references)

- New constants: `constants/constants_transaction_actions.go`
  (`ActionKindCreatePath`, …, `ErrActionUnknownKind`, JSON schema
  templates).
- New tables: `store/transaction_action.go` with
  `AppendAction / ListActionsBySeq / MarkActionReverted`.
- Engine: `gitmap-v27/txn/actions/{plan,apply,reverse}.go`, one dispatch
  table keyed by `ActionKind`.
- CLI surface unchanged — `gitmap-v27 revert` still drives everything via
  spec 28's flags. `--show-txn` is extended to print the action list.

## 10. Out of scope (still)

- Cross-machine portable transaction archives.
- Forward replay of an action list (the schema supports it; no caller
  exists yet).
- Action-level partial revert (`revert --action <seq>`). The unit of
  user-visible undo remains the whole transaction.
- Conflict resolution beyond `--force`: no three-way merge.

## 11. Memory anchor

- `mem://features/reversible-actions-schema` (this spec, condensed).
- `mem://features/transaction-revert` (parent journal — keep in sync).
