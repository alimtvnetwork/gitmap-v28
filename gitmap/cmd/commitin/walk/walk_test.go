package walk

import (
	"strings"
	"testing"
)

// TestWalkFirstParentReturnsHydratedCommitsInOrder feeds a fake git
// runner two SHAs and verifies they come back oldest→newest, fully
// hydrated, with files split correctly.
func TestWalkFirstParentReturnsHydratedCommitsInOrder(t *testing.T) {
	restore := SetGitRunnerForTest(fakeRunner)
	defer restore()
	got, err := WalkFirstParent("/fake/repo")
	if err != nil {
		t.Fatalf("walk: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("got %d commits, want 2", len(got))
	}
	if got[0].Sha != "aaa" || got[0].OrderIndex != 1 {
		t.Fatalf("first commit wrong: %+v", got[0])
	}
	if got[1].Sha != "bbb" || got[1].OrderIndex != 2 {
		t.Fatalf("second commit wrong: %+v", got[1])
	}
	if len(got[0].Files) != 2 || got[0].Files[0] != "main.go" {
		t.Fatalf("file list wrong for aaa: %v", got[0].Files)
	}
	if got[0].OriginalMessage != "first\n\nbody-line" {
		t.Fatalf("message rejoin wrong: %q", got[0].OriginalMessage)
	}
}

// TestWalkFirstParentEmptyRepoReturnsNilSlice covers the new-repo path
// where `git rev-list HEAD` errors with "unknown revision".
func TestWalkFirstParentEmptyRepoReturnsNilSlice(t *testing.T) {
	restore := SetGitRunnerForTest(func(_, sub string, _ ...string) (string, error) {
		if sub == "rev-list" {
			return "fatal: ambiguous argument 'HEAD': unknown revision", sentinelEmptyError{}
		}
		return "", nil
	})
	defer restore()
	got, err := WalkFirstParent("/fake/empty")
	if err != nil {
		t.Fatalf("expected nil error on empty repo, got %v", err)
	}
	if got != nil {
		t.Fatalf("expected nil slice, got %v", got)
	}
}

// sentinelEmptyError mimics a git non-zero exit so isEmptyRepoError fires.
type sentinelEmptyError struct{}

func (sentinelEmptyError) Error() string { return "exit status 128" }

// fakeRunner returns canned output for the three commands the walker
// issues. SHAs are 3-char placeholders to keep the test compact.
func fakeRunner(_, sub string, args ...string) (string, error) {
	switch sub {
	case "rev-list":
		return "aaa\nbbb", nil
	case "show":
		return fakeShow(args)
	}
	return "", nil
}

// fakeShow dispatches between the metadata format and the file-list
// format based on the args the hydrator passes.
func fakeShow(args []string) (string, error) {
	if hasArg(args, "--name-only") {
		if hasArg(args, "aaa") {
			return "main.go\nREADME.md", nil
		}
		return "src/lib.go", nil
	}
	if hasArg(args, "aaa") {
		return strings.Join([]string{"alice", "alice@x", "2024-01-02T03:04:05+00:00", "2024-01-02T03:04:06+00:00", "first", "body-line"}, "\x1f"), nil
	}
	return strings.Join([]string{"bob", "bob@x", "2024-02-03T04:05:06+02:00", "2024-02-03T04:05:07+02:00", "second", ""}, "\x1f"), nil
}

func hasArg(args []string, want string) bool {
	for _, a := range args {
		if a == want {
			return true
		}
	}
	return false
}
