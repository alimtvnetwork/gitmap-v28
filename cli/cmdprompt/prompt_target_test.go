package cmdprompt

import (
	"testing"
)

func TestPromptTargetSuite(t *testing.T) {
	tempDir := t.TempDir()

	targetRes := ResolvePromptTarget(tempDir)
	if targetRes.IsFailure() || targetRes.IsEmpty() {
		t.Fatalf("ResolvePromptTarget failed: %v", targetRes.AppError())
	}
	targets := targetRes.Data

	filtered := FilterPromptExclusions(targets, "non-existent")
	if len(filtered) != len(targets) {
		t.Fatal("unexpected filter result")
	}

	opts := ParsePromptArgs([]string{"install-prompts", "/tmp/repo", "--dry-run"})
	if opts.Action != "install" || !opts.IsDryRun || len(opts.Targets) != 1 {
		t.Fatalf("unexpected parsed options: %+v", opts)
	}
}
