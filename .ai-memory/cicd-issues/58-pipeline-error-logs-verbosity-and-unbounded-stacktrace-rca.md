# RCA: Pipeline Error Logs Verbosity, Unbounded Stack Trace Slicing, and Runner Cleanup Capture

## 1. Symptom

When executing `gitmap pipeline error-logs` on a failed workflow run (e.g. CI run `#35354190330`), the generated `GITMAP PIPELINE ERROR REPORT` exploded into an 8.1 MB (8,106,055 bytes) wall of repetitive, highly verbose text. Specifically:

1. **Massive Irrelevant Stack Trace**:
   Section [1/8] (`Lint Script Unit Tests` | `Run lint-script unit tests`, which had a 4-line Python assertion failure) displayed an enormous 1,677-line "Stack Trace" containing completely unrelated log output from subsequent jobs, including:
   - Go panic frames from a different job (`Full Suite Guard`).
   - 100+ passing package test results (`ok github.com/... 0.012s`).
   - `go install` dependency downloads (`go: downloading github.com/...`).
   - Runner image provisioning and git safe directory setup.
   - 2,370 lines of terminal progress bar ticker updates (`Checking boolean & enum compliance: [ 8/2370 ] 0.3%`).
   - GitHub Actions post-job git cleanup commands (`Removing includeIf entries`, `git-credentials-...`).

2. **Cross-Job Stack Trace Leakage**:
   The identical 1,677-line stack trace block was copied into every other failed section across all 8 failure cards.

3. **Runner Post-Step Git Cleanup in Error Details**:
   The error details for failed steps included 25 lines of runner git credential cleanup (`[command]/usr/bin/git version`, `Removing includeIf entries`, `includeif.gitdir:`, etc.).

4. **Duplicated Report Output**:
   The clipboard output contained `Combined Pipeline Section Failures` twice back-to-back, multiplying the 2.8 MB sections into an 8.1 MB clipboard payload.

## 2. Root Cause

1. **Unbounded Slicing in `extractStackTraceFromLog`**:
   In `cli/cmdpipeline/pipeline_error_extract.go:345`, finding `"goroutine "` or `"Stack Trace:"` triggered `return strings.TrimSpace(rawLogs[gIdx:])`. This naively sliced from the start of the panic to the very end of the multi-job log file (1,677 lines / ~350 KB), with zero frame boundary checks, zero line caps, and no stopping condition when the panic ended.
2. **Indiscriminate Cross-Job Stack Trace Broadcasting**:
   In `assembleJobItems` (`pipeline_error_extract.go:368-370`) and `buildSectionFailureFromRunJob` (`pipeline_error_extract.go:610-612`), any failed job that did not have its own stack trace was automatically assigned the workflow-level `extractStackTraceFromLog(rawLogs)`. Because `Full Suite Guard` had a Go panic, all unrelated jobs (Python unit tests, boolean linters, relative path checks) inherited `Full Suite Guard`'s 1,677-line block.
3. **Exit Code Lines Re-Arming Error Context Capture**:
   In `pipeline_error_extract.go:45`, `failureMarkers` included `"Process completed with exit code"`. When GitHub Actions emitted `Process completed with exit code 1.`, `recordErrorLine` treated it as an error line and reset `*ctxRem = 25`. Because GitHub Actions immediately runs post-job cleanup after a step fails, the next 25 lines of runner git credential and safe directory manipulation were captured as "error details".
4. **Missing Runner & Tool Noise Filtering**:
   `isIgnoredLogLine` only checked 4 phrases (`post job cleanup`, `safe.directory`, `removing ssh command`, `terminate orphan process`). It permitted runner provisioning, checkout setup, `go: downloading` lines, and terminal progress bar tickers to pass into error lines.
5. **Redundant Report Generation in Clipboard Formatter**:
   `buildClipboardErrorReport` (`cli/cmdpipeline/pipeline_logs.go:613`) appended `p.CombinedErrors` and then appended `p.ErrorLogs`. However, `updateCompactedErrorLogs` had already prepended `p.CombinedErrors` into `p.ErrorLogs`, resulting in duplicate section output.

## 3. Resolution

1. **Bounded, Pattern-Aware Stack Trace Extraction**:
   Created `cli/cmdpipeline/pipeline_stacktrace.go` with `extractBoundedStackLines`, `processStackScanLine`, and `sanitizeStackLine`. Stack traces are strictly bounded:
   - Go panics start at `goroutine \d+ [` or `panic:` and extract only stack frames (`.go:`, indent, function names, `created by`).
   - Extraction terminates immediately upon encountering non-stack lines (`FAIL`, `PASS`, `ok `, `? `, `##[`, `[command]`, `===`, `---`, or exit codes) or reaching the 25-frame cap.
   - Timestamps and runner prefixes (`job\tstep\t`) are cleaned.
2. **Job-Scoped Stack Trace Extraction**:
   Implemented `extractJobStackTrace` in `cli/cmdpipeline/pipeline_stacktrace.go`, filtering `rawLogs` lines to the specific `jobName`. In `assembleJobItems` and `buildSectionFailureFromRunJob`, jobs only receive a stack trace if their own job lines contained a panic. Unrelated jobs retain an empty stack trace.
3. **Terminal Exit Code Context Limiting**:
   Introduced `resolveContextLimit` and `isTerminalStepError`. Lines containing `Process completed with exit code`, `exit status `, or `FAILED (failures=` set `*ctxRem = 0`, halting context capture so post-step runner cleanups are excluded.
4. **Comprehensive Runner & Tool Noise Filtering**:
   Expanded `isIgnoredLogLine` with `isRunnerCleanupNoise`, `isRunnerSetupNoise`, and `isToolProgressNoise` to filter:
   - Git credential, includeIf, submodule, and safe directory runner commands.
   - Hosted runner provisioning, checkout setup, and deprecation warnings.
   - `go: downloading` tool install output and linter progress tickers (`Checking boolean & enum compliance:`).
   - Updated `isKeepLogLine` to filter both `isOkLogLine` and `isIgnoredLogLine`.
5. **De-duplicated Clipboard Report Output**:
   Updated `buildClipboardErrorReport` in `cli/cmdpipeline/pipeline_logs.go` to append `p.ErrorLogs` directly without prepending `p.CombinedErrors` again.
6. **Comprehensive Unit Testing**:
   Added unit tests in `cli/cmdpipeline/pipeline_compact_test.go`:
   - `TestExtractBoundedStackTrace_StopsAtNonStackLine`
   - `TestAssembleJobItems_DoesNotCrossContaminateStackTrace`
   - `TestIsIgnoredLogLine_FiltersRunnerAndCleanupNoise`
   - `TestResolveContextLimit_TerminalStepError`
   - `TestMultiJobLogParsing_BoundedAndNoiseFree`
   - `TestCachedLogFile35354190330_CompactPayloadSize` (verifying actual run `#35354190330` payload is < 50 KB instead of 8.1 MB).

## 4. Prevention & Learnings

- **Never slice log buffers with `[idx:]` without bounds**: Always parse line-by-line with strict start/stop predicates and hard line limits.
- **Never broadcast diagnostic artifacts across jobs**: Keep errors, summaries, and stack traces strictly scoped to the job/step in which they occurred.
- **Treat CI/CD exit codes as terminal boundaries**: Exit code lines signal the end of a step; never use them as context starters for following lines.
- **Keep string formatters non-redundant**: Ensure payload builder fields are not printed multiple times in composite reports.
