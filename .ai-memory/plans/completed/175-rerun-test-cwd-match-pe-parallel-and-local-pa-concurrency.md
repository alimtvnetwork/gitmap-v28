# Completed Plan 175: Rerun Test CWD Match, PE Parallel Log Fetching & Local PA Concurrency

- **Spec Reference:** [02-spec/21-app/175-rerun-test-cwd-match-pe-parallel-and-local-pa-concurrency.md](../../../02-spec/21-app/175-rerun-test-cwd-match-pe-parallel-and-local-pa-concurrency.md)
- **CI/CD Issue Reference:** [.ai-memory/cicd-issues/87-rerun-test-cwd-match-pe-parallel-and-local-pa-concurrency-rca.md](../../cicd-issues/87-rerun-test-cwd-match-pe-parallel-and-local-pa-concurrency-rca.md)
- **Originating Trigger:** CI/CD test failure in `TestResolveClosestActiveProject_CwdMatch` on commit `bc56920`, slow `gitmap pe` telemetry, and local machine PA concurrency restoration.
- **Execution Lifecycle:** Completed in 1 continuous orchestration loop.
- **Target Release:** `v6.355.1`

---

## Consolidated Subtasks & Verified Deliverables

### Subtask 01: Rerun Test CWD Match and Sequence Fix
- **Traceability ID:** Task-01
- **Target Files:** `cli/cmdagy/agy_rerun_test.go`
- **Delivered Changes:** Updated `TestResolveClosestActiveProject_CwdMatch` in `agy_rerun_test.go` so target `"1"` tests sequence 1 resolution (`other-1`, seq 1), and target `"."` tests CWD match (`match-1`, seq 2). Verified with guideline autofixer.
- **Status:** Complete & Verified

### Subtask 02: PE Parallel Log and Metadata Fetch
- **Traceability ID:** Task-02
- **Target Files:** `cli/cmdpipeline/pipeline_logs.go`, `cli/cmdpipeline/pipeline_error_extract.go`, `cli/cmdpipeline/pipeline_failure_tree.go`
- **Delivered Changes:**
  - Added `CorrelateFailedJobsWithRunJobs` in `cli/cmdpipeline/pipeline_error_extract.go`.
  - Added `queryReleaseAndPRsConcurrently` in `cli/cmdpipeline/pipeline_logs.go` to fetch release tags and open PR counts in parallel.
  - Added `fetchRunLogsAndJobsConcurrently` in `cli/cmdpipeline/pipeline_logs.go` to fetch raw failed logs and run jobs in parallel for each failed run.
  - Raised `resolveFetchConcurrency` from 4 to 8 workers.
  - Added `fetchWorkflowJobsParallel` in `cli/cmdpipeline/pipeline_failure_tree.go` to fetch failing workflow jobs in parallel for commit failure trees.
- **Status:** Complete & Verified

### Subtask 03: Restore Local PA Concurrency and SSH Isolation
- **Traceability ID:** Task-03
- **Target Files:** `cli/cmdssh/ssh_pull_fleet.go`
- **Delivered Changes:** Removed hardcoded `--parallel 2` from `executeLocalVMPull` in `cli/cmdssh/ssh_pull_fleet.go`. Local VM pull-all operations now run with full hardware concurrency (`runtime.NumCPU()`), while remote SSH commands and delegated sessions strictly retain half-priority throttling (`--parallel 2`).
- **Status:** Complete & Verified
