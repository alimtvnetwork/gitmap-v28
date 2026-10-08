package cmdsync

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

// SyncSingleRepo executes the safe 6-stage synchronization ceremony on a single target repo.
func SyncSingleRepo(sourceRoot string, proj ProjectConfig, opts SyncOptions) RepoSyncResult {
	res := RepoSyncResult{
		Repo:   proj.Name,
		Status: "OK",
		Action: "No Changes",
	}

	repoPath := proj.Path
	branch, branchErr := gitutil.CurrentBranch(repoPath)
	if branchErr != nil || branch == "" {
		branch = "main"
	}

	// 1. Pull latest upstream changes
	_, _ = gitutil.ExecGitWithTimeout(30*time.Second, repoPath, "pull", "origin", branch, "--no-rebase")

	// 2. Query pre-change tag
	preTag := gitutil.GetLatestTag(repoPath)
	if preTag == "" {
		preTag = "v0.1.0"
	}
	res.PreTag = preTag
	res.PostTag = preTag

	// 3. Safety backup branch
	timestamp := time.Now().UTC().Format("20060102150405")
	backupBranch := fmt.Sprintf("backup/sync-%s", timestamp)

	if !opts.DryRun {
		createSafetyBackupBranch(repoPath, branch, backupBranch, !opts.NoPush)
	}

	// 4. Mirror assets adhering to boundary invariants
	stats, mirrorErr := MirrorAssets(sourceRoot, repoPath)
	if mirrorErr != nil {
		res.Status = "FAIL"
		res.Action = "Mirror Failed"
		res.Error = mirrorErr.Error()
		return res
	}

	res.FilesAdded = stats.Added
	res.FilesUpdated = stats.Updated
	res.FilesRemoved = stats.Removed

	// 5. Inspect dirty working tree status
	statusBytes, _ := gitutil.ExecGitWithTimeout(10*time.Second, repoPath, "status", "--porcelain")
	statusStr := strings.TrimSpace(string(statusBytes))

	if statusStr == "" {
		res.Action = "Already Synced"
		return res
	}

	if opts.DryRun {
		res.Action = "Dry Run"
		res.PostTag = bumpPatchTag(preTag)
		return res
	}

	// 6. Stage and commit atomically
	_, addErr := gitutil.ExecGitWithTimeout(30*time.Second, repoPath, "add", "-A")
	if addErr != nil {
		res.Status = "FAIL"
		res.Action = "Git Add Failed"
		res.Error = addErr.Error()
		return res
	}

	commitMsg := "chore(sync): mirror canonical prompts, skills, specs"
	_, commitErr := gitutil.ExecGitWithTimeout(30*time.Second, repoPath, "commit", "-m", commitMsg)
	if commitErr != nil {
		res.Status = "FAIL"
		res.Action = "Commit Failed"
		res.Error = commitErr.Error()
		return res
	}

	if opts.NoRelease {
		return pushNoRelease(repoPath, branch, res, opts.NoPush)
	}

	// Release ceremony with version bump
	postTag := bumpPatchTag(preTag)
	res.PostTag = postTag

	_, tagErr := gitutil.ExecGitWithTimeout(15*time.Second, repoPath, "tag", postTag)
	if tagErr != nil {
		res.Action = "Committed (Tag Failed)"
		return res
	}

	return pushReleaseWithTags(repoPath, branch, res, opts.NoPush)
}

func handlePushError(res RepoSyncResult, err error) RepoSyncResult {
	msg := err.Error()
	if strings.Contains(msg, "Permission to") || strings.Contains(msg, "denied") {
		res.Status = "SKIPPED"
		res.Action = "Skipped (External Repo)"
		res.Error = "Push permission denied"
		return res
	}

	res.Status = "FAIL"
	res.Action = "Push Failed"
	res.Error = msg
	return res
}

func bumpPatchTag(tag string) string {
	clean := strings.TrimSpace(tag)
	hasPrefixV := strings.HasPrefix(clean, "v")
	numPart := strings.TrimPrefix(clean, "v")

	parts := strings.Split(numPart, ".")
	if len(parts) < 3 {
		return "v0.1.1"
	}

	re := regexp.MustCompile(`^(\d+)`)
	m := re.FindStringSubmatch(parts[2])
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
	_, _ = gitutil.ExecGitWithTimeout(15*time.Second, repoPath, "checkout", "-b", backupBranch)
	if pushBackup {
		_, _ = gitutil.ExecGitWithTimeout(30*time.Second, repoPath, "push", "origin", backupBranch)
	}
	_, _ = gitutil.ExecGitWithTimeout(15*time.Second, repoPath, "checkout", branch)
}

func pushNoRelease(repoPath, branch string, res RepoSyncResult, noPush bool) RepoSyncResult {
	if noPush {
		res.Action = "Committed"
		return res
	}
	_, pushErr := gitutil.ExecGitWithTimeout(45*time.Second, repoPath, "push", "origin", branch)
	if pushErr != nil {
		return handlePushError(res, pushErr)
	}
	res.Action = "Committed & Pushed"
	return res
}

func pushReleaseWithTags(repoPath, branch string, res RepoSyncResult, noPush bool) RepoSyncResult {
	if noPush {
		res.Action = "Released Locally"
		return res
	}
	outBytes, pushErr := gitutil.ExecGitWithTimeout(60*time.Second, repoPath, "push", "origin", branch, "--tags")
	if pushErr != nil {
		errStr := string(outBytes) + " " + pushErr.Error()
		return handlePushError(res, fmt.Errorf("%s", errStr))
	}
	res.Action = "Released & Pushed"
	return res
}
