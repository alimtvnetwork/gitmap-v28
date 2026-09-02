package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/gitutil"
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
	items := []RemediationItem{
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
	if found := FindRemediationItem(items, "non-existent"); found != nil {
		t.Fatalf("expected nil for non-existent repo, got %v", found)
	}
}

func assertFoundRepo(t *testing.T, items []RemediationItem, query, expectedName string) {
	t.Helper()
	found := FindRemediationItem(items, query)
	if found == nil {
		t.Fatalf("expected item for query %q, got nil", query)
	}
	if found.RepoName != expectedName {
		t.Fatalf("expected repo name %q, got %q", expectedName, found.RepoName)
	}
}

func TestBatchRemediationSaveLoadRemove(t *testing.T) {
	orig := getRemediationStateFile()
	defer func() { _ = os.Remove(orig) }()

	items := []RemediationItem{
		{RepoName: "repo-a", RepoPath: "/a", SummaryReason: "+1 modified"},
		{RepoName: "repo-b", RepoPath: "/b", SummaryReason: "+2 untracked"},
	}
	if err := SaveRemediationState(items); err != nil {
		t.Fatalf("save remediation state: %v", err)
	}

	loaded := LoadRemediationState()
	if len(loaded) != 2 {
		t.Fatalf("expected 2 loaded items, got %d", len(loaded))
	}

	RemoveRemediationItem("repo-a")
	afterRemove := LoadRemediationState()
	if len(afterRemove) != 1 || afterRemove[0].RepoName != "repo-b" {
		t.Fatalf("expected 1 remaining item repo-b, got %v", afterRemove)
	}

	RemoveRemediationItem("repo-b")
	if remaining := LoadRemediationState(); len(remaining) != 0 {
		t.Fatalf("expected empty state after removing all, got %v", remaining)
	}
}

func TestParseRecipeIndex(t *testing.T) {
	recipes := []gitutil.RemediationRecipe{
		{Title: "Opt 1"},
		{Title: "Opt 2"},
		{Title: "Opt 3"},
	}
	testRecipeMap(t, recipes, map[string]int{
		"1": 0, "stash": 0, "s": 0,
		"2": 1, "wip": 1, "w": 1,
		"3": 2, "discard": 2, "clean": 2, "d": 2,
		"unknown": -1,
	})
}

func testRecipeMap(t *testing.T, recipes []gitutil.RemediationRecipe, cases map[string]int) {
	t.Helper()
	for input, expected := range cases {
		idx := parseRecipeIndex(input, recipes)
		if idx != expected {
			t.Errorf("parseRecipeIndex(%q) = %d, want %d", input, idx, expected)
		}
	}
}

func TestIsReconcileAllRequested(t *testing.T) {
	if !isReconcileAllRequested([]string{"--all"}) {
		t.Errorf("expected true for --all")
	}
	if !isReconcileAllRequested([]string{"-a"}) {
		t.Errorf("expected true for -a")
	}
	if isReconcileAllRequested([]string{"codelane"}) {
		t.Errorf("expected false for codelane")
	}
}

func TestReconcileWorkflowE2E(t *testing.T) {
	orig := getRemediationStateFile()
	defer func() { _ = os.Remove(orig) }()

	tempDir := t.TempDir()
	repoDir := initDummyGitRepoWithRemote(t, tempDir)

	diag := gitutil.InspectDirtyState(repoDir)
	recipes := gitutil.GenerateRemediationRecipes(repoDir, diag)
	item := RemediationItem{
		RepoName:      "sample-repo",
		RepoPath:      repoDir,
		SummaryReason: diag.SummaryReason,
		Recipes:       recipes,
	}
	_ = SaveRemediationState([]RemediationItem{item})

	err := runReconcileCmd([]string{"sample-repo", "discard"})
	if err != nil {
		t.Fatalf("runReconcileCmd failed: %v", err)
	}

	remaining := LoadRemediationState()
	if len(remaining) != 0 {
		t.Fatalf("expected 0 remaining items, got %d", len(remaining))
	}
}

func initDummyGitRepoWithRemote(t *testing.T, baseDir string) string {
	t.Helper()
	remoteDir := filepath.Join(baseDir, "remote.git")
	runCmdIn(baseDir, "git", "init", "--bare", remoteDir)

	repoDir := filepath.Join(baseDir, "local")
	runCmdIn(baseDir, "git", "clone", remoteDir, repoDir)
	runCmdIn(repoDir, "git", "config", "user.name", "Tester")
	runCmdIn(repoDir, "git", "config", "user.email", "tester@example.com")

	_ = os.WriteFile(filepath.Join(repoDir, "tracked.txt"), []byte("tracked"), 0644)
	runCmdIn(repoDir, "git", "add", "tracked.txt")
	runCmdIn(repoDir, "git", "commit", "-m", "init")
	runCmdIn(repoDir, "git", "push", "origin", "HEAD")

	_ = os.WriteFile(filepath.Join(repoDir, "untracked.txt"), []byte("dirty"), 0644)
	return repoDir
}

func runCmdIn(dir, name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	_ = cmd.Run()
}
