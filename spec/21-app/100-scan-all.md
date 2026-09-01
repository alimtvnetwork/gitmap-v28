# 100 — `gitmap scan all` (Bulk Re-Scan)

> Status: **planned** — targets v3.33.0
> Depends on: spec 90 (`ScanFolder` table, shipped v3.7.0)

## Goal

Provide a single command that re-scans **every directory the user has ever
scanned**, without forcing them to remember or re-type each root. Useful for:

- Refreshing the entire `gitmap` worldview after a long break.
- Pulling in new repos cloned manually under existing roots.
- Producing an up-to-date snapshot before a `release` / `export`.

## CLI Surface

| Form               | Notes                                  |
|--------------------|----------------------------------------|
| `gitmap scan all`  | Canonical form.                        |
| `gitmap scan a`    | Short alias (mirrors `s` for `scan`).  |

`scan all` is implemented as a **subcommand dispatch** inside `cmd/scan.go`,
not as a `--all` flag. This keeps existing `scan` flags (`--exclude`,
`--mode`, `--output`) free of conflicts and matches the user's preferred
verbal phrasing.

All other `scan` flags (`--mode`, `--output`, `--output-path`,
`--github-desktop`, `--quiet`) are accepted and forwarded to each
per-root scan **unchanged**. `--output-path` is intentionally **rejected**
in `all` mode (each root writes to its own `.gitmap/output/`).

## Source of Roots

Read from the `ScanFolder` table (spec 90). Query:

```sql
SELECT ScanFolderId, AbsolutePath, Label, LastScannedAt
FROM   ScanFolder
ORDER  BY LastScannedAt DESC;
```

No new state. No deriving roots from `Repo.AbsolutePath`. If `ScanFolder` is
empty (fresh install), print:

```
No previously-scanned roots found. Run 'gitmap scan <dir>' first.
```

…and exit `0` (informational, not an error).

## Execution Model — Parallel Worker Pool

- Default concurrency: **N=4** workers (constant `ScanAllDefaultWorkers`).
- Override: `--workers <n>` (1–16, clamped).
- Each worker pulls one root off a buffered channel and runs the existing
  `runScanForRoot(root, opts)` helper end-to-end (artifact write + DB upsert
  + `EnsureScanFolder` LastScannedAt bump).
- Per-root logs are **buffered** and flushed atomically when that root
  finishes, so parallel output never interleaves mid-line. Each block is
  prefixed with the root path and a short status icon (`✓` / `⚠` / `✗`).
- A final summary line is printed after all workers drain:

```
Scanned 7 root(s) · 3 missing · 142 repos discovered · 4.8s
```

### Why parallel

The user explicitly chose parallel. SQLite uses `SetMaxOpenConns(1)` (see
core memory), so DB writes serialize naturally. Filesystem walking and
`git remote -v` invocations are the slow parts and benefit most from
concurrency.

### Locking

The existing advisory `gitmap.lock` (see `mem://tech/process-synchronization`)
is acquired **once** at the top of `scan all`, not per-root. Workers share
it. This prevents a second `gitmap` process from racing the same DB.

## Missing-Root Handling

For each root, before dispatching to a worker:

```go
if _, err := os.Stat(root); errors.Is(err, fs.ErrNotExist) {
    missing = append(missing, root)
    continue
}
```

Missing roots are **never modified during scan** — the `ScanFolder` row
stays intact in case the drive comes back. After all workers drain, if
`len(missing) > 0`:

```
⚠ 3 root(s) skipped — directory not found:
    D:\old-projects
    \\nas\archive
    /Volumes/External/work

Prune them from the database? [y/N]
```

- **Interactive TTY**: prompt y/N. On `y`, run
  `SQLDetachReposFromScanFolder` then `SQLDeleteScanFolderByID` for each
  missing row, inside a single transaction.
- **Non-TTY** (CI, piped): print the list + the hint
  `Run 'gitmap sf rm <path>' to remove them` and exit `0`.

A `--prune-missing` flag bypasses the prompt and prunes unconditionally
(useful for cron jobs).

## Exit Codes

| Code | Meaning                                                       |
|------|---------------------------------------------------------------|
| 0    | All present roots scanned successfully (missing roots warned).|
| 1    | One or more roots failed mid-scan (FS / git / DB error).      |
| 2    | Lock contention — another `gitmap` is already running.        |
| 3    | `ScanFolder` empty — informational, not a hard error.         |

Exit `0` even when only some roots succeed is **deliberate**: bulk mode
should not fail loudly when one drive is unplugged. Per-root failures are
visible in the summary count.

## Constants (per `mem://style/code-constraints`)

Add to `gitmap/constants/constants_cli.go`:

```go
const (
    CmdScanAll              = "all"
    CmdScanAllAlias         = "a"
    FlagScanAllWorkers      = "--workers"
    FlagScanAllPruneMissing = "--prune-missing"
    ScanAllDefaultWorkers   = 4
    ScanAllMaxWorkers       = 16
)
```

Add to `gitmap/constants/constants_messages.go`:

```go
const (
    MsgScanAllNoRoots      = "No previously-scanned roots found. Run 'gitmap scan <dir>' first."
    MsgScanAllSummary      = "Scanned %d root(s) · %d missing · %d repos discovered · %s"
    MsgScanAllMissingHint  = "Run 'gitmap sf rm <path>' to remove them"
    MsgScanAllPrunePrompt  = "Prune them from the database? [y/N] "
)
```

No magic strings anywhere in the new code.

## Error Handling (per Code Red policy)

Every per-root failure is logged to `os.Stderr` with the standardized
format `[scan-all] <root>: <err>` and aggregated into the final summary.
Errors are **never swallowed** — even a single failed `git remote -v`
inside one repo is surfaced in the per-root buffered log, and the root
counts as `⚠` rather than `✓`.

## Testing

- Unit: `TestScanAll_EmptyScanFolder` (exit 3, friendly message).
- Unit: `TestScanAll_AllRootsMissing` (no scans run, prune prompt skipped
  in non-TTY).
- Integration: spin up 3 temp dirs with fake repos, run `scan all`,
  assert each gets its own `.gitmap/output/gitmap.json`.
- Integration: corrupt one root mid-scan (chmod 000), assert exit 1 and
  the other roots still complete.
- Race: run `scan all` concurrently in two processes, assert second exits
  with code 2 (lock held).

## Open Questions

None — all four design questions resolved on 2026-04-20.

## Out of Scope

- Distributed / remote scanning (no SSH-into-other-machines).
- Per-root config overrides (use `gitmap sf` to manage labels).
- Rolling back partial scans on failure (each root is independent).
