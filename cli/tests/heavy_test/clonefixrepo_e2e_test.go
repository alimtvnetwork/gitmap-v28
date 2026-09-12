package heavy_test

import (
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmd"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

const (
	testCloneFixGitInitCmd = "init"
	testCloneFixRepoName   = "gitmap-v28"
	testCloneFixFolderName = "gitmap"
	testCloneFixRemoteURL  = "https://github.com/alimtvnetwork/gitmap-v28.git"
	testCloneFixRemoteKey  = "remote.origin.url"
)

func TestResolveCloneFixRepoNameUsesRemote(t *testing.T) {
	dir := t.TempDir()
	runTestGit(t, dir, testCloneFixGitInitCmd)
	runTestGit(t, dir, constants.GitConfigCmd, testCloneFixRemoteKey, testCloneFixRemoteURL)

	got := cmd.ResolveCloneFixRepoName(dir)
	if got != testCloneFixRepoName {
		t.Fatalf("resolveCloneFixRepoName() = %q, want %s", got, testCloneFixRepoName)
	}
}

func TestResolveCloneFixRepoNameFallsBackToFolder(t *testing.T) {
	dir := filepath.Join(t.TempDir(), testCloneFixFolderName)

	got := cmd.ResolveCloneFixRepoName(dir)
	if got != testCloneFixFolderName {
		t.Fatalf("resolveCloneFixRepoName() = %q, want %s", got, testCloneFixFolderName)
	}
}

func runTestGit(t *testing.T, dir string, args ...string) error {
	t.Helper()
	c := exec.Command(constants.GitBin, args...)
	c.Dir = dir
	if out, err := c.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, out)
	}

	return nil
}
