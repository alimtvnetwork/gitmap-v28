package cmd

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
