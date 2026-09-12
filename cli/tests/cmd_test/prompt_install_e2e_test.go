package cmd_test

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdprompt"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

func TestPromptE2ESuite(t *testing.T) {
	tempDir := t.TempDir()

	// 1. Dry run execution
	res := cmdprompt.ExecuteSinglePromptInstall(tempDir, true)
	if res.IsFailed() {
		t.Fatalf("expected dry-run to succeed, got %+v", res)
	}

	// 2. Write metadata and read
	meta := model.PromptArchitectMetadata{
		Version:     "v2.0.0",
		InstalledAt: "2026-08-26T17:30:00Z",
		Status:      "active",
	}

	if err := cmdprompt.WritePromptArchitectMetadata(tempDir, meta); err != nil {
		t.Fatalf("WritePromptArchitectMetadata failed: %v", err)
	}

	read, errRead := cmdprompt.ReadPromptArchitectMetadata(tempDir)
	if errRead != nil || read.Version != "v2.0.0" {
		t.Fatalf("unexpected read metadata: %+v", read)
	}
}
