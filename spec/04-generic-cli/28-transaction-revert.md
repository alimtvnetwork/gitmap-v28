# Transaction journal & revert (v4.18.0+)

## 1. Goal

Every gitmap-v28 operation that **mutates the on-disk filesystem** is recorded
as a row in a SQLite `Transaction` table the moment it begins, and a
matching reverse-action is captured before the change is applied. A user
can then roll back any of the last 50 transactions with a single command.

The transaction journal lives in the same SQLite DB as the rest of
gitmap-v28 state, so it is automatically anchored to the binary path
(see `mem://tech/database-location`) and never crosses install
boundaries.

## 2. Scope (in / out)

**Recorded as transactions** — these mutate the filesystem and are
the only operations a user ever wants to undo by accident:

| Op family            | Sub-commands                                                                 |
|----------------------|------------------------------------------------------------------------------|
| Clone family         | `clone`, `clone-next` (`cn`), `clone-from`, `clone-now` (`cnow`), `clone-pick` (`cp`), `clone-fix-repo` (`cfr` / `cfrp`) |
| Move / merge         | `mv`, `merge-both`, `merge-left`, `merge-right`                              |
| Repo rewrite         | `fix-repo` (`fr`)                                                            |
| History rewrite      | `history-purge` (`hp`), `history-pin` (`hpin`)                               |

**Not recorded** (read-only or DB-only):

- `scan`, `rescan`, `scan-clone`, `clone-multi --dry-run`
- `vscodepm` sync / `vscode-workspace`, `vpath`
- `release`, `release-alias` (`as` / `ra` / `rap`) — DB-only
- `audit-legacy`, `regoldens`, `doctor`, `version`, `help`, `docs`

## 3. Storage layout

### 3.1 SQLite tables

```
Transaction(
  TransactionId   INTEGER PRIMARY KEY AUTOINCREMENT,
  Kind            TEXT NOT NULL,            -- "clone" | "mv" | "merge" | "fix-repo" | "history-purge" | "history-pin"
  Status          TEXT NOT NULL,            -- "pending" | "committed" | "aborted" | "reverted"
  Argv            TEXT NOT NULL,            -- raw os.Args joined by \x1f
  Cwd             TEXT NOT NULL,
  CreatedAt       INTEGER NOT NULL,         -- unix seconds
  CommittedAt     INTEGER,                  -- nullable
  RevertedAt      INTEGER,                  -- nullable
  ReverseSummary  TEXT NOT NULL,            -- human-readable "delete <path>" / "rename A <- B"
  RepoSlug        TEXT NOT NULL,            -- best-effort, "" when unknown
  GitSha          TEXT NOT NULL             -- HEAD sha at op-start, "" when not a git repo
);

TransactionFile(
  TransactionFileId INTEGER PRIMARY KEY AUTOINCREMENT,
  TransactionId     INTEGER NOT NULL REFERENCES "Transaction"(TransactionId) ON DELETE CASCADE,
  RelPath           TEXT NOT NULL,          -- relative to the affected repo root
  AbsPath           TEXT NOT NULL,          -- snapshot at op-start
  BackupPath        TEXT NOT NULL,          -- absolute path under .gitmap/txn/<txn-id>/data/<repo>/<git-sha>/files/
  ByteSize          INTEGER NOT NULL,
  Sha256            TEXT NOT NULL,
  Action            TEXT NOT NULL           -- "delete" | "edit" | "rename"
);
```

### 3.2 Filesystem backup layout

```
<gitmap-binary-dir>/.gitmap/
└── txn/
    └── <txn-id>/
        └── data/
            └── <repo-slug-or-name>/
                └── <git-sha-or-_unknown>/
                    └── files/
                        ├── path/inside/repo/foo.go
                        └── ...
```

- `<git-sha>` is `git rev-parse HEAD` taken at transaction-begin time
  (so a revert restores the bytes the user actually had on disk, not
  whatever the remote HEAD points to today). Falls back to `_unknown`
  for non-git directories.
- Backups are **only created** for `delete` and `edit` actions.
  `fix-repo --replace` is **explicitly out of scope** (see §6.1).
- The relative directory shape inside `files/` mirrors the repo, so
  a recursive `cp -r` is the entire revert action for an edit/delete.

## 4. Lifecycle

```
        Begin                 record-files          Commit
   user ───────► Transaction ─────────────► do FS ────────►  status=committed
                  status=pending                              CommittedAt=now
                                            FAIL
                                            ───────► Abort   status=aborted

        Revert
   user ───────► restore files in reverse order ─────► status=reverted
                                                       RevertedAt=now
```

Pre-mutation order is mandatory:

1. `txn.Begin(kind)` → INSERT row, returns `TransactionID`.
2. For every file the op will delete OR edit-in-place:
   `txn.SnapshotFile(absPath, repoRoot)` → copies bytes into the
   backup dir, INSERTs a `TransactionFile` row with sha256 + size.
3. The op performs its FS mutation.
4. `txn.Commit()` flips `Status='committed'` and stamps `CommittedAt`.
   Any failure between (1) and (4) → `txn.Abort()` (status='aborted').

## 5. CLI surface — extend `gitmap-v28 revert`

The legacy `gitmap-v28 revert <version>` flow (rolls back the gitmap-v28
binary itself) is preserved as the default when no transaction flag
is present. New flags switch into transaction mode:

| Flag                     | Behavior                                                                 |
|--------------------------|--------------------------------------------------------------------------|
| `--list-txn`             | Print the last 50 transactions (id, kind, status, argv tail, repo).      |
| `--show-txn <id>`        | Print one transaction including every captured file.                     |
| `--txn <id>`             | Revert the named transaction.                                            |
| `--last-txn`             | Revert the most recent committed transaction.                            |
| `--prune-txn`            | Drop everything beyond the 50-row cap right now (also runs auto on Begin). |
| `--force`                | Skip the confirm prompt. Required for clone-revert when the dest dir has been edited since commit. |

Help text lives at `gitmap-v28/helptext/revert.md` (extended, not replaced).

## 6. Per-op reverse semantics

| Op                | Snapshot taken      | Revert action                                          |
|-------------------|---------------------|--------------------------------------------------------|
| `clone*`          | none (dest is new)  | `os.RemoveAll(dest)` + drop matching `Repo` row        |
| `mv`              | full src tree       | restore src from `files/`, delete dst                  |
| `merge-*`         | full target tree    | restore target from `files/`                           |
| `fix-repo`        | per-edited file     | overwrite each edited file with its `files/` copy      |
| `history-purge`   | full repo           | restore repo from `files/`                             |
| `history-pin`     | full repo           | restore repo from `files/`                             |

### 6.1 fix-repo `--replace` exemption

Per current scope: `fix-repo` runs that pass `--replace` (the in-place
rewriter) do **not** snapshot the per-file bytes. They still record a
`Transaction` row (kind=`fix-repo`, ReverseSummary=`replace mode — file
contents NOT recoverable, re-run fix-repo with the prior version`) so
the audit trail stays complete.

## 7. Retention

- Hard cap: **50 transactions**.
- `txn.Begin` calls `txn.PruneOldest(50)` first, which `DELETE`s rows
  with the smallest `TransactionId` values until count ≤ cap and
  recursively `os.RemoveAll`s their backup dirs.
- `gitmap-v28 revert --prune-txn` lets users force a prune cycle.

## 8. Concurrency

- The existing `gitmap.lock` advisory file lock guards every Begin
  / Commit / Revert call so two gitmap-v28 processes never interleave
  transaction rows.
- SQLite `SetMaxOpenConns(1)` (project-wide rule) makes intra-process
  serialization free.

## 9. Out-of-scope (deferred)

- `fix-repo --replace` per-file backups (see §6.1).
- Non-FS mutations (alias create/rename, release publish, scan upserts).
- Cross-machine restore / portable transaction archives.
- A full UI under `gitmap-v28 txn ...` — left as a follow-up; the `revert`
  flag surface is sufficient for v1.

## 10. Implementation references

- Schema: `gitmap-v28/constants/constants_transaction.go`
- Tables: `gitmap-v28/store/transaction.go`, `gitmap-v28/store/transaction_file.go`
- Engine: `gitmap-v28/txn/begin.go`, `gitmap-v28/txn/snapshot.go`, `gitmap-v28/txn/revert.go`, `gitmap-v28/txn/prune.go`
- CLI: `gitmap-v28/cmd/revert.go` (extended), `gitmap-v28/cmd/revert_txn.go` (new)
- Help: `gitmap-v28/helptext/revert.md` (extended)
- Memory: `mem://features/transaction-revert`