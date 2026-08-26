package cmd_test

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func TestPromptE2ESuite(t *testing.T) {
	tempDir := t.TempDir()

	// 1. Dry run execution
	res := cmd.ExecuteSinglePromptInstall(tempDir, true)
	if !res.IsSuccess {
		t.Fatalf("expected dry-run to succeed, got %+v", res)
	}

	// 2. Write metadata and read
	meta := model.PromptArchitectMetadata{
		Version:     "v2.0.0",
		InstalledAt: "2026-08-26T17:30:00Z",
		Status:      "active",
	}
	if err := cmd.WritePromptArchitectMetadata(tempDir, meta); err != nil {
		t.Fatalf("WritePromptArchitectMetadata failed: %v", err)
	}

	read, errRead := cmd.ReadPromptArchitectMetadata(tempDir)
	if errRead != nil || read.Version != "v2.0.0" {
		t.Fatalf("unexpected read metadata: %+v", read)
	}
}
