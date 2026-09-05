package cmd

import (
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
