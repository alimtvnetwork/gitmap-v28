# Component & CLI Specification: Pipeline PE Unit Test Traceback & Heatmap Display

> **Spec ID:** 220-pipeline-pe-unit-test-traceback-and-heatmap  
> **Status:** Approved / Active  
> **Component:** CLI Pipeline (`cli/cmdpipeline/`) & Heatmap (`gitmap pe`)  

---

## 1. Technical Components & Modification Plan

### 1.1 Stack Trace Parser Enhancement (`cli/cmdpipeline/pipeline_stacktrace.go`)
- **New Stack Start Markers:**
  - `Traceback (most recent call last):` (Python unittest/pytest)
  - `panic:` (Go runtime)
  - `goroutine ` (Go runtime)
  - `Stack Trace:` (Generic / C# / Java / Node)
- **New Stack Frame Recognizers:**
  - Python stack frame: `File "...", line \d+, in \w+`
  - Indented code lines (`    <expression>`)
  - Exception lines: `\w+Error: .*`, `Exception: .*`
  - Go test failure frames: `\w+_test.go:\d+: .*`
- **Updated Stack Terminators:**
  - Ensure `isStackTerminator` does NOT prematurely terminate on Python unittest banners like `======================================================================` or `----------------------------------------------------------------------` before the traceback body is captured.
  - Terminate on subsequent job boundaries (`##[`) or step completion lines.

### 1.2 Error Line & Noise Filter (`cli/cmdpipeline/pipeline_error_extract.go`)
- **Noise Suppression (`isToolProgressNoise`):**
  - Add filter for `Ran \d+ tests in \d+(\.\d+)?s` (e.g., `Ran 18 tests in 75.725s`).
  - Add filter for `^=+$` and `^-+$` decorative divider lines when they are not part of an active stack frame.
- **Summary Elevation (`isStrongerSummary`):**
  - Elevate `FAIL: <test_name>` (e.g. `FAIL: test_gofmt_check_clean_repo (...)`) above generic runner lines like `FAIL: Step 'Run lint-script unit tests' (step #4) failure` or `Process completed with exit code 1`.
  - Ensure assertion failures (`AssertionError: 1 != 0`) are captured in `ErrorLines`.

### 1.3 Cache Poisoning Guard (`cli/cmdpipeline/pipeline_query.go`)
- **Log Fallback Guard:**
  - When `gh run view <id> --log-failed` fails during an in-progress or early run state, the generated fallback (`Step '...' (step #...) failure`) must NOT be saved as the permanent `.log` file if the run is still active.
  - If a cached log on disk contains only synthetic fallback lines, check if the run is now completed and attempt to refresh the log with real logs from GitHub CLI.

### 1.4 Heatmap PE Formatting & AI Ingestion (`cli/cmdpipeline/pipeline_details.go`, `pipeline_ai_errors.go`)
- Format failed unit tests cleanly in the rendered terminal card:
  - Header: `│ Failed Test: test_gofmt_check_clean_repo`
  - Traceback / Stack Trace: Rendered with yellow/cyan highlights and proper indentation
  - Details: Shows assertion failure details clearly so both humans and AI can immediately diagnose root causes.

---

## 2. CLI Interface & Expected Output

### 2.1 Before (Truncated / Generic)
```text
  ┌─ [1/1] CI #37254814139 ──────────────────────────
  │ Job:   Lint Script Unit Tests
  │ Step:  Run lint-script unit tests
  │ Error: FAIL: Step 'Run lint-script unit tests' (step #4) failure
  └──────────────────────────────────────────────────────────
```

### 2.2 After (Enhanced High-Fidelity Output)
```text
  ┌─ [1/1] CI #37254814139 ──────────────────────────
  │ When Run:  2026-10-05 02:17:27 UTC (5m ago)
  │ Duration:  4m 37s
  │ Branch:    main | Commit: 0ebee163e9da547d48c5b77dcac39e14a7e9535d
  │ Saved Log: data/pipeline/alimtvnetwork-gitmap-v28/37254814139.log
  │ URL:       https://github.com/alimtvnetwork/gitmap-v28/actions/runs/37254814139
  │ Job:   Lint Script Unit Tests
  │ Step:  Run lint-script unit tests
  │ Error: FAIL: test_gofmt_check_clean_repo (__main__.TestGoFormatCheck.test_gofmt_check_clean_repo)
  │ Stack Trace:
  │   Traceback (most recent call last):
  │     File ".github/scripts/tests/test_ci_scripts.py", line 171, in test_gofmt_check_clean_repo
  │       self.assertEqual(res.returncode, 0)
  │   AssertionError: 1 != 0
  └──────────────────────────────────────────────────────────
```

---

## 3. Verification & Test Plan
1. **Unit Test Verification (`pipeline_error_extract_test.go`):**
   - Test parsing of Python `unittest` output containing `Traceback`, `AssertionError`, and `Ran 18 tests in ...`.
   - Verify `Ran 18 tests` is filtered out.
   - Verify traceback lines are captured in `job.StackTrace`.
   - Verify `FAIL: test_...` is selected as `job.FailureSummary`.
2. **Go Formatting Verification:**
   - Execute `gofmt -w` on the 5 unformatted files.
   - Execute `python .github/scripts/go-format-check.py --check-only` -> Must return code 0.
3. **CI Script Test Suite:**
   - Execute `python .github/scripts/tests/test_ci_scripts.py` -> All 18 tests pass.
