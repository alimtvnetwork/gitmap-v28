# Specification: Pipeline PE Unit Test Traceback Extraction & Heatmap Modernization

> **Spec ID:** 220-pipeline-pe-unit-test-traceback-and-heatmap  
> **Status:** Approved / Active  
> **Date:** 2026-10-05  
> **Scope:** `cli/cmdpipeline/`, `.github/scripts/tests/`, and CLI heatmap error display  

---

## 1. User Request (Verbatim)

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

## 2. Root Cause Analysis (4-Part RCA)

### 2.1 Problem Statement
When running `gitmap pe` (pipeline error log viewer) on workflow runs with failing unit tests (such as Python `unittest` in GitHub Actions run `37254814139`), the terminal output and heatmap only display a single generic step line:
`FAIL: Step 'Run lint-script unit tests' (step #4) failure`
It completely fails to extract the unit test failure details:
1. Which specific test failed (`FAIL: test_gofmt_check_clean_repo (...)`)
2. The Python traceback call (`Traceback (most recent call last): ... File "...", line 171`)
3. The assertion failure (`AssertionError: 1 != 0`)
4. Noise lines like `Ran 18 tests in 75.725s` are retained or misclassified instead of being filtered out.
Additionally, the actual unit test failed in CI because 5 Go files in `cli/` were not formatted with `gofmt`, and `REPO_ROOT` in `test_ci_scripts.py` was improperly resolved to `.github/`.

### 2.2 Root Cause Analysis
1. **Premature Fallback Log Caching (`cli/cmdpipeline/pipeline_query.go`):** When `gitmap pe` queries an active or newly-failing run, `gh run view --log-failed` fails temporarily while the run is still in-progress. The code then generates a 1-line fallback `Step '...' failure` and permanently writes it to `.log` cache. Subsequent invocations read this 1-line dummy file from disk and never query the actual logs.
2. **Missing Python Traceback Markers in Stack Trace Parser (`cli/cmdpipeline/pipeline_stacktrace.go`):** `findStackTraceStart` only checked `goroutine `, `panic:`, and `Stack Trace:`. It had zero recognition of `Traceback (most recent call last):` or `FAIL: test_`. `isStackFrame` only checked Go runtime markers and had no recognition of Python stack frames (`File "...", line ...`). `isStackTerminator` immediately killed parsing upon seeing `===` or `---` which are standard Python unittest separators!
3. **Log Context & Noise Filter Deficiencies (`cli/cmdpipeline/pipeline_error_extract.go`):** `isToolProgressNoise` lacked filters for test runner execution summaries like `^Ran \d+ tests in`. `isTerminalStepError` stopped context capture prematurely upon seeing `failed (failures=`.
4. **Unformatted Go Files in Codebase:** 5 files (`cli/cmd/help.go`, `cli/cmdagent/agent_cleanup.go`, `cli/cmdagent/agent_diagnose.go`, `cli/cmdcursor/cursor_sync.go`, `cli/cmddb/cmddb_reset.go`) were committed without `gofmt`, causing `python .github/scripts/go-format-check.py --check-only` to exit with code 1.

### 2.3 Corrective Action
1. Enhance `pipeline_stacktrace.go` to recognize Python tracebacks, pytest failures, and Go test failures, extracting multi-line stack frames including `Traceback (most recent call last):`, file/line references, and exception types.
2. In `pipeline_error_extract.go`, filter out test runner noise (`Ran \d+ tests in ...`) while preserving `FAIL: <test_name>` and assertion lines.
3. In `pipeline_query.go`, prevent fallback synthetic logs from permanently poisoning the run log cache so completed runs always fetch full logs.
4. Format all 5 dirty Go files with `gofmt -w` and fix `REPO_ROOT` path resolution in `.github/scripts/tests/test_ci_scripts.py`.
5. Update `pipeline_details.go` and `pipeline_ai_errors.go` to display unit test tracebacks and failed test information prominently in the heatmap PE cards and clipboard logs.

### 2.4 Prevention
Centralize unit test log extraction in dedicated parser functions with comprehensive regression tests in `pipeline_error_extract_test.go`.

---

## 3. Architecture & Functional Specification

### 3.1 Unit Test Failure Extraction Pipeline
```text
Raw GitHub Actions Log
       │
       ▼
┌──────────────────────────────────────────────┐
│ parseLogLine (Strips ANSI, splits job/step)  │
└──────────────────────────────────────────────┘
       │
       ├─► Is Tool Noise? (Ran \d+ tests in...) ──► Filter / Drop
       │
       ├─► Python Traceback / Go Test Failure?
       │   ├─► StackTrace Extractor (captures Traceback through Exception)
       │   └─► FailureSummary (captures FAIL: test_name)
       │
       ▼
┌──────────────────────────────────────────────┐
│ FailedJobItem                                │
│ - FailureSummary: FAIL: test_gofmt_check...  │
│ - ErrorLines: [AssertionError: 1 != 0, ...]  │
│ - StackTrace: Traceback (most recent call):  │
│               File "...", line 171           │
└──────────────────────────────────────────────┘
       │
       ▼
┌──────────────────────────────────────────────┐
│ Heatmap PE Terminal & Clipboard Formatter    │
│ - Displays Failed Test Header                │
│ - Displays Formatted Stack Trace             │
│ - Supplies High-Fidelity Context to AI & Dev │
└──────────────────────────────────────────────┘
```

### 3.2 Non-Negotiable Acceptance Criteria
1. When parsing unit test logs containing `Traceback (most recent call last):`, `job.StackTrace` must contain the complete traceback block.
2. `FAIL: test_name` must be elevated as the primary `FailureSummary` over generic `Process completed with exit code 1`.
3. `Ran \d+ tests in ...` lines must be discarded as non-actionable noise.
4. `python .github/scripts/go-format-check.py --check-only` must pass cleanly with exit code 0.
5. All new parser behaviors must be backed by unit tests in `pipeline_error_extract_test.go`.
