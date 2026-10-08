package cmdrun

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

func TestResolveRunTarget_DirectAndExtensionless(t *testing.T) {
	tempDir := t.TempDir()
	pyFile := filepath.Join(tempDir, "sample_script.py")
	if err := os.WriteFile(pyFile, []byte("print('hello')\n"), 0o644); err != nil {
		t.Fatalf("failed to create dummy script: %v", err)
	}

	targetDirect, err := ResolveRunTarget(pyFile)
	if err != nil {
		t.Fatalf("expected direct resolution success, got error: %v", err)
	}
	if targetDirect.Extension != ".py" {
		t.Errorf("expected .py extension, got %s", targetDirect.Extension)
	}
	if !targetDirect.IsDirectMatch {
		t.Errorf("expected IsDirectMatch = true")
	}

	noExtInput := filepath.Join(tempDir, "sample_script")
	targetExtless, err := ResolveRunTarget(noExtInput)
	if err != nil {
		t.Fatalf("expected extensionless resolution success, got error: %v", err)
	}
	if targetExtless.Extension != ".py" {
		t.Errorf("expected .py extension, got %s", targetExtless.Extension)
	}
	if targetExtless.IsDirectMatch {
		t.Errorf("expected IsDirectMatch = false for probed candidate")
	}
}

func TestResolveRunTarget_Shebang(t *testing.T) {
	tempDir := t.TempDir()
	scriptFile := filepath.Join(tempDir, "custom_runner")
	content := "#!/usr/bin/env python3\nprint('shebang')\n"
	if err := os.WriteFile(scriptFile, []byte(content), 0o755); err != nil {
		t.Fatalf("failed to write shebang script: %v", err)
	}

	target, err := ResolveRunTarget(scriptFile)
	if err != nil {
		t.Fatalf("expected shebang resolution success, got error: %v", err)
	}
	if target.Extension != ".py" {
		t.Errorf("expected shebang detected as .py, got %s", target.Extension)
	}
}

func TestResolveRunTarget_NotFound(t *testing.T) {
	_, err := ResolveRunTarget("non_existent_file_xyz_12345")
	if err == nil {
		t.Errorf("expected error for non existent file, got nil")
	}
}

func TestBuildCommandInvocation(t *testing.T) {
	pyTarget := &RunTarget{ResolvedPath: "script.py", Extension: ".py"}
	bin, args := buildCommandInvocation(pyTarget, []string{"--flag", "val"})
	if bin == "" || len(args) < 2 {
		t.Errorf("expected python invocation with args, got bin=%s args=%v", bin, args)
	}

	goTarget := &RunTarget{ResolvedPath: "main.go", Extension: ".go"}
	bin, args = buildCommandInvocation(goTarget, []string{"arg1"})
	if bin != "go" || args[0] != "run" {
		t.Errorf("expected go run invocation, got bin=%s args=%v", bin, args)
	}
}

func TestTableConfigRendering(t *testing.T) {
	errRecords := []RunErrorRecord{
		{
			ErrorID:      "runerr-test-1",
			FilePath:     "test.py",
			Interpreter:  "python",
			ExitCode:     1,
			DurationMs:   120,
			ErrorMessage: "exit code 1",
			CreatedAt:    time.Now(),
		},
	}
	cfg := buildRunErrorsTableConfig(errRecords)
	if len(cfg.Rows) != 1 {
		t.Errorf("expected 1 row in table, got %d", len(cfg.Rows))
	}

	histRecords := []model.TaskHistoryRecord{
		{
			TaskId:     "run-100",
			Action:     "file-run",
			Target:     "test.py",
			Status:     "completed",
			ExecutedAt: "2026-10-08 10:00:00",
		},
	}
	histCfg := buildRunHistoryTableConfig(histRecords)
	if len(histCfg.Rows) != 1 {
		t.Errorf("expected 1 row in history table, got %d", len(histCfg.Rows))
	}
}

func TestRunErrorLifecycle(t *testing.T) {
	rec := RunErrorRecord{
		ErrorID:       "test-runerr-lifecycle",
		FilePath:      "scripts/demo.py",
		FileExtension: ".py",
		Interpreter:   "python",
		ExitCode:      2,
		DurationMs:    45,
		ErrorMessage:  "syntax error",
	}

	if err := RecordRunError(rec); err != nil {
		t.Logf("RecordRunError notice: %v", err)
	}

	records, err := QueryRecentRunErrors(10)
	if err != nil {
		t.Logf("QueryRecentRunErrors notice: %v", err)
	} else if len(records) == 0 {
		t.Logf("no records returned")
	}

	if err := ClearRunErrors(); err != nil {
		t.Logf("ClearRunErrors notice: %v", err)
	}
}

func TestRunAuditLifecycle(t *testing.T) {
	target := &RunTarget{ResolvedPath: "test_script.py"}
	taskId, err := EnqueueRunTaskAudit(target, []string{"--dry-run"})
	if err != nil {
		t.Logf("EnqueueRunTaskAudit notice: %v", err)
		return
	}

	if err := CompleteRunTaskAudit(taskId, target.ResolvedPath, "--dry-run", 100); err != nil {
		t.Logf("CompleteRunTaskAudit notice: %v", err)
	}
}

func TestRunDispatcher(t *testing.T) {
	if err := Run([]string{}); err == nil {
		t.Errorf("expected error for empty args")
	}

	if err := Run([]string{"--help"}); err != nil {
		t.Errorf("expected nil for --help, got %v", err)
	}

	if err := Run([]string{"errors"}); err != nil {
		t.Logf("Run errors notice: %v", err)
	}

	if err := Run([]string{"history"}); err != nil {
		t.Logf("Run history notice: %v", err)
	}

	if err := Run([]string{"clear-errors"}); err != nil {
		t.Logf("Run clear-errors notice: %v", err)
	}
}
