package cmdpull

import (
	"os"
	"testing"
)

func TestExtractUnresolvedReposDirty(t *testing.T) {
	states := getMockDirtyStates()
	unresolved := ExtractUnresolvedRepos(states)
	if len(unresolved) != 2 {
		t.Fatalf("expected 2 unresolved repos, got %d", len(unresolved))
	}
	assertUnresolvedDirty(t, unresolved[0], "dirty-repo-1")
	assertUnresolvedDirty(t, unresolved[1], "dirty-repo-2")
}

func getMockDirtyStates() []*PullRepoState {
	return []*PullRepoState{
		{RepoName: "clean-repo", RepoPath: "/path/to/clean", Changes: "up-to-date", IsDirty: false},
		{RepoName: "dirty-repo-1", RepoPath: "/path/to/dirty1", Changes: "dirty", IsDirty: false},
		{RepoName: "dirty-repo-2", RepoPath: "/path/to/dirty2", Changes: "failed", IsDirty: true},
	}
}

func assertUnresolvedDirty(t *testing.T, item UnresolvedPullRepo, expectedName string) {
	if item.RepoName != expectedName || !item.IsDirty {
		t.Errorf("expected %s to be marked dirty, got %+v", expectedName, item)
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

func TestInteractivePromptEligibleCI(t *testing.T) {
	origCI := os.Getenv("CI")
	defer os.Setenv("CI", origCI)
	os.Setenv("CI", "true")
	if isInteractivePromptEligible() {
		t.Error("expected isInteractivePromptEligible to return false when CI=true")
	}
}

func TestInteractivePromptEligibleNonInt(t *testing.T) {
	origNonInt := os.Getenv("GITMAP_NON_INTERACTIVE")
	defer os.Setenv("GITMAP_NON_INTERACTIVE", origNonInt)
	os.Setenv("GITMAP_NON_INTERACTIVE", "1")
	if isInteractivePromptEligible() {
		t.Error("expected isInteractivePromptEligible to return false when GITMAP_NON_INTERACTIVE=1")
	}
}

func TestResolveDirtyItemTarget(t *testing.T) {
	itemWithName := UnresolvedPullRepo{RepoName: "my-repo", RepoPath: "/path/to/my-repo"}
	if resolveDirtyItemTarget(itemWithName) != "my-repo" {
		t.Errorf("expected my-repo, got %s", resolveDirtyItemTarget(itemWithName))
	}
	itemWithoutName := UnresolvedPullRepo{RepoName: "", RepoPath: "/path/to/my-other-repo"}
	if resolveDirtyItemTarget(itemWithoutName) != "my-other-repo" {
		t.Errorf("expected my-other-repo, got %s", resolveDirtyItemTarget(itemWithoutName))
	}
}

func TestResolveRepoDirtyFlag(t *testing.T) {
	sDirty := &PullRepoState{IsDirty: true}
	if !resolveRepoDirtyFlag(sDirty) {
		t.Errorf("expected true for s.IsDirty")
	}
	sClean := &PullRepoState{IsDirty: false, Changes: "up-to-date"}
	if resolveRepoDirtyFlag(sClean) {
		t.Errorf("expected false for clean state")
	}
}
