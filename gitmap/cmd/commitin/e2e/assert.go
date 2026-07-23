package e2e

import (
	"os/exec"
	"strings"
	"testing"
)

// CommitNode is one row from `git log --format=%H%x09%an%x09%ae%x09%aI%x09%cI%x09%s`
// against the destination repo. Used for ordering + identity asserts.
type CommitNode struct {
	Sha         string
	AuthorName  string
	AuthorEmail string
	AuthorDate  string // ISO-8601 strict (matches %aI)
	CommitDate  string // ISO-8601 strict (matches %cI)
	Subject     string
}

// LogFirstParent returns the destination repo's first-parent log
// oldest→newest. Empty repo (no HEAD yet) returns nil, nil.
func (r *Repo) LogFirstParent(t *testing.T) []CommitNode {
	t.Helper()
	cmd := exec.Command("git", "log", "--first-parent", "--reverse",
		"--format=%H%x09%an%x09%ae%x09%aI%x09%cI%x09%s", "HEAD")
	cmd.Dir = r.Path
	out, err := cmd.Output()
	if err != nil {
		// Treat "unknown revision HEAD" (empty repo) as no commits.
		if strings.Contains(err.Error(), "exit status 128") {
			return nil
		}
		t.Fatalf("git log: %v", err)
	}
	return parseLogLines(string(out))
}

// AssertCommitCount fatals when the destination repo's commit count
// differs from `want`. Empty `want=0` is allowed.
func (r *Repo) AssertCommitCount(t *testing.T, want int) {
	t.Helper()
	got := len(r.LogFirstParent(t))
	if got != want {
		t.Fatalf("commit count: got %d, want %d", got, want)
	}
}

// AssertHasSubject fatals when no commit on first-parent HEAD has the
// exact `subject` line. Used to spot-check message-rule pipelines.
func (r *Repo) AssertHasSubject(t *testing.T, subject string) {
	t.Helper()
	for _, c := range r.LogFirstParent(t) {
		if c.Subject == subject {
			return
		}
	}
	t.Fatalf("no commit with subject %q on HEAD first-parent", subject)
}

// parseLogLines splits the tab-separated `git log` payload into nodes.
// Trailing empty line from `git log` is skipped silently.
func parseLogLines(raw string) []CommitNode {
	lines := strings.Split(strings.TrimRight(raw, "\n"), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return nil
	}
	out := make([]CommitNode, 0, len(lines))
	for _, ln := range lines {
		parts := strings.SplitN(ln, "\t", 6)
		if len(parts) < 6 {
			continue
		}
		out = append(out, CommitNode{
			Sha: parts[0], AuthorName: parts[1], AuthorEmail: parts[2],
			AuthorDate: parts[3], CommitDate: parts[4], Subject: parts[5],
		})
	}
	return out
}
