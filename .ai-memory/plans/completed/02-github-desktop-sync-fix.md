# Completed Plan: 02-github-desktop-sync-fix

## User Request (Verbatim)
> In Ubuntu, if we do run the GitKraken, it does not add all the projects to GitHub Desktop. It adds some, but then again crashes or I don't know what happens. You can try to investigate the GitKraken errors to see if there is any error logs and make sure you do not swallow any error or just skip any error. Remember that, try to fix it.

---

## 1. Executive Summary & Root Cause Analysis
- **Root Cause 1 (Launcher Detachment & Exit 0 Spoofing):** On Linux, GitHub Desktop's launcher script (`/usr/lib/github-desktop/resources/app/static/github`) executes Node/Electron in the background (`&`) and exits immediately with code `0`.
- **Root Cause 2 (Electron IPC Flooding & Modal Collision):** GitMap iterated through 70+ repositories in tight unthrottled loops, launching 70+ detached processes in milliseconds. Electron's deep link handler (`x-github-client://openLocalRepo/...`) received these simultaneous requests, causing its React UI renderer to attempt opening multiple `AddRepository` modal popups simultaneously:
  `warn: [ui] Attempted to add a popup of already existing type - AddRepository.`
  This resulted in UI collisions, dropped repository additions (~95% lost), and client crashes.
- **Root Cause 3 (Missing Canonical Subcommand):** Invocations executed `github <path>` instead of `github open <path>`.
- **Root Cause 4 (Swallowed Errors):** In `cli/cmdide/ide_sync.go` and `cli/cmdide/ide_add.go`, errors were discarded via `exec.Command(cli, p).Run() == nil`, masking failures and reporting 100% false success.
- **Root Cause 5 (Candidate Path Omission):** `cli/desktop/resolve.go` only checked user-local `$HOME/.local/...` and missed `/usr/bin/github`, `/usr/bin/github-desktop`, `/usr/local/bin/github`, `/usr/lib/github-desktop/...`, and `github-desktop` in `exec.LookPath`.

---

## 2. Completed Subtasks & Verified Deliverables

### Subtask 01: CLI Resolution & Desktop Package Fix
- **Owner:** Worker 01
- **Files Modified:**
  - `cli/constants/constants_messages.go`: Added `GitHubDesktopAltBin` and `GitHubDesktopVerbOpen`.
  - `cli/desktop/resolve.go`: Added `github-desktop` probe to `exec.LookPath` and expanded `linuxCandidates()` with system and package paths (`/usr/bin/github`, `/usr/bin/github-desktop`, `/usr/local/bin/github`, `/usr/local/bin/github-desktop`, `/usr/lib/github-desktop/resources/app/static/github`, `/snap/bin/github-desktop`, Flatpak paths).
  - `cli/desktop/desktop.go`: Changed command dispatch to `exec.Command(cli, constants.GitHubDesktopVerbOpen, repoPath)`, captured `CombinedOutput()`, enriched failure error messages, and added 250ms inter-iteration pacing delay in `addAll`.
- **Status:** `DONE`

### Subtask 02: Pacing & Error Surfacing in IDE Sync & Desktop Sync
- **Owner:** Worker 02
- **Files Modified:**
  - `cli/cmdide/ide_types.go`: Added `IDESyncFailure` struct, `DesktopFailures []IDESyncFailure` to `IDESyncSummary`, and `HasErrors bool` / `ErrorDetails map[string]string` to `IDEActionResult`.
  - `cli/cmdide/ide_sync.go`: Passed canonical `"open"` argument, captured `CombinedOutput()`, populated `DesktopFailures` on failure without swallowing errors, added 250ms pacing delay between items in `syncDesktopPaths`, and rendered diagnostics via `printDesktopFailureDetails`.
  - `cli/cmdide/ide_add.go`: Passed `"open"` argument with `CombinedOutput()`, surfaced failures in `ErrorDetails["desktop"]` with `HasErrors = true`, and rendered failure diagnostics in `printActionResult`.
  - `cli/cmddesktopsync/desktopsync_ops.go`: Passed `"open"` argument, trimmed output on error, and applied 250ms pacing delay between items in `syncAll`.
- **Status:** `DONE`

---

## 3. Specifications & Verification Evidence
- Architecture Spec: `02-spec/21-app/02-github-desktop-sync-fix/01-architecture-spec.md`
- Component Spec: `02-spec/21-app/02-github-desktop-sync-fix/02-component-spec.md`
- FastGate pre-commit lint: `PASS`
- Unit tests: `go test -C cli -count=1 ./cmdide/... ./cmddesktopsync/... ./desktop/...` -> `PASS`
