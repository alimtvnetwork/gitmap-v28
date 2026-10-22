package gitutil

import (
	"os/exec"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func TestBuildGitDirArgs_EmptyPath(t *testing.T) {
	args := []string{"status", "--short"}
	res := buildGitDirArgs("", args...)
	if len(res) != len(args) {
		t.Fatalf("buildGitDirArgs empty path len = %d, want %d", len(res), len(args))
	}

	if res[0] != "status" || res[1] != "--short" {
		t.Fatalf("buildGitDirArgs empty path unexpected output: %v", res)
	}
}

func TestBuildGitDirArgs_NonEmptyPath(t *testing.T) {
	repoPath := "/custom/repo/path"
	res := buildGitDirArgs(repoPath, "rev-parse", "HEAD")
	if len(res) != 4 {
		t.Fatalf("buildGitDirArgs non-empty path len = %d, want 4", len(res))
	}

	hasExpectedPrefix := res[0] == constants.GitDirFlag && res[1] == repoPath
	if !hasExpectedPrefix {
		t.Fatalf("buildGitDirArgs prefix = [%s, %s], want [-C, %s]", res[0], res[1], repoPath)
	}
}

func initTestGitRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	cmd := exec.Command(constants.GitBin, "init", dir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to init temp git repo: %v", err)
	}

	return dir
}

func setTestGitConfig(t *testing.T, dir, key, val string) {
	t.Helper()
	cmd := exec.Command(constants.GitBin, constants.GitDirFlag, dir, constants.GitConfigCmd, key, val)
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to set git config %s=%s: %v", key, val, err)
	}
}

func addTestOriginRemote(t *testing.T, dir, url string) {
	t.Helper()
	cmd := exec.Command(constants.GitBin, constants.GitDirFlag, dir, "remote", "add", constants.GitOrigin, url)
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to add origin remote: %v", err)
	}
}

func TestGetRepoConfig_ValidKey(t *testing.T) {
	dir := initTestGitRepo(t)
	setTestGitConfig(t, dir, "test.user.name", "gitmap-tester")
	val, err := GetRepoConfig(dir, "test.user.name")
	if err != nil {
		t.Fatalf("GetRepoConfig error: %v", err)
	}

	if val != "gitmap-tester" {
		t.Fatalf("GetRepoConfig = %q, want %q", val, "gitmap-tester")
	}
}

func TestGetRepoConfig_NonExistentKey(t *testing.T) {
	dir := initTestGitRepo(t)
	_, err := GetRepoConfig(dir, "nonexistent.config.key")
	hasErr := err != nil
	if !hasErr {
		t.Fatal("GetRepoConfig expected error for non-existent key, got nil")
	}
}

func TestGetRemoteOriginURL_NonGitDir(t *testing.T) {
	dir := t.TempDir()
	val := GetRemoteOriginURL(dir)
	hasEmptyVal := len(val) == 0
	if !hasEmptyVal {
		t.Fatalf("GetRemoteOriginURL on non-git dir = %q, want empty string", val)
	}
}

func TestGetRemoteOriginURL_ConfiguredRemote(t *testing.T) {
	dir := initTestGitRepo(t)
	testURL := "https://github.com/example/repo.git"
	setTestGitConfig(t, dir, constants.GitRemoteOrigin, testURL)
	val := GetRemoteOriginURL(dir)
	if val != testURL {
		t.Fatalf("GetRemoteOriginURL = %q, want %q", val, testURL)
	}
}

func TestReadRemoteOriginURL(t *testing.T) {
	dir := initTestGitRepo(t)
	testURL := "https://github.com/example/repo.git"
	setTestGitConfig(t, dir, constants.GitRemoteOrigin, testURL)
	readVal, err := ReadRemoteOriginURL(dir)
	if err != nil || readVal != testURL {
		t.Fatalf("ReadRemoteOriginURL = (%q, %v), want (%q, nil)", readVal, err, testURL)
	}
}

func TestSetRemoteOriginURL(t *testing.T) {
	dir := initTestGitRepo(t)
	addTestOriginRemote(t, dir, "https://old.url")
	newURL := "https://github.com/example/updated.git"
	if err := SetRemoteOriginURL(dir, newURL); err != nil {
		t.Fatalf("SetRemoteOriginURL failed: %v", err)
	}

	if got := GetRemoteOriginURL(dir); got != newURL {
		t.Fatalf("GetRemoteOriginURL after SetRemoteOriginURL = %q, want %q", got, newURL)
	}
}

func TestRemoteURL_Delegation(t *testing.T) {
	dir := initTestGitRepo(t)
	expectedURL := "https://github.com/example/delegated.git"
	setTestGitConfig(t, dir, constants.GitRemoteOrigin, expectedURL)
	gotURL, err := RemoteURL(dir)
	if err != nil {
		t.Fatalf("RemoteURL error: %v", err)
	}

	if gotURL != expectedURL {
		t.Fatalf("RemoteURL = %q, want %q", gotURL, expectedURL)
	}
}
