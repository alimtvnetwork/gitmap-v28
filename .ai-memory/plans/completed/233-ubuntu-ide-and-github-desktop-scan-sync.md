# Completed Plan: 233-ubuntu-ide-and-github-desktop-scan-sync

## Metadata
- **Plan Slug:** `233-ubuntu-ide-and-github-desktop-scan-sync`
- **Specification:** [02-spec/21-app/233-ubuntu-ide-and-github-desktop-scan-sync/01-architecture-spec.md](../../02-spec/21-app/233-ubuntu-ide-and-github-desktop-scan-sync/01-architecture-spec.md)
- **Component & CLI Spec:** [02-spec/21-app/233-ubuntu-ide-and-github-desktop-scan-sync/02-component-and-cli-spec.md](../../02-spec/21-app/233-ubuntu-ide-and-github-desktop-scan-sync/02-component-and-cli-spec.md)
- **Root Cause Analysis:** [02-spec/22-app-issues/71-ubuntu-ide-github-desktop-scan-omission.md](../../02-spec/22-app-issues/71-ubuntu-ide-github-desktop-scan-omission.md)
- **Status:** `COMPLETED`
- **Completed Date:** `2026-10-06`
- **Version Milestone:** `v6.497.0`

---

## User Request (Verbatim)
```text
Serious issue in the gitmap. When we do the scan, it should actually add the repos to the GitHub Desktop and VS Code. Apparently, both of those are not happening in the Ubuntu. I want you to do the end-to-end test, try to understand what went wrong and how you can fix it, find the root cause, and fix this stuff. Also write the root cause of it, why it happened, how it happened. Add all these projects to GitHub Desktop and also the VS Code. First do it using scripts and then try to add it using gitmap. gitmap scan should be doing it. If the repo is already there, then it's not cloned again, but it will try to be checked if each one of them are added to the IDE in VS Code, Antigravity, Cursor, GitHub Desktop, and so on. User can also skip. User can use a flag like excuse sync and provide GitHub Desktop, VS Code, things like that. Users can also remove from this addition. There should be a command to do that. Think of the flag, think of the command that will do all these things on top of the repos. Add, sync, all these things. Write the command, write the process, fix the code, and then finally make a release bump, minor bump, and make a release.
```

---

## Deliverables Summary

### 1. Grounded 4-Part Root Cause Analysis (Issue 71)
- **Vector 1 (Auto-Mkdir Omission):** `cli/vscodepm/path.go` checked `!dirExists(extDir)` and returned `ErrExtensionMissing` without auto-creating `~/.config/Code/User/globalStorage/alefragnani.project-manager`. Swallowed as a soft skip during scans.
- **Vector 2 (GitHub Desktop Linux CLI Blindspot):** `cli/desktop/resolve.go` returned `nil` for Linux candidates, preventing fallback resolution when not on standard path.
- **Vector 3 (Scan Toggle Default False):** `cli/cmdscan/flags.go` defaulted `ghDesktopFlag` to false.
- **Vector 4 (Missing IDE Hooks):** `cli/cmdscan/scan.go` lacked Cursor and Antigravity synchronization during scan passes.
- Authored [02-spec/22-app-issues/71-ubuntu-ide-github-desktop-scan-omission.md](../../02-spec/22-app-issues/71-ubuntu-ide-github-desktop-scan-omission.md) and [.ai-memory/issues/71-ubuntu-ide-github-desktop-scan-omission.md](../issues/71-ubuntu-ide-github-desktop-scan-omission.md).

### 2. Standalone Multi-IDE Python Synchronizer
- Created [03-ai-scripts/40-ubuntu-ide-desktop-sync.py](../../03-ai-scripts/40-ubuntu-ide-desktop-sync.py).
- Proactively auto-creates missing IDE storage hierarchies with `0o755` permissions.
- Idempotently synchronizes VS Code Project Manager, Cursor Project Manager, Google Antigravity project descriptors, and GitHub Desktop.
- Verified on Ubuntu host: 74 repositories discovered and synchronized across IDE targets (74 in VS Code, 74 in Cursor, 74 in Antigravity).

### 3. GitMap IDE Command Suite (`cli/cmdide/`)
- Implemented full CLI command group:
  * `gitmap ide add <path> [--target <ide>]`
  * `gitmap ide sync [dir] [--target <ide>] [--exclude <ide>]`
  * `gitmap ide remove <path> [--target <ide>]` (alias: `rm`)
  * `gitmap ide list` (alias: `ls`)
  * `gitmap ide status`
  * `gitmap ide help`
- Wired into root utility dispatcher (`cli/cmd/rootutility.go`), usage groups (`cli/cmd/rootusage_groups.go`), and bridge helpers (`cli/cmd/clihelpers.go`).

### 4. GitMap Scan Integration & User Flags
- Auto-created extension directories via `EnsureProjectsJSONPath()` in `cli/vscodepm/path.go` and wired into `cli/vscodepm/sync.go`.
- Added Linux installation candidate resolution to `cli/desktop/resolve.go`.
- Added `--sync-ide`, `--skip-sync` (`--no-ide-sync`), and `--exclude-sync` flags in `cli/cmdscan/flags.go`.
- Integrated `scan.ideSync` phase in `cli/cmdscan/scan.go` via decoupled hook `SyncRecordsToIDEsFn`.

---

## Verification Evidence
- `go vet ./cmdide ./cmdscan ./desktop ./vscodepm ./store ./cmd`: PASS (exit code 0).
- `go test -v ./cmdide`: PASS (all 4 unit tests passed, 0.013s).
- `go test -v ./desktop`: PASS (all unit tests passed, 0.002s).
- `go test -v ./vscodepm`: PASS (all unit tests passed, 0.008s).
- `python3 03-ai-scripts/40-ubuntu-ide-desktop-sync.py --skip-sync desktop --json`: PASS (74 repositories verified and synchronized).
