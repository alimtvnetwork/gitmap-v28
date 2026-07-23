package e2e

import (
	"testing"
	"time"
)

// TestHarnessSmoke proves the three harness primitives (NewRepo,
// Commit, LogFirstParent) work against the real `git` binary and that
// commit-date replication round-trips. This is NOT a commit-in
// pipeline test — Steps 9–12 cover that. This test only guards the
// harness itself so a regression here doesn't masquerade as a
// pipeline failure in the higher-level suites.
func TestHarnessSmoke(t *testing.T) {
	src := NewRepo(t, "src")
	when := time.Date(2024, 1, 15, 12, 30, 0, 0, time.UTC)
	sha := src.Commit("readme.md", "hello\n", "initial commit", when)
	if len(sha) != 40 {
		t.Fatalf("commit sha looks malformed: %q", sha)
	}
	log := src.LogFirstParent(t)
	if len(log) != 1 {
		t.Fatalf("got %d log entries, want 1", len(log))
	}
	if log[0].Sha != sha {
		t.Fatalf("log sha %q != commit sha %q", log[0].Sha, sha)
	}
	if log[0].Subject != "initial commit" {
		t.Fatalf("subject = %q, want %q", log[0].Subject, "initial commit")
	}
	// Date round-trip: %aI is RFC3339 strict — must contain our timestamp.
	if log[0].AuthorDate[:19] != "2024-01-15T12:30:00" {
		t.Fatalf("author date = %q, want prefix 2024-01-15T12:30:00", log[0].AuthorDate)
	}
	if log[0].CommitDate[:19] != "2024-01-15T12:30:00" {
		t.Fatalf("commit date = %q, want prefix 2024-01-15T12:30:00", log[0].CommitDate)
	}
	src.MustExist("readme.md")
}

// TestRunWithoutInputsReturnsParseSurface proves the e2e.Run helper
// wires orchestrator.Run cleanly: an empty inputs slice should still
// produce a deterministic non-zero exit code (input unusable) rather
// than panic. Real pipeline coverage starts in Step 9.
func TestRunWithoutInputsReturnsParseSurface(t *testing.T) {
	src := NewRepo(t, "dst")
	raw := NewRawArgs(src.Path) // no inputs → unusable
	res := Run(t, raw)
	if res.ExitCode == 0 {
		t.Fatalf("expected non-zero exit for empty inputs, got 0\nstderr=%s", res.Stderr)
	}
}
