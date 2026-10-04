# Specification 176: TestResolveFetchConcurrency Alignment with 8-Worker Concurrency Ceiling

## Status
- **Status:** active
- **Domain:** app / cli / pipeline / tests
- **Author:** Lead AI Architect
- **Date:** 2026-09-27

---

## 1. Context & Motivation

During execution of commit `6d4f04ef` (Application Spec 175), `resolveFetchConcurrency` in `cli/cmdpipeline/pipeline_logs.go` was upgraded to expand worker pool capacity from 4 to 8 concurrent workers (`if total <= 8 { return total } return 8`) to accelerate `gitmap pe` log retrieval across multi-job pipelines.

However, `TestResolveFetchConcurrency` in `cli/cmdpipeline/pipeline_errorlogs_test.go` still asserted the legacy ceiling of 4:
```go
if resolveFetchConcurrency(10) != 4 {
    t.Errorf("expected 4, got %d", resolveFetchConcurrency(10))
}
```
Because `resolveFetchConcurrency(10)` now yields 8, the test assertion failed across all cross-platform runners in CI/CD (macOS, Ubuntu, and Windows).

---

## 2. Requirements & Modifications

### 2.1 Unit Test Assertion Update
- File: `cli/cmdpipeline/pipeline_errorlogs_test.go`
- Update `TestResolveFetchConcurrency`:
  - When `total = 0`, expected concurrency is `0`.
  - When `total = 2`, expected concurrency is `2`.
  - When `total = 10`, expected concurrency is `8` (matching the current maximum worker limit).

### 2.2 Version Bump & Packaging
- Bump `cli/constants/constants.go` to `6.355.2`.
- Record changes in `changelog.md`.
- Deploy locally to `%USERPROFILE%\AppData\Local\gitmap-cli\gitmap.exe`.
- Deploy to remote node `w1` (`node-w1`).
- Push tag `v6.355.2` and commit to `main`.
