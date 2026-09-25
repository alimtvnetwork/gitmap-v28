# Plan 162: Pipeline Stage Timings & Antigravity Queue Inspector

## Metadata
- **Spec:** [162-pipeline-stage-timings-and-agy-queue-inspector.md](../../../02-spec/21-app/162-pipeline-stage-timings-and-agy-queue-inspector.md)
- **Status:** Completed
- **Completed Date:** 2026-09-25
- **Total Subtasks:** 5

## Context & Objectives
1. Record each pipeline job/stage in repository SQLite split DB (`data/pipeline/<repo-slug>/sql.db`) including job name, start/end timestamps, single stage duration, and combined aggregate stage timings.
2. Provide formal Mermaid ER diagrams and specification documentation in `02-spec/` (`02-spec/09-pipeline/` and `02-spec/21-app/10-pipeline-and-repo-split-db/`).
3. Add Antigravity IDE inspection capabilities:
   - `gitmap agy active`: Inspect running conversations (`not_fully_idle != 0`).
   - `gitmap agy queues`: Inspect pending prompt queues across workspace directories.
   - Header indicators in `gitmap agy status`.

## Subtasks & Outcomes

### Subtask 162-01: Pipeline Jobs and Segments Schema & DB Ops
- **Status:** Completed
- **Changes:**
  - Added `PipelineJob` table and index in `cli/pipelinedb/pipeline_split_schema_job.go`.
  - Added CRUD and aggregation helpers: `RecordPipelineJobs`, `RecordPipelineSegments`, `QueryPipelineJobs`, `QueryPipelineSegments`, `QueryStageSummary`.
  - Added tests in `pipeline_stage_test.go` and `pipeline_segment_test.go`.

### Subtask 162-02: Pipeline Stage Timings Sync & CLI Display
- **Status:** Completed
- **Changes:**
  - Added `syncRunJobsAndSegments` during pipeline sync in `cli/cmdpipeline/pipeline_segments.go`.
  - Added `pipeline stages [runId]` with aliases `stage`, `jobs`, `job`, `sj`.
  - Computed Combined Stage Sum, Wall-Clock Duration, and Concurrency Speedup Multiplier.
  - Added help text in `cli/cmdpipeline/pipeline_help_menu.go`.

### Subtask 162-03: Pipeline Spec and ERD Diagram
- **Status:** Completed
- **Changes:**
  - Created Mermaid ERD in `02-spec/09-pipeline/pipeline-split-db-erd.mmd`.
  - Updated `02-spec/21-app/10-pipeline-and-repo-split-db/01-pipeline-split-db-architecture.md` with full schema and entity relationships.
  - Created specification in `02-spec/21-app/162-pipeline-stage-timings-and-agy-queue-inspector.md`.

### Subtask 162-04: Antigravity Active Running & Queued Prompts Inspector
- **Status:** Completed
- **Changes:**
  - Implemented `gitmap agy active` querying `conversation_summaries.db` (`(killed IS NULL OR killed = 0) AND not_fully_idle != 0`).
  - Implemented `gitmap agy queues` scanning `.ai-memory/temp/agy-prompt-queue.json` across workspaces.
  - Integrated active and queue counts in `gitmap agy status` header.
  - Added `--json` support and unit tests in `cli/cmdagy/agy_active_queues_test.go`.

### Subtask 162-05: Verification and CI Gates
- **Status:** Completed
- **Changes:**
  - Decomposed all modified/new files to strictly $\le 100$ lines.
  - Verified 100% test pass in `pipelinedb`, `cmdpipeline`, and `cmdagy`.
  - Zero nested ifs, affirmative booleans, relative markdown links.
