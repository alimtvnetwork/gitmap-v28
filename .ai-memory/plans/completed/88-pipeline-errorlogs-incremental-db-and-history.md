# Consolidated Milestone: Pipeline Error Logs, Incremental DB Caching & Historical Analysis

> **Task Slug:** `88-pipeline-errorlogs-incremental-db-and-history`
> **Status:** Completed
> **Prompt Version:** 2.1.0
> **Started:** In response to incomplete error section reporting in `gitmap pipeline error-logs` and user requirement for positional `-N` inspection, incremental SQLite caching, `--last-failures`, and 5-run history table.
> **Steps / Loops to Complete:** Completed in 2 phases (Phase 1: 2-Agent Planning & Spec Generation; Phase 2: 2-Agent Parallel Execution, Local CI Quality Gates across 38 gates).
> **Quality Gate Verification:** 38/38 gates passed cleanly with exit code 0.

---

## Executive Summary & Implemented Features

1. **Universal Command Aliases**:
   - Registered `gitmap pipelines` as root alias for `gitmap pipeline`.
   - Registered subcommands: `errorlogs`, `errorslogs`, `error-log`, `errors-log`, `errors-logs`, `errors`, `last-failed-logs`.
2. **Infallible Error Section Correlation**:
   - Correlated raw failed logs with `gh run view <runId> --json jobs` to guarantee 0 failing steps or jobs are omitted.
   - Formatted and combined all failed section error logs into `SectionFailures` and `CombinedErrors`.
3. **Positional Negative Indexing (`-N`)**:
   - `gitmap pipeline errors -2`: Inspects the 2nd latest pipeline run.
   - If failed: shows run metadata header and combined error logs.
   - If passing: shows clean status, duration, historical ETA, and commit info.
4. **Incremental SQLite Database Caching**:
   - `gitmap pipelines last-failed-logs`: Scans up to 20 runs, queries SQLite for cached `RunId`s, and downloads logs ONLY for new/missing failed runs (zero re-download guarantee).
   - Added SQLite composite indexes for `RepoSlug` and `RunId`.
5. **Database Identification**:
   - Explicitly displays relative pipeline database path (`./data/pipeline_db/pipeline_<slug>.db`) in terminal output.
6. **Recent Pipeline Failure Summaries**:
   - Compact 5-run history summary table rendered in `gitmap pipeline errors` / `errorlogs`.
7. **`--last-failures <N>` Flag**:
   - Displays cached failure logs from SQLite for the last N failed runs.

---

## Granular Subtasks Archive

### 01-task-routing-and-flags.md

# Subtask 01: Universal Command Routing & Flag Parsing
> Target Files: gitmap/cmd/rootutility.go, gitmap/cmd/pipeline.go, gitmap/cmd/pipeline_flags.go
1. Register gitmap pipelines as alias for gitmap pipeline in rootutility.go.
2. Register errorslogs, errors-log, errors-logs, errors, last-failed-logs in rootutility.go and pipeline.go.
3. Parse negative positional index (-N) and --last-failures <N> flags in pipeline_flags.go.

---

### 02-task-infallible-error-extract.md

# Subtask 02: Infallible Error Extraction & Section Correlation
> Target Files: gitmap/cmd/pipeline_error_extract.go, gitmap/cmd/pipeline_segments.go
1. Correlate raw failed logs with gh run view <runId> --json jobs so zero failing steps or jobs are omitted.
2. Combine all failed section error logs for the current pipeline run into SectionFailures and CombinedErrors.
3. Format each failed section with Workflow, Job, Step, error summary, and log lines.

---

### 03-task-incremental-cache-db.md

# Subtask 03: Incremental SQLite Database Caching
> Target Files: gitmap/pipelinedb/pipeline_split_ops.go, gitmap/pipelinedb/pipeline_split_schema.go, gitmap/cmd/pipeline_sync_cache.go
1. Add SQLite composite indexes for RepoSlug and RunId.
2. Check SQLite for already cached RunIds before downloading logs via gh run view.
3. Download logs only for new/missing failed runs (incremental sync).
4. Format and display relative database file path.

---

### 04-task-positional-and-history.md

# Subtask 04: Positional Negative Indexing & Failure History
> Target Files: gitmap/cmd/pipeline_logs.go, gitmap/cmd/pipeline_history.go
1. gitmap pipeline errors -2: failing case shows combined error logs; passing case shows status, duration, ETA, and commit metadata.
2. gitmap pipeline errors: shows current failure logs + summary table of last 5 pipeline runs.
3. gitmap pipeline errors --last-failures <N>: displays cached failure logs from SQLite.

---

### 05-task-tests-and-ci-verification.md

# Subtask 05: Unit Tests & Quality Gate CI Verification
> Target Files: gitmap/cmd/pipeline_errorlogs_test.go, gitmap/cmd/pipeline_history_test.go
1. Unit tests for aliases, negative indexing, incremental caching, and history rendering.
2. Run local CI runner: python 03-ai-scripts/06-cicd-local-runner.py ensuring all 38 gates pass.
