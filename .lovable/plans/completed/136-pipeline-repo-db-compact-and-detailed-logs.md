# Plan 136: Pipeline Repo-DB Compact & Detailed Error Logs Architecture (Consolidated Milestone)

> **Originating Request:** The user requested a compact error version by default that eliminates passing `ok lines` from server error logs, while preserving detailed full logs with `--detailed` / `--verbose` / `--v` / `-v`. Furthermore, all logs must be isolated into the repository-scoped SQLite DB (`data/pipeline_db/pipeline_<repoSlug>.db`), strictly never in the main DB (`gitmap.db`), and structured across 3 connected tables with a master run record table and two sub-tables for detailed error logs and compact error logs.
> **Execution Budget & Loop Trace:** Completed across 2 phases and 3 subtasks within a self-looping budget of N=200 steps.

---

## 1. Executive Summary of Implementation

1. **Repository DB Isolation**:
   - Error logs and raw diagnostic bodies are stored exclusively within the repository-scoped split database (`data/pipeline_db/pipeline_<repoSlug>.db`).
   - The main database (`store.DB` / `gitmap.db`) is strictly excluded from storing error log contents.
2. **3 Connected Tables Architecture in Repository DB**:
   - **Master Table (`PipelineRun`)**:
     - Stores all run-level execution metadata: `RunId` (PRIMARY KEY/UNIQUE), `RepoSlug`, `WorkflowName`, `Status`, `Conclusion`, `Branch`, `Sha`, `EtaSeconds`, `DurationSeconds`, `RunUrl`, `IsSuccess`, `Notes`, `Comments`, `CreatedAt`, `UpdatedAt`.
   - **Detailed Sub-Table (`PipelineDetailErrorLog`)**:
     - Linked to master table via `RunId REFERENCES PipelineRun(RunId) ON DELETE CASCADE`.
     - Stores the full, uncompressed logs fetched from the server, preserving all diagnostic lines verbatim.
   - **Compact Sub-Table (`PipelineCompactErrorLog`)**:
     - Linked to master table via `RunId REFERENCES PipelineRun(RunId) ON DELETE CASCADE`.
     - Stores the compact error logs with all passing `ok lines` stripped (`ok\t`, `ok `, `?\t`, `? `, `--- PASS:`, `=== RUN`, `PASS`, `✔ ok`, `✔ Macro`).
     - Tracks `FilteredOkCount` indicating how many noise lines were eliminated.
   - Retained `PipelineErrorLog` for backwards compatibility with legacy queries.
3. **Server Reading & Dual-Table Recording**:
   - In `cli/cmdpipeline/pipeline_recorder.go`, `pipeline_sync_cache.go`, and `pipeline_persist.go`:
   - When reading failed logs from the server (GitHub Actions):
     - Calculates compact logs and filtered line count via `FilterCompactLogText`.
     - Ensures master record exists in `PipelineRun`.
     - Records uncompressed log into `PipelineDetailErrorLog`.
     - Records filtered compact log into `PipelineCompactErrorLog`.
4. **CLI Command Selection & Flag Matrix**:
   - Default: `gitmap pipeline error-logs`, `errors`, `errorlogs`, `err`, `last-failed-logs` returns compact version (no `ok lines`).
   - Flag-enabled: `gitmap pipeline error-logs --detailed` (or `--verbose`, `--v`, `-v`, `-V`) returns detailed full version (all lines preserved).
   - Historical cached failures (`RenderLastCachedFailures`) query `PipelineCompactErrorLog` by default, and `PipelineDetailErrorLog` when `--detailed` is specified.

---

## 2. Granular Subtask Trace & Resolution

### Subtask 01: PipelineDB 3 Connected Tables Schema & Operations
- Updated `cli/pipelinedb/pipeline_split_schema.go`:
  - Added `sqlCreatePipelineDetailErrorLog` with foreign key referencing `PipelineRun(RunId)`.
  - Added `sqlCreatePipelineCompactErrorLog` with foreign key referencing `PipelineRun(RunId)`.
- Created `cli/pipelinedb/pipeline_compact_record.go`:
  - Defined `PipelineCompactErrorRecord` entity.
- Updated `cli/pipelinedb/pipeline_split_ops.go`:
  - Implemented `RecordDetailErrorLog`, `RecordCompactErrorLog`, `RecordDualErrorLog`.
  - Implemented `HasDetailErrorLog`, `HasCompactErrorLog`.
  - Implemented `QueryDetailedErrors`, `QueryCompactErrors`.
  - Implemented `QueryDetailedErrorLogsByRunId`, `QueryCompactErrorLogsByRunId`.
  - Implemented `GetRunWorkflowName`.
  - Updated `Clear()` and `Reset()` to truncate/drop all 3 tables with foreign key safety.
- Created unit tests in `cli/pipelinedb/pipeline_split_db_test.go`:
  - Verified 3-table persistence, foreign key relationships, detail/compact queries, and clearing.

### Subtask 02: Server Fetch Dual-Table Storage & Main DB Error Log Ban
- Updated `cli/cmdpipeline/pipeline_recorder.go`:
  - Updated `persistSingleFailedRunLog` to record both detailed and compact records via `pipeDb.RecordDualErrorLog`.
  - Verified `recordInMasterDB` never stores error log bodies in the main DB.
- Updated `cli/cmdpipeline/pipeline_sync_cache.go`:
  - Updated `saveParsedFailedJobs` and `saveFallbackErrorLog` to persist dual detailed and compact records.
- Updated `cli/cmdpipeline/pipeline_persist.go`:
  - Added `persistLogToRepoSplitDb` to immediately populate `PipelineDetailErrorLog` and `PipelineCompactErrorLog` in the repo split DB whenever logs are fetched from the server and cached.

### Subtask 03: Compact/Detailed Query Wiring, CLI Flags & Verification
- Updated `cli/cmdpipeline/pipeline_error_extract.go`:
  - Added `FilterCompactLogText(rawText string) (string, int)` to strip `ok lines` and compute count.
- Updated `cli/cmdpipeline/pipeline_history.go`:
  - Enhanced `renderSingleCachedFailureCard` to render compact errors by default and detailed errors when `isDetailed` is true.
  - Updated `RenderLastCachedFailures` to accept optional `isDetailed ...bool`.
  - Updated `HandlePipelineHistoryErrors` to pass `hasDetailedArg(args)` to `RenderLastCachedFailures`.
- Added unit tests in `cli/cmdpipeline/pipeline_compact_test.go`:
  - Verified `TestFilterCompactLogText` strips passing ok lines.
  - Verified `TestDualTableRepoSplitDbStorage` creates master `PipelineRun`, detailed `PipelineDetailErrorLog`, and compact `PipelineCompactErrorLog`.
- Verified quality gates across repository linters:
  - `python linter-scripts/check-nested-ifs.py`: 0 violations across 2,859 files.
  - `python linter-scripts/check-boolean-guidelines.py`: 0 violations across 2,859 files.
  - `python linter-scripts/check-relative-paths.py`: 0 violations across 6,813 files.
  - `gofmt -w cli/`: Clean.
  - `go vet ./...`: 0 errors.

---

## 3. Files Modified & Created

| File | Status | Description |
|---|---|---|
| `.agents/skills/pipeline-compact-error-logs/SKILL.md` | Modified | Updated native skill with 3-table repo DB schema and CLI flag matrix |
| `cli/pipelinedb/pipeline_compact_record.go` | New | Defined `PipelineCompactErrorRecord` struct |
| `cli/pipelinedb/pipeline_split_schema.go` | Modified | Added `PipelineDetailErrorLog` and `PipelineCompactErrorLog` tables |
| `cli/pipelinedb/pipeline_split_db.go` | Modified | Updated `InitSchema` with detailed and compact tables |
| `cli/pipelinedb/pipeline_split_ops.go` | Modified | Added CRUD and query operations for detailed and compact logs |
| `cli/pipelinedb/pipeline_split_db_test.go` | Modified | Added 3 connected tables unit tests |
| `cli/cmdpipeline/pipeline_error_extract.go` | Modified | Added `FilterCompactLogText` helper |
| `cli/cmdpipeline/pipeline_recorder.go` | Modified | Updated failed run persistence to write dual detailed/compact logs |
| `cli/cmdpipeline/pipeline_sync_cache.go` | Modified | Updated sync caching to write dual detailed/compact logs |
| `cli/cmdpipeline/pipeline_persist.go` | Modified | Added `persistLogToRepoSplitDb` on server log fetching |
| `cli/cmdpipeline/pipeline_history.go` | Modified | Updated history cached failure card rendering for detail vs compact mode |
| `cli/cmdpipeline/pipeline_compact_test.go` | Modified | Added unit tests for filter and dual-table storage |
| `.lovable/plans/01-index.md` | Modified | Updated plans index with completed milestone 136 |
