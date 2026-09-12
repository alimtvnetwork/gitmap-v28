package cmdpipeline

import (
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
	if detailRes.IsFailure() || detailRes.Count() != 1 {
		t.Fatalf("expected 1 detail log, got %d (err: %v)", detailRes.Count(), detailRes.AppError())
	}
	details := detailRes.Data
	if !strings.Contains(details[0].RawLogs, "PASS: Test0") {
		t.Errorf("detail log missing raw PASS line: %s", details[0].RawLogs)
	}

	compactRes := pipeDb.QueryCompactErrorLogsByRunId(runId)
	if compactRes.IsFailure() || compactRes.Count() != 1 {
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
