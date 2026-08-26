package cmd

import (
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func TestPromptVersionMetadata(t *testing.T) {
	tempDir := t.TempDir()

	meta := model.PromptArchitectMetadata{
		Version:     "v2.0.0",
		InstalledAt: "2026-08-26T17:30:00Z",
		Status:      "active",
	}

	if errWrite := WritePromptArchitectMetadata(tempDir, meta); errWrite != nil {
		t.Fatalf("WritePromptArchitectMetadata failed: %v", errWrite)
	}

	readMeta, errRead := ReadPromptArchitectMetadata(tempDir)
	if errRead != nil || readMeta.Version != "v2.0.0" {
		t.Fatalf("ReadPromptArchitectMetadata unexpected: %+v (err: %v)", readMeta, errRead)
	}

	if !IsPromptArchitectInstalled(readMeta) {
		t.Fatal("expected IsPromptArchitectInstalled to return true")
	}

	sanitized := SanitizeTargetDirectory(tempDir)
	if !filepath.IsAbs(sanitized) {
		t.Fatalf("expected absolute path, got %s", sanitized)
	}
}
