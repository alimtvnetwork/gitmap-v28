# Subtask 91-05: AUM Search Benchmarks, Documentation & E2E Isolation

Spec Reference: [02-spec/21-app/141-ssh-join-common-scan-and-agy-rop.md](../../../../02-spec/21-app/141-ssh-join-common-scan-and-agy-rop.md)
Parent Plan: [.ai-memory/plans/completed/91-ssh-join-common-scan-and-agy-rop.md](../../completed/91-ssh-join-common-scan-and-agy-rop.md)

## Objective
Benchmark GitMap Go AUM search/find/list vs Python search side-by-side, document the results in `readme.md`, and implement live local VM / E2E test suites with `//go:build e2e` guard so CI/CD does not run or fail on external VM dependencies.

## Functional Requirements
1. **Benchmark Engine Calibration:**
   - In `cli/cmdautomation/benchmark.go`, correct the Python script flags (`--pattern` and `--path`) to accurately measure `03-ai-scripts/12-fast-cached-grep.py`.
   - Measure real Go multi-core parallel file traversal + content match vs Python single-process / multi-threaded traversal.
2. **README Documentation:**
   - Add a dedicated benchmark comparison section in `readme.md` showing Go vs Python execution time, memory usage, and throughput.
3. **E2E Isolation Guard:**
   - Ensure all live VM cluster E2E tests (probing `w1`, `w2`, `w3`) and heavy performance benchmarks have `//go:build e2e` build tags.
   - Normal `go test ./...` in GitHub Actions CI/CD runs fast unit tests without requiring live VMs.

## Files to Create/Modify
- `cli/cmdautomation/benchmark.go` [MODIFY]
- `cli/cmdautomation/search_benchmark_test.go` [MODIFY]
- `cli/cmdssh/ssh_e2e_live_test.go` [NEW]
- `readme.md` [MODIFY]
