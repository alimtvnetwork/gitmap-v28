package cmdpipeline

import (
	"os"
	"strings"
	"testing"
)

func TestParsePipelineAIDelay(t *testing.T) {
	tests := []struct {
		args        []string
		wantDelay   int
		wantSubArgs []string
	}{
		{
			args:        []string{"status"},
			wantDelay:   20,
			wantSubArgs: []string{"status"},
		},
		{
			args:        []string{"status", "-t", "35", "--json"},
			wantDelay:   35,
			wantSubArgs: []string{"status", "--json"},
		},
		{
			args:        []string{"status", "-t", "5"},
			wantDelay:   20,
			wantSubArgs: []string{"status"},
		},
		{
			args:        []string{"eta", "--delay", "50"},
			wantDelay:   50,
			wantSubArgs: []string{"eta"},
		},
		{
			args:        []string{"etc", "--time", "65"},
			wantDelay:   65,
			wantSubArgs: []string{"etc"},
		},
	}

	for _, tt := range tests {
		gotDelay, gotSubArgs := parsePipelineAIDelay(tt.args)
		if gotDelay != tt.wantDelay {
			t.Errorf("parsePipelineAIDelay(%v) delay = %d, want %d", tt.args, gotDelay, tt.wantDelay)
		}

		if len(gotSubArgs) != len(tt.wantSubArgs) {
			t.Errorf("parsePipelineAIDelay(%v) subArgs = %v, want %v", tt.args, gotSubArgs, tt.wantSubArgs)
		}
	}
}

func TestExtractPipelineAISubcmd(t *testing.T) {
	sub, args := extractPipelineAISubcmd([]string{"errors", "--json"})
	if sub != "errors" || len(args) != 1 || args[0] != "--json" {
		t.Errorf("expected extractPipelineAISubcmd(errors) to yield errors and [--json], got: %s, %v", sub, args)
	}

	sub2, args2 := extractPipelineAISubcmd([]string{"-t", "30"})
	if sub2 != "status" || len(args2) != 2 {
		t.Errorf("expected flag prefix to default to status, got: %s", sub2)
	}
}

func TestResolvePipelineAINextCommand_Errors(t *testing.T) {
	payload := PipelineStatusPayload{
		IsRunning:  true,
		EtaSeconds: 120,
		HasErrors:  true,
	}

	resolvePipelineAINextCommand(&payload)
	if payload.NextAiCommand != "gitmap pipeline fix agy" {
		t.Errorf("expected NextAiCommand to be 'gitmap pipeline fix agy', got: %s", payload.NextAiCommand)
	}
	if !payload.IsStopWaiting {
		t.Errorf("expected IsStopWaiting to be true when errors present")
	}
	if payload.RecommendedAction != "fix_errors" {
		t.Errorf("expected RecommendedAction to be 'fix_errors', got: %s", payload.RecommendedAction)
	}
}

func TestResolvePipelineAINextCommand_RunningClean(t *testing.T) {
	payload := PipelineStatusPayload{
		IsRunning:  true,
		EtaSeconds: 75,
		HasErrors:  false,
	}

	resolvePipelineAINextCommand(&payload)
	if payload.NextAiCommand != "gitmap pipeline-ai status -t 75" {
		t.Errorf("expected wait command, got: %s", payload.NextAiCommand)
	}
	if payload.IsStopWaiting {
		t.Errorf("expected IsStopWaiting to be false when running cleanly")
	}
	if payload.RecommendedAction != "wait" {
		t.Errorf("expected RecommendedAction to be 'wait', got: %s", payload.RecommendedAction)
	}
}

func TestHasAnyFailingJobOrStep(t *testing.T) {
	cleanJobs := []ghJobItem{
		{Status: "completed", Conclusion: "success"},
		{Status: "in_progress", Conclusion: ""},
	}
	if hasAnyFailingJobOrStep(cleanJobs) {
		t.Errorf("expected false for cleanJobs")
	}

	failedStepJobs := []ghJobItem{
		{
			Status: "in_progress",
			Steps: []ghStepItem{
				{Name: "Lint", Status: "completed", Conclusion: "failure"},
			},
		},
	}
	if !hasAnyFailingJobOrStep(failedStepJobs) {
		t.Errorf("expected true when step has failed")
	}
}

func TestTruncateErrorLines(t *testing.T) {
	text := "line1\nline2\nline3\nline4\nline5"
	truncated := truncateErrorLines(text, 3)
	lines := strings.Split(truncated, "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
}

func TestRunPipelineAI_SkipDelay(t *testing.T) {
	_ = os.Setenv("GITMAP_SKIP_DELAY", "1")
	defer os.Unsetenv("GITMAP_SKIP_DELAY")

	err := runPipelineAI([]string{"status", "-t", "25", "--json"})
	if err != nil {
		t.Fatalf("runPipelineAI failed with GITMAP_SKIP_DELAY=1: %v", err)
	}
}
