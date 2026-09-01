# 117 — Update Awareness: `list --update`, `update apply/all`, `hd`

## Purpose

Close the loop after `gitmap scan`: surface which scanned repos have upgrades
available, act on them individually or in bulk, and expose the status through
a single-screen dashboard and the existing `stats` output.

## Command Signatures

    gitmap list --update [flags]        # alias: lu, alt: gitmap update list
    gitmap update apply <repo> [flags]  # alias: ua
    gitmap update all [flags]           # alias: uall
    gitmap hd [flags]                   # alias for: gitmap help --dashboard

## Data Model

Reuses `VersionProbe`, `ScanFolder`, and `Releases` tables. Adds one seed
value to `TaskType`:

```sql
INSERT OR IGNORE INTO TaskType (Name) VALUES ('Upgrade');
```

`TaskTypeUpgrade` is enqueued when `update apply` or `update all` fails on a
row. The task's `SourceCommand` is `update apply` and `CommandArgs` is the
resolved repo slug or path, so `do-pending` can replay it verbatim.

`constants/constants_pending_task.go` gains:

```go
const TaskTypeUpgrade = "Upgrade"
```

`SQLSeedTaskTypes` is extended to include `('Upgrade')`.

## `list --update` Behavior

1. Load every row from `VersionProbe` joined with `ScanRecord`.
2. Compute `behind = index(latest) - index(current)` using the semver-sorted
   tag list. Rows with unknown `current` render `behind = -` and count toward
   the `unknown` bucket.
3. Filter by `--source`, sort by `behind DESC, repo ASC`, apply `--limit`.
4. When `--stale-days N` is set, re-probe any row whose last probe is older
   than N days before rendering.
5. Exit `0` when rendered, `2` when the scan cache is missing.

Output columns: `repo`, `current`, `latest`, `behind`, `source`, `path`. JSON
schema is `spec/08-json-schemas/list-update.schema.json`.

## `update apply` Behavior

1. Resolve `<repo>` against `ScanRecord` (slug, nickname, or path).
2. Confirm the plan unless `-y`; skip execution entirely when `--dry-run`.
3. If `--stash`, run `git stash push -u` first; restore with `git stash pop`
   on success only.
4. Strategy dispatch:
   - `fetch-checkout` (default): `git fetch --tags` then `git checkout <tag>`
     where `<tag>` is `--tag` or the row's `latest`.
   - `source-release`: only valid when `source == release`; invokes the
     repo's `gitmap release-self --apply` equivalent.
5. On success, insert a `CompletedTask` row.
6. On failure, insert a `PendingTask` with `TaskTypeId = Upgrade`,
   `CommandArgs = <resolved-repo>`, `FailureReason = <error>`. Exit `1`.

Exit codes: `0` upgraded, `1` failed and enqueued, `2` already up-to-date or
unknown repo.

## `update all` Behavior

1. Build the queue from `list --update` (same filters via `--only`).
2. Confirm once unless `-y`; abort entirely on `n`.
3. Spawn a worker pool of size `--parallel` (clamp `[1, 16]`).
4. Each worker runs `update apply` per row, honoring `--strategy`, `--stash`,
   `--dry-run`.
5. `--stop-on-error` cancels the queue after the first failure; the default
   drains the queue and reports partial success.
6. Summary line prints `ok / failed / skipped` and lists any newly created
   PendingTask ids.

Exit codes: `0` all ok, `1` partial, `2` nothing upgradable.

## `hd` Behavior

Read-only aggregate; no writes. Assembles:

| Field | Source |
|-------|--------|
| `version.installed` | `constants.Version` |
| `version.latest` | Latest gitmap release tag (cached; refreshed by `--refresh`) |
| `repos.scanned` | `SELECT COUNT(*) FROM ScanRecord` |
| `repos.upgradable` | Row count from `list --update` |
| `repos.unknown` | Probe rows with `IsAvailable = 0` and `Error != ''` |
| `pending.queued` | `SELECT COUNT(*) FROM PendingTask` |
| `pending.oldest` | `ORDER BY CreatedAt LIMIT 1` |
| `lastScanAt` | `SELECT MAX(FinishedAt) FROM ScanRun` |
| `recent` | Last N rows of `CompletedTask ORDER BY CompletedAt DESC` |

## `stats` Extension

`gitmap stats` gains an Upgrades block, printed after the existing command
table:

```
UPGRADES
  repos scanned:        142
  up-to-date:           118
  upgradable:            22   →  gitmap list --update
  unknown / no tags:      2
  last full scan:  2026-07-22 10:14
```

JSON output gains a top-level object:

```json
"upgrades": {
  "scanned":    142,
  "upToDate":   118,
  "upgradable":  22,
  "unknown":      2,
  "lastScanAt": "2026-07-22T10:14:00Z"
}
```

## `scan` Post-Run Hint

After every successful `gitmap scan`, if the newly-updated probe cache
reports `upgradable > 0`, the CLI prints one closing line:

```
  ℹ  22 repo(s) have upgrades available. Run:  gitmap list --update
```

This hint is suppressed by `--quiet`, or when `stderr` is not a terminal.

## Implementation Files (Go side)

| File | Responsibility |
|------|----------------|
| `cmd/listupdate.go` | `list --update` handler; also entry point for `update list` |
| `cmd/updateapply.go` | Single-repo upgrade path |
| `cmd/updateall.go` | Worker-pool wrapper around `updateapply` |
| `cmd/helpdashboard.go` | `hd` renderer, both text and JSON |
| `cmd/scan.go` | Emit the post-scan upgrade hint |
| `cmd/stats.go` | Compose the Upgrades block |
| `constants/constants_cli.go` | `CmdListUpdate`, `CmdUpdateApply`, `CmdUpdateAll`, `CmdHD` + aliases |
| `constants/constants_pending_task.go` | `TaskTypeUpgrade`; extend `SQLSeedTaskTypes` |
| `constants/constants_update_awareness.go` | New: messages, hint format, table headers |
| `store/upgradesview.go` | New: joined view feeding `list --update` and `stats` |

## Docs-Site Wiring

- `src/data/commands.ts`: add entries for `list --update`, `update apply`,
  `update all`, `hd`. Category `tools` for the update commands, `scanning`
  for `list --update` and `hd`.
- `gitmap/helptext/`: new files `list-update.md`, `update-apply.md`,
  `update-all.md`, `hd.md`. Cross-links added to `scan.md`, `stats.md`,
  `pending.md`, `do-pending.md`.
- `src/pages/FlagReference.tsx`: include the new flags.

## Acceptance Criteria

1. `gitmap list --update` renders the columned table and JSON schema listed
   above; exit `0` on success, `2` when scan cache is missing.
2. `gitmap update apply <repo>` upgrades, honors `-y`, `--dry-run`,
   `--stash`, and on failure enqueues a `PendingTask(Upgrade)` with
   `SourceCommand = "update apply"`.
3. `gitmap update all` walks the `list --update` queue with `--parallel`
   workers, supports `--only`, `--stop-on-error`, returns exit `0/1/2`.
4. `gitmap hd` prints the dashboard in text and JSON.
5. `gitmap stats` output contains the Upgrades block in both formats.
6. `gitmap scan` prints the trailing upgrade hint when `upgradable > 0` and
   stderr is a TTY.
7. `gitmap do-pending` replays `Upgrade` tasks via `gitmap update apply`.
8. Docs site lists all four new commands with usage, flags, and examples.
9. `gitmap help --json` payload includes the four new commands under the
   correct groups.

## Code Style

All functions ≤ 15 lines. Positive logic. Blank line before every return. No
magic strings. No switch statements.

## See Also

- [95-pending-task-workflow.md](95-pending-task-workflow.md)
- [19-list-versions.md](19-list-versions.md)
- [90-scan-folder-and-version-probe.md](90-scan-folder-and-version-probe.md)
