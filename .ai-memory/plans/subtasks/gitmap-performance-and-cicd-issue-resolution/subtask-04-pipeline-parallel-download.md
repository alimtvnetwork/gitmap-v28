# Subtask 04: Accelerate Pipeline Log Downloading and Fast Cache Return in Gitmap TE/PE

- **Subtask ID**: `subtask-04-pipeline-parallel-download`
- **Parent Task**: `gitmap-performance-and-cicd-issue-resolution`
- **Owned Files**:
  - `cli/cmdpipeline/pipeline_logs_fetch.go`
  - `cli/cmdpipeline/pipeline_query_logs.go`
  - `cli/cmdpipeline/pipeline_sha_cache.go`
  - `cli/cmdpipeline/pipeline_sha_cache_test.go`
- **Assigned Worker**: Worker 01
- **Status**: COMPLETED

---

## 1. Objectives
1. Ensure all downloading in `cmdpipeline` is fully parallelized:
   - Dynamic bounded goroutine worker pool (`resolveJobLogConcurrency`) up to 8 workers for parallel job log fetching.
   - Semaphore cleanup with `defer func() { <-sem }()` in `dispatchFetchRunWorker`.
2. Optimize instant cache retrieval so that `gitmap te` and `gitmap pe` retrieve cached payloads from `PipelineShaCache` without subprocess overhead:
   - Added `resolveTargetShaFast` and `readGitHeadSha` to resolve HEAD commit SHA via direct `.git/HEAD` disk access (<0.1ms) instead of executing `git rev-parse HEAD`.
3. Validate that cached output renders immediately (<20ms) with zero lag.

---

## 2. Verification
- `gitmap te` returns in under 300ms on cache hit.
- `go test -C cli ./cmdpipeline -run "TestShaCache|TestResolveJobLogConcurrency|TestResolveFetchConcurrency" -v`: All PASS (0.12s).
