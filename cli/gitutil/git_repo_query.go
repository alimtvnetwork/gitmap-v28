package gitutil

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// GetRepoHeadCommitSha retrieves the full HEAD commit SHA of the repository at repoPath.
func GetRepoHeadCommitSha(repoPath string) (string, error) {
	out, err := RunGitInDir(repoPath, constants.GitRevParse, constants.GitHEAD)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

// GetRepoRemoteURL retrieves the configured URL for the specified remoteName at repoPath.
func GetRepoRemoteURL(repoPath, remoteName string) (string, error) {
	key := fmt.Sprintf("remote.%s.url", remoteName)

	return GetRepoConfig(repoPath, key)
}

// GetRepoTrackingRemote retrieves the upstream tracking reference for the active branch.
func GetRepoTrackingRemote(repoPath string) (string, error) {
	out, err := RunGitInDir(repoPath, constants.GitRevParse, constants.GitAbbrevRef, "@{u}")
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

// HasUnpushedCommits reports whether the repository has commits not yet pushed to upstream.
func HasUnpushedCommits(repoPath string) (bool, error) {
	out, err := RunGitInDir(repoPath, "cherry", "-v")
	if err != nil {
		return false, err
	}

	hasCommits := len(strings.TrimSpace(string(out))) > 0

	return hasCommits, nil
}

// GetRepoDirtyFiles returns relative paths of uncommitted, modified, or untracked files in repoPath.
func GetRepoDirtyFiles(repoPath string) ([]string, error) {
	out, err := RunGitInDir(repoPath, constants.GitStatus, constants.GitPorcelainFlag)
	if err != nil {
		return nil, err
	}

	return parsePorcelainDirtyLines(string(out)), nil
}

func parsePorcelainDirtyLines(output string) []string {
	var files []string
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		clean := strings.TrimSpace(line)
		if len(clean) > 3 {
			files = append(files, strings.TrimSpace(clean[3:]))
		}
	}

	return files
}

// IsRepoDetached reports whether the repository is in a detached HEAD state.
func IsRepoDetached(repoPath string) (bool, error) {
	out, err := RunGitInDir(repoPath, constants.GitSymbolicRef, constants.GitQuietFlag, constants.GitHEAD)
	if err != nil {
		return true, nil
	}

	isDetached := len(strings.TrimSpace(string(out))) == 0

	return isDetached, nil
}

