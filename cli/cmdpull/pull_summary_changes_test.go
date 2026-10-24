package cmdpull

import (
	"testing"
)

func TestHasIncomingChanges_Nil(t *testing.T) {
	if hasIncomingChanges(nil) {
		t.Fatalf("expected hasIncomingChanges(nil) to be false")
	}
}

func TestHasIncomingChanges_Identical(t *testing.T) {
	state := &PullRepoState{OldSHA: "abc", NewSHA: "abc"}
	if hasIncomingChanges(state) {
		t.Fatalf("expected identical SHAs to return false")
	}
}

func TestHasIncomingChanges_Different(t *testing.T) {
	state := &PullRepoState{OldSHA: "abc", NewSHA: "def"}
	if !hasIncomingChanges(state) {
		t.Fatalf("expected different SHAs to return true")
	}
}

func TestFilterGitStatusLines(t *testing.T) {
	raw := " M foo.go\n\n?? bar.go\n M .gitmap/test\n"
	lines := filterGitStatusLines(raw)
	if len(lines) != 2 {
		t.Fatalf("expected 2 filtered lines, got %d: %v", len(lines), lines)
	}
	if lines[0] != " M foo.go" || lines[1] != "?? bar.go" {
		t.Fatalf("unexpected filtered lines: %v", lines)
	}
}

func TestIsIgnoredStatusLine(t *testing.T) {
	if !isIgnoredStatusLine("") {
		t.Fatalf("expected empty string to be ignored")
	}
	if !isIgnoredStatusLine(" M .gitmap/release/latest.json") {
		t.Fatalf("expected .gitmap status to be ignored")
	}
	if isIgnoredStatusLine(" M cli/cmdpull/pull.go") {
		t.Fatalf("expected standard status line to not be ignored")
	}
}
