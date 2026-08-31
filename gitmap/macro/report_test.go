package macro

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestFormatReport_JSON(t *testing.T) {
	rep := NewExecutionReport("test-report", 2, time.Now())
	rep.Steps = append(rep.Steps, StepExecution{
		StepNum:        1,
		CommandLine:    "echo 1",
		WorkingDir:     "/tmp",
		Status:         "success",
		ExitCode:       0,
		ElapsedSeconds: 0.1,
		Logs:           []string{"line 1", "line 2"},
		ErrorLogs:      []string{},
	})
	rep.Finalize(time.Now(), false)

	out, err := FormatReport(rep, false)
	if err != nil {
		t.Fatalf("FormatReport JSON failed: %v", err)
	}
	if !strings.Contains(out, `"macro": "test-report"`) {
		t.Errorf("expected JSON to contain macro name, got: %s", out)
	}
	if !strings.Contains(out, `"logs"`) || !strings.Contains(out, `"line 1"`) {
		t.Errorf("expected JSON to contain logs array, got: %s", out)
	}
	if !strings.Contains(out, `"status": "success"`) {
		t.Errorf("expected status success in JSON, got: %s", out)
	}
}

func TestFormatReport_YAML(t *testing.T) {
	rep := NewExecutionReport("test-yaml-report", 1, time.Now())
	rep.Steps = append(rep.Steps, StepExecution{
		StepNum:        1,
		CommandLine:    "git status",
		WorkingDir:     "/work",
		Status:         "success",
		ExitCode:       0,
		ElapsedSeconds: 0.05,
		Logs:           []string{"branch main up to date"},
		ErrorLogs:      []string{},
	})
	rep.Finalize(time.Now(), false)

	out, err := FormatReport(rep, true)
	if err != nil {
		t.Fatalf("FormatReport YAML failed: %v", err)
	}
	if !strings.Contains(out, "macro: test-yaml-report") {
		t.Errorf("expected YAML to contain macro name, got: %s", out)
	}
	if !strings.Contains(out, "branch main up to date") {
		t.Errorf("expected YAML to contain logs, got: %s", out)
	}
	if !strings.Contains(out, "status: success") {
		t.Errorf("expected status in YAML, got: %s", out)
	}
}

func TestSaveReportToFile(t *testing.T) {
	tmpDir := t.TempDir()
	targetFile := filepath.Join(tmpDir, "sub", "report.json")

	savedPath, err := SaveReportToFile(targetFile, `{"test": true}`)
	if err != nil {
		t.Fatalf("SaveReportToFile failed: %v", err)
	}
	if _, statErr := os.Stat(savedPath); statErr != nil {
		t.Fatalf("target file was not created: %v", statErr)
	}
}
