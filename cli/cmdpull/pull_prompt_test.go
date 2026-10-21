package cmdpull

import (
	"os"
	"testing"
)

func TestExtractUnresolvedReposDirty(t *testing.T) {
	states := []*PullRepoState{
		{
			RepoName: "clean-repo",
			RepoPath: "/path/to/clean",
			Changes:  "up-to-date",
			IsDirty:  false,
		},
		{
			RepoName: "dirty-repo-1",
			RepoPath: "/path/to/dirty1",
			Changes:  "dirty",
			IsDirty:  false,
		},
		{
			RepoName: "dirty-repo-2",
			RepoPath: "/path/to/dirty2",
			Changes:  "failed",
			IsDirty:  true,
		},
	}

	unresolved := ExtractUnresolvedRepos(states)
	if len(unresolved) != 2 {
		t.Fatalf("expected 2 unresolved repos, got %d", len(unresolved))
	}

	if unresolved[0].RepoName != "dirty-repo-1" || !unresolved[0].IsDirty {
		t.Errorf("expected dirty-repo-1 to be marked dirty, got %+v", unresolved[0])
	}
	if unresolved[1].RepoName != "dirty-repo-2" || !unresolved[1].IsDirty {
		t.Errorf("expected dirty-repo-2 to be marked dirty, got %+v", unresolved[1])
	}
}

func TestExtractUnresolvedReposClean(t *testing.T) {
	states := []*PullRepoState{
		{
			RepoName: "repo-a",
			RepoPath: "/path/to/a",
			Changes:  "up-to-date",
			Step:     PullStepTypeUpToDate,
		},
	}

	unresolved := ExtractUnresolvedRepos(states)
	if len(unresolved) != 0 {
		t.Fatalf("expected 0 unresolved repos, got %d", len(unresolved))
	}
}

func TestInteractivePromptEligibleNonInteractiveEnv(t *testing.T) {
	origCI := os.Getenv("CI")
	defer os.Setenv("CI", origCI)

	os.Setenv("CI", "true")
	if isInteractivePromptEligible() {
		t.Error("expected isInteractivePromptEligible to return false when CI=true")
	}

	os.Unsetenv("CI")
	origNonInt := os.Getenv("GITMAP_NON_INTERACTIVE")
	defer os.Setenv("GITMAP_NON_INTERACTIVE", origNonInt)

	os.Setenv("GITMAP_NON_INTERACTIVE", "1")
	if isInteractivePromptEligible() {
		t.Error("expected isInteractivePromptEligible to return false when GITMAP_NON_INTERACTIVE=1")
	}
}
