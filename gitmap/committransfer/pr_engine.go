package committransfer

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/pterm/pterm"
)

// ProcessPR optionally creates and merges a PR for the given commit based on PRMode.
func ProcessPR(targetDir string, originalSubject, cleanedBody, shortSHA, prMode, newSHA string) error {
	if !shouldCreatePR(originalSubject, prMode) {
		return nil
	}
	branchName := fmt.Sprintf("pr/replay-%s-%d", shortSHA, time.Now().Unix())
	if err := createAndPushBranch(targetDir, branchName); err != nil {
		return err
	}
	if err := createGHPR(targetDir, branchName, cleanedBody); err != nil {
		return err
	}
	pterm.Success.Printf("Created PR for %s\n", shortSHA)
	if err := mergeAndRestore(targetDir, branchName); err != nil {
		return err
	}
	pterm.Success.Printf("Merged PR for %s\n", shortSHA)
	return nil
}

func shouldCreatePR(subject, prMode string) bool {
	if prMode == "all" {
		return true
	}
	lower := strings.ToLower(subject)
	isTag := strings.Contains(lower, "tag:") || strings.Contains(lower, "release ") || strings.Contains(lower, "version ")
	isRelease := strings.Contains(lower, "chore(release):") || strings.Contains(lower, "release v")
	if prMode == "tags" && isTag {
		return true
	}
	return prMode == "release" && isRelease
}

func createAndPushBranch(dir, branchName string) error {
	if err := exec.Command("git", "-C", dir, "checkout", "-b", branchName).Run(); err != nil {
		return apperror.Wrap(err, "ProcessPR: failed to checkout new PR branch", nil)
	}
	if err := exec.Command("git", "-C", dir, "push", "-u", "origin", branchName).Run(); err != nil {
		return apperror.Wrap(err, "ProcessPR: failed to push PR branch", nil)
	}
	return nil
}

func createGHPR(dir, branchName, body string) error {
	title := strings.Split(body, "\n")[0]
	cmd := exec.Command("gh", "pr", "create", "--title", title, "--body", body, "--head", branchName)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return apperror.Wrap(err, "ProcessPR: failed to create PR via gh", nil)
	}
	return nil
}

func mergeAndRestore(dir, branchName string) error {
	cmd := exec.Command("gh", "pr", "merge", branchName, "--squash", "--delete-branch")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return apperror.Wrap(err, "ProcessPR: failed to merge PR via gh", nil)
	}
	if err := exec.Command("git", "-C", dir, "checkout", "-").Run(); err != nil {
		return apperror.Wrap(err, "ProcessPR: failed to checkout previous branch", nil)
	}
	if err := exec.Command("git", "-C", dir, "pull").Run(); err != nil {
		return apperror.Wrap(err, "ProcessPR: failed to pull latest after PR merge", nil)
	}
	return nil
}
