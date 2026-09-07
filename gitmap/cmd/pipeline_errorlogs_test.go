package cmd

import (
	"strings"
	"testing"
	"time"
)

func TestHandlePipelineErrorLogsWithTimeline(t *testing.T) {
	err := handlePipelineErrorLogs([]string{"-t", "--json"})
	if err != nil {
		t.Errorf("expected handlePipelineErrorLogs with -t and --json to succeed, got %v", err)
	}
}

func TestHandlePipelineErrorLogsWithCheckAndFix(t *testing.T) {
	errCheck := handlePipelineErrorLogs([]string{"--check", "--json"})
	if errCheck != nil {
		t.Errorf("expected handlePipelineErrorLogs with --check and --json to succeed, got %v", errCheck)
	}

	errFix := handlePipelineErrorLogs([]string{"--fix", "--json"})
	if errFix != nil {
		t.Errorf("expected handlePipelineErrorLogs with --fix and --json to succeed, got %v", errFix)
	}
}

func TestPipelineErrorLogsPayloadRerunETA(t *testing.T) {
	now := time.Now().UTC()
	runs := []ghRunItem{
		{
			DatabaseId: 201,
			Name:       "CI",
			Status:     "completed",
			Conclusion: "failure",
			CreatedAt:  now.Add(-60 * time.Second).Format(time.RFC3339),
			UpdatedAt:  now.Add(-30 * time.Second).Format(time.RFC3339),
		},
		{
			DatabaseId: 202,
			Name:       "CI",
			Status:     "completed",
			Conclusion: "success",
			CreatedAt:  now.Add(-200 * time.Second).Format(time.RFC3339),
			UpdatedAt:  now.Add(-100 * time.Second).Format(time.RFC3339), // 100s duration
		},
	}

	payload := buildErrorLogsPayload("alimtvnetwork/gitmap-v28", runs)
	payload.RerunEtaSeconds = calculateAverageDuration(runs, payload.WorkflowName)

	if payload.Conclusion != "failure" {
		t.Errorf("expected failure conclusion, got %s", payload.Conclusion)
	}
	if payload.RerunEtaSeconds != 100 {
		t.Errorf("expected rerun ETA 100s, got %d", payload.RerunEtaSeconds)
	}
}

func TestRunInternalCICDChecks(t *testing.T) {
	results := runInternalCICDChecks(false)
	if len(results) == 0 {
		t.Errorf("expected internal CI/CD checks to return results, got 0")
	}

	hasGofmtCheck := false
	for _, r := range results {
		if r.Name == "gofmt formatting" {
			hasGofmtCheck = true
			break
		}
	}
	if !hasGofmtCheck {
		t.Errorf("expected gofmt formatting probe to be present in results")
	}
}

func TestPipelineDispatcherErrorlogsAlias(t *testing.T) {
	err := runPipeline([]string{"errorlogs", "--json"})
	if err != nil {
		t.Errorf("expected runPipeline errorlogs --json to succeed, got %v", err)
	}

	errTimeline := runPipeline([]string{"errorlogs", "-t", "--json"})
	if errTimeline != nil {
		t.Errorf("expected runPipeline errorlogs -t --json to succeed, got %v", errTimeline)
	}
}

func TestParseFailedLogLines(t *testing.T) {
	raw := "macos-latest / test\trun tests\t2026-09-07T17:07:10.7297330Z --- FAIL: TestSample (0.00s)\n" +
		"macos-latest / test\trun tests\t2026-09-07T17:07:10.7299850Z     sample_test.go:12: Expected true to be false\n" +
		"macos-latest / test\trun tests\t2026-09-07T17:07:10.7300340Z FAIL\n" +
		"macos-latest / test\trun tests\t2026-09-07T17:07:10.7300760Z FAIL\tgithub.com/alimtvnetwork/gitmap-v28/gitmap/power\t0.020s\n" +
		"windows-latest / smoke\tcfr cg\t2026-09-07T17:12:01.2500124Z ##[error]cfr cg exited 10\n"

	jobs := ParseFailedLogLines(raw)
	if len(jobs) != 2 {
		t.Fatalf("expected 2 failed jobs, got %d", len(jobs))
	}

	if jobs[0].JobName != "macos-latest / test" || jobs[0].StepName != "run tests" {
		t.Errorf("unexpected job 0: %+v", jobs[0])
	}
	if jobs[0].FailureSummary != "sample_test.go:12: Expected true to be false" {
		t.Errorf("unexpected summary 0: %s", jobs[0].FailureSummary)
	}

	if jobs[1].JobName != "windows-latest / smoke" || jobs[1].StepName != "cfr cg" {
		t.Errorf("unexpected job 1: %+v", jobs[1])
	}
	if jobs[1].FailureSummary != "cfr cg exited 10" {
		t.Errorf("unexpected summary 1: %s", jobs[1].FailureSummary)
	}
}

func TestCollectFailedRuns(t *testing.T) {
	runs := []ghRunItem{
		{DatabaseId: 1, HeadSha: "abc", Conclusion: "failure"},
		{DatabaseId: 2, HeadSha: "abc", Conclusion: "failure"},
		{DatabaseId: 3, HeadSha: "abc", Conclusion: "success"},
		{DatabaseId: 4, HeadSha: "def", Conclusion: "failure"},
	}

	collected := collectFailedRuns(runs)
	if len(collected) != 2 {
		t.Fatalf("expected 2 collected runs matching first failure sha abc, got %d", len(collected))
	}
	if collected[0].DatabaseId != 1 || collected[1].DatabaseId != 2 {
		t.Errorf("unexpected collected runs: %+v", collected)
	}
}

func TestFormatAggregatedErrorLogs(t *testing.T) {
	failedRuns := []FailedRunItem{
		{
			WorkflowName: "Cross-Platform Build",
			RunId:        123,
			Url:          "https://github.com/example/123",
			FailedJobs: []FailedJobItem{
				{
					JobName:        "macos-latest",
					StepName:       "test",
					FailureSummary: "panic: nil pointer",
					ErrorLines:     []string{"panic: nil pointer", "exit status 2"},
				},
			},
		},
	}

	formatted := formatAggregatedErrorLogs(failedRuns)
	if !strings.Contains(formatted, "==> Failed Run: Cross-Platform Build (#123)") {
		t.Errorf("missing header in formatted logs: %s", formatted)
	}
	if !strings.Contains(formatted, "Summary: panic: nil pointer") {
		t.Errorf("missing summary in formatted logs: %s", formatted)
	}
}
