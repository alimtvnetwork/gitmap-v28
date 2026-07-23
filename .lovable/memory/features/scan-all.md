---
name: Scan All
description: gitmap scan all / scan a re-scans every root in ScanFolder table using a parallel worker pool with prune-missing prompt
type: feature
---

# Feature: `gitmap scan all` (a.k.a. `scan a`)

**Spec:** `spec/01-app/100-scan-all.md`
**Status:** planned for v3.33.0
**Depends on:** ScanFolder table (spec 90, shipped v3.7.0)

## Behavior summary

- **Aliases:** `gitmap scan all` (canonical) and `gitmap scan a` (short).
  Subcommand dispatch inside `cmd/scan.go` — **not** a `--all` flag, to
  avoid colliding with existing `scan` flags.
- **Source of roots:** `SELECT AbsolutePath FROM ScanFolder ORDER BY LastScannedAt DESC`.
  No derivation from `Repo.AbsolutePath`. If table is empty, exit 3 with
  friendly message — not a hard error.
- **Execution:** parallel worker pool, default N=4, `--workers` flag
  (1–16). Per-root logs buffered + flushed atomically so output never
  interleaves. SQLite `SetMaxOpenConns(1)` already serializes DB writes.
- **Locking:** acquire `gitmap.lock` once at top, shared by all workers.
  Second concurrent `scan all` exits with code 2.
- **Missing roots:** never auto-pruned during scan. After scan, if any
  missing, prompt `Prune them from the database? [y/N]` on TTY; print
  `gitmap sf rm <path>` hint on non-TTY. `--prune-missing` flag bypasses
  prompt.
- **Forwarded flags:** `--mode`, `--output`, `--github-desktop`, `--quiet`
  pass through to each per-root scan unchanged. `--output-path` is
  **rejected** in `all` mode (each root keeps its own `.gitmap/output/`).

## Exit codes

| 0 | All present roots scanned (missing-root warnings ok)        |
| 1 | One or more roots failed mid-scan                           |
| 2 | Lock contention (another gitmap running)                    |
| 3 | ScanFolder empty (informational)                            |

## Constants (no magic strings)

In `constants_cli.go`: `CmdScanAll`, `CmdScanAllAlias`,
`FlagScanAllWorkers`, `FlagScanAllPruneMissing`, `ScanAllDefaultWorkers=4`,
`ScanAllMaxWorkers=16`.

In `constants_messages.go`: `MsgScanAllNoRoots`, `MsgScanAllSummary`,
`MsgScanAllMissingHint`, `MsgScanAllPrunePrompt`.

## Why these decisions

- **ScanFolder DB table** (not derived parents) — zero new state, table is
  already populated on every scan since v3.7.0.
- **Parallel worker pool** — chosen by user. Filesystem walk + `git
  remote -v` are I/O bound, SQLite serialization handles writes.
- **Skip + offer to prune** — chosen by user. Safest default (drive may
  reappear), but offers cleanup for stale entries.
- **Subcommand `scan all` / `scan a`** — chosen by user. Matches their
  original phrasing, keeps the `scan` flag namespace clean.
