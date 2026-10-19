# Specification: GitMap Performance and CI/CD Issue Resolution

## 1. Overview & Objectives

This specification defines the architectural enhancements and surgical remediations for:
1. **Accelerating `gitmap te` / `gitmap pe` execution speed**:
   - Concurrent/parallel multi-job and multi-run log downloading.
   - Multi-core parallel chunked log line processing and regex filtering for rapid error isolation.
   - Ultra-fast cache display (<20ms) directly from SQLite `PipelineShaCache` without redundant CLI round-trips.
   - Mandatory elapsed duration reporting at the conclusion of every execution: `⚡ [gitmap te] Request completed in 42ms (from cache)` or `⏱️ [gitmap te] Request completed in 1.48s`.
2. **Remediating 100% of CI/CD pipeline failures**:
   - Resolving legacy reference violations (`gitmap-v6` -> `gitmap-v28`) in `02-spec/01-spec-authoring-guide/13-root-readme-conventions.md`. <!-- gitmap-legacy-ref-allow -->
   - Decomposing oversized Go source files (`cli/constants/constants_install.go`, `cli/cmdinstall/install_packages.go`, `cli/cmdinstall/install_packages_test.go`) so that every file strictly adheres to the 300-line cap.

---

## 2. Architecture & Design Principles

### 2.1 Parallel Pipeline Downloading Architecture
- **Job-level log fetching**: Concurrent downloads bounded by worker goroutines with `sync.WaitGroup` to avoid sequential blocking across multiple failing workflow jobs.
- **Run-level log fetching**: Concurrent metadata and log fetching for multi-workflow configurations.
- **Instant Cache Presentation**: Immediate deserialization and terminal rendering from `repodb/pipeline.db` (`PipelineShaCache` table) on identical commit SHAs.

### 2.2 Parallel Line-by-Line Log Processing Architecture
- For log payloads exceeding 100 lines:
  - Partition raw log lines into deterministic slices based on CPU core availability (`runtime.NumCPU()`).
  - Distribute chunk processing across worker goroutines:
    - Pre-filtering blank lines and known noise patterns (`isIgnoredLogLine`).
    - Stripping ANSI escape sequences (`stripAnsiIfPresent`).
    - Flagging lines with failure markers (`hasFailureMarker`).
  - Reduce sequential scanning overhead to sub-millisecond aggregation.

### 2.3 Elapsed Request Timing Invariant
- Every execution of `gitmap te` and `gitmap pe` must track start time (`startTime := time.Now()`).
- At the end of terminal output (and within structured JSON payload when requested), print the total request elapsed time:
  ```text
  ⏱️ [gitmap te] Request completed in 124ms (from cache)
  ```
  or for fresh downloads:
  ```text
  ⏱️ [gitmap te] Request completed in 2.18s
  ```

### 2.4 File Size & Coding Guideline Gates (300-Line Cap)
- `cli/constants/constants_install.go` (755 lines) -> decomposed into:
  - `cli/constants/constants_install.go` (core CLI commands and basic tool names)
  - `cli/constants/constants_install_tools.go` (extended tool names and metadata)
  - `cli/constants/constants_install_pkgs.go` (package manager definitions and script templates)
- `cli/cmdinstall/install_packages.go` (403 lines) -> decomposed into:
  - `cli/cmdinstall/install_packages.go` (dispatch, orchestration, and execution)
  - `cli/cmdinstall/install_packages_resolve.go` (package name resolution and manager mapping)
- `cli/cmdinstall/install_packages_test.go` (431 lines) -> decomposed into:
  - `cli/cmdinstall/install_packages_test.go` (main test suite)
  - `cli/cmdinstall/install_packages_resolve_test.go` (resolver test suite)

---

## 3. Verification & Acceptance Criteria

1. `python .github/scripts/check-legacy-refs.py` exits with code 0 (zero legacy references found).
2. `python .github/scripts/tests/test_ci_scripts.py` passes completely without test failures.
3. `python 03-ai-scripts/52-file-size-check.py --diff HEAD~1 --limit 300` exits with code 0 (all modified Go files $\le 300$ lines).
4. `gitmap te` prints elapsed execution duration at the conclusion of both cached and live executions.
5. All tests in `cli/cmdpipeline` and `cli/cmdinstall` pass cleanly (`go test ./cli/cmdpipeline -v`, `go test ./cli/cmdinstall -v`).
