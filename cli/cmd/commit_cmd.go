// Package cmd provides CLI commands and execution dispatchers for gitmap.
package cmd

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// RunCommit is the exported entry point for commit and cm.
func RunCommit(args []string) error {
	return runCommit(args)
}

func runCommit(args []string) error {
	if isCommitPushHelpArg(args) {
		RenderCommitHelp()
		return nil
	}
	cleanArgs, hasPush, isDryRun := parseCommitFlags(args)
	if isDryRun {
		return previewCommitChanges()
	}
	return executeCommit(cleanArgs, hasPush)
}

func parseCommitFlags(args []string) ([]string, bool, bool) {
	var clean []string
	hasPush, isDryRun := false, false
	for _, a := range args {
		switch a {
		case "--push", "-p":
			hasPush = true
		case "--dry-run", "-n":
			isDryRun = true
		default:
			clean = append(clean, a)
		}
	}
	return clean, hasPush, isDryRun
}

func previewCommitChanges() error {
	printPaddedInfo("Dry run: previewing changes to be committed...")
	return execGitInheritCP("status", "--short")
}

func executeCommit(args []string, hasPush bool) error {
	printPaddedInfo("Staging all changes...")

	if err := execGitInheritCP("add", "-A"); err != nil {
		return apperror.WrapSimple(err, "git add failed:")
	}

	hasChanges, err := hasStagedChangesCP()
	if err != nil {
		return apperror.WrapSimple(err, "check git status failed:")
	}

	if !hasChanges {
		printPaddedInfo("Working tree clean, nothing to commit.")

		return handleCleanWorkingTree(hasPush)
	}

	if err := dispatchGitCommit(args); err != nil {
		return apperror.WrapSimple(err, "git commit failed:")
	}

	return handleOptionalPush(hasPush)
}

func handleCleanWorkingTree(hasPush bool) error {
	if !hasPush {
		printPaddedSuccess("Working tree clean.")

		return nil
	}

	unpushed := countUnpushedCommitsCP()
	if unpushed == 0 {
		printPaddedSuccess("Everything is up to date.")

		return nil
	}

	return pushUnpushedCommits(unpushed)
}

func pushUnpushedCommits(unpushed int) error {
	printPaddedInfo("Pushing %d unpushed commit(s) to remote...", unpushed)

	if err := execGitInheritCP("push"); err != nil {
		return apperror.WrapSimple(err, "git push failed:")
	}

	printPaddedSuccess("Pushed %d commit(s) to remote.", unpushed)

	return nil
}

func handleOptionalPush(hasPush bool) error {
	if !hasPush {
		printPaddedSuccess("Changes committed successfully.")

		return nil
	}

	printPaddedInfo("Pushing to remote...")

	if err := execGitInheritCP("push"); err != nil {
		return apperror.WrapSimple(err, "git push failed:")
	}

	printPaddedSuccess("Changes committed and pushed successfully.")

	return nil
}

func dispatchGitCommit(args []string) error {
	if len(args) == 0 {
		printPaddedInfo("Opening git commit editor...")
		return execGitInheritCP("commit")
	}
	if args[0] == "-m" {
		return execGitInheritCP(append([]string{"commit"}, args...)...)
	}
	msg := strings.Join(args, " ")
	printPaddedInfo("Committing: %s", msg)
	return execGitInheritCP("commit", "-m", msg)
}
