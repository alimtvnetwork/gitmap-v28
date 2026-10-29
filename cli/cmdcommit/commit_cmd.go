// Package cmd provides CLI commands and execution dispatchers for gitmap.
package cmdcommit

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
	if isCommitInfoSubcommand(args) {
		return RunCommitInfo(args[1:])
	}

	if IsCommitPushHelpArg(args) {
		RenderCommitHelp()

		return nil
	}

	return dispatchCommit(args)
}

func isCommitInfoSubcommand(args []string) bool {
	return len(args) > 0 && strings.EqualFold(args[0], "info")
}

func dispatchCommit(args []string) error {
	cleanArgs, meta := ExtractCommitFlags(args)
	if !cmdpull.IsGitRepoCWD() {
		return handleNonGitRepoCommit(cleanArgs, meta)
	}

	if meta.IsDryRun {
		return previewCommit(cleanArgs, meta)
	}

	return executeCommit(cleanArgs, meta)
}

func parseCommitFlags(args []string) ([]string, bool, bool) {
	clean, meta := ExtractCommitFlags(args)

	return clean, meta.HasPush, meta.IsDryRun
}

func executeCommit(args []string, meta CommitMetadata) error {
	if errTime := verifyTimeRestriction(meta); errTime != nil {
		return errTime
	}
	defer cleanupTimeSimulatedEnv(meta.IsAutoWindow)

	hasChanges, err := stageAndCheckStaged()
	if err != nil {
		return err
	}

	if !hasChanges {
		printPaddedInfo("Working tree clean, nothing to commit.")

		return handleCleanWorkingTree(meta.HasPush)
	}

	return commitAndOptionalPush(args, meta)
}

func stageAndCheckStaged() (bool, error) {
	printPaddedInfo("Staging all changes...")

	if err := ExecGitInheritCP("add", "-A"); err != nil {
		return false, apperror.WrapSimple(err, "git add failed:")
	}

	hasChanges, err := hasStagedChangesCP()
	if err != nil {
		return false, apperror.WrapSimple(err, "check git status failed:")
	}

	return hasChanges, nil
}

func commitAndOptionalPush(args []string, meta CommitMetadata) error {
	if err := dispatchGitCommit(args); err != nil {
		return apperror.WrapSimple(err, "git commit failed:")
	}

	RecordCurrentCommitSafely(meta)

	return handleOptionalPush(meta.HasPush)
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

	if err := ExecGitInheritCP("push"); err != nil {
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

	if err := ExecGitInheritCP("push"); err != nil {
		return apperror.WrapSimple(err, "git push failed:")
	}

	printPaddedSuccess("Changes committed and pushed successfully.")

	return nil
}

func dispatchGitCommit(args []string) error {
	if len(args) == 0 {
		printPaddedInfo("Opening git commit editor...")

		return ExecGitInheritCP("commit")
	}

	if args[0] == "-m" {
		return ExecGitInheritCP(append([]string{"commit"}, args...)...)
	}

	msg := strings.Join(args, " ")
	printPaddedInfo("Committing: %s", msg)

	return ExecGitInheritCP("commit", "-m", msg)
}
