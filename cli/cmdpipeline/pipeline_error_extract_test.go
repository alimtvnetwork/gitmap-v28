package cmdpipeline

import (
	"strings"
	"testing"
)

func TestExtractTauriBundlerAndCompilerWarnings(t *testing.T) {
	rawLogs := strings.Join([]string{
		"build\tBuild Desktop App\t2026-03-23T10:00:00.0000000Z Browserslist: browsers data (caniuse-lite) is 8 months old. Please run: npx update-browserslist-db",
		"build\tBuild Desktop App\t2026-03-23T10:00:01.0000000Z vite v7.3.1 building client environment for production...",
		"build\tBuild Desktop App\t2026-03-23T10:00:02.0000000Z transforming...",
		"build\tBuild Desktop App\t2026-03-23T10:00:05.0000000Z warning: unused import: `std::path::Path`",
		"build\tBuild Desktop App\t2026-03-23T10:00:05.0000000Z   --> src-tauri/src/main.rs:12:5",
		"build\tBuild Desktop App\t2026-03-23T10:00:05.0000000Z    |",
		"build\tBuild Desktop App\t2026-03-23T10:00:05.0000000Z 12 | use std::path::Path;",
		"build\tBuild Desktop App\t2026-03-23T10:00:05.0000000Z    |     ^^^^^^^^^^^^^^^",
		"build\tBuild Desktop App\t2026-03-23T10:00:05.0000000Z    |",
		"build\tBuild Desktop App\t2026-03-23T10:00:05.0000000Z    = note: `#[warn(unused_imports)]` on by default",
		"build\tBuild Desktop App\t2026-03-23T10:00:06.0000000Z Info Looking up installed tauri packages...",
		"build\tBuild Desktop App\t2026-03-23T10:00:07.0000000Z failed to bundle project: Failed to copy binary from \"/Users/runner/work/agm-desktop/agm-desktop/src-tauri/target/universal-apple-darwin/release/agm\" to \"/Users/runner/work/agm-desktop/agm-desktop/src-tauri/target/universal-apple-darwin/release/bundle/macos/agm.app/Contents/MacOS/agm\": \"/Users/runner/work/agm-desktop/agm-desktop/src-tauri/target/universal-apple-darwin/release/agm\" does not exist",
		"build\tBuild Desktop App\t2026-03-23T10:00:08.0000000Z ##[error]Process completed with exit code 1.",
	}, "\n")

	jobs := ParseFailedLogLines(rawLogs)
	if len(jobs) == 0 {
		t.Fatalf("expected at least 1 failed job item, got 0")
	}

	job := jobs[0]
	if !strings.Contains(job.FailureSummary, "failed to bundle project") {
		t.Errorf("expected FailureSummary to contain 'failed to bundle project', got: %s", job.FailureSummary)
	}
	if !strings.Contains(job.FailureSummary, "does not exist") {
		t.Errorf("expected FailureSummary to contain 'does not exist', got: %s", job.FailureSummary)
	}
	if strings.Contains(job.FailureSummary, "Process completed with exit code 1") {
		t.Errorf("FailureSummary should NOT be generic exit code when bundler error exists, got: %s", job.FailureSummary)
	}

	hasWarningLinePointer := false
	hasWarningCodeSnippet := false
	for _, w := range job.Warnings {
		if strings.Contains(w, "--> src-tauri/src/main.rs:12:5") {
			hasWarningLinePointer = true
		}
		if strings.Contains(w, "use std::path::Path;") {
			hasWarningCodeSnippet = true
		}
	}

	if !hasWarningLinePointer {
		t.Errorf("expected job.Warnings to capture line pointer '--> src-tauri/src/main.rs:12:5', got: %v", job.Warnings)
	}
	if !hasWarningCodeSnippet {
		t.Errorf("expected job.Warnings to capture code snippet 'use std::path::Path;', got: %v", job.Warnings)
	}
}

func TestIsStrongerSummaryBundlerPriority(t *testing.T) {
	bundlerErr := `failed to bundle project: Failed to copy binary: file does not exist`
	genericExit := `Error: Process completed with exit code 1.`

	if !isStrongerSummary(bundlerErr, genericExit) {
		t.Errorf("bundler error should be stronger summary than generic exit code")
	}

	if isStrongerSummary(genericExit, bundlerErr) {
		t.Errorf("generic exit code should NOT override existing bundler error summary")
	}
}

func TestIsToolProgressNoiseTauri(t *testing.T) {
	noiseLines := []string{
		"Browserslist: browsers data (caniuse-lite) is 8 months old. Please run: npx update-browserslist-db",
		"transforming...",
		"vite v7.3.1 building client environment for production...",
		"Info Looking up installed tauri packages to check mismatched versions...",
		"Running beforeBuildCommand `npm run build`",
	}

	for _, line := range noiseLines {
		if !isToolProgressNoise(line) {
			t.Errorf("expected line to be recognized as tool progress noise: %s", line)
		}
	}
}

func TestIsToolProgressNoise_RanTests(t *testing.T) {
	lines := []string{
		"Ran 18 tests in 75.725s",
		"Ran 1 test in 0.002s",
		"Ran 142 tests in 3.14s",
	}
	for _, l := range lines {
		if !isToolProgressNoise(l) {
			t.Errorf("expected line to be recognized as tool progress noise: %s", l)
		}
	}
}

func TestIsStrongerSummary_UnitTestPriority(t *testing.T) {
	unitTestErr := "FAIL: test_user_authentication (tests.test_auth.TestAuth.test_user_authentication)"
	genericExit := "Error: Process completed with exit code 1."
	stepNotice := "Step 'Run Unit Tests' (step #4) failure"

	if !isStrongerSummary(unitTestErr, genericExit) {
		t.Errorf("unit test failure should be stronger than generic exit code")
	}
	if !isStrongerSummary(unitTestErr, stepNotice) {
		t.Errorf("unit test failure should be stronger than step failure notice")
	}
	if isStrongerSummary(genericExit, unitTestErr) {
		t.Errorf("generic exit code should NOT override unit test failure")
	}
	if isStrongerSummary(stepNotice, unitTestErr) {
		t.Errorf("step notice should NOT override unit test failure")
	}
}

func TestExtractPythonStackTrace(t *testing.T) {
	raw := "======================================================================\n" +
		"FAIL: test_something (tests.test_foo.TestFoo.test_something)\n" +
		"----------------------------------------------------------------------\n" +
		"Traceback (most recent call last):\n" +
		"  File \"/workspace/tests/test_foo.py\", line 42, in test_something\n" +
		"    self.assertEqual(a, b)\n" +
		"AssertionError: 1 != 2\n" +
		"----------------------------------------------------------------------\n" +
		"Ran 18 tests in 75.725s\n"

	stack := extractStackTraceFromLog(raw)
	if !strings.Contains(stack, "Traceback (most recent call last):") {
		t.Fatalf("missing traceback header: %s", stack)
	}
	if !strings.Contains(stack, "File \"/workspace/tests/test_foo.py\", line 42") {
		t.Fatalf("missing file frame: %s", stack)
	}
	if !strings.Contains(stack, "self.assertEqual(a, b)") {
		t.Fatalf("missing code frame: %s", stack)
	}
	if !strings.Contains(stack, "AssertionError: 1 != 2") {
		t.Fatalf("missing assertion error: %s", stack)
	}
	if strings.Contains(stack, "Ran 18 tests in") {
		t.Fatalf("stack trace captured beyond terminator: %s", stack)
	}
}

func TestExtractPythonStackTraceWithAssertionDiff(t *testing.T) {
	raw := "Traceback (most recent call last):\n" +
		"  File \"tests/test_diff.py\", line 10, in test_diff\n" +
		"    self.assertEqual(list_a, list_b)\n" +
		"AssertionError: Lists differ: [1] != [2]\n" +
		"--- expected\n" +
		"+++ actual\n" +
		"- [1]\n" +
		"+ [2]\n" +
		"----------------------------------------------------------------------\n" +
		"Ran 1 test in 0.001s\n"

	stack := extractStackTraceFromLog(raw)
	if !strings.Contains(stack, "--- expected") || !strings.Contains(stack, "+++ actual") {
		t.Fatalf("expected assertion diff lines to be captured in stack trace: %s", stack)
	}
	if !strings.Contains(stack, "- [1]") || !strings.Contains(stack, "+ [2]") {
		t.Fatalf("expected diff values to be captured: %s", stack)
	}
}

func TestIsOnlyFallbackLogs(t *testing.T) {
	fallbackOnly := "Job 1\tStep 1\tFAIL: Step 'Run Unit Tests' (step #4) failure\n" +
		"Job 1\tJob Execution\tFAIL: Job 'Job 1' failure: timed out\n"

	if !isOnlyFallbackLogs(fallbackOnly) {
		t.Errorf("expected fallbackOnly to be identified as only fallback logs")
	}

	realLog := "Job 1\tStep 1\tFAIL: Step 'Run Unit Tests' (step #4) failure\n" +
		"Traceback (most recent call last):\n" +
		"  File \"test.py\", line 1\n" +
		"AssertionError: fail\n"

	if isOnlyFallbackLogs(realLog) {
		t.Errorf("expected realLog with traceback to NOT be identified as only fallback logs")
	}
}

func TestParseFailedLogLines_PythonUnitTestE2E(t *testing.T) {
	rawLogs := strings.Join([]string{
		"Lint Script Unit Tests\tRun lint-script unit tests\t2026-10-05T02:17:50.7182265Z ##[group]Run python3 .github/scripts/tests/test_ci_scripts.py",
		"Lint Script Unit Tests\tRun lint-script unit tests\t2026-10-05T02:17:50.7183234Z python3 .github/scripts/tests/test_ci_scripts.py",
		"Lint Script Unit Tests\tRun lint-script unit tests\t2026-10-05T02:17:50.7478692Z ##[endgroup]",
		"Lint Script Unit Tests\tRun lint-script unit tests\t2026-10-05T02:19:06.5336430Z ..........F.......",
		"Lint Script Unit Tests\tRun lint-script unit tests\t2026-10-05T02:19:06.5337201Z ======================================================================",
		"Lint Script Unit Tests\tRun lint-script unit tests\t2026-10-05T02:19:06.5338034Z FAIL: test_gofmt_check_clean_repo (__main__.TestGoFormatCheck.test_gofmt_check_clean_repo)",
		"Lint Script Unit Tests\tRun lint-script unit tests\t2026-10-05T02:19:06.5338863Z ----------------------------------------------------------------------",
		"Lint Script Unit Tests\tRun lint-script unit tests\t2026-10-05T02:19:06.5339507Z Traceback (most recent call last):",
		"Lint Script Unit Tests\tRun lint-script unit tests\t2026-10-05T02:19:06.5340450Z   File \"/home/runner/work/gitmap-v28/gitmap-v28/.github/scripts/tests/test_ci_scripts.py\", line 171, in test_gofmt_check_clean_repo",
		"Lint Script Unit Tests\tRun lint-script unit tests\t2026-10-05T02:19:06.5341123Z     self.assertEqual(res.returncode, 0)",
		"Lint Script Unit Tests\tRun lint-script unit tests\t2026-10-05T02:19:06.5341478Z AssertionError: 1 != 0",
		"Lint Script Unit Tests\tRun lint-script unit tests\t2026-10-05T02:19:06.5341775Z ----------------------------------------------------------------------",
		"Lint Script Unit Tests\tRun lint-script unit tests\t2026-10-05T02:19:06.5342119Z Ran 18 tests in 75.725s",
		"Lint Script Unit Tests\tRun lint-script unit tests\t2026-10-05T02:19:06.5342351Z FAILED (failures=1)",
		"Lint Script Unit Tests\tRun lint-script unit tests\t2026-10-05T02:19:06.5486180Z ##[error]Process completed with exit code 1.",
	}, "\n")

	jobs := ParseFailedLogLines(rawLogs)
	if len(jobs) == 0 {
		t.Fatalf("expected at least 1 failed job item, got 0")
	}

	job := jobs[0]
	if !strings.Contains(job.FailureSummary, "FAIL: test_gofmt_check_clean_repo") {
		t.Errorf("expected FailureSummary to be unit test failure line, got: %s", job.FailureSummary)
	}

	if !strings.Contains(job.StackTrace, "Traceback (most recent call last):") {
		t.Errorf("expected StackTrace to contain Traceback header, got: %s", job.StackTrace)
	}
	if !strings.Contains(job.StackTrace, "test_gofmt_check_clean_repo") {
		t.Errorf("expected StackTrace to contain test name, got: %s", job.StackTrace)
	}
	if !strings.Contains(job.StackTrace, "AssertionError: 1 != 0") {
		t.Errorf("expected StackTrace to contain AssertionError, got: %s", job.StackTrace)
	}

	for _, line := range job.ErrorLines {
		if strings.Contains(line, "Ran 18 tests in") {
			t.Errorf("expected 'Ran 18 tests in' noise to be excluded from ErrorLines, but found: %s", line)
		}
	}
}

func TestExtractStackTraceCapturesTestFailureHeader(t *testing.T) {
	raw := strings.Join([]string{
		"======================================================================",
		"FAIL: test_validate_output (tests.test_cli.TestCLI.test_validate_output)",
		"----------------------------------------------------------------------",
		"Traceback (most recent call last):",
		"  File \"tests/test_cli.py\", line 88, in test_validate_output",
		"    self.assertEqual(actual, expected)",
		"AssertionError: 'foo' != 'bar'",
		"----------------------------------------------------------------------",
		"Ran 5 tests in 1.23s",
		"FAILED (failures=1)",
	}, "\n")

	stack := extractStackTraceFromLog(raw)
	assertStackHeaderAndFrames(t, stack)
}

func assertStackHeaderAndFrames(t *testing.T, stack string) {
	if !strings.Contains(stack, "FAIL: test_validate_output") {
		t.Fatalf("expected stack to capture FAIL header, got: %s", stack)
	}
	if !strings.Contains(stack, "Traceback (most recent call last):") {
		t.Fatalf("expected stack to contain Traceback header, got: %s", stack)
	}
	if !strings.Contains(stack, "AssertionError: 'foo' != 'bar'") {
		t.Fatalf("expected stack to contain AssertionError, got: %s", stack)
	}
	if strings.Contains(stack, "Ran 5 tests in") {
		t.Fatalf("stack trace captured past terminator: %s", stack)
	}
}

func TestExtractStackTrace_AssertionDiffWithBlankLines(t *testing.T) {
	raw := strings.Join([]string{
		"Traceback (most recent call last):",
		"  File \"tests/test_diff.py\", line 45, in test_multiline",
		"    self.assertEqual(first, second)",
		"AssertionError: Multi-line strings differ:",
		"",
		"--- expected",
		"+++ actual",
		"",
		"- alpha",
		"- beta",
		"+ alpha",
		"+ gamma",
		"----------------------------------------------------------------------",
		"Ran 1 test in 0.05s",
	}, "\n")

	stack := extractStackTraceFromLog(raw)
	assertAssertionDiffFrames(t, stack)
}

func assertAssertionDiffFrames(t *testing.T, stack string) {
	if !strings.Contains(stack, "--- expected") || !strings.Contains(stack, "+++ actual") {
		t.Fatalf("expected diff headers in stack, got: %s", stack)
	}
	if !strings.Contains(stack, "- beta") || !strings.Contains(stack, "+ gamma") {
		t.Fatalf("expected diff lines not to be clipped by blank lines, got: %s", stack)
	}
	if strings.Contains(stack, "Ran 1 test in") {
		t.Fatalf("stack trace captured beyond terminator: %s", stack)
	}
}

func TestIsTestFailureSummary_Expansions(t *testing.T) {
	validCases := []string{
		"FAIL: test_something (suite.TestClass)",
		"FAIL: custom_check (suite.TestClass)",
		"FAIL:\ttest_tab_delimited",
		"ERROR: test_db_timeout (tests.test_db.TestDB)",
		"ERROR: test_setup_failure",
		"FAILED tests/test_api.py::test_create - AssertionError",
		"--- FAIL: TestPipelineRunner (0.05s)",
	}
	for _, c := range validCases {
		if !isTestFailureSummary(c) {
			t.Errorf("expected isTestFailureSummary to return true for: %s", c)
		}
	}
	assertInvalidTestSummaries(t)
}

func assertInvalidTestSummaries(t *testing.T) {
	invalidCases := []string{
		"Process completed with exit code 1.",
		"Step 'Run Unit Tests' (step #4) failure",
		"Ran 18 tests in 75.725s",
		"FAILED (failures=1)",
		"PASS",
	}
	for _, c := range invalidCases {
		if isTestFailureSummary(c) {
			t.Errorf("expected isTestFailureSummary to return false for: %s", c)
		}
	}
}

func TestPrintFailedJobItemToBuilder_RendersTracebackAndErrors(t *testing.T) {
	var sb strings.Builder
	j := FailedJobItem{
		JobName:        "Test Runner",
		StepName:       "Run Tests",
		FailureSummary: "FAIL: test_feature (tests.TestFeature)",
		ErrorLines: []string{
			"FAIL: test_feature (tests.TestFeature)",
			"    Detailed assertion message",
		},
		StackTrace: "FAIL: test_feature\nTraceback (most recent call last):\n  File \"test.py\", line 10\nAssertionError",
	}
	printFailedJobItemToBuilder(&sb, j)
	out := sb.String()
	assertJobItemRenderOutput(t, out)
}

func assertJobItemRenderOutput(t *testing.T, out string) {
	if !strings.Contains(out, "[Job: Test Runner | Step: Run Tests]") {
		t.Errorf("expected job header, got: %s", out)
	}
	if !strings.Contains(out, "Detailed assertion message") {
		t.Errorf("expected error lines, got: %s", out)
	}
	if !strings.Contains(out, "Stack Trace:") || !strings.Contains(out, "AssertionError") {
		t.Errorf("expected stack trace in output, got: %s", out)
	}
}
