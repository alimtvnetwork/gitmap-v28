package osclean

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCleanTerminalHistory_DryRun(t *testing.T) {
	opts := TerminalCleanOptions{
		IsDryRun:     true,
		ShouldReseed: true,
	}

	res := CleanTerminalHistory(opts)
	if res.IsFailure() {
		t.Fatalf("unexpected failure: %v", res.AppError())
	}

	summary := res.Value
	if !summary.IsDryRun {
		t.Errorf("expected summary.IsDryRun to be true")
	}

	if len(summary.Shells) == 0 {
		t.Errorf("expected at least 1 shell stats returned")
	}
}

func TestIsShellSelected(t *testing.T) {
	if !isShellSelected("powershell", nil) {
		t.Errorf("expected empty only to select powershell")
	}
	if !isShellSelected("powershell", []string{"pwsh"}) {
		t.Errorf("expected pwsh alias to match powershell")
	}
	if !isShellSelected("bash", []string{"bash"}) {
		t.Errorf("expected bash to match bash")
	}
	if isShellSelected("bash", []string{"powershell"}) {
		t.Errorf("expected powershell only not to match bash")
	}
}

func TestProcessTerminalHistoryFile_DryRun(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test_history.txt")
	content := "echo 1\necho 2\necho 3\n"
	if err := os.WriteFile(testFile, []byte(content), 0600); err != nil {
		t.Fatalf("failed creating test file: %v", err)
	}

	stats := ShellCleanStats{
		Shell: "test",
		Label: "Test Shell",
	}

	processTerminalHistoryFile(&stats, testFile, []string{"gitmap"}, true, true, false)

	if stats.FilesCleared != 1 {
		t.Errorf("expected 1 file cleared, got %d", stats.FilesCleared)
	}
	if stats.BytesFreed != int64(len(content)) {
		t.Errorf("expected %d bytes freed, got %d", len(content), stats.BytesFreed)
	}
	if !stats.IsReseeded {
		t.Errorf("expected IsReseeded true")
	}

	// Verify file content was preserved in dry run
	data, err := os.ReadFile(testFile)
	if err != nil || string(data) != content {
		t.Errorf("expected content to be preserved in dry run")
	}
}

func TestProcessTerminalHistoryFile_Live(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test_history.txt")
	content := "echo 1\necho 2\n"
	if err := os.WriteFile(testFile, []byte(content), 0600); err != nil {
		t.Fatalf("failed creating test file: %v", err)
	}

	stats := ShellCleanStats{
		Shell: "test",
		Label: "Test Shell",
	}

	canonical := []string{"gitmap", "gitmap status"}
	processTerminalHistoryFile(&stats, testFile, canonical, false, true, false)

	if stats.FilesCleared != 1 {
		t.Errorf("expected 1 file cleared, got %d", stats.FilesCleared)
	}
	if stats.ReseedCount != 2 {
		t.Errorf("expected 2 reseeded commands, got %d", stats.ReseedCount)
	}

	data, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("failed reading reseeded file: %v", err)
	}
	expected := "gitmap\ngitmap status\n"
	if string(data) != expected {
		t.Errorf("expected content %q, got %q", expected, string(data))
	}
}
