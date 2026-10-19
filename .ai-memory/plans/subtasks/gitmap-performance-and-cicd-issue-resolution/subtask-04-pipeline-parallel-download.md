# Subtask 04: Accelerate Pipeline Log Downloading and Fast Cache Return in Gitmap TE/PE

- **Subtask ID**: `subtask-04-pipeline-parallel-download`
- **Parent Task**: `gitmap-performance-and-cicd-issue-resolution`
- **Owned Files**:
  - `cli/cmdpipeline/pipeline_logs_fetch.go`
  - `cli/cmdpipeline/pipeline_query_logs.go`
  - `cli/cmdpipeline/pipeline_sha_cache.go`
- **Assigned Worker**: Worker 01
- **Status**: PENDING

---

## 1. Objectives
1. Ensure all downloading in `cmdpipeline` is fully parallelized:
   - Parallel download across failed jobs using bounded goroutines.
   - Parallel download across multiple failed workflow runs when applicable.
2. Optimize instant cache retrieval so that `gitmap te` and `gitmap pe` retrieve cached payloads from `PipelineShaCache` without any redundant external process or network stalls.
3. Validate that cached output renders immediately with zero lag.

---

## 2. Verification
- `gitmap te` returns in under 300ms on cache hit.
- `go test ./cli/cmdpipeline/... -v`
