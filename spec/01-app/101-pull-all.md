# 101 — `gitmap pull all` (CWD-Scoped Bulk Pull + run-script)

> Status: **planned** — targets v3.34.0
> Depends on: spec 90 (`ScanFolder` table, v3.7.0), spec 100 (`scan all` worker-pool model)

## Goal

When the user is `cd`'d into a directory that is (or sits inside) a known
scan root, `gitmap pull all` updates **every repo registered under that
root** in parallel. For repos that ship a `run.ps1` (Windows) or `run.sh`
(Unix), the script is executed **instead of** `git pull` — the script is
trusted to handle its own update / build / deploy logic.

## CLI Surface

| Form              | Notes                                                |
|-------------------|------------------------------------------------------|
| `gitmap pull all` | Canonical form. CWD must be a known scan root.       |
| `gitmap p all`    | Inherits the existing `p` alias for `pull`.          |
| `gitmap pull a`   | Short alias, mirrors `scan a`.                       |

Implemented as a **subcommand dispatch** inside `cmd/pull.go`, parallel to
the `scan all` pattern in spec 100. Existing `pull` flags (`--group`,
`--verbose`) are accepted; `--all` is **rejected** in this mode (the CWD
already defines the scope, and `pull --all` already means "every tracked
repo" globally — keep them distinct).

## Scope Resolution — ScanFolder.AbsolutePath

```go
cwd, _ := filepath.EvalSymlinks(mustGetwd())
sf, err := db.FindScanFolderForPath(cwd) // exact match OR ancestor
if err != nil { /* exit 3, friendly message */ }
repos := db.SQLSelectReposByScanFolderId(sf.ScanFolderId)
```

`FindScanFolderForPath` matches CWD against `ScanFolder.AbsolutePath`
either exactly or as an ancestor (`strings.HasPrefix(cwd, sf.AbsolutePath + sep)`).
If multiple ancestors match, pick the **longest** (most specific) — a
nested scan root wins over its parent.

If CWD is not under any `ScanFolder` row, exit 3 with:

```
Current directory is not a registered scan root.
Run 'gitmap scan .' here first, or 'gitmap pull --all' to pull every tracked repo.
```

This is informational, not a hard error. **No fall-back to path-prefix
matching** — the user explicitly chose the strict ScanFolder semantic.

## Per-Repo Update Strategy — run-script Replaces git pull

For each repo in `repos`:

```go
script := detectRunScript(repo.AbsolutePath) // "run.ps1" on Windows, "run.sh" on Unix
if script != "" {
    runScript(repo.AbsolutePath, script)   // INSTEAD of git pull
} else {
    gitPull(repo.AbsolutePath)
}
```

### `detectRunScript`

| GOOS       | Lookup order              | Interpreter                           |
|------------|---------------------------|---------------------------------------|
| `windows`  | `run.ps1`                 | `pwsh.exe` → fallback `powershell.exe`|
| `linux`/`darwin` | `run.sh`            | `bash` (must be executable)           |

A Windows repo's `run.ps1` is **not** invoked on Unix and vice-versa —
this matches the user's chosen "Windows: run.ps1 · Unix: run.sh" rule.
Cross-OS execution is out of scope for v3.34.0.

### Script execution rules

- Working directory: the repo root.
- Stdin: `/dev/null` (script must be non-interactive).
- Stdout/stderr: captured into the per-repo log buffer (see below).
- Timeout: `PullAllScriptTimeout = 10 * time.Minute` (override via
  `--script-timeout <duration>`).
- Exit code ≠ 0 marks the repo as ✗ in the summary; pull continues for
  the rest (per chosen "continue on failure" semantics).
- Environment: parent env + `GITMAP_PULL_ALL=1` so scripts can detect
  bulk-mode and skip interactive prompts.

### Why "instead of", not "after"

The user explicitly chose *replaces*. Rationale: a repo that ships
`run.ps1` typically does its own `git pull` (often with extra logic like
`--rebase`, submodule sync, or LFS fetch) and then builds/deploys.
Invoking `git pull` first would either duplicate work or fight the
script's intent.

## Concurrency — Parallel N=4, Continue on Failure

Same worker-pool topology as spec 100 (`scan all`). One repo failing
**does not** abort the rest. Per-repo logs are buffered and flushed
atomically when each worker finishes. Final summary:

```
Pulled 12 repo(s) · 9 git pull · 3 run.ps1 · 1 failed · 8.3s
```

Constants:

- `PullAllDefaultWorkers = 4`
- `PullAllMaxWorkers     = 16`
- Override via `--workers <n>` (1–16).

### Locking

Acquire `gitmap.lock` once at the top, shared by all workers
(`mem://tech/process-synchronization`). A second concurrent `pull all`
exits with code 2.

## Exit Codes

| Code | Meaning                                                          |
|------|------------------------------------------------------------------|
| 0    | All repos updated successfully.                                  |
| 1    | At least one repo failed (`git pull` non-zero or script exit≠0). |
| 2    | Lock contention — another `gitmap` is already running.           |
| 3    | CWD is not under any registered scan root (informational).       |

## Constants (per `mem://style/code-constraints`)

`gitmap/constants/constants_cli.go`:

```go
const (
    CmdPullAll              = "all"
    CmdPullAllShort         = "a"
    FlagPullAllWorkers      = "--workers"
    FlagPullAllScriptTimeout = "--script-timeout"
    PullAllDefaultWorkers   = 4
    PullAllMaxWorkers       = 16
    PullAllScriptTimeout    = 10 * time.Minute
    RunScriptWindows        = "run.ps1"
    RunScriptUnix           = "run.sh"
    RunScriptEnvFlag        = "GITMAP_PULL_ALL=1"
)
```

`gitmap/constants/constants_messages.go`:

```go
const (
    MsgPullAllNotInScanRoot = "Current directory is not a registered scan root.\nRun 'gitmap scan .' here first, or 'gitmap pull --all' to pull every tracked repo."
    MsgPullAllSummary       = "Pulled %d repo(s) · %d git pull · %d %s · %d failed · %s"
    MsgPullAllScriptHit     = "[%s] running %s (replaces git pull)"
    MsgPullAllScriptTimeout = "[%s] script %s timed out after %s"
)
```

No magic strings. No raw `"run.ps1"` / `"run.sh"` literals outside the
constants file.

## Error Handling (Code Red policy)

Every per-repo failure logged to `os.Stderr` with format
`[pull-all] <repo>: <err>`. Script timeouts log
`[pull-all] <repo>: script timed out after <dur>`. Errors aggregated into
the final summary count — never swallowed. Use `errors.Is` for the
`context.DeadlineExceeded` check on script timeout.

## Security

- `detectRunScript` calls `filepath.Clean` on the resolved script path
  and verifies it's still inside `repo.AbsolutePath` (path-traversal guard
  per `mem://tech/security-hardening`).
- Script must be a regular file (`fi.Mode().IsRegular()`), not a symlink
  pointing outside the repo.
- On Unix, `run.sh` must have the executable bit set; otherwise log a
  warning and fall back to `git pull`.

## Testing

- Unit: `TestPullAll_CwdNotInScanRoot` → exit 3, friendly message.
- Unit: `TestDetectRunScript_WindowsPrefersPs1`, `TestDetectRunScript_UnixSh`.
- Integration: 3 fake repos under one scan root — one with `run.ps1`,
  one with `run.sh`, one without. Assert the right strategy fires per OS.
- Integration: `run.ps1` exits 1 → repo marked ✗, other repos still pulled.
- Integration: `run.ps1` sleeps past `--script-timeout 2s` → marked ✗
  with timeout message.
- Race: two concurrent `pull all` → second exits 2.

## Out of Scope

- `run.py` / `run.js` / arbitrary script languages — Python/Node-flavoured
  follow-up for a later spec.
- Per-repo override config (`gitmap.repo.json`) for non-default script names.
- Rollback of partial pulls on failure (each repo is independent).
- Submodule recursion control — defer to whatever the script or
  `git pull --recurse-submodules` already does.
