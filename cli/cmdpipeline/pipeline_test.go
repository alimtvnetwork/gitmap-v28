package cmdpipeline

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestParseSlugFromGitURL(t *testing.T) {
	tests := []struct {
		url  string
		want string
	}{
		{"https://github.com/alimtvnetwork/gitmap-v28.git", "alimtvnetwork/gitmap-v28"},
		{"git@github.com:alimtvnetwork/gitmap-v28.git", "alimtvnetwork/gitmap-v28"},
		{"https://github.com/owner/custom-repo", "owner/custom-repo"},
		{"owner/repo", "owner/repo"},
	}

	for _, tt := range tests {
		got := parseSlugFromGitURL(tt.url)
		if got != tt.want {
			t.Errorf("parseSlugFromGitURL(%q) = %q, want %q", tt.url, got, tt.want)
		}
	}
}

func TestBuildStatusPayload(t *testing.T) {
	runs := []ghRunItem{
		{
			DatabaseId: 101,
			Name:       "CI",
			Status:     "in_progress",
			Conclusion: "",
			CreatedAt:  "2026-08-30T16:00:00Z",
			HeadBranch: "main",
			HeadSha:    "123456",
			Url:        "https://example.com",
		},
	}

	p := buildStatusPayload("owner/repo", "v1.0.0", 2, runs)
	if !p.IsRunning {
		t.Errorf("expected IsRunning true, got false")
	}

	if p.PendingPRs != 2 {
		t.Errorf("expected 2 pending PRs, got %d", p.PendingPRs)
	}

	if p.LastTagRelease != "v1.0.0" {
		t.Errorf("expected v1.0.0 tag, got %s", p.LastTagRelease)
	}
}

func TestBuildErrorLogsPayload(t *testing.T) {
	runs := []ghRunItem{
		{
			DatabaseId: 202,
			Name:       "Build",
			Status:     "completed",
			Conclusion: "failure",
			Url:        "https://example.com/202",
		},
	}

	p := buildErrorLogsPayload("owner/repo", runs)
	if p.Conclusion != "failure" {
		t.Errorf("expected failure conclusion, got %s", p.Conclusion)
	}

	if p.RunId != 202 {
		t.Errorf("expected RunId 202, got %d", p.RunId)
	}
}

func TestWriteOrRenderErrorLogs_File(t *testing.T) {
	tmpDir := t.TempDir()
	outFile := filepath.Join(tmpDir, "error.log")

	payload := PipelineErrorLogsPayload{
		Repo:         "owner/repo",
		WorkflowName: "Test",
		Status:       "completed",
		Conclusion:   "failure",
		ErrorLogs:    "Fatal error: undefined constant",
	}

	err := writeOrRenderErrorLogs(ErrorLogOutputParams{
		Payload:  payload,
		IsJSON:   false,
		FilePath: outFile,
	})

	if err != nil {
		t.Fatalf("writeOrRenderErrorLogs failed: %v", err)
	}

	content, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	if !strings.Contains(string(content), "Fatal error") {
		t.Fatalf("unexpected content: %s", string(content))
	}
}

func TestPipelineHelp(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printPipelineHelp()

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	out := buf.String()

	if !strings.Contains(out, "gitmap pipeline [command]") {
		t.Fatalf("pipeline help missing usage string: %s", out)
	}
}

func TestCalculateETA_HistoricalAverage(t *testing.T) {
	now := time.Now().UTC().Format(time.RFC3339)
	runs := []ghRunItem{
		{
			Name:      "Release",
			Status:    "in_progress",
			CreatedAt: now,
		},
		{
			Name:       "Release",
			Status:     "completed",
			CreatedAt:  "2026-08-31T08:50:00Z",
			UpdatedAt:  "2026-08-31T08:51:40Z", // 100s
			Conclusion: "success",
		},
		{
			Name:       "Release",
			Status:     "completed",
			CreatedAt:  "2026-08-31T08:40:00Z",
			UpdatedAt:  "2026-08-31T08:41:20Z", // 80s
			Conclusion: "success",
		},
	}

	eta := calculateETA(runs)
	// Average is (100 + 80) / 2 = 90s, minus ~0s elapsed => >= 80s
	if eta < 20 {
		t.Fatalf("expected ETA >= 20, got %d", eta)
	}
}

func TestIsErrorLogsSubcmdPe(t *testing.T) {
	if !isErrorLogsSubcmd("pe") {
		t.Errorf("expected isErrorLogsSubcmd('pe') to be true")
	}
	if !isErrorLogsSubcmd("errors") {
		t.Errorf("expected isErrorLogsSubcmd('errors') to be true")
	}
}

func TestIsPipelineClearAction(t *testing.T) {
	if !isPipelineClearAction("clear") {
		t.Errorf("expected isPipelineClearAction('clear') to be true")
	}
	if !isPipelineClearAction("reset") {
		t.Errorf("expected isPipelineClearAction('reset') to be true")
	}
}

func TestTargetShaFiltering(t *testing.T) {
	runs := []ghRunItem{
		{HeadSha: "sha-clean", Status: "completed", Conclusion: "success", Name: "CI"},
		{HeadSha: "sha-stale", Status: "completed", Conclusion: "failure", Name: "CI"},
	}
	failed := collectRunsMatchingSha(runs, "sha-clean")
	if len(failed) != 1 || failed[0].Conclusion != "success" {
		t.Fatalf("expected 1 clean run for sha-clean, got %d", len(failed))
	}
}

func TestPipelineHelpMentionsPe(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	printPipelineHelp()
	_ = w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	if !strings.Contains(buf.String(), "gitmap pe") {
		t.Fatalf("pipeline help missing 'gitmap pe': %s", buf.String())
	}
}

func TestPipelineHelpMentionsPd(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	printPipelineHelp()
	_ = w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	if !strings.Contains(buf.String(), "gitmap pd") {
		t.Fatalf("pipeline help missing 'gitmap pd': %s", buf.String())
	}
}

func TestBuildPipelineDetailsPayload(t *testing.T) {
	run := ghRunItem{
		DatabaseId: 12345, Name: "Cross-Platform Build",
		HeadBranch: "main", HeadSha: "abc1234567890",
		Status: "completed", Conclusion: "failure",
	}
	jobs := []ghJobItem{
		{
			DatabaseId: 101, Name: "Ubuntu 22.04", Status: "completed", Conclusion: "success",
			Steps: []ghStepItem{{Name: "Build the app", Conclusion: "success"}},
		},
		{
			DatabaseId: 102, Name: "macOS aarch64", Status: "completed", Conclusion: "failure",
			Steps: []ghStepItem{{Name: "go test ./...", Conclusion: "failure"}},
		},
	}
	p := buildPipelineDetailsPayload("test/repo", run, jobs, true)
	if p.TotalCount != 2 || p.PassedCount != 1 || p.FailedCount != 1 {
		t.Fatalf("unexpected summary counts: %+v", p)
	}
	if !p.IsFromCache || len(p.Jobs) != 2 || p.Jobs[1].Step != "go test ./..." {
		t.Fatalf("unexpected details payload: %+v", p)
	}
}
