# Completed Plan: Fleet Update JSON Communication, Running Prompts Backup/Restore E2E & Deploy Polish

**Plan Version:** 1.0.0  
**Status:** Completed  
**Spec Reference:** [02-spec/21-app/179-fleet-update-json-prompts-backup-and-deploy-polish.md](../../../02-spec/21-app/179-fleet-update-json-prompts-backup-and-deploy-polish.md)

---

## Completed Tasks Summary

### Task 01: Fleet Update JSON Communication, Aliases & High-Quality Display
- Added aliases `ua`, `update all`, `update-all`, `updateall`, and `update -all` in `cli/cmd/rootutility.go` and `cli/cmdupdate/update_fleet.go`.
- Resolved duplicate `deploy` routing conflict by removing the shadowing entry in `cli/cmd/rootutility.go` so `gitmap deploy` routes cleanly to `runSSHDeploy` in `cli/cmd/roottooling.go`.
- Implemented structured JSON communication protocol over SSH in `cli/cmdupdate/update_fleet.go` and `cli/cmdssh/ssh_update_remote.go`:
  - Suppressed all 6 PowerShell streams (`*>$null`, `$WarningPreference='SilentlyContinue'`, `$ErrorActionPreference='SilentlyContinue'`, `$env:GITMAP_UPDATING='1'`) on remote Windows nodes and redirected shell output on Unix nodes.
  - Returned structured, single-line JSON summaries containing `current_version`, `previous_version`, `success`, and `details`.
  - Added `extractJSONSubstring` in `ParseFleetUpdateTelemetry` to guard against arbitrary shell banners or MOTDs.
  - Added `cleanTelemetryDetails` and `sanitizeTableRowDetail` to guarantee fallback output never leaks multi-line noise, 404s, or probing logs into table rows.
  - Upgraded `renderFleetUpdateSummary` and `buildFleetUpdateRows` to display a clean, single-line terminal table with columns: `ALIAS`, `IP`, `STATUS`, `DURATION`, `DETAILS`.

### Task 02: Running Prompts End-to-End Discovery, Backup & Restore
- In `cli/cmdagy/agy_running_prompts_io.go`, updated `enqueueSingleRestoredItem` to use `resolveQueueFilePath(item.ProjectPath)` instead of hardcoded paths.
- Verified end-to-end prompt discovery with `gitmap agy running-prompts ls` finding active and enqueued prompts across projects.
- Verified snapshot functionality with `gitmap backup-running-prompts` writing batches into SQLite Split-DB (`data/backup-prompts/sql.db`).
- Verified prompt injection with `gitmap restore-running-prompts` enqueuing backed-up prompts into project queues.

### Task 03: Deploy & Sync Modes Verification
- Verified `gitmap deploy`, `gitmap deploy-right`, and `gitmap deploy-left` CLI argument and flag parsing (`--sync`, `--sync-right`, `--sync-left`, `--overwrite`, `--skip`, `--json`, `--parallel`, `--dry-run`).
- Added `dispatchSCDeployOps` in `cli/cmd/rootcore.go` so `gitmap sc deploy`, `gitmap sc deploy-right`, and `gitmap sc deploy-left` dispatch directly to `cmdssh.RunSSHDeployCLI`.

### Task 04: High-Quality Help Text & UI Documentation
- Created `cli/cmdupdate/update_help_menu.go` with `RenderUpdateRichHelp` using `termhelp.HelpMenu`.
- Created `cli/cmdssh/ssh_deploy_rich_help.go` with `RenderDeployRichHelp` using `termhelp.HelpMenu`.
- Upgraded `RenderRunningPromptsHelp` in `cli/cmdagy/agy_running_prompts_render.go` to use `termhelp.HelpMenu`.
- Registered `update`, `ua`, `deploy`, `running-prompts` in `cli/cmd/rich_help_dispatcher.go`.
- Authored comprehensive markdown documentation:
  - `cli/helptext/update.md`: Fleet update, local self-update, JSON protocol, flags, examples.
  - `cli/helptext/running-prompts.md`: Running prompts lifecycle, Split-DB backup, queue restore, examples.
  - `cli/helptext/deploy.md`: Directional sync, parallel workers, conflict resolution, examples.

---

## Verification Outcomes
- `check-nested-ifs.py`: 0 violations across changed files.
- `check-boolean-guidelines.py`: 0 violations across changed files.
- `check-relative-paths.py`: 0 violations across all 8,164 files.
- `go-format-check.py`: 0 unformatted files across all 3,645 Go files.
- `go build`: Compiled successfully with zero errors.
