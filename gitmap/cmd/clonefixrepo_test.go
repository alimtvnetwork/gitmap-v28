package cmd

import (
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

const (
	TestCloneFixGitInitCmd = "init"
	TestCloneFixRepoName   = "gitmap-v28"
	TestCloneFixFolderName = "gitmap"
	TestCloneFixRemoteURL  = "https://github.com/alimtvnetwork/gitmap-v28.git"
	TestCloneFixRemoteKey  = "remote.origin.url"
)

func TestResolveCloneFixRepoNameUsesRemote(t *testing.T) {
	dir := t.TempDir()
	runTestGit(t, dir, TestCloneFixGitInitCmd)
	runTestGit(t, dir, constants.GitConfigCmd, TestCloneFixRemoteKey, TestCloneFixRemoteURL)

	got := resolveCloneFixRepoName(dir)
	if got != TestCloneFixRepoName {
		t.Fatalf("resolveCloneFixRepoName() = %q, want %s", got, TestCloneFixRepoName)
	}
}

func TestResolveCloneFixRepoNameFallsBackToFolder(t *testing.T) {
	dir := filepath.Join(t.TempDir(), TestCloneFixFolderName)

	got := resolveCloneFixRepoName(dir)
	if got != TestCloneFixFolderName {
		t.Fatalf("resolveCloneFixRepoName() = %q, want %s", got, TestCloneFixFolderName)
	}
}

func runTestGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command(constants.GitBin, args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, out)
	}
}
