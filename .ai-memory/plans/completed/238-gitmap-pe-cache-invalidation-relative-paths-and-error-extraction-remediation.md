# Master Execution Plan: GitMap PE Cache Invalidation, Relative Paths & Error Extraction Remediation

> **Plan Version:** 1.0.0  
> **Status:** Active  
> **Parent Task:** Task-238  
> **Associated Spec:** `02-spec/21-app/238-gitmap-pe-cache-invalidation-relative-paths-and-error-extraction-remediation/01-architecture-spec.md`  

---

## 1. Goal & Objectives

1. Fix `cli/cmdpipeline/pipeline_cache_eval.go` to invalidate cache when active runs (`in_progress` or `queued`) exist, and enforce target affinity so specific commit requests never serve cached data from other commits.
2. Fix all absolute path leakages across `cli/cmdpipeline/pipeline_logs.go`, `cli/cmdpipeline/pipeline_details.go`, and `cli/cmdpipeline/pipeline_error_extract.go` using `FormatRelativeDbPath`.
3. Support deep negative offset targets in `queryWorkflowRunsForTarget` by resolving commit SHAs from git history.
4. Clean redundant tab prefixes (`job\tstep\t`) in failure summaries.
5. Verify unit tests in `cli/cmdpipeline/`, run linters (`check-nested-ifs.py`, `check-enum-and-boolean.py`, `check-relative-paths.py`), execute minor version bump (`v6.507.0`), commit, and release.

---

## 2. Subtask Breakdown

- **Subtask 01:** `01-pipeline-cache-invalidation-and-target-sha-fix.md`
  - Refactor `pipeline_cache_eval.go` to reject cache hits for active runs and ensure strict target affinity.
  - Refactor `queryWorkflowRunsForTarget` in `pipeline_logs.go` to resolve deep offsets.
  - Update unit tests in `pipeline_cache_eval_test.go`.

- **Subtask 02:** `02-pipeline-relative-paths-and-tab-prefix-hygiene.md`
  - Sanitize log paths to relative format in `pipeline_logs.go`, `pipeline_details.go`, and `pipeline_error_extract.go`.
  - Strip redundant `job\tstep\t` prefixes in section failure formatters.
  - Run linters: `check-nested-ifs.py`, `check-enum-and-boolean.py`, `check-relative-paths.py`.

- **Subtask 03:** `03-cicd-verification-minor-bump-and-release.md`
  - Run unit tests in `cli/cmdpipeline/...`.
  - Bump minor version to `v6.507.0`.
  - Atomic commit via `gitmap cpf`.
  - Release ceremony via `29-release-orchestrator.py`.
  - Verify CI/CD telemetry.

---

## 3. Worker Allocation (A = 2, H = 2)

- **Worker 01:** Owns Subtask 01 (`cli/cmdpipeline/pipeline_cache_eval.go`, `cli/cmdpipeline/pipeline_cache_eval_test.go`).
- **Worker 02:** Owns Subtask 02 (`cli/cmdpipeline/pipeline_logs.go`, `cli/cmdpipeline/pipeline_details.go`, `cli/cmdpipeline/pipeline_error_extract.go`).
- **Lead Orchestrator:** Owns Subtask 03 (testing, minor bump `v6.507.0`, atomic commit, release ceremony, CI monitoring).

---

## 4. Execution Tracking Ledger

| Subtask ID | File | Owner | Status | Evidence |
| :--- | :--- | :--- | :--- | :--- |
| Subtask-01 | `.ai-memory/plans/subtasks/238-gitmap-pe-cache-invalidation-relative-paths-and-error-extraction-remediation/01-pipeline-cache-invalidation-and-target-sha-fix.md` | Worker 01 | COMPLETED | Cache invalidation on active runs + strict target affinity verified |
| Subtask-02 | `.ai-memory/plans/subtasks/238-gitmap-pe-cache-invalidation-relative-paths-and-error-extraction-remediation/02-pipeline-relative-paths-and-tab-prefix-hygiene.md` | Worker 02 | COMPLETED | FormatRelativeDbPath across all outputs + cleanDisplayErrorText verified |
| Subtask-03 | `.ai-memory/plans/subtasks/238-gitmap-pe-cache-invalidation-relative-paths-and-error-extraction-remediation/03-cicd-verification-minor-bump-and-release.md` | Lead | COMPLETED | Full pipeline test suite PASS, minor bump to v6.507.0, release ceremony executed |
