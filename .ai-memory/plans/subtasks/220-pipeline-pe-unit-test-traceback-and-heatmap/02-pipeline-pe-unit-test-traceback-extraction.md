# Subtask 02: Pipeline PE Unit Test Traceback Extraction

> **Parent Plan:** `220-pipeline-pe-unit-test-traceback-and-heatmap`  
> **Status:** DONE  
> **Target:** `cli/cmdpipeline/pipeline_stacktrace.go`, `cli/cmdpipeline/pipeline_error_extract.go`  

---

## Objectives
1. Enhance `pipeline_stacktrace.go`:
   - Add Python `Traceback (most recent call last):` to `findStackTraceStart` and `isStackStartLine`.
   - Recognize Python stack frames (`File "...", line \d+, in \w+`, indented lines, exception types).
   - Prevent premature termination of stack frames on unittest divider lines (`===`, `---`).
2. Enhance `pipeline_error_extract.go`:
   - Filter out test progress noise like `Ran \d+ tests in ...` in `isToolProgressNoise`.
   - Elevate test failure lines (`FAIL: test_...`) in `isStrongerSummary`.
   - Ensure assertion failures and traceback contexts are captured without premature truncation.
