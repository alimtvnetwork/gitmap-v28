package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOSFixMirrors(t *testing.T) {
	tempDir := t.TempDir()
	sourcesFile := filepath.Join(tempDir, "sources.list")

	malaysiaMirror := "deb http://my.archive.ubuntu.com/ubuntu/ jammy main\n"
	if err := os.WriteFile(sourcesFile, []byte(malaysiaMirror), 0644); err != nil {
		t.Fatalf("failed to write test sources: %v", err)
	}

	if !HasRegionalMirrorGlitch(sourcesFile) {
		t.Fatal("expected glitch detected for my.archive.ubuntu.com")
	}

	if err := FixRegionalMirrors(sourcesFile); err != nil {
		t.Fatalf("FixRegionalMirrors failed: %v", err)
	}

	fixed, _ := os.ReadFile(sourcesFile)
	if string(fixed) != "deb http://archive.ubuntu.com/ubuntu/ jammy main\n" {
		t.Fatalf("unexpected fixed content: %s", string(fixed))
	}
}
