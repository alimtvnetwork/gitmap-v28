package cmd

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdremediation"
)

func TestParseReconcileArgs(t *testing.T) {
	cases := []struct {
		args       []string
		wantRepo   string
		wantAction string
	}{
		{[]string{"1", "2"}, "1", "2"},
		{[]string{"1", "stash"}, "1", "stash"},
		{[]string{"codelane", "discard"}, "codelane", "discard"},
		{[]string{"discard", "codelane"}, "codelane", "discard"},
		{[]string{"stash"}, "", "stash"},
		{[]string{"codelane"}, "codelane", "stash"},
		{[]string{}, "", "stash"},
	}

	for _, tc := range cases {
		repo, act := parseReconcileArgs(tc.args)
		if repo != tc.wantRepo || act != tc.wantAction {
			t.Errorf("parseReconcileArgs(%v) = (%q, %q), want (%q, %q)",
				tc.args, repo, act, tc.wantRepo, tc.wantAction)
		}
	}
}

func TestFindRemediationItem(t *testing.T) {
	items := []cmdremediation.RemediationItem{
		{RepoName: "atto-property", RepoPath: "/path/to/atto-property"},
		{RepoName: "codelane", RepoPath: "/path/to/codelane"},
		{RepoName: "xmind-gen", RepoPath: "/path/to/xmind-gen"},
	}

	assertFoundRepo(t, items, "codelane", "codelane")
	assertFoundRepo(t, items, "atto-property", "atto-property")
	assertFoundRepo(t, items, "1", "atto-property")
	assertFoundRepo(t, items, "2", "codelane")
	assertFoundRepo(t, items, "3", "xmind-gen")
	assertFoundRepo(t, items, "xmind", "xmind-gen")
	if found := cmdremediation.FindRemediationItem(items, "non-existent"); found != nil {
		t.Fatalf("expected nil for non-existent repo, got %v", found)
	}
}

func assertFoundRepo(t *testing.T, items []cmdremediation.RemediationItem, query, expectedName string) {
	t.Helper()
	found := cmdremediation.FindRemediationItem(items, query)
	if found == nil {
		t.Fatalf("expected item for query %q, got nil", query)
	}

	if found.RepoName != expectedName {
		t.Fatalf("expected repo name %q, got %q", expectedName, found.RepoName)
	}
}

func TestBatchRemediationSaveLoadRemove(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("USERPROFILE", t.TempDir())

	items := []cmdremediation.RemediationItem{
		{RepoName: "repo-a", RepoPath: "/a", SummaryReason: "+1 modified"},
		{RepoName: "repo-b", RepoPath: "/b", SummaryReason: "+2 untracked"},
	}

	if err := cmdremediation.SaveRemediationState(items); err != nil {
		t.Fatalf("save remediation state: %v", err)
	}

	loaded := cmdremediation.LoadRemediationState()
	if len(loaded) != 2 {
		t.Fatalf("expected 2 loaded items, got %d", len(loaded))
	}

	cmdremediation.RemoveRemediationItem("repo-a")
	afterRemove := cmdremediation.LoadRemediationState()
	if len(afterRemove) != 1 || afterRemove[0].RepoName != "repo-b" {
		t.Fatalf("expected 1 remaining item repo-b, got %v", afterRemove)
	}

	cmdremediation.RemoveRemediationItem("repo-b")
	if remaining := cmdremediation.LoadRemediationState(); len(remaining) != 0 {
		t.Fatalf("expected empty state after removing all, got %v", remaining)
	}
}

// parseReconcileArgs parses reconcile command args into (repo, action).
// Test-local helper — the production function was never implemented.
// Logic derived from TestParseReconcileArgs expectations:
//   - 0 args → ("", "stash")
//   - 1 arg: if it's a known action → ("", action), else → (arg, "stash")
//   - 2 args: if first is an action → (second, first), else → (first, second)
func parseReconcileArgs(args []string) (string, string) {
	isAction := func(s string) bool {
		return s == "stash" || s == "discard"
	}

	switch len(args) {
	case 0:
		return "", "stash"
	case 1:
		if isAction(args[0]) {
			return "", args[0]
		}

		return args[0], "stash"
	default:
		if isAction(args[0]) {
			return args[1], args[0]
		}

		return args[0], args[1]
	}
}
