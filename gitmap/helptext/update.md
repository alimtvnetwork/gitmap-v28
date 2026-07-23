# gitmap update

Self-update gitmap from the source repository. Pulls latest, rebuilds, and deploys.

## Alias

None

## Usage

    gitmap update [--repo-path <path>] [--verbose] [--report-errors json [--report-errors-file <path>]] [--debug-repo-detect] [--debug-windows]

## Flags

| Flag | Description |
|------|-------------|
| `--repo-path <path>` | Override the source repository path for this run |
| `--verbose` | Enable verbose logging to file |
| `--report-errors json` | Append a JSON-Lines entry for every non-fatal failure during the build/deploy phase (e.g. `npm install` or `npm run build` failing) so CI can branch on them without parsing prose. |
| `--report-errors-file <path>` | Write the JSONL report to this path. When omitted, the file is auto-created at `<TMP>/gitmap-update-report-YYYYMMDD-HHMMSS.jsonl`. |
| `--debug-repo-detect` | Print marker checks (`gitmap/main.go`, `package.json`, `vite` dep, `node_modules`, prebuilt `dist/` locations, npm on PATH) and the resulting decision (`use-prebuilt-*`, `auto-build`, `skip-no-build-script`, `skip-not-a-vite-repo`, `use-legacy-source`, `no-docs-source`). When combined with `--report-errors json`, entries are mirrored under `stage="repo-detect"`. |
| `--debug-windows` | Print a `[debug-windows]` dump on every phase of the self-update handoff: phase name, GOOS, self executable, self/parent PIDs, resolution source (`config`/`sibling`/`PATH`), resolved cleanup target, target-exists check, child argv, relevant env vars (`GITMAP_DEBUG_WINDOWS`, `GITMAP_UPDATE_CLEANUP_DELAY_MS`, `GITMAP_DEBUG_REPO_DETECT`, `GITMAP_REPORT_ERRORS`, `GITMAP_REPORT_ERRORS_FILE`, `PATH`, `GITMAP_DEPLOY_PATH`), spawned child PID, and the path to the durable handoff log file. **As of v3.90.0** the dump also includes (a) the exact shell-quoted spawn command line that Phase 3 will execute (copy-paste safe in PowerShell/cmd/bash/zsh) plus an explicit "no `git` subprocess is launched" note, and (b) a pre-flight enumeration of every `filepath.Glob` pattern and the matching `os.Remove`/`os.RemoveAll` targets the deployed binary will operate on, so the planned filesystem changes are visible before any deletion happens. The flag is propagated through Phase 2 and Phase 3 via both argv and the `GITMAP_DEBUG_WINDOWS=1` env bridge so the dump runs on both sides of the detached cleanup spawn. Despite the name, it works on Unix too. |
| `--debug-windows-json[=<path>]` | **(v3.91.0)** Mirror every `[debug-windows]` event to a structured NDJSON file. Default path: `output/gitmap-debug-windows-<timestamp>.jsonl`; pass `--debug-windows-json=/path/to/trace.jsonl` to override, or set `GITMAP_DEBUG_WINDOWS_JSON=<path>`. One JSON object per line with a stable envelope (`ts`, `event`, `pid`, `ppid`, `goos`, `self`, `version`) plus event-specific fields. Events: `header`, `footer`, `handoff`, `child_pid`, `note`, `command_plan`, `cleanup_plan`. The opened path is auto-forwarded to the Phase 3 cleanup child via env+argv, so both phases append to the **same file** for one consolidated trace per handoff. Sink is off by default; `--debug-windows` alone keeps the v3.90 console-only behaviour. File-open failures degrade silently to console-only. |

## Handoff log file

The Phase 3 cleanup handoff (and the cleanup child itself) **always** write structured events to a durable log file at:

    <TMP>/gitmap-update-handoff-YYYYMMDD.log

This happens regardless of `--verbose` or `--debug-windows` so failures are recoverable even when stdout/stderr is swallowed by an intermediate launcher (run.ps1 wrappers, hidden Windows process attrs, detached spawns). Each line is a single key=value record:

    2026-04-24T12:34:56Z pid=12345 ppid=12000 goos=windows phase=phase-3 event=resolve source=config target=C:\bin\gitmap.exe
    2026-04-24T12:34:56Z pid=12345 ppid=12000 goos=windows phase=phase-3 event=start_ok target=C:\bin\gitmap.exe pid=23456
    2026-04-24T12:34:58Z pid=23456 ppid=12345 goos=windows phase=cleanup event=start self=C:\bin\gitmap.exe
    2026-04-24T12:35:00Z pid=23456 ppid=12345 goos=windows phase=cleanup event=done removed=3

The path is also printed once on every `gitmap update` run via `→ Handoff log file: ...` so you always know where to look.

## Prerequisites

- Git must be installed
- Source repository must be accessible

## Examples

### Example 1: Update to a newer version

    gitmap update

**Output:**

    ■ Checking for updates...
    Current version: v2.19.0
    Latest version:  v2.22.0
    v2.19.0 → v2.22.0
    ■ Pulling latest source...
    ■ Building gitmap.exe...
    ■ Deploying to E:\bin-run\gitmap.exe...
    ✓ Updated to v2.22.0
    → Run 'gitmap changelog --latest' to see what's new

### Example 2: Already up to date

    gitmap update

**Output:**

    ■ Checking for updates...
    Current version: v2.22.0
    Latest version:  v2.22.0
    ✓ Already up to date (v2.22.0)

### Example 3: Update with custom repo path

    gitmap update --repo-path C:\Projects\gitmap-v28

**Output:**

    → Repo path: C:\Projects\gitmap-v28
    ■ Pulling latest source...
    ■ Building gitmap.exe...
    ✓ Updated to v2.49.1

### Example 4: Update with network error

    gitmap update

**Output:**

    ■ Checking for updates...
    ✗ Failed to pull latest: network timeout
    → Check your internet connection and try again

### Example 5: No source repo linked — clone into new path

    gitmap update

**Output:**

    ⚠ The saved source repository path no longer exists on disk.

    Enter the new path to the gitmap source repo: D:\gitmap

    ■ Path does not exist. Cloning gitmap source into D:\gitmap...
    Cloning into 'D:\gitmap'...
    ✓ Cloned successfully.
    → Repo path: D:\gitmap
    ■ Pulling latest source...
    ■ Building gitmap.exe...
    ✓ Updated to v2.56.1

### Example 6: No source repo linked and no path provided

    gitmap update

**Output:**

    ✗ Source repository path not found.

    This binary was installed without a linked source repo, so 'update'
    cannot locate the code to pull and rebuild.

    How to fix:

      Option 1 — Re-install via the one-liner (recommended):
        irm https://raw.githubusercontent.com/.../install.ps1 | iex

      Option 2 — Clone the repo and build from source:
        git clone https://github.com/.../gitmap-v28.git C:\gitmap-src
        cd C:\gitmap-src
        .\run.ps1

      Option 3 — Download the latest release manually:
        https://github.com/.../gitmap-v28/releases/latest

      Option 4 — Use --repo-path to specify it manually:
        gitmap update --repo-path C:\gitmap-src

      After building from source, 'gitmap update' will work automatically.

## Updater Fallback

If no source repo is available and `gitmap-updater` is installed, `gitmap update`
automatically delegates to it. The updater checks GitHub releases and downloads
the latest version without needing a local source checkout.

    gitmap update

**Output (with updater installed):**

    → No source repo found. Delegating to gitmap-updater...

    ■ Checking for updates...
    Current version: v2.49.0
    Latest version:  v2.49.1
    v2.49.0 → v2.49.1
    ■ Downloading installer for v2.49.1...
    ■ Running installer...
    ✓ Update complete.

## Troubleshooting

If you installed gitmap from a GitHub release (e.g. via the one-liner installer),
the binary does not have a source repo path embedded. You have three choices:

1. **Install `gitmap-updater`** — it handles updates via GitHub releases automatically.
2. **Use `--repo-path`** to point at a local clone for a one-time update.
3. **Clone and rebuild** from source so future updates work automatically.

## Reporting non-fatal failures (`--report-errors json`)

The auto-build step (`npm install` and `npm run build` for the docs site) is
intentionally non-fatal — a failed build never aborts `gitmap update`. To make
those failures visible to CI without scraping logs, pass `--report-errors json`.

Each failure is appended as one JSON object per line (JSONL) to the report
file. The CLI prints the file path before the update starts and a summary line
after it finishes.

### Example

    gitmap update --report-errors json --report-errors-file C:\logs\update.jsonl

**Output (with one auto-build failure):**

    -> Error report (json): C:\logs\update.jsonl
    ...
    !! Auto-build failed - 'gitmap hd' will fail
    ...
    -> Wrote 1 non-fatal failure entry to C:\logs\update.jsonl

### Entry schema

Each line contains a single object with these fields:

| Field | Type | Notes |
|-------|------|-------|
| `timestamp` | string (ISO-8601 UTC) | When the failure was recorded |
| `stage` | string | Stable identifier — currently `docs-npm-install` or `docs-npm-build` |
| `command` | string | The exact command that failed |
| `exitCode` | integer | Process exit code (0 if the command "succeeded" but post-conditions failed) |
| `cwd` | string | Working directory at time of failure |
| `message` | string | Human-readable summary |
| `paths` | object | Stage-specific paths (e.g. `repoRoot`, `packageJson`, `expectedDist`) |
| `os` | string | `windows` or `unix` |

Without `--report-errors json`, the env vars are not set and the scripts skip
the writer entirely — there is zero overhead for normal updates.

## See Also

- [version](version.md) — Check current version
- [doctor](doctor.md) — Diagnose installation issues
- [changelog](changelog.md) — View release notes for new version

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter update
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
