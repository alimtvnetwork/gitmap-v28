# Plan 175: Pipeline Errors Performance Optimization & CI Failure Fixes

## 1. Overview
Remediated the 22 pipeline failures from GitHub Actions CI run #34923716481 and implemented high-performance parallel fetching, file-first staging, live progress indicators, and `--no-output-log` flag support for `gitmap pipeline errors`.

## 2. Root Cause Analysis (RCA)
1. **Undefined `context` in `cli/cmdssh/sshjoin_add_cmd.go`**:
   - *Root Cause*: `SJAddCmd.RunE` invoked `context.Background()` when `cmd.Context()` was nil, but `"context"` was missing from the imports block.
   - *Fix*: Added `"context"` to the `import (...)` block.
2. **Nested `if` in `cli/cmd/cluster.go:151`**:
   - *Root Cause*: `dispatchInvertedClusterHelp` contained `if isMatched` nested inside `if hasInverted`, exceeding nesting depth limit.
   - *Fix*: Flattened using an affirmative guard clause `if !hasInverted { return nil, false }`.
3. **Generate Drift**:
   - *Root Cause*: New constants added in `cli/constants/constants_cli.go` (`CmdServersClientsAlias`, `CmdClientsAlias`) were not regenerated into `cli/completion/allcommands_generated.go`.
   - *Fix*: Regenerated via `go generate ./...` (`go run ./internal/gencommands`).
4. **Pipeline Errors Latency & Missing Features**:
   - *Root Cause*: Sequential `gh run view` subprocess calls for failed runs; repeated subprocess queries for immutable finished run jobs; missing file-first staging and user progress indicators.
   - *Fix*:
     - Implemented parallel worker pool (`fetchAllFailedRunsParallel`) with bounded concurrency (4 workers) for failed run item fetching.
     - Implemented repository-scoped filesystem directory (`resolvePipelineDirForRepo`) and file-first staging.
     - Added job caching (`readCachedPipelineJobs` / `writeCachedPipelineJobs`) to eliminate repeated subprocess queries for cached runs.
     - Added live progress indicator (`"Reading pipeline logs..."`) on initial fetch.
     - Added `--no-output-log` / `-n` flag via affirmative boolean `HasSuppressOutputLog` to stage logs to disk without dumping to terminal.

## 3. Subtasks Executed
- **Subtask 01**: CI Build & Linter Fixes (`sshjoin_add_cmd.go`, `cluster.go`, `allcommands_generated.go`).
- **Subtask 02**: Pipeline Parallel Logs, File Staging, Live Progress & Performance (`pipeline_flags.go`, `pipeline.go`, `pipeline_logs.go`, `pipeline_persist.go`, `pipeline_segments.go`, `pipeline_query.go`, `pipeline_errorlogs_test.go`).

## 4. Verification
- `check-boolean-guidelines.py`: PASS (0 violations across 2969 files).
- `check-nested-ifs.py`: PASS (0 violations across 2969 files).
- `check-gosec-diff.py`: PASS (0 new findings).
- `go vet`: PASS across all modified packages.
