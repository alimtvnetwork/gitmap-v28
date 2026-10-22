package gitutil

import (
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// GetRepoUnpushedCommitCount returns the number of unpushed commits ahead of upstream tracking.
func GetRepoUnpushedCommitCount(repoPath string) result.IntResult {
	res := RunGitInDir(repoPath, "rev-list", "--count", "@{u}..HEAD")
	if res.IsFailure() {
		return result.Fail[int](res.AppError())
	}

	clean := strings.TrimSpace(string(res.Data))
	count, parseErr := strconv.Atoi(clean)
	if parseErr != nil {
		return result.Fail[int](apperror.WrapSimple(parseErr, "GetRepoUnpushedCommitCount"))
	}

	return result.Ok(count)
}

// GetRepoMergeBase finds the best common ancestor commit between two references.
func GetRepoMergeBase(repoPath, ref1, ref2 string) result.StringResult {
	res := RunGitInDir(repoPath, "merge-base", ref1, ref2)
	if res.IsFailure() {
		return result.Fail[string](res.AppError())
	}

	return result.Ok(strings.TrimSpace(string(res.Data)))
}
