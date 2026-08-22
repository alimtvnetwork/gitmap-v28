# CI/CD Issue 07: installctx `ctxExplainEnabled` Concurrency Race Guard

- **Stage**: CI `go test -race` / Test Matrix
- **Status**: ✅ Resolved
- **Affected Packages**: `github.com/alimtvnetwork/gitmap-v28/gitmap/cmd`

## 1. Problem Description

During concurrent `go test -race` runs across package `cmd`, unsynchronized reads and writes to the global boolean `ctxExplainEnabled` across `installctx.go`, `installctxexplain.go`, `installctxmac.go`, `installctx_harness_test.go`, and `installctx_linux_e2e_test.go` triggered race detector data race warnings.

## 2. Root Cause Analysis

- `ctxExplainEnabled` was declared as a bare package-level `bool` (`var ctxExplainEnabled bool`).
- Test helpers (`withExplain`, `withHome`) and templating routines read and wrote `ctxExplainEnabled` concurrently without synchronization locks or mutex guards when parallel tests executed.

## 3. Solution

1. **`sync.RWMutex` Protection & Accessors**:
   - In `gitmap/cmd/installctx.go`, protected `ctxExplainEnabled` with `ctxExplainMu sync.RWMutex` via `setCtxExplainEnabled(on bool)` and `isCtxExplainEnabled() bool`.
2. **Standardized Call Sites**:
   - Updated `installctxexplain.go`, `installctxmac.go`, `installctx_harness_test.go`, and `installctx_linux_e2e_test.go` to use `isCtxExplainEnabled()` and `setCtxExplainEnabled()`.

## 4. What NOT to Repeat

- **Never use bare global state for feature toggles in CLI commands**: Always wrap package globals in mutexes (`sync.RWMutex`) or atomic primitives (`atomic.Bool`) with dedicated getters/setters to guarantee race safety under `go test -race`.
