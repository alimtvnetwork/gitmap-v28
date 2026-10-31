package cmdsync

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsync/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/lazyregex"
)

var leadingDigitsRe = lazyregex.New(`^(\d+)`)

// SyncSingleRepo executes the safe 6-stage synchronization ceremony on a single target repo.
func SyncSingleRepo(sourceRoot string, proj ProjectConfig, opts SyncOptions) RepoSyncResult {
	res := newInitialSyncResult(proj.Name)
	branch := resolveCurrentRepoBranch(proj.Path)
	if !opts.DryRun {
		pullLatestUpstream(proj.Path, branch)
	}

	preTag := resolvePreTag(proj.Path)
	res.PreTag = preTag
	res.PostTag = preTag

	performSafetyBackupIfMutating(proj.Path, branch, opts)

	stats, mirrorErr := MirrorAssetsWithScope(sourceRoot, proj.Path, opts.Scope, opts.DryRun)
	if mirrorErr != nil {
		return failSyncResult(res, constants.ActionMirrorFailed, mirrorErr)
	}

	applyMirrorStats(&res, stats)

	return processWorkingTreeChanges(proj.Path, branch, preTag, res, opts)
}

func newInitialSyncResult(repoName string) RepoSyncResult {
	return RepoSyncResult{
		Repo:   repoName,
		Status: constants.StatusOk,
		Action: constants.ActionNoChanges,
	}
}

func resolveCurrentRepoBranch(repoPath string) string {
	branchRes := gitutil.CurrentBranch(repoPath)
	branch := branchRes.Data
	if branchRes.IsFailure() || branch == "" {
		return constants.DefaultBranch
	}

	return branch
}

func pullLatestUpstream(repoPath, branch string) {
	_, _ = gitutil.ExecGitWithTimeout(8*time.Second, repoPath, constants.GitCmdPull, constants.GitRemoteOrigin, branch, constants.FlagNoRebase)
}

func resolvePreTag(repoPath string) string {
	preTag := gitutil.GetLatestTag(repoPath)
	if preTag == "" {
		return constants.DefaultInitialTag
	}

	return preTag
}

func performSafetyBackupIfMutating(repoPath, branch string, opts SyncOptions) {
	if !opts.DryRun {
		timestamp := time.Now().UTC().Format("20060102150405")
		backupBranch := fmt.Sprintf("backup/sync-%s", timestamp)
		createSafetyBackupBranch(repoPath, branch, backupBranch, !opts.NoPush)
	}
}

func applyMirrorStats(res *RepoSyncResult, stats MirrorStats) {
	res.FilesAdded = stats.Added
	res.FilesUpdated = stats.Updated
	res.FilesRemoved = stats.Removed
	res.FilesPreserved = stats.Preserved
	res.FilesConflicts = stats.Conflicts
	res.FilesUnchanged = stats.InSync
}

func failSyncResult(res RepoSyncResult, action string, err error) RepoSyncResult {
	res.Status = constants.StatusFail
	res.Action = action
	res.Error = err.Error()

	return res
}

func processWorkingTreeChanges(repoPath, branch, preTag string, res RepoSyncResult, opts SyncOptions) RepoSyncResult {
	statusBytes, _ := gitutil.ExecGitWithTimeout(10*time.Second, repoPath, constants.GitCmdStatus, constants.FlagPorcelain)
	statusStr := strings.TrimSpace(string(statusBytes))
	if statusStr == "" {
		res.Action = constants.ActionAlreadySynced

		return res
	}

	if opts.DryRun {
		res.Action = constants.ActionDryRun
		res.PostTag = bumpPatchTag(preTag)

		return res
	}

	return commitAndReleaseChanges(repoPath, branch, preTag, res, opts)
}

func commitAndReleaseChanges(repoPath, branch, preTag string, res RepoSyncResult, opts SyncOptions) RepoSyncResult {
	_, addErr := gitutil.ExecGitWithTimeout(30*time.Second, repoPath, constants.GitCmdAdd, constants.FlagAllShort)
	if addErr != nil {
		return failSyncResult(res, constants.ActionGitAddFailed, addErr)
	}

	_, commitErr := gitutil.ExecGitWithTimeout(30*time.Second, repoPath, constants.GitCmdCommit, constants.FlagMsgShort, constants.DefaultCommitMsg)
	if commitErr != nil {
		return failSyncResult(res, constants.ActionCommitFailed, commitErr)
	}

	if opts.NoRelease {
		return pushNoRelease(repoPath, branch, res, opts.NoPush)
	}

	return releaseCeremonyWithTag(repoPath, branch, preTag, res, opts.NoPush)
}

func releaseCeremonyWithTag(repoPath, branch, preTag string, res RepoSyncResult, noPush bool) RepoSyncResult {
	postTag := bumpPatchTag(preTag)
	res.PostTag = postTag

	_, tagErr := gitutil.ExecGitWithTimeout(15*time.Second, repoPath, constants.GitCmdTag, postTag)
	if tagErr != nil {
		res.Action = constants.ActionTagFailed

		return res
	}

	return pushReleaseWithTags(repoPath, branch, res, noPush)
}

func handlePushError(res RepoSyncResult, err error) RepoSyncResult {
	msg := err.Error()
	if strings.Contains(msg, "Permission to") || strings.Contains(msg, "denied") {
		res.Status = constants.StatusSkipped
		res.Action = constants.ActionSkippedExternal
		res.Error = "Push permission denied"

		return res
	}

	res.Status = constants.StatusFail
	res.Action = constants.ActionPushFailed
	res.Error = msg

	return res
}

func bumpPatchTag(tag string) string {
	clean := strings.TrimSpace(tag)
	hasPrefixV := strings.HasPrefix(clean, "v")
	numPart := strings.TrimPrefix(clean, "v")

	parts := strings.Split(numPart, ".")
	if len(parts) < 3 {
		return constants.DefaultFallbackTag
	}

	m := leadingDigitsRe.FindStringSubmatch(parts[2])
	if len(m) > 1 {
		parts[2] = incrementPatchSegment(m[1], parts[2])
	}

	res := strings.Join(parts, ".")
	if hasPrefixV {
		res = "v" + res
	}

	return res
}

func incrementPatchSegment(digits, fallback string) string {
	patchNum, err := strconv.Atoi(digits)
	if err != nil {
		return fallback
	}

	return strconv.Itoa(patchNum + 1)
}

func createSafetyBackupBranch(repoPath, branch, backupBranch string, pushBackup bool) {
	_, _ = gitutil.ExecGitWithTimeout(15*time.Second, repoPath, constants.GitCmdCheckout, constants.FlagBranchShort, backupBranch)
	if pushBackup {
		_, _ = gitutil.ExecGitWithTimeout(30*time.Second, repoPath, constants.GitCmdPush, constants.GitRemoteOrigin, backupBranch)
	}

	_, _ = gitutil.ExecGitWithTimeout(15*time.Second, repoPath, constants.GitCmdCheckout, branch)
}

func pushNoRelease(repoPath, branch string, res RepoSyncResult, noPush bool) RepoSyncResult {
	if noPush {
		res.Action = constants.ActionCommitted

		return res
	}

	_, pushErr := gitutil.ExecGitWithTimeout(45*time.Second, repoPath, constants.GitCmdPush, constants.GitRemoteOrigin, branch)
	if pushErr != nil {
		return handlePushError(res, pushErr)
	}

	res.Action = constants.ActionCommittedAndPushed

	return res
}

func pushReleaseWithTags(repoPath, branch string, res RepoSyncResult, noPush bool) RepoSyncResult {
	if noPush {
		res.Action = constants.ActionReleasedLocally

		return res
	}

	_, pushErr := gitutil.ExecGitWithTimeout(60*time.Second, repoPath, constants.GitCmdPush, constants.GitRemoteOrigin, branch, res.PostTag)
	if pushErr != nil {
		return retryPushBranchAndTag(repoPath, branch, res)
	}

	res.Action = constants.ActionReleasedAndPushed

	return res
}

func retryPushBranchAndTag(repoPath, branch string, res RepoSyncResult) RepoSyncResult {
	branchBytes, bErr := gitutil.ExecGitWithTimeout(45*time.Second, repoPath, constants.GitCmdPush, constants.GitRemoteOrigin, branch)
	if bErr != nil {
		errStr := string(branchBytes) + " " + bErr.Error()

		return handlePushError(res, apperror.NewExecutionError(errStr))
	}

	tagBytes, tErr := gitutil.ExecGitWithTimeout(30*time.Second, repoPath, constants.GitCmdPush, constants.GitRemoteOrigin, res.PostTag)
	if tErr != nil {
		errStr := string(tagBytes) + " " + tErr.Error()

		return handlePushError(res, apperror.NewExecutionError(errStr))
	}

	res.Action = constants.ActionReleasedAndPushed

	return res
}
