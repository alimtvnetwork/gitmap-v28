# Subtask 06: Always Display Elapsed Execution Time in Gitmap TE/PE Output and Verify All Quality Gates

- **Subtask ID**: `subtask-06-pipeline-timing-output`
- **Parent Task**: `gitmap-performance-and-cicd-issue-resolution`
- **Owned Files**:
  - `cli/cmdpipeline/pipeline_logs_output.go`
  - `cli/cmdpipeline/pipeline_logs_entry.go`
  - `cli/cmdpipeline/pipeline_timing_test.go`
- **Assigned Worker**: Worker 03
- **Status**: PENDING

---

## 1. Objectives
1. Record start timestamp `time.Now()` at the initiation of `handlePipelineErrorLogs` / `executePipelineErrorLogs`.
2. Compute elapsed execution duration and print it at the very conclusion of all terminal outputs:
   - For cache hits:
     `⚡ [gitmap te] Request completed in 45ms (from cache)`
   - For live queries:
     `⏱️ [gitmap te] Request completed in 2.15s`
3. Include elapsed duration in JSON output payload as `elapsedDuration` and `elapsedSeconds`.
4. Ensure the timing is ALWAYS displayed, including when logs are copied to clipboard or rendered.
5. Add unit tests in `pipeline_timing_test.go` verifying the timing message is produced.

---

## 2. Verification
- `gitmap te` outputs elapsed duration banner.
- `go test ./cli/cmdpipeline/... -v`
