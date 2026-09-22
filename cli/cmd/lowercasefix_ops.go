// Package cmd — lowercasefix_ops.go implements scanning and two-step git renaming.
package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// ExecuteLowerCaseFix orchestrates scanning, renaming, and committing.
func ExecuteLowerCaseFix(opts LowerCaseFixOptions) error {
	cwd, _ := os.Getwd()
	isGit := isGitRepository()
	pairs, totalScanned, err := findRenameCandidates(cwd, opts)
	if err != nil {
		return err
	}
	for i := range pairs {
		pairs[i].IsGitTracked = isGit && isFileTrackedInGit(pairs[i].OldPath)
	}

	renderRenameHeader(opts, isGit, len(pairs), cwd)
	if len(pairs) == 0 {
		fmt.Printf("%s No uppercase files matching patterns found to rename.\n\n", constants.ColorGreen+"✓"+constants.ColorReset)

		return nil
	}

	if opts.IsDryRun {
		for i, p := range pairs {
			prefix := "git mv"
			if !p.IsGitTracked {
				prefix = "fs rename"
			}
			fmt.Printf("  [%d/%d] %s → %s (%s preview)\n", i+1, len(pairs), p.OldBase, p.NewBase, prefix)
		}
		renderDryRunNotice(len(pairs))

		return nil
	}

	return applyRenamesAndCommit(pairs, opts, totalScanned, isGit)
}

func applyRenamesAndCommit(pairs []RenamePair, opts LowerCaseFixOptions, scanned int, isGit bool) error {
	renamedCount := 0
	for i, p := range pairs {
		err := performSafeGitRename(p)
		renderFileStepLog(i+1, len(pairs), p, err == nil, err)
		if err != nil {
			return err
		}
		renamedCount++
	}

	commitSHA := ""
	if isGit && !opts.IsNoCommit {
		sha, commitErr := commitRenames(pairs, opts.CommitMessage)
		if commitErr != nil {
			return commitErr
		}
		commitSHA = sha
	}

	renderRenameSummary(RenameSummary{
		TotalScanned: scanned,
		TotalMatched: len(pairs),
		TotalRenamed: renamedCount,
		IsGitRepo:    isGit,
		CommitSHA:    commitSHA,
	})

	return nil
}
