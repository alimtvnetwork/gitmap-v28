package cmd

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpending"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsends"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

func processSingleRepoSend(repo model.ScanRecord, opts cmdpending.SendsOptions, finalMsg string) cmdpending.RepoSendResultRecord {
	dir := repo.AbsolutePath
	if dir == "" {
		dir = repo.RelativePath
	}

	porcelain, _ := currentSendsGitExecutor(dir, "status", "--porcelain")
	untracked, modified, staged, _ := cmdpending.ParsePorcelainStatusLines(porcelain)
	filesStaged := untracked + modified + staged

	branch, _ := currentSendsGitExecutor(dir, "rev-parse", "--abbrev-ref", "HEAD")
	branch = strings.TrimSpace(branch)
	if branch == "" {
		branch = repo.Branch
	}
	if branch == "" {
		branch = "main"
	}

	relPath := repo.RelativePath
	if relPath == "" {
		relPath = repo.RepoName
	}

	isTargetAll := strings.EqualFold(opts.Target, "all")
	isClean := filesStaged == 0

	if isClean {
		return handleCleanRepoSend(dir, repo.RepoName, relPath, branch, isTargetAll, opts)
	}

	if opts.IsDryRun {
		return cmdpending.RepoSendResultRecord{
			RepoName:     repo.RepoName,
			RelativePath: relPath,
			Status:       "dry-run-simulated",
			Branch:       branch,
			FilesStaged:  filesStaged,
			IsSuccess:    true,
		}
	}

	repoCommitMsg := resolveRepoCommitMessage(opts, finalMsg, porcelain)

	return executeRealRepoCommitAndPush(dir, repo.RepoName, relPath, branch, filesStaged, repoCommitMsg, opts.IsPushed)
}

func resolveRepoCommitMessage(opts cmdpending.SendsOptions, finalMsg, porcelain string) string {
	if !isSemanticGroupingApplicable(opts) {
		return finalMsg
	}

	categories := cmdsends.CategorizeChangedFiles(porcelain)
	if len(categories) > 0 {
		return cmdsends.FormatMarkdownCommitBody(finalMsg, categories)
	}

	return finalMsg
}

func isSemanticGroupingApplicable(opts cmdpending.SendsOptions) bool {
	isTargetAll := strings.EqualFold(opts.Target, "all")
	if isTargetAll {
		return true
	}

	return isWIPTrigger(opts.RawMessage)
}

func isWIPTrigger(msg string) bool {
	clean := strings.TrimSpace(strings.ToLower(msg))
	isWIPExact := clean == "wip" || clean == "save changes"
	if isWIPExact {
		return true
	}
	hasWIPPrefix := strings.HasPrefix(clean, "wip:") || strings.HasPrefix(clean, "wip ") || strings.HasPrefix(clean, "wip-")
	if hasWIPPrefix {
		return true
	}

	return strings.Contains(clean, "wip") || strings.Contains(clean, "save changes")
}

func handleCleanRepoSend(dir, repoName, relPath, branch string, isTargetAll bool, opts cmdpending.SendsOptions) cmdpending.RepoSendResultRecord {
	if isTargetAll {
		return cmdpending.RepoSendResultRecord{
			RepoName:     repoName,
			RelativePath: relPath,
			Status:       "clean-skipped",
			Branch:       branch,
			FilesStaged:  0,
			IsSuccess:    true,
		}
	}

	unpushedCount := checkRepoUnpushedCount(dir)
	rec, hasHandled := tryHandleUnpushedCleanRepo(dir, repoName, relPath, branch, unpushedCount, opts.IsPushed, opts.IsDryRun)
	if hasHandled {
		return rec
	}

	return cmdpending.RepoSendResultRecord{
		RepoName:     repoName,
		RelativePath: relPath,
		Status:       "already-clean",
		Branch:       branch,
		FilesStaged:  0,
		IsSuccess:    true,
	}
}

func tryHandleUnpushedCleanRepo(dir, repoName, relPath, branch string, unpushedCount int, isPushed, isDryRun bool) (cmdpending.RepoSendResultRecord, bool) {
	if unpushedCount <= 0 || !isPushed {
		return cmdpending.RepoSendResultRecord{}, false
	}

	return handleUnpushedCleanRepo(dir, repoName, relPath, branch, isDryRun)
}

func handleUnpushedCleanRepo(dir, repoName, relPath, branch string, isDryRun bool) (cmdpending.RepoSendResultRecord, bool) {
	if isDryRun {
		return cmdpending.RepoSendResultRecord{
			RepoName:     repoName,
			RelativePath: relPath,
			Status:       "dry-run-simulated",
			Branch:       branch,
			FilesStaged:  0,
			IsSuccess:    true,
		}, true
	}

	_, errPush := currentSendsGitExecutor(dir, "push")
	if errPush != nil {
		return cmdpending.RepoSendResultRecord{
			RepoName:     repoName,
			RelativePath: relPath,
			Status:       "failed",
			Branch:       branch,
			FilesStaged:  0,
			IsSuccess:    false,
			ErrorMessage: fmt.Sprintf("git push failed: %v", errPush),
		}, true
	}

	headSha, _ := currentSendsGitExecutor(dir, "rev-parse", "--short", "HEAD")

	return cmdpending.RepoSendResultRecord{
		RepoName:     repoName,
		RelativePath: relPath,
		Status:       "pushed",
		Branch:       branch,
		HeadSHA:      strings.TrimSpace(headSha),
		FilesStaged:  0,
		IsSuccess:    true,
	}, true
}
