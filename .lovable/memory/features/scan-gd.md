---
name: Scan GD
description: gitmap scan gd / scan github-desktop registers every repo under the current scan root in GitHub Desktop, sequentially, idempotent
type: feature
---

# Feature: `gitmap scan gd` (bulk GitHub Desktop registration)

**Spec:** `spec/01-app/102-scan-gd.md`
**Status:** planned for v3.35.0
**Depends on:** ScanFolder table (spec 90), CWD resolver (spec 101), desktop-sync infra (spec 11)

## Behavior summary

- **Aliases:** `gitmap scan gd` (canonical short), `gitmap scan github-desktop` (long), `gitmap s gd`. Subcommand dispatch in `cmd/scan.go`. **Rejects** all existing `scan` flags (no filesystem walk, no DB upsert).
- **Scope (strict):** match CWD against `ScanFolder.AbsolutePath` (exact or longest ancestor). No fallback. If CWD isn't a scan root, exit 3 with hint to run `gitmap scan .` or `gitmap scan <dir> --github-desktop`.
- **Operation:** read repos from DB by `ScanFolderId`, call `desktop.RegisterRepo` for each. Idempotent — Desktop dedupes by path internally.
- **Concurrency:** **sequential**, not parallel. The Desktop IPC backend on Windows doesn't benefit from parallelism, and 100 repos completes in seconds anyway. `--workers` flag is rejected.
- **Failure model:** continue on failure, collect failures, print summary `✓ N registered · ⚠ M failed`. Exit 1 if any failed.
- **Platform:** Windows-only today. macOS/Linux invocation exits 4 with friendly message. (Spec 11 owns Desktop integration scope.)
- **Coexists with `--github-desktop` flag:** `scan <dir> --github-desktop` registers as a side-effect of a fresh scan; `scan gd` registers existing DB repos without re-walking. Both stay shipped, both documented.

## Exit codes

| 0 | All repos registered                                             |
| 1 | One or more registrations failed                                 |
| 2 | Lock contention                                                  |
| 3 | CWD not under any registered scan root (informational)           |
| 4 | Platform unsupported (non-Windows)                               |

## Constants (no magic strings)

`constants_cli.go`: `CmdScanGd="gd"`, `CmdScanGithubDesktop="github-desktop"`.

`constants_messages.go`: `MsgScanGdHeader`, `MsgScanGdSummaryOk`, `MsgScanGdSummaryFailed`, `MsgScanGdNotInScanRoot`, `MsgScanGdUnsupportedOS`.

## Why these decisions

- **Strict ScanFolder scope** — chosen by user. Predictable; matches `pull all` (spec 101). Avoids accidentally registering repos imported via `clone` that live outside the current root.
- **Keep `--github-desktop` flag, no deprecation** — chosen by user. Two distinct verbs serve two distinct workflows (discovery-time vs. retrofit). Deprecation would break muscle memory for no real gain.
- **Continue + summarize on failure** — chosen by user. Same ethos as `scan all` / `pull all`. Surfaces all problems in one pass.
- **Sequential execution** — implementation choice (not asked of user). Desktop IPC on Windows uses a single channel; parallelism adds complexity without speed.
