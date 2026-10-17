# Subtask 01: Parallel Log Download & Fast ANSI/Regex Extraction

- **Task ID:** subtask-286-01
- **Parent Task:** 286-pipeline-pe-parallel-download-sha-cache-and-speedup
- **Status:** DONE
- **Owner:** Worker 01

## Objectives
1. In `cli/cmdpipeline/pipeline_query_logs.go`:
   - Refactor `fetchJobLevelFailedLogs` to fetch failing job logs concurrently using a bounded worker pool (semaphore of size 4).
   - Collect and preserve exact job ordering using pre-allocated slice indexed by job order.
2. In `cli/cmdpipeline/pipeline_extract_clean.go`:
   - In `parseLogLine(raw string)`, implement fast ANSI check `stripAnsiIfPresent(raw)`: check `strings.IndexByte(raw, 0x1b) != -1` before calling `ansiRegex.ReplaceAllString`.
3. In `cli/cmdpipeline/pipeline_extract_scan.go`:
   - Optimize line scanning: check for failure markers and error prefixes efficiently.
4. Verify line counts of all modified Go files remain strictly <= 300 lines with nesting depth <= 1.
