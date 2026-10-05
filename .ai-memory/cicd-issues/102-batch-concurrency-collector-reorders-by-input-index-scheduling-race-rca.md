# RCA-102: Batch Concurrency Collector Reorders By Input Index Scheduling Race Condition

> **Issue ID:** RCA-102  
> **Workflow:** Cross-Platform Build (`cross-platform.yml`)  
> **Job:** `ubuntu-latest / go build + test` (Run ID: `37279833279`, Job ID: `111665073367`)  
> **Step:** `go test ./... (cached)` (step #8)  
> **Date:** 2026-10-05  
> **Author:** MD ALIM UL KARIM  

---

## Part 1: Symptom

During the `Cross-Platform Build` workflow on GitHub Actions `ubuntu-latest`, step `go test ./... (cached)` failed with exit code 1:

```text
--- FAIL: TestE2E_BatchConcurrency_CollectorReordersByInputIndex (0.00s)
    clonenextbatchconcurrent_e2e_test.go:140: completion order matched input order — randomization stub failed, reorder-by-index assertion would be trivial
FAIL
FAIL	github.com/alimtvnetwork/gitmap-v28/cli/cmdclone	0.256s
```

The test asserted that the worker completion order must differ from the input order (`assertCompletionOrderRandomized`). However, on high-performance multi-core Linux runners, the completion order matched the input order sequentially, causing `assertCompletionOrderRandomized` to fail immediately.

---

## Part 2: Root Cause

In `cli/cmdclone/clonenextbatchconcurrent_e2e_test.go`:
```go
gates := make([]chan struct{}, n)
for i := range gates {
    gates[i] = make(chan struct{})
}
go func() {
    for i := n - 1; i >= 0; i-- {
        close(gates[i])
        runtime.Gosched()
    }
}()
```
The test attempted to reverse the worker completion order by launching an uncoordinated background goroutine that looped through `i := n - 1; i >= 0; i--` closing `gates[i]` with only `runtime.Gosched()`.

Because channel closes and `runtime.Gosched()` execute in microseconds, the background goroutine closed all 20 channels before worker goroutines were scheduled and had dequeued jobs from the `jobs` channel. When the workers began dequeuing jobs (`repo-0`, `repo-1`, ..., `repo-19`), every channel `gates[idx]` was already closed. Thus, every worker ran uninhibited in its natural enqueue order, completing sequentially (`repo-0` through `repo-19`).

---

## Part 3: Surgical Fix

Replaced the uncoordinated background goroutine and `runtime.Gosched()` with deterministic cascading channel coordination:
1. `gates[n-1]` is closed upfront so worker `n-1` (the last input repo) is guaranteed to execute first without waiting.
2. Inside `processOneBatchRepoFn`, each worker `idx` waits on `<-gates[idx]`, locks the mutex, appends to `completionOrder`, unlocks, and then if `idx > 0`, closes `gates[idx-1]`.
3. This creates an unbroken, deterministic chain: Worker 19 finishes and releases Worker 18; Worker 18 finishes and releases Worker 17; ...; Worker 1 finishes and releases Worker 0.
4. Completion order is mathematically guaranteed to be exactly `[repo-(n-1), repo-(n-2), ..., repo-0]`, eliminating all scheduling races, sleep variance, and platform differences.

Target File: `cli/cmdclone/clonenextbatchconcurrent_e2e_test.go`

---

## Part 4: Verification

1. Executed targeted unit test uncached:
   ```bash
   go test -count=1 -v -run TestE2E_BatchConcurrency_CollectorReordersByInputIndex ./cmdclone
   ```
   **Result:** `PASS (0.118s)`.

2. Executed repeated stress run with 10 iterations:
   ```bash
   go test -count=10 -v -run TestE2E_BatchConcurrency_CollectorReordersByInputIndex ./cmdclone
   ```
   **Result:** 10/10 PASS with zero flakes.

3. Verified full `cmdclone` test suite:
   ```bash
   go test -short ./cmdclone
   ```
   **Result:** `ok github.com/alimtvnetwork/gitmap-v28/cli/cmdclone 1.278s`.

4. Verified repository linters:
   - `python .github/scripts/go-format-check.py --check-only` (3953 files clean)
   - `python linter-scripts/check-relative-paths.py` (8704 files clean)
   - `python linter-scripts/check-enum-and-boolean.py` (clean)
