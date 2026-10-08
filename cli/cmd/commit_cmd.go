// Package cmd provides CLI commands and execution dispatchers for gitmap.
package cmd

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpull"
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

	return dispatchCommit(args)
}

func dispatchCommit(args []string) error {
	cleanArgs, hasPush, isDryRun := parseCommitFlags(args)
	if !cmdpull.IsGitRepoCWD() {
		return handleNonGitRepoCommit(cleanArgs, hasPush, isDryRun)
	}

	if isDryRun {
		return previewCommitChanges()
	}

	return executeCommit(cleanArgs, hasPush)
}

func parseCommitFlags(args []string) ([]string, bool, bool) {
	var clean []string
	hasPush, isDryRun := false, false

	for _, a := range args {
		clean, hasPush, isDryRun = evaluateCommitFlag(a, clean, hasPush, isDryRun)
	}

	return clean, hasPush, isDryRun
}

func evaluateCommitFlag(arg string, clean []string, hasPush, isDryRun bool) ([]string, bool, bool) {
	switch arg {
	case "--push", "-p":
		return clean, true, isDryRun
	case "--dry-run", "-n":
		return clean, hasPush, true
	default:
		return append(clean, arg), hasPush, isDryRun
	}
}

func previewCommitChanges() error {
	printPaddedInfo("Dry run: previewing changes to be committed...")

	return execGitInheritCP("status", "--short")
}

func executeCommit(args []string, hasPush bool) error {
	hasChanges, err := stageAndCheckStaged()
	if err != nil {
		return err
	}

	if !hasChanges {
		printPaddedInfo("Working tree clean, nothing to commit.")

		return handleCleanWorkingTree(hasPush)
	}

	return commitAndOptionalPush(args, hasPush)
}

func stageAndCheckStaged() (bool, error) {
	printPaddedInfo("Staging all changes...")

	if err := execGitInheritCP("add", "-A"); err != nil {
		return false, apperror.WrapSimple(err, "git add failed:")
	}

	hasChanges, err := hasStagedChangesCP()
	if err != nil {
		return false, apperror.WrapSimple(err, "check git status failed:")
	}

	return hasChanges, nil
}

func commitAndOptionalPush(args []string, hasPush bool) error {
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
