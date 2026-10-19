# Subtask 05: Implement Parallel Line Processing and Chunked Log Parsing in Gitmap TE/PE

- **Subtask ID**: `subtask-05-pipeline-parallel-lines`
- **Parent Task**: `gitmap-performance-and-cicd-issue-resolution`
- **Owned Files**:
  - `cli/cmdpipeline/pipeline_extract_scan.go`
  - `cli/cmdpipeline/pipeline_extract_parallel.go`
  - `cli/cmdpipeline/pipeline_extract_clean.go`
- **Assigned Worker**: Worker 02
- **Status**: PENDING

---

## 1. Objectives
1. Implement parallel chunked log line processing in `cli/cmdpipeline/pipeline_extract_parallel.go`:
   - For logs with $> 100$ lines, chunk the lines and process them concurrently across worker goroutines.
   - Run noise-filtering, ANSI stripping, and failure marker detection across CPU cores.
   - Ensure thread-safe, deterministic aggregation into structured failed job items.
2. Maintain backward compatibility with existing `ParseFailedLogLines` and `extractCleanErrorLines`.
3. Keep all files strictly $\le 300$ lines with nesting depth $\le 1$.

---

## 2. Verification
- Unit tests comparing parallel line processing against sequential baseline for 100% parity.
- `go test ./cli/cmdpipeline/... -v`
