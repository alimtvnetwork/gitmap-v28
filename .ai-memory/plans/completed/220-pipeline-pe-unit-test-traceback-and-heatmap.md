# Plan: Pipeline PE Unit Test Traceback Extraction & Heatmap Modernization

> **Plan ID:** 220-pipeline-pe-unit-test-traceback-and-heatmap  
> **Status:** Completed  
> **Created:** 2026-10-05  
> **Completed:** 2026-10-05  
> **Associated Spec:** `02-spec/21-app/220-pipeline-pe-unit-test-traceback-and-heatmap/`  

---

## User Request (Verbatim)

```text
When we do the heat map PE, PE is not showing the error message. That's one problem. I'm giving you the terminal version that is shown and what was in the error. Both of these are very important because when error happens with a unit test, you need to put the unit test information correctly with all the lines that is there. Probably not the run 18 tests, this type of line. But yes, the traceback call, which test is failed, this type of information is very important for the logger and for the AI to fix the error right away. Make sure that you respect that and you update the heat map regarding this heat map PE, and also fix this test, make a minor bump in the version, and make a release.

# Actionable Items Must Follow Non-Negotiable

1. Write spec under 02-spec/21-app/<slug>/ and enqueue plan task in .ai-memory/plans/<slug>.md (subtasks in .ai-memory/plans/subtasks/<slug>/) first
2. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String
3. Update the heat map regarding the heat map PE issue.
4. Fix the unit test to ensure error messages are correctly displayed.
5. Include the traceback call and failed test information for logging and AI error fixing.
6. Make a minor version bump.
7. Make a release.

## Must follow and spawn agent using

@[.agents/skills/execute-parent-task-with-n-steps-v6]
```

---

## Execution Summary

### Subtask 01: Job-Level Log Fetcher & DB Cache Self-Healing
- **Owner:** Worker 01
- **Status:** COMPLETED
- **Files Changed:** `cli/cmdpipeline/pipeline_query.go`, `cli/cmdpipeline/pipeline_persist.go`
- **Accomplishments:**
  - Implemented `fetchFailedJobLogs` calling `gh run view --job <databaseId> --repo <repo> --log-failed` to fetch actual logs for failing jobs when run-level log archives are unavailable or in progress.
  - Implemented `fetchJobLevelFailedLogs` and `normalizeJobLogOutput` to preserve job names and line prefixes.
  - Replaced early return in `persistLogToRepoSplitDb` with non-fallback verification (`!isCorruptOrFallbackErrorLog`), allowing the SQLite Split-DB to update when real logs arrive.

### Subtask 02: Traceback Parser & Heatmap Formatter
- **Owner:** Worker 02
- **Status:** COMPLETED
- **Files Changed:** `cli/cmdpipeline/pipeline_stacktrace.go`, `cli/cmdpipeline/pipeline_error_extract.go`, `cli/cmdpipeline/pipeline_history.go`, `cli/cmdpipeline/pipeline_error_extract_test.go`
- **Accomplishments:**
  - Enhanced `findStackTraceStart` to scan backward for `FAIL: <test>` or `ERROR: <test>` headers preceding `Traceback (most recent call last):`, capturing test failure metadata in the stack trace block.
  - Updated `isStackStopLine` so empty lines within Python tracebacks and multi-line assertion diffs are not prematurely truncated.
  - Expanded `isTestFailureSummary` to match `FAIL:`, `ERROR:`, `FAILED`, and `--- FAIL:`.
  - Updated `printFailedJobItemToBuilder` in `pipeline_history.go` to format `j.StackTrace` and `j.ErrorLines` in the commit inspector / heatmap summary cards.
  - Added regression unit tests in `pipeline_error_extract_test.go` verifying header capture and diff preservation.

### Subtask 03: Targeted Verification & Quality Audit
- **Owner:** Lead
- **Status:** COMPLETED
- **Accomplishments:**
  - Verified `python .github/scripts/go-format-check.py --check-only` passes across 3,951 Go files (exit code 0).
  - Verified `python linter-scripts/check-forbidden-strings.py` and `python linter-scripts/check-relative-paths.py` pass cleanly.
  - Verified `python 03-ai-scripts/05-guideline-autofixer.py cli/cmdpipeline --check-only` reports 100% compliance.
  - Verified targeted unit tests in `cli/cmdpipeline` pass in 0.144s.

### Subtask 04: Minor Version Bump & Release Ceremony
- **Owner:** Lead
- **Status:** COMPLETED
- **Accomplishments:**
  - Executed minor version bump from `v6.477.0` to `v6.478.0`.
  - Updated `version.json`, `package.json`, `cli/constants/constants.go`, `.gitmap/release/latest.json`, `readme.md`, `what-to-read.md`, and `changelog.md`.
  - Created release notes `.ai-memory/release/release-notes-v6.478.0.md`.
