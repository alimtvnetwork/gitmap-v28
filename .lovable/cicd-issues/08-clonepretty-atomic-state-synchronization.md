# CI/CD Issue 08: `clonepretty` Global Flag Atomic Synchronization

- **Stage**: CI `go test -race` / Test Matrix
- **Status**: ✅ Resolved
- **Affected Packages**: `github.com/alimtvnetwork/gitmap-v28/gitmap/cmd`

## 1. Problem Description

`clonepretty.go` and `clonespinner.go` declared bare package-level booleans (`cloneDryRunFlag`, `cloneSpinnerOff`, `isCloneAssumeYes`). During concurrent clone test executions under `go test -race`, concurrent reads in `runCloneCommandPretty`, `newCloneCommand`, and `startCloneSpinner` alongside mutating test helpers (`SetCloneDryRun`, `SetCloneSpinnerOff`, `SetCloneAssumeYes`) triggered race detector warnings.

## 2. Root Cause Analysis

- `cloneDryRunFlag`, `cloneSpinnerOff`, and `isCloneAssumeYes` were declared as plain non-synchronized `bool` variables.
- Parallel tests mutated and checked these toggles simultaneously.

## 3. Solution

1. **Convert to `atomic.Bool`**:
   - Replaced `bool` with `atomic.Bool` for `cloneDryRunFlag`, `cloneSpinnerOff`, and `isCloneAssumeYes`.
   - Updated `SetCloneDryRun(on bool)` -> `cloneDryRunFlag.Store(on)`, `IsCloneDryRun() bool` -> `cloneDryRunFlag.Load()`, `SetCloneAssumeYes(on bool)` -> `isCloneAssumeYes.Store(on)`, and `SetCloneSpinnerOff(off bool)` -> `cloneSpinnerOff.Store(off)`.
2. **Updated Call Sites**:
   - Used `.Load()` in `runCloneCommandPretty`, `newCloneCommand`, and `startCloneSpinner`.

## 4. What NOT to Repeat

- **Never use bare boolean variables for CLI behavior toggles and hooks across test suites**: Always use `sync/atomic.Bool` for scalar flags or `sync.RWMutex` for complex state so that `go test -race` runs without races.
