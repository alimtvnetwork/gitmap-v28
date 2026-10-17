# Subtask 02: SQLite PipelineShaCache & TTL Leak Remediation

- **Task ID:** subtask-286-02
- **Parent Task:** 286-pipeline-pe-parallel-download-sha-cache-and-speedup
- **Status:** DONE
- **Owner:** Worker 02

## Objectives
1. In `cli/pipelinedb/pipeline_split_schema.go` and `cli/pipelinedb/pipeline_split_schema_exec.go`:
   - [x] Defined `sqlCreatePipelineShaCache` DDL and added to `pipelineSchemaQueries()`.
2. In `cli/pipelinedb/pipeline_split_queries.go`:
   - [x] Implemented `RecordShaCache(commitSha string, runId uint64, repoSlug string, payloadJson string) error`.
   - [x] Implemented `QueryShaCache(commitSha, repoSlug string) (string, error)`.
3. In `cli/cmdpipeline/pipeline_sha_cache.go` (new file, <= 300 lines):
   - [x] Implemented instant SHA cache hit handler `TryRenderFromShaCache(repo string, flags PipelineErrorFlags) (bool, error)`.
   - [x] Implemented `SavePayloadToShaCache(repo string, payload PipelineErrorLogsPayload)`.
4. In `cli/cmdpipeline/pipeline_logs_entry.go`:
   - [x] Called `TryRenderFromShaCache` at the beginning of `executePipelineErrorLogs`.
   - [x] Called `SavePayloadToShaCache` when fresh payload is generated.
5. In `cli/cmdpipeline/pipeline_cache_eval.go`:
   - [x] Fixed stale TTL cache leak: strictly invalidate `within_ttl` if `headSha != ""` and `!matchCommitSha(headSha, latest.Sha)`.
6. In `cli/cmdpipeline/pipeline_logs_repo.go`:
   - [x] Added in-memory TTL caching for `queryReleaseAndPRsConcurrently(repo)` (10-minute cache) to eliminate remote `gh` CLI stalls.
7. Verified line counts of all modified Go files remain strictly <= 300 lines with nesting depth <= 1.

