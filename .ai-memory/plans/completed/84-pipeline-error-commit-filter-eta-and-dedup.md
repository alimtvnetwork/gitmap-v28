# Plan 84: Pipeline Error Commit Navigation, Warning Ingestion, De-Duplication & Rolling ETA Engine

> **Execution Summary:**
> - **Origin:** Initiated from user request to enhance GitMap's pipeline error diagnostics (`gitmap pe`, `gitmap pipeline error-logs`, `gitmap pd`) with commit SHA and negative offset navigation, preserve compiler/tool warnings and pre-failure context, eliminate duplicate error lines and compact linker command noise, implement historical ETA moving averages backed by rolling JSON cache, and provide updated documentation.
> - **Workflow:** 3-Phase Parent Task N-Step Loop (Planning, Architecture & Verification).
> - **Outcome:** 100% completed with zero linter errors, clean `go vet`, and comprehensive unit tests.

---

## Consolidated Subtasks & Deliverables

### Subtask 01: CLI Flag Parsing & Commit Target Routing
- **Files Modified:**
  - `cli/cmdpipeline/pipeline_flags.go`
  - `cli/cmdpipeline/pipeline_commit_groups.go`
  - `cli/cmdpipeline/pipeline_history.go`
- **Accomplishments:**
  - Extended `PipelineErrorFlags` to detect and parse target commits and relative offsets (`-1`, `-2`, `-1n`, `-2N`, `HEAD~1`, `~3`).
  - Added `ParseNegativeIndex` supporting multiple git-style negative offset conventions.
  - Implemented `parseCommitTarget`, `evaluateCommitCandidate`, `isSkipTokenForCommit`, and `isCommitHexSha`.
  - Harmonized negative offset parsing in `pipeline_commit_groups.go` and `pipeline_history.go` to delegate to `ParseNegativeIndex`.

### Subtask 02: Compiler Warning Ingestion & Pre-Failure Context Capture
- **Files Modified:**
  - `cli/cmdpipeline/pipeline.go`
  - `cli/cmdpipeline/pipeline_error_extract.go`
- **Accomplishments:**
  - Added `Warnings []string` to `SectionFailure` and `FailedJobItem` domain structs.
  - Updated `scanLogLinesIntoMap` and `processLogLine` to buffer pre-failure compiler warnings (matching `warning:`, `note:`, `##[warning]`) and attach them to failed job items upon encountering failure.
  - Preserved warnings across section aggregation (`buildSectionFailureFromRunJob` and `mergeMatchedFailure`).

### Subtask 03: De-Duplication of Repeated Error Lines & Linker Noise Compacting
- **Files Modified:**
  - `cli/cmdpipeline/pipeline_error_extract.go`
  - `cli/cmdpipeline/pipeline_logs.go`
  - `cli/cmdpipeline/pipeline_details.go`
- **Accomplishments:**
  - Implemented `compactLinkerCommandLine` to detect massive compiler/linker invocations (e.g. `link.exe`, `collect2`) and compact hundreds of repetitive `.lib`, `.rlib`, `/LIBPATH:` arguments into `... [<N> library and object files omitted] ...`, making underlying fatal notes (`CVT1100`, `LNK1123`) immediately visible.
  - Added `filterOutSummaryLine` to suppress duplicate occurrences of `FailureSummary` under `Details:`.
  - Updated `renderSingleSectionFailureRow` to render `renderSectionWarnings` and `renderSectionErrorLinesDedup`, suppressing redundant full stack traces in high-level combined sections.
  - Updated `renderFailedJobSection` to render `renderJobCardWarnings` and `renderJobCardDetails`.
  - Updated `renderSingleDetailsFailure` in `pipeline_details.go` for consistency.

### Subtask 04: Accurate Historical ETA Decision Engine & JSON Rolling Summary Cache
- **Files Modified/Created:**
  - `cli/cmdpipeline/pipeline_eta_cache.go` (new)
  - `cli/cmdpipeline/pipeline_status.go`
  - `cli/cmdpipeline/pipeline_dynamic_timeline.go`
- **Accomplishments:**
  - Created `cli/cmdpipeline/pipeline_eta_cache.go` managing a rolling cache of the last 20 pipeline run durations in `ResolveEtaCacheFilePath(repo)`.
  - Implemented moving average calculation for workflow run durations, tracking success runs and failure runs >= 45s to avoid skewing from instant pre-flight failures.
  - Integrated `GetCachedWorkflowETA` and `UpdatePipelineEtaCache` into `calculateAverageDuration` in `pipeline_status.go` and `fallbackWorkflowDuration` in `pipeline_dynamic_timeline.go`, preventing inaccurate 3-minute hardcoded fallbacks.

### Subtask 05: Documentation, Help Text Parity & Unit Testing
- **Files Modified/Created:**
  - `cli/cmdpipeline/pipeline_help_menu.go`
  - `cli/cmdpipeline/pipeline_logs.go`
  - `cli/helptext/pipeline.md`
  - `cli/cmdpipeline/pipeline_error_commit_filter_test.go` (new)
- **Accomplishments:**
  - Updated `gitmap pe` usage, flags, and help menus with commit SHA and negative offset syntax.
  - Documented commit SHA navigation, negative offsets (`-1`, `-1n`, `HEAD~1`), warning preservation, linker compacting, and rolling `eta_history.json` in `cli/helptext/pipeline.md`.
  - Added comprehensive automated unit tests covering flag parsing, negative offset variants, linker command compacting, summary deduplication, warning extraction, and ETA cache persistence.
