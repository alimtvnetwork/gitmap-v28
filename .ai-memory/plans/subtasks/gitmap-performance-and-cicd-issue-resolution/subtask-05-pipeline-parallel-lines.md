# Subtask 05: Implement Parallel Line Processing and Chunked Log Parsing in Gitmap TE/PE

- **Subtask ID**: `subtask-05-pipeline-parallel-lines`
- **Parent Task**: `gitmap-performance-and-cicd-issue-resolution`
- **Owned Files**:
  - `cli/cmdpipeline/pipeline_extract_scan.go`
  - `cli/cmdpipeline/pipeline_extract_parallel.go`
  - `cli/cmdpipeline/pipeline_extract_parallel_test.go`
  - `cli/cmdpipeline/pipeline_extract_clean.go`
- **Assigned Worker**: Worker 02
- **Status**: COMPLETED

---

## 1. Objectives
1. Implement parallel chunked log line processing in `cli/cmdpipeline/pipeline_extract_parallel.go`:
   - For logs with $> 100$ lines, chunk the lines across worker goroutines bounded by `runtime.NumCPU()`.
   - Concurrently perform ANSI stripping, timestamp removal, annotation parsing, and failure marker detection.
   - Deterministically aggregate into `jobMap` and `order` slices matching sequential state machine logic.
2. Updated `cli/cmdpipeline/pipeline_extract_scan.go`:
   - Auto-dispatch to `scanLogLinesParallel` when `len(lines) > 100 && runtime.NumCPU() > 1`.
   - Preserve 100% backward compatibility and exact deterministic output for `ParseFailedLogLines`.
3. Keep all files strictly $\le 300$ lines with nesting depth $\le 1$.

---

## 2. Verification
- `pipeline_extract_parallel.go` (239 lines, $\le 300$).
- `pipeline_extract_scan.go` (212 lines, $\le 300$).
- `pipeline_extract_parallel_test.go` (138 lines, $\le 300$).
- `go test -C cli ./cmdpipeline -run "TestParseFailedLogLines|TestParallel" -v`: All 11 tests PASS (0.13s).
