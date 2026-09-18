# Subtask 04: Multi-Project Parallel Pipeline Errors Scanning & Batching

## Objective
Implement multi-project parallel scanning and batching for pipeline error fixes with a default limit of 3 projects, persistent offset tracking for subsequent batches, and AGY task enqueueing.

## Target Files
- `cli/cmdagy/agy_fix_multi_project.go` (new)
- `cli/cmdagy/agy_fix_pipeline.go`
- `cli/cmdpipeline/pipeline_multi_errors.go` (new)

## Status
Completed: 2026-09-18

## Verification
- Implemented concurrent scanning: `ScanFailingProjectCandidates` in `cli/cmdagy/agy_fix_pipeline_batch_scan.go`.
- Implemented cursor persistence: `LoadPipelineFixBatchCursor`, `SavePipelineFixBatchCursor`, `ResetPipelineFixBatchCursor` in `cli/cmdagy/agy_fix_pipeline_batch_cursor.go`.
- Implemented batch pagination: `RunAgyFixMultiProjectBatch` in `cli/cmdagy/agy_fix_pipeline_batch_exec.go` with default limit of 3 projects per run.
- Subsequent runs automatically advance the cursor to process the next batch of projects.
- Added `--reset-batch` flag to restart from project 1.
- All files $\le 92$ lines, functions $\le 15$ lines.
- `go vet ./cmdagy/...`: PASS (code 0).
