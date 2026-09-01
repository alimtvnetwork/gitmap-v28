# Issue: Update Sync Lock Loop (v2.3.9 → v2.3.11)

**Status**: ✅ Resolved

## Problem

`gitmap-v28 update` looped 20 times trying to sync the active PATH binary, failing each time due to file locks.

## Root Causes

1. **Parent held lock during sync**: `runUpdate()` used `cmd.Run()` to wait for handoff — parent process (active PATH binary) was still alive during PATH sync step. Windows blocked `Copy-Item`.
2. **Copy-first instead of rename-first**: `run.ps1` PATH sync tried `Copy-Item` (overwrite) as primary strategy. Windows blocks overwrite of running `.exe` but allows rename.
3. **`Read-Host` in non-interactive session**: Generated update script ended with `Read-Host` which fails in non-interactive PowerShell via `exec.Command`.
4. **Syntax error from incomplete refactoring**: Switching between `cmd.Run()` and `cmd.Start()` lost a closing brace, causing build failure.

## Key Insight

The parent holds a lock on `E:\bin-run\gitmap.exe`. The handoff runs from `gitmap-update-<pid>.exe` — a DIFFERENT file. So `cmd.Run()` (foreground) is correct because the lock conflict is between the PARENT and the PATH SYNC step, not between parent and worker. Rename-first resolves this by renaming (not overwriting) the locked binary.

## Solutions Applied

1. **`cmd.Run()` (foreground/blocking)** — parent waits for worker, terminal stays stable. Handoff copy is a different binary so no lock conflict for deploy.
2. **Rename-first PATH sync in `-Update` mode** — rename active binary to `.old` first, then copy new one.
3. **Removed `Read-Host`** — update script exits cleanly without user input.
4. **Diagnostic logs** — prints active exe path and handoff copy path at update start.
5. **Unique temp script names** — `gitmap-update-*.ps1` to avoid stale collisions.

## Prevention Rules — Do Not Repeat

1. **Always use `cmd.Run()` in `runUpdate()`** — foreground execution keeps terminal stable
2. **NEVER use `cmd.Start()` + `os.Exit(0)` in `runUpdate()`** — async detach breaks the session
3. **Always use rename-first for PATH sync during update** — copy-overwrite is unreliable on Windows
4. **Never add interactive prompts to generated scripts** — they run in non-interactive PowerShell
5. **After switching between `cmd.Run()` and `cmd.Start()`, verify closing brace** — mechanical error that breaks build
6. **Any update-flow change must update ALL of:** `cmd/update.go`, `run.ps1`, `spec/01-app/09-build-deploy.md`, `spec/03-general/02-powershell-build-deploy.md`, `spec/03-general/03-self-update-mechanism.md`, `spec/02-app-issues/`
