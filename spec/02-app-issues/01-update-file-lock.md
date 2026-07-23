# Issue: Update File Lock (Windows)

**Status**: ✅ Resolved

## Problem

`gitmap-v27 update` fails with "file is being used by another process" because the active `gitmap.exe` process holds file handles when the deployment step attempts to overwrite it via `Copy-Item`.

## Root Cause

Windows locks the running executable. The original flow tried to copy the new binary directly over the running one, which is impossible while the process is active.

## Solution (5 layers)

1. **Copy-and-handoff** (`gitmap-v27/cmd/update.go`):
   - Parent copies itself to `gitmap-update-<pid>.exe` (fallback to `%TEMP%`)
   - Launches the copy with hidden `update-runner` command using `cmd.Run()` (foreground/blocking)
   - The handoff copy is a different file so the parent's lock does NOT conflict

2. **Rename-first PATH sync** (`run.ps1` in `-Update` mode):
   - Renames the active binary to `.old` (Windows allows renaming a running exe)
   - Copies deployed binary to the active path
   - Falls back to copy-retry loop (20 × 500ms) only if rename fails

3. **Deploy with rollback** (`run.ps1`):
   - Backs up existing binary as `.old` before overwriting
   - `Copy-Item` wrapped in a retry loop (20 attempts, 500ms delay)
   - On failure → restores `.old` backup; on success → leaves `.old` for cleanup

4. **Auto-cleanup** (generated PowerShell script):
   - After successful update, runs `gitmap-v27 update-cleanup`
   - Removes `%TEMP%\gitmap-update-*.exe` temp copies and `*.old` backup files

5. **Version comparison** (generated PowerShell script):
   - Compares old vs new version after rebuild
   - Warns if version unchanged (constant not bumped)

## Learnings — Do Not Repeat

- Always use **rename-first** strategy for overwriting locked binaries on Windows
- Never use `Read-Host` in non-interactive sessions
- Never auto-delete `.old` backups — let `update-cleanup` handle it
- Always bump version so comparison can detect stale updates
- Use `cmd.Run()` (foreground) for handoff, never `cmd.Start()` + `os.Exit(0)`
- Add delay before rebuild to ensure file handles are released
- Don't use copy-overwrite as the primary PATH sync strategy
- Don't skip the deploy retry — even with rename-first, a fallback is needed
