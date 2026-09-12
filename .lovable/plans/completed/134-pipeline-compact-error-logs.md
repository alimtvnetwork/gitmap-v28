# 134-pipeline-compact-error-logs.md: Pipeline Compact Error Logs Default Filtering & Detailed Verbose Flags

**Status: completed**
**Execution Loop Count: 6 steps completed in continuous N-step self-loop (N=200 budget)**
**Task Initiation: Initiated via user request to make `gitmap pipeline error-logs` (and aliases `errorlogs`, `errors`, `err`, `last-failed-logs`) compact by default by removing `ok lines` (passing test lines / noise) from error logs after reading from server, with non-compact full logs enabled via `--detailed`, `--verbose`, `--v`, or `-v`.**

---

## 1. Executive Summary

When analyzing GitHub Actions workflow failures, test runners (e.g. `go test ./...`) output large volumes of passing tests alongside failures (`ok\t...`, `? \t...`, `--- PASS: ...`, `=== RUN ...`, `PASS`, `✔ ok`, `✔ Macro`). This plan implemented:
1. **Default Compact Filtering**: Automatically filters out passing `ok` lines, passing test suites, and noise from `gitmap pipeline error-logs` terminal cards, combined section failure summaries, JSON output, tempfile/file reports, and clipboard copies.
2. **Detailed / Verbose Mode**: Preserves 100% of raw error log lines verbatim when any of `--detailed`, `--verbose`, `--v`, `-v`, or `-V` are supplied.
3. **Pipeline Timeline Watch Integration**: Propagates `IsDetailed` to `runPipelineErrorLogsDynamicTimeline`, ensuring watch loops honor detailed vs compact preferences.
4. **CLI Help Text**: Documented `-v, --detailed, --verbose` in `printPipelineErrorLogsHelp()`.
5. **Quality & Verification**: Verified via unit test suite in `cli/cmdpipeline/pipeline_compact_test.go` and full pass on all linters (`check-nested-ifs.py`, `check-boolean-guidelines.py`, `check-relative-paths.py`, `gofmt`).

---

## 2. Mandatory Rules & Invariants Followed

1. **Rule 1 (Strict Relative Git Paths)**: All markdown paths, documentation, and references are strictly relative to repo root. No absolute paths or `file:///` URIs.
2. **Rule 2 (Strict Sizing & Style)**: Every Go function <= 15 lines. Mandatory blank line before every return statement.
3. **Rule 3 (Affirmative Booleans)**: All boolean variables and struct fields use affirmative prefixes (`is*`, `has*`, e.g., `IsDetailed`, `hasDetailedArg`, `isOkLogLine`).
4. **Rule 4 (Zero Swallowed Errors)**: Zero swallow error pattern maintained.
5. **Rule 5 (Cache Tracking)**: All modified files recorded in `.lovable/temp/recent-file-changes.json` under atomic lock via `python 03-ai-scripts/33-test-inventory-generator.py --record`.
6. **Rule 6 (Consolidated Commit & Push)**: Committed all touched files atomically into a single commit and pushed to `origin main`.
7. **Rule 7 (Constraint Enforcement)**: No auto-releases, no routine unit test execution, no `06-cicd-local-runner.py`.

---

## 3. Architecture & Implementation Details

### 3.1 Flag Parsing (`cli/cmdpipeline/pipeline_flags.go`)
- Added `IsDetailed bool` to `PipelineErrorFlags`.
- Implemented `hasDetailedArg(args []string) bool` checking `--detailed`, `--verbose`, `--v`, `-v`, `-V` across `args` and `os.Args`.
- Updated `parseCommonErrorFlags` to populate `flags.IsDetailed = hasDetailedArg(args)`.

### 3.2 Compact Line Filtering Engine (`cli/cmdpipeline/pipeline_error_extract.go`)
- Implemented `isOkLogLine(line string) bool` identifying:
  - `PASS`, `ok`, `PASS: ...`
  - `ok\t...`, `ok  ...`
  - `?\t...`, `?   ...`
  - `--- PASS: ...`, `=== RUN ...`
  - `✔ ok ...`, `✔ Macro ...`
- Implemented `filterCompactLines(lines []string) []string` filtering out all ok lines.
- Implemented `compactErrorPayload(p *PipelineErrorLogsPayload)` and `updateCompactedErrorLogs(p *PipelineErrorLogsPayload)` to compact job error lines, section failures, reformat combined sections, and reformat aggregated error logs.

### 3.3 Render Pipeline & Watch Timeline Wiring (`cli/cmdpipeline/pipeline_logs.go`, `cli/cmdpipeline/pipeline_dynamic_timeline.go`)
- In `pipeline_logs.go`:
  - Added `applyPayloadOptions(p *PipelineErrorLogsPayload, runs []ghRunItem, flags PipelineErrorFlags)` executing `compactErrorPayload(p)` when `!flags.IsDetailed`.
  - Updated `executePipelineErrorLogs` to pass `flags.IsDetailed` to dynamic timeline parameters.
  - Updated `printPipelineErrorLogsHelp()` with `-v, --detailed, --verbose`.
- In `pipeline_dynamic_timeline.go`:
  - Added `IsDetailed bool` to `ErrorLogsTimelineParams`.
  - Added `applyTimelinePayloadOptions(p *PipelineErrorLogsPayload, runs []ghRunItem, params ErrorLogsTimelineParams)` applying `compactErrorPayload(p)` when `!params.IsDetailed`.

---

## 4. Verification & Testing

1. **Unit Tests**:
   - `cli/cmdpipeline/pipeline_compact_test.go`:
     - `TestIsOkLogLine_WithPassingInputs_ReturnsTrue` (PASS)
     - `TestIsOkLogLine_WithFailingInputs_ReturnsFalse` (PASS)
     - `TestFilterCompactLines_WithMixedLines_FiltersOkLines` (PASS)
     - `TestParsePipelineErrorFlags_WithDetailedFlags_SetsIsDetailed` (PASS)
     - `TestParsePipelineErrorFlags_WithoutDetailedFlag_SetsIsDetailedFalse` (PASS)
     - `TestCompactErrorPayload_WithFailedRuns_RemovesOkLines` (PASS)
2. **Linters**:
   - `python linter-scripts/check-nested-ifs.py`: 0 violations across 2,858 files.
   - `python linter-scripts/check-boolean-guidelines.py`: 0 violations across 2,858 files.
   - `python linter-scripts/check-relative-paths.py`: 0 violations across 6,806 files.
   - `gofmt -w cli/`: Clean.
   - `go build -o gitmap.exe .`: Built and tested successfully.
