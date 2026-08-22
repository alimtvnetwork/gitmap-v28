# CI/CD Issue 11: `cmdFaithfulExiter` Mutex Isolation and `maybeExitOnCmdFaithfulMismatch` Race

- **Stage**: CI `go test -race` / Test Matrix
- **Status**: ✅ Resolved
- **Affected Packages**: `github.com/alimtvnetwork/gitmap-v28/gitmap/cmd`

## 1. Problem Description

`go test -race ./cmd/...` exited abruptly with code 3 during test suite execution without formatting test failure results. In CI output, the divergence warning appeared twice:
```
[FAIL] verify-cmd-faithful: scripts-fixer — 1 divergence(s); displayed cmd: did not match executor argv (non-fatal unless --verify-cmd-faithful-exit-on-mismatch is set)
verify-cmd-faithful: FAIL - one or more cmd: lines did not match the executor's argv (see per-row reports above); exiting with code 3 because --verify-cmd-faithful-exit-on-mismatch was set
```
Followed by `FAIL github.com/alimtvnetwork/gitmap-v28/gitmap/cmd 20.869s`.

## 2. Root Cause Analysis

- `clonetermverifyexit_test.go` tests `--verify-cmd-faithful-exit-on-mismatch` by calling `setCmdFaithfulExitOnMismatch(true)`.
- While the unit test was running, concurrent clone tests executed `maybeExitOnCmdFaithfulMismatch()`.
- If `withRecordingExiter` completed and its `t.Cleanup` restored `cmdFaithfulExiter` to `os.Exit` while a concurrent test was checking `cmdFaithfulExitOnMismatchEnabled()`, the concurrent test called `os.Exit(3)`, immediately aborting the test runner process.

## 3. Solution

1. **Synchronize `maybeExitOnCmdFaithfulMismatch` under `cmdFaithfulExiterMu.RLock`**:
   - Updated `maybeExitOnCmdFaithfulMismatch()` in `clonetermverifyexit.go` to hold `cmdFaithfulExiterMu.RLock()` throughout its entire check and exit invocation.
2. **Hold Exclusive `cmdFaithfulExiterMu.Lock()` During Exit-Check Tests**:
   - Updated `withRecordingExiter` in `clonetermverifyexit_test.go` to hold `cmdFaithfulExiterMu.Lock()` across the entire test lifecycle until `t.Cleanup` finishes resetting state and restoring the original handler.
3. **Capture Stdout in `TestCommitCodingGuidelinesNoCommitNoPushPrintsBothNotes`**:
   - Injected `Stdout: &stdout` in `codingguidelines_test.go` to prevent status messages from polluting `os.Stdout`.

## 4. What NOT to Repeat

- **Never mutate package-level process exit hooks or exit-triggering flags without holding an exclusive lock**: Any test that swaps `os.Exit` for a test recorder MUST lock out all concurrent tests that could execute the exit condition.
