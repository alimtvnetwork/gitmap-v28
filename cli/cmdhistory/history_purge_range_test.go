package cmdhistory

import (
	"strings"
	"testing"
)

func TestPurgeHelpNamesRangeAndDryRun(t *testing.T) {
	if !strings.Contains(purgeHelpText, "--dry-run") {
		t.Fatal("help missing dry-run")
	}
	if !strings.Contains(purgeHelpText, "--range") {
		t.Fatal("help missing range")
	}
	if !strings.Contains(purgeHelpText, ".cursor/skills/execute-parent-task-with-n-steps-v7") {
		t.Fatal("help missing example path")
	}
}

func TestExpandCommitRangeOneCommit(t *testing.T) {
	shas, err := expandCommitRange("../..", "f110b0526^..f110b0526")
	if err != nil {
		t.Fatal(err)
	}
	if len(shas) != 1 {
		t.Fatalf("count = %d", len(shas))
	}
	if !strings.HasPrefix(shas[0], "f110b0526") {
		t.Fatalf("sha = %s", shas[0])
	}
}

func TestExpandCommitRangeRejectsSha(t *testing.T) {
	_, err := expandCommitRange(".", "abc123")
	if err == nil {
		t.Fatal("expected an error")
	}
}
