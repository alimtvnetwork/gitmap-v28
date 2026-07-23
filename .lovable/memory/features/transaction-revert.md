# Memory: features/transaction-revert
Updated: now

## Operating rules (apply on every FS-mutating op)

Every operation in spec/04-generic-cli/28-transaction-revert.md §2
(clone family, mv, merge-*, fix-repo, history-purge, history-pin)
MUST call `txn.Begin(kind)` BEFORE touching the filesystem and
`txn.Commit()` on success. Any error path between Begin and Commit
MUST call `txn.Abort(err)`. Files about to be deleted or edited in
place MUST be snapshotted via `txn.SnapshotFile(absPath, repoRoot)`
BEFORE the mutation. Snapshots land at
`<gitmap-binary-dir>/.gitmap/txn/<txn-id>/data/<repo>/<git-sha>/files/<rel>`.

## Exemption (intentional, do not "fix")

`fix-repo --replace` records the row but skips per-file snapshots.
ReverseSummary is set to a warning string explaining the user must
re-run fix-repo with the prior version. Spec §6.1.

## Hard cap

50 transactions, enforced by `txn.PruneOldest(50)` called from
`txn.Begin`. Surplus rows AND their backup directories are removed
atomically (DB row delete cascades to TransactionFile, then
`os.RemoveAll(<txn-id>/)`). NEVER raise this cap silently.

## CLI surface (extends `gitmap revert`, does NOT replace)

Legacy `gitmap revert <version>` (binary-version rollback) stays the
default when no transaction flag is present. New flags switch into
transaction mode: `--list-txn`, `--show-txn <id>`, `--txn <id>`,
`--last-txn`, `--prune-txn`, `--force`.

## Schema

Tables: `Transaction`, `TransactionFile` (PascalCase, INTEGER PK
AUTOINCREMENT). `SchemaVersionCurrent` is bumped in the same commit
that adds the tables (21 → 22).

## NEVER

1. Mutate the FS without an active `txn` handle when the op is in scope.
2. Hand-edit `Transaction` / `TransactionFile` rows from ad-hoc SQL.
3. Snapshot for `fix-repo --replace` runs.
4. Raise the 50-row cap without a spec edit.
