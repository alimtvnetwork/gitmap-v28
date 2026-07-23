# 31 — update-cleanup Phase 3 observability gap

## Summary

`gitmap-v27 update` can still appear to "repeat the same cleanup error with no logs" even after the Phase 3 deployed-binary handoff was introduced.

The actual failure point is **inside the detached `update-cleanup` child**, after the Phase 3 parent has already reported a successful spawn. The parent logs `start_ok`, but several inner cleanup branches still write **only to child stderr**:

- `filepath.Glob` failures in `updatecleanup_remove.go`
- `os.Remove` retry exhaustion in `updatecleanup_remove.go`
- drive-root shim skip/remove failures in `updatecleanup_extra.go`
- `os.RemoveAll` failures for `*.gitmap-tmp-*` swap dirs in `updatecleanup_extra.go`

On Windows, that child is launched hidden via `cmd.Start()` + `HideWindow: true`, so those stderr-only diagnostics are exactly the lines most likely to be lost or missed. The user then sees:

1. Phase 3 handoff started successfully
2. Cleanup still failed
3. No durable record of *which exact inner operation* failed

That makes the error feel silent/repeating even though the child did emit some stderr.

## Root cause

**Observability gap, not primary control-flow failure.**

The Phase 3 architecture is correct: the deployed binary must perform cleanup because the handoff copy still holds the file lock. The regression is that the observability upgrade was incomplete.

The repo already had three logging channels:

1. Console stderr / stdout (`--debug-windows`)
2. Durable handoff log (`updatehandofflog.go`)
3. NDJSON sink (`updatedebugwindows_json.go`)

But only the outer lifecycle events (`start`, `done`, `start_fail`, etc.) were mirrored durably. Several **inner cleanup branches** still used direct `fmt.Fprintf(os.Stderr, ...)` without corresponding `logHandoffEvent(...)` and without structured JSON notes. As a result, the system recorded that cleanup *ran*, but not why an inner removal step failed.

## Files involved

- `gitmap-v27/cmd/updatecleanup.go`
- `gitmap-v27/cmd/updatecleanup_remove.go`
- `gitmap-v27/cmd/updatecleanup_extra.go`
- `gitmap-v27/cmd/updatehandoff_phase3.go`
- `gitmap-v27/cmd/updatehandofflog.go`
- `gitmap-v27/cmd/updatedebugwindows.go`
- `gitmap-v27/cmd/updatedebugwindows_json.go`

## Solution

1. Keep the Phase 3 deployed-binary handoff architecture.
2. Mirror every cleanup-inner failure/skip/retry branch to the durable handoff log.
3. Emit matching `--debug-windows-json` events so detached-child failures survive even when console stderr is swallowed.
4. Add targeted automated tests for the log-line formatter and the child-env forwarding path so the observability contract does not regress silently.

## Expected result after fix

When cleanup fails again, the user gets at least one durable artifact that names the exact failing branch, for example:

- `phase=cleanup event=remove_retry path=... attempt=3 err="Access is denied"`
- `phase=cleanup event=remove_fail path=... err="Access is denied"`
- `phase=cleanup event=swap_remove_fail path=... err="The process cannot access the file..."`
- `phase=cleanup event=drive_root_skip path=... reason=size_guard bytes=7340032`

So even if the hidden child console output disappears, the failure remains forensically recoverable.

## Validation

- `go test ./cmd -run 'TestFormatHandoffLogLine|TestBuildCleanupChildEnv|TestBuildCleanupChildArgs'`
- Confirm `GITMAP_UPDATE_CLEANUP_DELAY_MS`, `GITMAP_DEBUG_WINDOWS`, and `GITMAP_DEBUG_WINDOWS_JSON` are forwarded into the cleanup child env.
- Native Windows manual check with `gitmap-v27 update --debug-windows --debug-windows-json` should now produce both console diagnostics and durable per-branch child failure events.
