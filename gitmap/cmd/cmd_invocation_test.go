package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInstallCommandInvocations(t *testing.T) {
	// 1. gitmap install (no args) -> should succeed cleanly with instructions
	err := runInstall([]string{})
	if err != nil {
		t.Errorf("expected runInstall with no args to return nil usage guide, got %v", err)
	}

	// 2. gitmap install --list
	err = runInstall([]string{"--list"})
	if err != nil {
		t.Errorf("expected runInstall --list to succeed, got %v", err)
	}

	// 3. gitmap install ls
	err = runInstall([]string{"ls"})
	if err != nil {
		t.Errorf("expected runInstall ls to succeed, got %v", err)
	}
}

func TestPipelineCommandInvocations(t *testing.T) {
	// 1. gitmap pipeline (default status)
	err := runPipeline([]string{})
	if err != nil {
		t.Errorf("expected runPipeline() to succeed, got %v", err)
	}

	// 2. gitmap pipeline status --json
	err = runPipeline([]string{"status", "--json"})
	if err != nil {
		t.Errorf("expected runPipeline status --json to succeed, got %v", err)
	}

	// 3. gitmap pipeline waittime
	err = runPipeline([]string{"waittime"})
	if err != nil {
		t.Errorf("expected runPipeline waittime to succeed, got %v", err)
	}

	// 4. gitmap pipeline error-logs
	err = runPipeline([]string{"error-logs"})
	if err != nil {
		t.Errorf("expected runPipeline error-logs to succeed, got %v", err)
	}

	// 5. gitmap pipeline error-logs --json
	err = runPipeline([]string{"error-logs", "--json"})
	if err != nil {
		t.Errorf("expected runPipeline error-logs --json to succeed, got %v", err)
	}

	// 6. gitmap pipeline error-logs --tempfile
	tempFile := "test-ci-err.json"
	err = runPipeline([]string{"error-logs", "--json", "--tempfile", tempFile})
	if err != nil {
		t.Errorf("expected runPipeline error-logs --tempfile to succeed, got %v", err)
	}
	defer os.Remove(filepath.Join(resolveTempDir(), tempFile))

	// 7. gitmap pipeline help
	err = runPipeline([]string{"help"})
	if err != nil {
		t.Errorf("expected runPipeline help to succeed, got %v", err)
	}

	// 8. gitmap pipeline logs
	err = runPipeline([]string{"logs"})
	if err != nil {
		t.Errorf("expected runPipeline logs to succeed, got %v", err)
	}
}

func TestTopLevelErrorLogsAndLogs(t *testing.T) {
	// Top-level error-logs invocation
	err := runPipeline([]string{"error-logs"})
	if err != nil {
		t.Errorf("expected top-level error-logs to succeed, got %v", err)
	}

	// Top-level logs invocation
	err = runPipeline([]string{"logs"})
	if err != nil {
		t.Errorf("expected top-level logs to succeed, got %v", err)
	}

	// Top-level waittime invocation
	err = runPipeline([]string{"waittime"})
	if err != nil {
		t.Errorf("expected top-level waittime to succeed, got %v", err)
	}
}

func TestHelpFlagTrigger(t *testing.T) {
	// Ensure hasHelpFlag catches all variants including positional 'help'
	cases := [][]string{
		{"help"},
		{"--help"},
		{"-h"},
		{"install", "help"},
		{"clone", "--help"},
	}

	for _, c := range cases {
		if !hasHelpFlag(c) {
			t.Errorf("expected hasHelpFlag to return true for %v", c)
		}
	}

	nonHelp := [][]string{
		{"vscode"},
		{"status"},
		{"--json"},
	}

	for _, c := range nonHelp {
		if hasHelpFlag(c) {
			t.Errorf("expected hasHelpFlag to return false for %v", c)
		}
	}
}

func TestBuildErrorLogsPayloadWithFallback(t *testing.T) {
	_ = os.MkdirAll(".gitmap", 0755)
	_ = os.WriteFile(".gitmap/last_error.log", []byte("sample local test error"), 0644)
	defer os.Remove(".gitmap/last_error.log")

	payload := buildErrorLogsPayload("test-repo", []ghRunItem{})
	if !strings.Contains(payload.ErrorLogs, "sample local test error") {
		t.Errorf("expected local error fallback in payload, got %s", payload.ErrorLogs)
	}
}
