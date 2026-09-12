package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractTargetArg(t *testing.T) {
	if got := extractTargetArg([]string{}); got != "." {
		t.Fatalf("expected ., got %s", got)
	}

	if got := extractTargetArg([]string{"--flag", "myfolder"}); got != "myfolder" {
		t.Fatalf("expected myfolder, got %s", got)
	}
}

func TestResolveExplorerTarget(t *testing.T) {
	target, err := resolveExplorerTarget([]string{"."})
	if err != nil {
		t.Fatalf("resolveExplorerTarget failed: %v", err)
	}

	cwd, _ := os.Getwd()
	if filepath.Clean(target) != filepath.Clean(cwd) {
		t.Fatalf("expected %s, got %s", cwd, target)
	}
}

func TestCheckExplorerTargetType(t *testing.T) {
	tempDir := t.TempDir()
	isDir, err := checkExplorerTargetType(tempDir)
	if err != nil || !isDir {
		t.Fatalf("expected tempDir to be isDir=true, err=%v", err)
	}

	sampleFile := filepath.Join(tempDir, "sample.txt")
	_ = os.WriteFile(sampleFile, []byte("test"), 0644)
	isDirFile, fileErr := checkExplorerTargetType(sampleFile)
	if fileErr != nil || isDirFile {
		t.Fatalf("expected sampleFile to be isDir=false, err=%v", fileErr)
	}
}

func TestBuildExplorerCmd(t *testing.T) {
	tempDir := t.TempDir()
	cmd := buildExplorerCmd(tempDir, true)
	if cmd == nil || cmd.Path == "" {
		t.Fatal("expected non-nil explorer cmd")
	}
}
