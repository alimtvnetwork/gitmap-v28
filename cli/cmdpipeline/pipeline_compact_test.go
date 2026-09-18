package cmdpipeline

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
)

func TestIsOkLogLine_WithPassingInputs_ReturnsTrue(t *testing.T) {
	cases := []string{
		"PASS",
		"ok",
		"PASS: TestSomething",
		"ok\tgithub.com/alimtvnetwork/gitmap-v28/cli/cmd\t0.12s",
		"ok  github.com/alimtvnetwork/gitmap-v28/cli/cmd 0.12s",
		"?\tgithub.com/alimtvnetwork/gitmap-v28/cli/pkg\t[no test files]",
		"?   github.com/alimtvnetwork/gitmap-v28/cli/pkg [no test files]",
		"--- PASS: TestSample (0.05s)",
		"=== RUN   TestSample",
		"✔ ok: macro verified",
		"✔ Macro exported successfully",
		"test proxy::opencode_sync::canonical_family_tests::canonical_families_match_ids_are_unique_globally_after_normalization ... ok",
		"test proxy::opencode_sync::tests::apply_sync_keeps_non_conflicting_user_fields ... ok",
		"test proxy::opencode_sync::tests::build_variants_flash_resolve_to_real_ids ... ok",
		"test proxy::opencode_sync::tests::test_base_url_matches_different_urls ... ok",
		"test proxy::opencode_sync::tests::catalog_marks_gemini_31_flash_lite_as_non_variant ... ignored",
		"test result: ok. 12 passed; 0 failed; 0 ignored; 0 measured; 0 filtered out; finished in 0.05s",
	}
	assertAllOkLinesTrue(t, cases)
}

func assertAllOkLinesTrue(t *testing.T, cases []string) {
	for _, c := range cases {
		if !isOkLogLine(c) {
			t.Errorf("expected isOkLogLine(%q) to be true, got false", c)
		}
	}
}

func TestIsOkLogLine_WithFailingInputs_ReturnsFalse(t *testing.T) {
	cases := []string{
		"--- FAIL: TestSample (0.05s)",
		"FAIL\tgithub.com/alimtvnetwork/gitmap-v28/cli/cmd\t0.50s",
		"panic: runtime error: invalid memory address",
		"syntax error: unexpected token",
		"exit status 1",
		"Error: compilation failed",
		"##[error]Process completed with exit code 1.",
		"test proxy::mappers::tool_result_compressor::tests::test_sanitize_tool_result_blocks ... FAILED",
	}
	assertAllFailingLinesFalse(t, cases)
}

func assertAllFailingLinesFalse(t *testing.T, cases []string) {
	for _, c := range cases {
		if isOkLogLine(c) {
			t.Errorf("expected isOkLogLine(%q) to be false, got true", c)
		}
	}
}

func TestFilterCompactLines_WithMixedLines_FiltersOkLines(t *testing.T) {
	input := []string{
		"=== RUN   TestFailure",
		"    sample_test.go:42: assertion failed",
		"--- FAIL: TestFailure (0.02s)",
		"=== RUN   TestPass",
		"--- PASS: TestPass (0.01s)",
		"FAIL",
		"ok\tgithub.com/alimtvnetwork/gitmap-v28/cli/cmd\t0.30s",
	}
	result := filterCompactLines(input)
	assertCompactedLines(t, result)
}

func assertCompactedLines(t *testing.T, result []string) {
	if len(result) != 3 {
		t.Fatalf("expected 3 filtered lines, got %d: %v", len(result), result)
	}
	if !strings.Contains(result[0], "assertion failed") {
		t.Errorf("unexpected first line: %s", result[0])
	}
	if !strings.Contains(result[1], "--- FAIL: TestFailure") {
		t.Errorf("unexpected second line: %s", result[1])
	}
	if result[2] != "FAIL" {
		t.Errorf("unexpected third line: %s", result[2])
	}
}

func TestParsePipelineErrorFlags_WithDetailedFlags_SetsIsDetailed(t *testing.T) {
	testDetailedFlagVariant(t, "--detailed")
	testDetailedFlagVariant(t, "--verbose")
	testDetailedFlagVariant(t, "--v")
	testDetailedFlagVariant(t, "-v")
	testDetailedFlagVariant(t, "-V")
}

func testDetailedFlagVariant(t *testing.T, flag string) {
	flags := ParsePipelineErrorFlags([]string{"pipeline", "error-logs", flag})
	if !flags.IsDetailed {
		t.Errorf("expected IsDetailed to be true with flag %s", flag)
	}
}

func TestParsePipelineErrorFlags_WithoutDetailedFlag_SetsIsDetailedFalse(t *testing.T) {
	flags := ParsePipelineErrorFlags([]string{"pipeline", "error-logs"})
	if flags.IsDetailed {
		t.Errorf("expected IsDetailed to be false by default")
	}
}

func TestCompactErrorPayload_WithFailedRuns_RemovesOkLines(t *testing.T) {
	payload := buildSampleErrorPayload()
	compactErrorPayload(&payload)
	assertPayloadCompacted(t, payload)
}

func buildSampleErrorPayload() PipelineErrorLogsPayload {
	return PipelineErrorLogsPayload{
		FailedRuns: []FailedRunItem{
			{
				WorkflowName: "CI",
				RunId:        12345,
				FailedJobs: []FailedJobItem{
					{
						JobName:  "Test",
						StepName: "Unit Tests",
						ErrorLines: []string{
							"=== RUN   TestOne",
							"--- PASS: TestOne (0.01s)",
							"--- FAIL: TestTwo (0.02s)",
							"ok\tpkg/math\t0.05s",
						},
					},
				},
			},
		},
		SectionFailures: []SectionFailure{
			{
				WorkflowName: "CI",
				RunId:        12345,
				JobName:      "Test",
				StepName:     "Unit Tests",
				ErrorLines: []string{
					"=== RUN   TestOne",
					"--- PASS: TestOne (0.01s)",
					"--- FAIL: TestTwo (0.02s)",
					"ok\tpkg/math\t0.05s",
				},
			},
		},
	}
}

func assertPayloadCompacted(t *testing.T, payload PipelineErrorLogsPayload) {
	jobLines := payload.FailedRuns[0].FailedJobs[0].ErrorLines
	if len(jobLines) != 1 || !strings.Contains(jobLines[0], "--- FAIL: TestTwo") {
		t.Errorf("expected only FAIL line in job error lines, got: %v", jobLines)
	}

	secLines := payload.SectionFailures[0].ErrorLines
	if len(secLines) != 1 || !strings.Contains(secLines[0], "--- FAIL: TestTwo") {
		t.Errorf("expected only FAIL line in section error lines, got: %v", secLines)
	}
}

func TestFilterCompactLogText(t *testing.T) {
	raw := "=== RUN TestA\n--- PASS: TestA (0.01s)\n--- FAIL: TestB (0.02s)\nok pkg/a 0.05s\nProcess completed with exit code 1."
	compact, count := FilterCompactLogText(raw)
	if count != 3 {
		t.Errorf("expected 3 filtered ok lines, got %d", count)
	}
	if strings.Contains(compact, "PASS") || strings.Contains(compact, "ok pkg/a") {
		t.Errorf("compact text contains ok lines: %q", compact)
	}
	if !strings.Contains(compact, "FAIL: TestB") {
		t.Errorf("compact text missing FAIL line: %q", compact)
	}
}

func TestDualTableRepoSplitDbStorage(t *testing.T) {
	repo := "test-owner/test-dual-table"
	pipeDb, err := pipelinedb.OpenPipelineSplitDb(repo)
	if err != nil {
		t.Fatalf("failed to open split db: %v", err)
	}
	defer pipeDb.Close()
	_ = pipeDb.Reset()

	run := ghRunItem{
		DatabaseId: 88001, Name: "CI", Status: "completed", Conclusion: "failure",
		HeadBranch: "main", HeadSha: "sha88001",
	}
	recordSingleSplitRun(pipeDb, PipelineStatusPayload{Repo: repo}, run)
	clean := "--- FAIL: Test1\nok pkg 0.01s\nError: boom"
	raw := "=== RUN Test1\n--- PASS: Test0\n--- FAIL: Test1\nok pkg 0.01s\nError: boom"
	persistSingleFailedRunLog(pipeDb, repo, run, clean, raw)

	verifySplitDbDualLogs(t, pipeDb, 88001)
}

func verifySplitDbDualLogs(t *testing.T, pipeDb *pipelinedb.PipelineSplitDb, runId uint64) {
	detailRes := pipeDb.QueryDetailedErrorLogsByRunId(runId)
	if detailRes.IsCountOtherThan(1) {
		t.Fatalf("expected 1 detail log, got %d (err: %v)", detailRes.Count(), detailRes.AppError())
	}
	details := detailRes.Data
	if !strings.Contains(details[0].RawLogs, "PASS: Test0") {
		t.Errorf("detail log missing raw PASS line: %s", details[0].RawLogs)
	}

	compactRes := pipeDb.QueryCompactErrorLogsByRunId(runId)
	if compactRes.IsCountOtherThan(1) {
		t.Fatalf("expected 1 compact log, got %d (err: %v)", compactRes.Count(), compactRes.AppError())
	}
	compacts := compactRes.Data
	if strings.Contains(compacts[0].ErrorText, "ok pkg") {
		t.Errorf("compact log contains ok lines: %s", compacts[0].ErrorText)
	}
	if compacts[0].FilteredOkCount != 1 {
		t.Errorf("expected FilteredOkCount 1, got %d", compacts[0].FilteredOkCount)
	}
}

func TestCleanAnnotationError_WithFileAndLine_FormatsLocation(t *testing.T) {
	raw := "::error file=cli/cmd/join.go,line=34,col=2::[gocritic] appendAssign: append result not assigned"
	got := cleanAnnotationError(raw)
	want := "cli/cmd/join.go:34:2: [gocritic] appendAssign: append result not assigned"
	if got != want {
		t.Errorf("cleanAnnotationError() = %q, want %q", got, want)
	}
}

func TestCleanAnnotationError_WithoutAnnotation_ReturnsOriginal(t *testing.T) {
	raw := "regular error line without annotation"
	got := cleanAnnotationError(raw)
	if got != raw {
		t.Errorf("cleanAnnotationError() = %q, want %q", got, raw)
	}
}

func TestFormatSectionMetadata_WithScriptAndFile_DisplaysThem(t *testing.T) {
	sec := SectionFailure{
		WorkflowName:   "CI",
		RunId:          12345,
		JobName:        "gocritic",
		StepName:       "Gocritic diff",
		FailureSummary: "cli/cmd/join.go:34:2: [gocritic] appendAssign",
		ErrorLines: []string{
			"  │ File:     cli/cmd/join.go:34:2",
			"  │ Script:   .github/scripts/check-single-linter-diff.py",
		},
	}
	var sb strings.Builder
	formatSectionMetadata(&sb, sec)
	out := sb.String()
	if !strings.Contains(out, "Script:   .github/scripts/check-single-linter-diff.py") {
		t.Errorf("expected script in metadata, got: %s", out)
	}
	if !strings.Contains(out, "File:     cli/cmd/join.go:34:2") {
		t.Errorf("expected file in metadata, got: %s", out)
	}
}

func TestExtractBoundedStackTrace_StopsAtNonStackLine(t *testing.T) {
	raw := "goroutine 1 [running]:\n" +
		"testing.tRunner()\n" +
		"\t/opt/go/testing.go:100 +0x12\n" +
		"panic(0x123)\n" +
		"\t/opt/go/panic.go:50 +0x4\n" +
		"FAIL\tgithub.com/example/pkg\t0.1s\n" +
		"ok\tgithub.com/example/other\t0.05s\n" +
		"go: downloading github.com/foo/bar\n"

	stack := extractStackTraceFromLog(raw)
	if strings.Contains(stack, "FAIL") || strings.Contains(stack, "downloading") {
		t.Errorf("stack trace captured lines past termination: %s", stack)
	}
	if !strings.Contains(stack, "goroutine 1 [running]:") || !strings.Contains(stack, "panic(0x123)") {
		t.Errorf("stack trace missing frames: %s", stack)
	}
}

func TestAssembleJobItems_DoesNotCrossContaminateStackTrace(t *testing.T) {
	rawLogs := "JobA\tStep1\tFAIL: test_python_assert\n" +
		"JobA\tStep1\tAssertionError: 1 != 0\n" +
		"JobB\tStep2\tgoroutine 5 [running]:\n" +
		"JobB\tStep2\t\t/opt/go/panic.go:10 +0x1\n" +
		"JobB\tStep2\tFAIL\tgithub.com/pkg\t0.1s\n"

	jobMap := make(map[string]*FailedJobItem)
	var order []string
	scanLogLinesIntoMap(rawLogs, jobMap, &order)
	items := assembleJobItems(jobMap, order, rawLogs)

	assertNoCrossContamination(t, items)
}

func assertNoCrossContamination(t *testing.T, items []FailedJobItem) {
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
	if len(items[0].StackTrace) != 0 {
		t.Errorf("expected JobA to have empty stack trace, got: %s", items[0].StackTrace)
	}
	if !strings.Contains(items[1].StackTrace, "goroutine 5") {
		t.Errorf("expected JobB to retain its stack trace, got: %s", items[1].StackTrace)
	}
}

func TestIsIgnoredLogLine_FiltersRunnerAndCleanupNoise(t *testing.T) {
	noise := []string{
		"[command]/usr/bin/git config --local safe.directory",
		"Removing includeIf entries pointing to credentials",
		"includeif.gitdir:/home/runner/.git.path",
		"git-credentials-12345.config",
		"Checking boolean & enum compliance: [ 8/2370 ] 0.3%",
		"go: downloading github.com/pelletier/go-toml v1.9.5",
	}
	for _, n := range noise {
		if !isIgnoredLogLine(n) {
			t.Errorf("expected %q to be ignored", n)
		}
	}
}

func TestResolveContextLimit_TerminalStepError(t *testing.T) {
	if resolveContextLimit("Process completed with exit code 1.") != 0 {
		t.Errorf("expected 0 for exit code line")
	}
	if resolveContextLimit("exit status 1") != 0 {
		t.Errorf("expected 0 for exit status line")
	}
	if resolveContextLimit("FAILED (failures=1)") != 0 {
		t.Errorf("expected 0 for failures line")
	}
	if resolveContextLimit("--- FAIL: TestSomething") != 25 {
		t.Errorf("expected 25 for non-terminal failure marker")
	}
}

func TestMultiJobLogParsing_BoundedAndNoiseFree(t *testing.T) {
	rawLogs := "Lint Script Unit Tests\tRun tests\t2026-09-18T14:06:24Z FAIL: test_gofmt\n" +
		"Lint Script Unit Tests\tRun tests\t2026-09-18T14:06:24Z AssertionError: 1 != 0\n" +
		"Lint Script Unit Tests\tRun tests\t2026-09-18T14:06:24Z FAILED (failures=1)\n" +
		"Lint Script Unit Tests\tRun tests\t2026-09-18T14:06:24Z Process completed with exit code 1.\n" +
		"Lint Script Unit Tests\tRun tests\t2026-09-18T14:06:24Z [command]/usr/bin/git version\n" +
		"Lint Script Unit Tests\tRun tests\t2026-09-18T14:06:24Z Adding repository directory\n" +
		"Full Suite Guard\tRun suite\t2026-09-18T14:11:41Z goroutine 364 [running]:\n" +
		"Full Suite Guard\tRun suite\t2026-09-18T14:11:41Z \t/opt/go/testing.go:100 +0x10\n" +
		"Full Suite Guard\tRun suite\t2026-09-18T14:11:41Z panic(0x123)\n" +
		"Full Suite Guard\tRun suite\t2026-09-18T14:11:41Z \t/opt/go/panic.go:20 +0x5\n" +
		"Full Suite Guard\tRun suite\t2026-09-18T14:11:41Z FAIL\tgithub.com/pkg\t10s\n" +
		"Full Suite Guard\tRun suite\t2026-09-18T14:11:41Z ok  \tgithub.com/other\t0.01s\n" +
		"Full Suite Guard\tRun suite\t2026-09-18T14:11:41Z go: downloading github.com/foo\n" +
		"Boolean & Enum Linter\tRun python\t2026-09-18T14:07:45Z Checking boolean & enum compliance: [ 8/2370 ] 0.3%\n" +
		"Boolean & Enum Linter\tRun python\t2026-09-18T14:07:45Z ❌ FAILED: Found 2 violation(s):\n" +
		"Boolean & Enum Linter\tRun python\t2026-09-18T14:07:45Z   - /repo/file.go:35: Nested 'if' detected\n"

	jobs := ParseFailedLogLines(rawLogs)
	assertMultiJobCleanOutput(t, jobs)
}

func assertMultiJobCleanOutput(t *testing.T, jobs []FailedJobItem) {
	if len(jobs) < 2 {
		t.Fatalf("expected at least 2 jobs, got %d", len(jobs))
	}
	for _, l := range jobs[0].ErrorLines {
		if strings.Contains(l, "[command]") || strings.Contains(l, "Adding repository") {
			t.Errorf("Job 0 captured runner cleanup lines: %s", l)
		}
	}
	if len(jobs[0].StackTrace) != 0 {
		t.Errorf("Job 0 leaked stack trace from Job 1: %s", jobs[0].StackTrace)
	}
	if strings.Contains(jobs[1].StackTrace, "ok  ") || strings.Contains(jobs[1].StackTrace, "downloading") {
		t.Errorf("Job 1 stack trace contains post-termination lines: %s", jobs[1].StackTrace)
	}
}

func TestCachedLogFile35354190330_CompactPayloadSize(t *testing.T) {
	logPath := filepath.Join(os.Getenv("LOCALAPPDATA"), "gitmap-cli", "data", "pipeline", "35354190330.log")
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Skip("skipping local log file test if not present")
	}

	jobs := ParseFailedLogLines(string(content))
	sections := buildSectionsForRun(jobs)
	combined := formatCombinedSectionFailures(sections)

	assertCompactSectionsSize(t, combined, sections)
}

func buildSectionsForRun(jobs []FailedJobItem) []SectionFailure {
	var sections []SectionFailure
	run := FailedRunItem{WorkflowName: "CI", RunId: 35354190330, FailedJobs: jobs}
	for _, j := range jobs {
		sections = append(sections, buildSectionFailureFromRunJob(run, j))
	}

	return sections
}

func assertCompactSectionsSize(t *testing.T, combined string, sections []SectionFailure) {
	if len(combined) > 50000 {
		t.Errorf("combined sections size unexpectedly large: %d bytes (expected < 50KB)", len(combined))
	}
	if len(sections) > 0 && strings.Contains(sections[0].StackTrace, "goroutine") {
		t.Errorf("Section 0 (Lint Script) leaked goroutine from Full Suite Guard")
	}
}
