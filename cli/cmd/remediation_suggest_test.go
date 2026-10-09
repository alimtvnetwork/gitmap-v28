package cmd

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdremediation"
)

func TestFindRemediationItem_PathVariations(t *testing.T) {
	items := []cmdremediation.RemediationItem{
		{RepoName: "scripts-fixer", RepoPath: "D:\\work\\scripts-fixer"},
		{RepoName: "gitmap", RepoPath: "D:\\work\\gitmap"},
	}
	assertFoundRepo(t, items, ".\\scripts-fixer\\", "scripts-fixer")
	assertFoundRepo(t, items, "./scripts-fixer/", "scripts-fixer")
	assertFoundRepo(t, items, "scripts-fixer", "scripts-fixer")
	assertFoundRepo(t, items, "scripts-fixer\\", "scripts-fixer")
}
func TestFindRemediationSuggestions(t *testing.T) {
	items := []cmdremediation.RemediationItem{
		{RepoName: "scripts-fixer", RepoPath: "/path/to/scripts-fixer"},
		{RepoName: "gitmap-engine", RepoPath: "/path/to/gitmap-engine"},
	}
	suggs := cmdremediation.FindRemediationSuggestions(items, "script-fix")
	if len(suggs) == 0 {
		t.Fatalf("expected suggestions for 'script-fix', got empty")
	}
	if suggs[0] != "scripts-fixer" {
		t.Errorf("top suggestion = %q, want 'scripts-fixer'", suggs[0])
	}
}
