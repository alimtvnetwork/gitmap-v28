package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpending"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func executeRealRepoCommitAndPush(dir, repoName, relPath, branch string, filesStaged int, finalMsg string, isPushed bool) cmdpending.RepoSendResultRecord {
	_, errAdd := currentSendsGitExecutor(dir, "add", "-A")
	if errAdd != nil {
		return cmdpending.RepoSendResultRecord{
			RepoName:     repoName,
			RelativePath: relPath,
			Status:       string(constants.SendResultStatusFailed),
			Branch:       branch,
			FilesStaged:  filesStaged,
			IsSuccess:    false,
			ErrorMessage: fmt.Sprintf("git add failed: %v", errAdd),
		}
	}

	_, errCommit := currentSendsGitExecutor(dir, "commit", "-m", finalMsg)
	if errCommit != nil {
		return cmdpending.RepoSendResultRecord{
			RepoName:     repoName,
			RelativePath: relPath,
			Status:       string(constants.SendResultStatusFailed),
			Branch:       branch,
			FilesStaged:  filesStaged,
			IsSuccess:    false,
			ErrorMessage: fmt.Sprintf("git commit failed: %v", errCommit),
		}
	}

	headShaOut, _ := currentSendsGitExecutor(dir, "rev-parse", "--short", "HEAD")
	headSha := strings.TrimSpace(headShaOut)

	if isPushed {
		return pushCommittedRepo(dir, repoName, relPath, branch, headSha, filesStaged)
	}

	return cmdpending.RepoSendResultRecord{
		RepoName:     repoName,
		RelativePath: relPath,
		Status:       string(constants.SendResultStatusCommitted),
		Branch:       branch,
		HeadSHA:      headSha,
		FilesStaged:  filesStaged,
		IsSuccess:    true,
	}
}

func pushCommittedRepo(dir, repoName, relPath, branch, headSha string, filesStaged int) cmdpending.RepoSendResultRecord {
	_, errPush := currentSendsGitExecutor(dir, "push", "origin", branch)
	if errPush != nil {
		return cmdpending.RepoSendResultRecord{
			RepoName:     repoName,
			RelativePath: relPath,
			Status:       string(constants.SendResultStatusFailed),
			Branch:       branch,
			HeadSHA:      headSha,
			FilesStaged:  filesStaged,
			IsSuccess:    false,
			ErrorMessage: fmt.Sprintf("git push failed: %v", errPush),
		}
	}

	return cmdpending.RepoSendResultRecord{
		RepoName:     repoName,
		RelativePath: relPath,
		Status:       string(constants.SendResultStatusPushed),
		Branch:       branch,
		HeadSHA:      headSha,
		FilesStaged:  filesStaged,
		IsSuccess:    true,
	}
}

func checkRepoUnpushedCount(dir string) int {
	revCount, errCount := currentSendsGitExecutor(dir, "rev-list", "--count", "@{u}..HEAD")
	if errCount != nil {
		return 0
	}

	count, _ := strconv.Atoi(strings.TrimSpace(revCount))

	return count
}
