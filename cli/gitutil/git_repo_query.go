package gitutil

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// GetRepoHeadCommitSha retrieves the full HEAD commit SHA of the repository at repoPath.
func GetRepoHeadCommitSha(repoPath string) result.StringResult {
	res := RunGitInDir(repoPath, constants.GitRevParse, constants.GitHEAD)
	if res.IsFailure() {
		return result.Fail[string](res.AppError())
	}

	return result.Ok(strings.TrimSpace(string(res.Data)))
}

// GetRepoRemoteURL retrieves the configured URL for the specified remoteName at repoPath.
func GetRepoRemoteURL(repoPath, remoteName string) result.StringResult {
	key := fmt.Sprintf("remote.%s.url", remoteName)

	return GetRepoConfig(repoPath, key)
}

// GetRepoTrackingRemote retrieves the upstream tracking reference for the active branch.
func GetRepoTrackingRemote(repoPath string) result.StringResult {
	res := RunGitInDir(repoPath, constants.GitRevParse, constants.GitAbbrevRef, "@{u}")
	if res.IsFailure() {
		return result.Fail[string](res.AppError())
	}

	return result.Ok(strings.TrimSpace(string(res.Data)))
}

// HasUnpushedCommits reports whether the repository has commits not yet pushed to upstream.
func HasUnpushedCommits(repoPath string) result.BoolResult {
	res := RunGitInDir(repoPath, "cherry", "-v")
	if res.IsFailure() {
		return result.Fail[bool](res.AppError())
	}

	hasCommits := len(strings.TrimSpace(string(res.Data))) > 0

	return result.Ok(hasCommits)
}

// GetRepoDirtyFiles returns relative paths of uncommitted, modified, or untracked files in repoPath.
func GetRepoDirtyFiles(repoPath string) result.StringSliceResult {
	res := RunGitInDir(repoPath, constants.GitStatus, constants.GitPorcelainFlag)
	if res.IsFailure() {
		return result.FailSlice[string](res.AppError())
	}

	files := parsePorcelainDirtyLines(string(res.Data))

	return result.OkSlice(files)
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
func IsRepoDetached(repoPath string) result.BoolResult {
	res := RunGitInDir(repoPath, constants.GitSymbolicRef, constants.GitQuietFlag, constants.GitHEAD)
	if res.IsFailure() {
		return result.Ok(true)
	}

	isDetached := len(strings.TrimSpace(string(res.Data))) == 0

	return result.Ok(isDetached)
}
