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
	resolveGitTracking(pairs, isGit)
	renderRenameHeader(opts, isGit, len(pairs), totalScanned, cwd)
	if len(pairs) == 0 {
		renderZeroMatchMessage(opts)

		return nil
	}
	if opts.IsDryRun {
		renderDryRunPreviews(pairs)

		return nil
	}
	isProceed, confirmErr := checkPreflightConfirmation(pairs, opts, isGit)
	if confirmErr != nil || !isProceed {
		return confirmErr
	}

	return applyRenamesAndCommit(pairs, opts, totalScanned, isGit)
}

func resolveGitTracking(pairs []RenamePair, isGit bool) {
	for i := range pairs {
		pairs[i].IsGitTracked = isGit && isFileTrackedInGit(pairs[i].OldPath)
	}
}

func checkPreflightConfirmation(pairs []RenamePair, opts LowerCaseFixOptions, isGit bool) (bool, error) {
	if !isGit || opts.IsYes {
		return true, nil
	}

	isConfirmed, err := promptPreflightConfirmation(pairs, opts)
	if err != nil {
		return false, err
	}
	if !isConfirmed {
		renderCancelledMessage()

		return false, nil
	}

	return true, nil
}

func renderZeroMatchMessage(opts LowerCaseFixOptions) {
	filterStr := resolveFilterStr(opts)
	fmt.Printf("%s✓ No uppercase files matching filter [%s] found to rename.%s\n", constants.ColorGreen, filterStr, constants.ColorReset)
	fmt.Printf("  All matching files in this directory are already lowercase. No changes needed.\n\n")
}

func renderDryRunPreviews(pairs []RenamePair) {
	for i, p := range pairs {
		prefix := resolveRenamePrefix(p.IsGitTracked)
		matched := resolveMatchedInfo(p.MatchedBy)
		fmt.Printf("  [%d/%d] %s → %s%s (%s preview)\n", i+1, len(pairs), p.OldBase, p.NewBase, matched, prefix)
	}
	renderDryRunNotice(len(pairs))
}

func resolveRenamePrefix(isGitTracked bool) string {
	if isGitTracked {
		return "git mv"
	}

	return "fs rename"
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

	commitSHA, commitErr := maybeCommitRenames(pairs, opts, isGit)
	if commitErr != nil {
		return commitErr
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

func maybeCommitRenames(pairs []RenamePair, opts LowerCaseFixOptions, isGit bool) (string, error) {
	if !isGit || opts.IsNoCommit {
		return "", nil
	}

	return commitRenames(pairs, opts.CommitMessage)
}
