# 102 — `gitmap scan gd` (Bulk GitHub Desktop Registration)

> Status: **planned** — targets v3.35.0
> Depends on: spec 90 (`ScanFolder` table), spec 101 (CWD→ScanFolder resolver), spec 11 (`desktop-sync` infra)

## Goal

Register every repo under the **current scan root** in GitHub Desktop in
one shot, without re-running the full scanner. Useful after a fresh
machine setup or after enabling Desktop on a project that was previously
scanned with `--github-desktop` off.

## CLI Surface

| Form                          | Notes                                               |
|-------------------------------|-----------------------------------------------------|
| `gitmap scan gd`              | Canonical short form.                               |
| `gitmap scan github-desktop`  | Long form, matches the existing `--github-desktop` flag name. |
| `gitmap s gd`                 | Inherits the existing `s` alias for `scan`.         |

Implemented as a **subcommand dispatch** inside `cmd/scan.go`, parallel
to `scan all` (spec 100). Existing `scan` flags are **rejected** (this
mode does not walk the filesystem; it reads from the DB).

### Relationship to `scan --github-desktop`

These are two distinct verbs and **both stay shipped**:

| Command                       | What it does                                      |
|-------------------------------|---------------------------------------------------|
| `scan <dir> --github-desktop` | Walk filesystem, discover repos, register each as a side-effect of the scan. Used during initial discovery. |
| `scan gd`                     | Read existing repos from DB (under current scan root) and register each in Desktop. **No filesystem walk, no DB upsert.** Used after the fact, when `--github-desktop` was off during the original scan. |

The docs entry for `scan` will get a "See also: `scan gd`" link clarifying
this. No deprecation. No silent redirect.

## Scope Resolution — CWD → ScanFolder (Strict)

Same resolver as spec 101 (`pull all`):

```go
cwd, _ := filepath.EvalSymlinks(mustGetwd())
sf, err := db.FindScanFolderForPath(cwd) // exact OR longest ancestor
if err != nil { /* exit 3 */ }
repos := db.SQLSelectReposByScanFolderId(sf.ScanFolderId)
```

If CWD is not under any registered `ScanFolder`, exit 3 with:

```
Current directory is not a registered scan root.
Run 'gitmap scan .' here first, or 'gitmap scan <dir> --github-desktop' to scan + register in one pass.
```

**No fall-back** to "every repo in the database" — the user explicitly
chose the strict scope. Same rationale as `pull all`.

## Registration Logic

Reuse the existing `desktop.RegisterRepo(absPath)` helper from spec 11
(`desktop-sync`). For each repo:

```go
for _, r := range repos {
    if err := desktop.RegisterRepo(r.AbsolutePath); err != nil {
        failures = append(failures, repoFailure{Path: r.AbsolutePath, Err: err})
        continue
    }
    registered++
}
```

GitHub Desktop dedupes by path internally, so re-registering an
already-known repo is a no-op — the command is **idempotent**. Safe to
re-run.

### Sequential, not parallel

Unlike `scan all` (spec 100) and `pull all` (spec 101), Desktop
registration is **sequential**. Reasons:

- The Desktop registration backend on Windows uses a single named pipe
  / registry write that does not benefit from parallelism.
- Total time is dominated by the per-call IPC handshake (~50ms),
  not by I/O. 100 repos completes in a few seconds even sequentially.
- Avoids the locking complexity of parallel writes against Desktop's
  internal store.

A `--workers` flag is **not** accepted (rejected with a friendly hint
to remove it).

## Failure Model — Continue + Summarize

Per the user's chosen "continue, summarize at end" semantics:

```
Registering 142 repo(s) under D:\projects in GitHub Desktop...
✓ 139 registered
⚠ 3 failed:
    D:\projects\old-repo (path no longer exists)
    D:\projects\broken (registration error: <details>)
    \\nas\shared\experiment (path inaccessible)

Done in 4.2s. Exit 1 (some failures).
```

- Each failure logged to `os.Stderr` with format
  `[scan-gd] <path>: <err>` (Code Red policy, no swallowing).
- Final exit code: **0** if all succeeded, **1** if any failed.
- A repo whose `AbsolutePath` no longer exists on disk is **not** auto-
  pruned from the DB — that's the job of `gitmap prune`. We just skip
  and report.

## Exit Codes

| Code | Meaning                                                          |
|------|------------------------------------------------------------------|
| 0    | All repos under the scan root registered successfully.           |
| 1    | One or more repos failed to register.                            |
| 2    | Lock contention — another `gitmap` is already running.           |
| 3    | CWD is not under any registered scan root (informational).       |
| 4    | Platform unsupported (GitHub Desktop integration is Windows-only today). |

## Locking

Acquire `gitmap.lock` once at the top
(`mem://tech/process-synchronization`). Sequential execution means no
inner contention.

## Constants (per `mem://style/code-constraints`)

`gitmap/constants/constants_cli.go`:

```go
const (
    CmdScanGd          = "gd"
    CmdScanGithubDesktop = "github-desktop"
)
```

`gitmap/constants/constants_messages.go`:

```go
const (
    MsgScanGdHeader        = "Registering %d repo(s) under %s in GitHub Desktop..."
    MsgScanGdSummaryOk     = "✓ %d registered"
    MsgScanGdSummaryFailed = "⚠ %d failed:"
    MsgScanGdNotInScanRoot = "Current directory is not a registered scan root.\nRun 'gitmap scan .' here first, or 'gitmap scan <dir> --github-desktop' to scan + register in one pass."
    MsgScanGdUnsupportedOS = "GitHub Desktop integration is currently Windows-only. (See spec 11.)"
)
```

No magic strings. No raw `"gd"` / `"github-desktop"` literals outside
`constants_cli.go`.

## Error Handling (Code Red)

- Each `desktop.RegisterRepo` failure: log to `os.Stderr` immediately,
  append to `failures` slice for the summary block.
- Wrap errors with `fmt.Errorf("scan-gd register %s: %w", path, err)` so
  callers can `errors.Is(err, desktop.ErrPipeClosed)` upstream.
- Never swallow.

## Testing

- Unit: `TestScanGd_CwdNotInScanRoot` → exit 3, friendly message.
- Unit: `TestScanGd_RejectsScanFlags` → `--exclude` / `--mode` /
  `--workers` all return parse errors with helpful hints.
- Integration: 5 fake repos under one scan root, mock `desktop.RegisterRepo`
  to fail for 2 of them. Assert exit 1, summary lists the 2 failures,
  the 3 successes still got registered.
- Integration (Windows-only): real Desktop registration of a temp repo,
  verify it appears in `%APPDATA%\GitHub Desktop\repositories.json`.
- Cross-platform: macOS/Linux invocation → exit 4 with friendly message.

## Out of Scope

- Bulk **un-registration** from Desktop (covered by future
  `gitmap unregister gd` if requested).
- Group / glob filters within the scan root — use existing
  `gitmap group` to slice further first.
- Non-Windows Desktop integration (tracked separately; not blocking
  this command — it just exits 4 cleanly).
