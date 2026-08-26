package cmd

import (
	"testing"
)

func TestPromptTargetSuite(t *testing.T) {
	tempDir := t.TempDir()

	targets, err := ResolvePromptTarget(tempDir)
	if err != nil || len(targets) == 0 {
		t.Fatalf("ResolvePromptTarget failed: %v", err)
	}

	filtered := FilterPromptExclusions(targets, "non-existent")
	if len(filtered) != len(targets) {
		t.Fatal("unexpected filter result")
	}

	opts := parsePromptArgs([]string{"install-prompts", "/tmp/repo", "--dry-run"})
	if opts.Action != "install" || !opts.IsDryRun || len(opts.Targets) != 1 {
		t.Fatalf("unexpected parsed options: %+v", opts)
	}
}
