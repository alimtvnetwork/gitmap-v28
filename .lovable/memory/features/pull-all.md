---
name: Pull All
description: gitmap pull all updates every repo under the current scan root in parallel; run.ps1/run.sh replaces git pull when present
type: feature
---

# Feature: `gitmap pull all` (CWD-scoped bulk pull)

**Spec:** `spec/01-app/101-pull-all.md`
**Status:** planned for v3.34.0
**Depends on:** ScanFolder table (spec 90), worker-pool pattern (spec 100)

## Behavior summary

- **Aliases:** `gitmap pull all` (canonical), `gitmap p all`, `gitmap pull a`.
  Subcommand dispatch in `cmd/pull.go`. Reject `--all` flag in this mode
  (keep distinct from existing `pull --all` global semantics).
- **Scope (strict):** match CWD against `ScanFolder.AbsolutePath` (exact or
  longest ancestor). **No** path-prefix fallback to `Repo.AbsolutePath`.
  If CWD isn't a registered scan root, exit 3 with hint to run
  `gitmap scan .` or `gitmap pull --all`.
- **Per-repo strategy:** if `run.ps1` (Windows) or `run.sh` (Unix) exists
  in the repo root, execute it **instead of** `git pull`. The script is
  trusted to handle its own pull/build/deploy. Otherwise plain `git pull`.
- **Cross-platform:** Windows looks at `run.ps1` only (`pwsh.exe` →
  `powershell.exe` fallback). Unix looks at `run.sh` only (must be
  executable, bash). No cross-OS script invocation.
- **Concurrency:** parallel worker pool, default N=4, `--workers` 1–16.
  Continue on failure — one repo's error never aborts the batch.
  Per-repo logs buffered + flushed atomically.
- **Script env:** `GITMAP_PULL_ALL=1` injected so scripts can skip
  interactive prompts. Stdin = `/dev/null`. Default timeout 10m,
  override via `--script-timeout <dur>`.
- **Locking:** acquire `gitmap.lock` once at top, shared by workers.
  Second concurrent `pull all` exits 2.

## Exit codes

| 0 | All repos updated                                                |
| 1 | One or more repos failed (git pull non-zero OR script exit ≠ 0)  |
| 2 | Lock contention                                                  |
| 3 | CWD not under any registered scan root (informational)           |

## Constants (no magic strings)

`constants_cli.go`: `CmdPullAll`, `CmdPullAllShort`, `FlagPullAllWorkers`,
`FlagPullAllScriptTimeout`, `PullAllDefaultWorkers=4`,
`PullAllMaxWorkers=16`, `PullAllScriptTimeout=10*time.Minute`,
`RunScriptWindows="run.ps1"`, `RunScriptUnix="run.sh"`,
`RunScriptEnvFlag="GITMAP_PULL_ALL=1"`.

`constants_messages.go`: `MsgPullAllNotInScanRoot`, `MsgPullAllSummary`,
`MsgPullAllScriptHit`, `MsgPullAllScriptTimeout`.

## Security

`detectRunScript` cleans the path, verifies containment inside repo
root, requires `IsRegular()` (no symlinks escaping). On Unix, requires
executable bit; otherwise warn + fall back to `git pull`.

## Why these decisions

- **Strict ScanFolder match (no prefix fallback)** — chosen by user.
  Predictable; avoids accidentally pulling repos imported via `clone`
  that happen to live under CWD but were never associated with a root.
- **run-script REPLACES git pull** — chosen by user. Scripts that ship
  `run.ps1` typically already pull + build + deploy; running `git pull`
  first would duplicate or fight the script's intent.
- **Windows: run.ps1 · Unix: run.sh** — chosen by user. Each platform
  has a native script; no fragile cross-OS interpreter chain.
- **Parallel N=4, continue on failure** — chosen by user. Same model as
  `scan all` (spec 100) and existing `pull --parallel`.
