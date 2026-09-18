# CI/CD Issue 11: `cmdFaithfulExiter` Mutex Isolation and `maybeExitOnCmdFaithfulMismatch` Race

- **Stage**: CI `go test -race` / Test Matrix
- **Status**: ✅ Resolved
- **Affected Packages**: `github.com/alimtvnetwork/gitmap-v28/gitmap/cmd`

## 1. Problem Description

`go test -race ./cmd/...` exited with failure / timeout during test suite execution. In CI output, the divergence warning appeared:
```
[FAIL] verify-cmd-faithful: scripts-fixer — 1 divergence(s); displayed cmd: did not match executor argv (non-fatal unless --verify-cmd-faithful-exit-on-mismatch is set)
verify-cmd-faithful: FAIL - one or more cmd: lines did not match the executor's argv (see per-row reports above); exiting with code 3 because --verify-cmd-faithful-exit-on-mismatch was set
```
Followed by `FAIL github.com/alimtvnetwork/gitmap-v28/gitmap/cmd`.

## 2. Root Cause Analysis

1. **Self-Deadlock on `sync.RWMutex`**:
   - `withRecordingExiter` acquired `cmdFaithfulExiterMu.Lock()`, and inside the same goroutine `maybeExitOnCmdFaithfulMismatch()` attempted to acquire `cmdFaithfulExiterMu.RLock()`.
   - In Go, a goroutine holding a write lock (`Lock()`) that calls `RLock()` self-deadlocks on runtime sema.
2. **Concurrency Race with Global Exiter**:
   - `cmdFaithfulExiter` was accessed directly by `maybeExitOnCmdFaithfulMismatch()`.
   - In `withRecordingExiter(t)`, safe getter/setter functions (`getCmdFaithfulExiter` / `setCmdFaithfulExiter`) synchronized under `cmdFaithfulExiterMu` must be used to atomically read and write the exiter function without holding nested locks across calls.

## 3. Solution

1. **Eliminate Nested RLock**:
   - `maybeExitOnCmdFaithfulMismatch()` uses `getCmdFaithfulExiter()` which cleanly acquires and releases `cmdFaithfulExiterMu.RLock()` to copy the function pointer, preventing self-deadlock.
2. **Thread-Safe Test State Management in `withRecordingExiter`**:
   - `withRecordingExiter` uses `setCmdFaithfulExiter(...)` and `resetCmdFaithfulState()`, restoring `setCmdFaithfulExiter(prev)` in `t.Cleanup`.
3. **Capture Stdout in `TestCommitCodingGuidelinesNoCommitNoPushPrintsBothNotes`**:
   - Injected `Stdout: &stdout` in `codingguidelines_test.go` to prevent skip notes from leaking to `os.Stdout`.

## 4. What NOT to Repeat

- **Never call `RLock()` while holding `Lock()` on the same `sync.RWMutex`**: Go's `RWMutex` is not re-entrant; always use granular getters/setters or separate synchronization primitives.
