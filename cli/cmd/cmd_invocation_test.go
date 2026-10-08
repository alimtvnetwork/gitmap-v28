package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstall"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpipeline"
)

func TestInstallCommandInvocations(t *testing.T) {
	// 1. gitmap install (no args) -> should succeed cleanly with instructions
	err := cmdinstall.RunInstall([]string{})
	if err != nil {
		t.Errorf("expected runInstall with no args to return nil usage guide, got %v", err)
	}

	// 2. gitmap install --list
	err = cmdinstall.RunInstall([]string{"--list"})
	if err != nil {
		t.Errorf("expected runInstall --list to succeed, got %v", err)
	}

	// 3. gitmap install ls
	err = cmdinstall.RunInstall([]string{"ls"})
	if err != nil {
		t.Errorf("expected runInstall ls to succeed, got %v", err)
	}
}

func TestPipelineCommandInvocations(t *testing.T) {
	origRunner := cmdpipeline.PipelineAgyFixRunner
	cmdpipeline.PipelineAgyFixRunner = nil
	defer func() { cmdpipeline.PipelineAgyFixRunner = origRunner }()
	// 1. gitmap pipeline (default status)
	err := cmdpipeline.RunPipeline([]string{})
	if err != nil {
		t.Errorf("expected runPipeline() to succeed, got %v", err)
	}

	// 2. gitmap pipeline status --json
	err = cmdpipeline.RunPipeline([]string{"status", "--json"})
	if err != nil {
		t.Errorf("expected runPipeline status --json to succeed, got %v", err)
	}

	// 3. gitmap pipeline waittime & eta
	err = cmdpipeline.RunPipeline([]string{"waittime"})
	if err != nil {
		t.Errorf("expected runPipeline waittime to succeed, got %v", err)
	}

	err = cmdpipeline.RunPipeline([]string{"eta"})
	if err != nil {
		t.Errorf("expected runPipeline eta to succeed, got %v", err)
	}

	// 4. gitmap pipeline error-logs
	err = cmdpipeline.RunPipeline([]string{"error-logs"})
	if err != nil {
		t.Errorf("expected runPipeline error-logs to succeed, got %v", err)
	}

	// 5. gitmap pipeline error-logs --json
	err = cmdpipeline.RunPipeline([]string{"error-logs", "--json"})
	if err != nil {
		t.Errorf("expected runPipeline error-logs --json to succeed, got %v", err)
	}

	// 6. gitmap pipeline error-logs --tempfile
	tempFile := "test-ci-err.json"
	err = cmdpipeline.RunPipeline([]string{"error-logs", "--json", "--tempfile", tempFile})
	if err != nil {
		t.Errorf("expected runPipeline error-logs --tempfile to succeed, got %v", err)
	}

	defer os.Remove(filepath.Join(cmdpipeline.ResolveTempDir(), tempFile))

	// 7. gitmap pipeline help
	err = cmdpipeline.RunPipeline([]string{"help"})
	if err != nil {
		t.Errorf("expected runPipeline help to succeed, got %v", err)
	}

	// 8. gitmap pipeline logs
	err = cmdpipeline.RunPipeline([]string{"logs"})
	if err != nil {
		t.Errorf("expected runPipeline logs to succeed, got %v", err)
	}
}

func TestTopLevelErrorLogsAndLogs(t *testing.T) {
	origRunner := cmdpipeline.PipelineAgyFixRunner
	cmdpipeline.PipelineAgyFixRunner = nil
	defer func() { cmdpipeline.PipelineAgyFixRunner = origRunner }()

	// Top-level error-logs invocation
	err := cmdpipeline.RunPipeline([]string{"error-logs"})
	if err != nil {
		t.Errorf("expected top-level error-logs to succeed, got %v", err)
	}

	// Top-level logs invocation
	err = cmdpipeline.RunPipeline([]string{"logs"})
	if err != nil {
		t.Errorf("expected top-level logs to succeed, got %v", err)
	}

	// Top-level waittime invocation
	err = cmdpipeline.RunPipeline([]string{"waittime"})
	if err != nil {
		t.Errorf("expected top-level waittime to succeed, got %v", err)
	}
}

func TestTopLevelPipelineErrorsShortcut(t *testing.T) {
	origRunner := cmdpipeline.PipelineAgyFixRunner
	cmdpipeline.PipelineAgyFixRunner = nil
	defer func() { cmdpipeline.PipelineAgyFixRunner = origRunner }()

	if err := cmdpipeline.RunPipelineErrors([]string{"--json"}); err != nil {
		t.Errorf("expected runPipelineErrors --json to succeed, got %v", err)
	}
	if err := cmdpipeline.RunPipeline([]string{"pe", "--json"}); err != nil {
		t.Errorf("expected runPipeline pe --json to succeed, got %v", err)
	}
	if err := cmdpipeline.RunPipelineErrors([]string{"clear", "-y"}); err != nil {
		t.Errorf("expected runPipelineErrors clear -y to succeed, got %v", err)
	}
}

func TestTopLevelPipelineDetailsShortcut(t *testing.T) {
	if err := cmdpipeline.RunPipelineDetails([]string{"--json"}); err != nil {
		t.Errorf("expected runPipelineDetails --json to succeed, got %v", err)
	}
	if err := cmdpipeline.RunPipeline([]string{"pd", "--json"}); err != nil {
		t.Errorf("expected runPipeline pd --json to succeed, got %v", err)
	}
	if err := cmdpipeline.RunPipeline([]string{"details", "--json"}); err != nil {
		t.Errorf("expected runPipeline details --json to succeed, got %v", err)
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

	payload := cmdpipeline.BuildErrorLogsPayload("test-repo", []cmdpipeline.GhRunItem{})
	if !strings.Contains(payload.ErrorLogs, "sample local test error") {
		t.Errorf("expected local error fallback in payload, got %s", payload.ErrorLogs)
	}
}
