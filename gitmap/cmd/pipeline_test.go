package cmd

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
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
			DatabaseID: 101,
			Name:       "CI",
			Status:     "in_progress",
			Conclusion: "",
			CreatedAt:  "2026-08-30T16:00:00Z",
			HeadBranch: "main",
			HeadSha:    "123456",
			URL:        "https://example.com",
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
			DatabaseID: 202,
			Name:       "Build",
			Status:     "completed",
			Conclusion: "failure",
			URL:        "https://example.com/202",
		},
	}

	p := buildErrorLogsPayload("owner/repo", runs)
	if p.Conclusion != "failure" {
		t.Errorf("expected failure conclusion, got %s", p.Conclusion)
	}
	if p.RunID != 202 {
		t.Errorf("expected RunID 202, got %d", p.RunID)
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
