# CI/CD RCA: TestResolveFetchConcurrency Assertion Out of Sync with Expanded Concurrency Ceiling

- **Issue ID:** CI-88
- **Date:** 2026-09-27
- **Branch:** `main`
- **Failed Commit:** `6d4f04efd5ea53ae7d5c625f8b351c6fecefed74`
- **Failed Runs:**
  - Cross-Platform Build (#36308712770) — macOS, Ubuntu, Windows (`go test ./...`)
  - CI (#36308712892) — Full Suite Guard (`go test ./...`)
- **Status:** Resolved / Fixed

---

## 1. Reproduction & Error Telemetry

### Exact Failure Log:
```text
● Job: macos-latest / go build + test | Step: go test ./... (no cache)
  --- FAIL: TestResolveFetchConcurrency (0.00s)
      pipeline_errorlogs_test.go:350: expected 4, got 8
  FAIL
  FAIL	github.com/alimtvnetwork/gitmap-v28/cli/cmdpipeline	41.312s
```

```text
● Job: ubuntu-latest / go build + test | Step: go test ./... (no cache)
  --- FAIL: TestResolveFetchConcurrency (0.00s)
      pipeline_errorlogs_test.go:351: expected 4, got 8
  FAIL
  FAIL	github.com/alimtvnetwork/gitmap-v28/cli/cmdpipeline	38.125s
```

---

## 2. Root Cause Analysis (4-Part)

### Part 1: Symptom
In GitHub Actions runs #36308712770 and #36308712892, cross-platform build jobs failed during `go test ./...` in package `github.com/alimtvnetwork/gitmap-v28/cli/cmdpipeline` on `TestResolveFetchConcurrency`. The test failed with `expected 4, got 8`.

### Part 2: Root Cause
In commit `6d4f04ef` (Application Spec 175), `resolveFetchConcurrency` in `cli/cmdpipeline/pipeline_logs.go` was upgraded to expand worker pool capacity from 4 to 8 concurrent workers (`if total <= 8 { return total } return 8`) to accelerate `gitmap pe` log retrieval across multi-job pipelines. However, `TestResolveFetchConcurrency` in `cli/cmdpipeline/pipeline_errorlogs_test.go` still asserted the legacy ceiling of 4 (`if resolveFetchConcurrency(10) != 4`), causing `resolveFetchConcurrency(10)` (which now returns 8) to fail the test assertion.

### Part 3: Resolution
Updated `TestResolveFetchConcurrency` in `cli/cmdpipeline/pipeline_errorlogs_test.go` so that `resolveFetchConcurrency(10)` expects the current concurrency ceiling of 8 (`if resolveFetchConcurrency(10) != 8`).

### Part 4: Prevention & Learnings
Whenever concurrency constants or bounds are tuned in core packages (such as `cmdpipeline`, `cloneconcurrency`, or `cmdpull`), all corresponding unit tests asserting boundary limits must be updated and checked in the same commit to prevent CI assertion drift.

---

## 3. Verification & Compliance
- **Modified Test File:** `cli/cmdpipeline/pipeline_errorlogs_test.go`
- **Rule Verification:** Positive booleans strictly maintained, functions <= 8–15 lines, no unmocked OS actions.
