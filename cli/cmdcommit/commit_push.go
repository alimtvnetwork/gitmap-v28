package cmdcommit

import (
	"fmt"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpull"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// IsCommitPushHelpArg returns true if the first argument is a help trigger word.
func IsCommitPushHelpArg(args []string) bool {
	if len(args) == 0 {
		return false
	}

	return args[0] == "help" || args[0] == "--help" || args[0] == "-h"
}

// RunCommitPush stages all changes, commits with the given message, and pushes.
func RunCommitPush(args []string) error {
	if IsCommitPushHelpArg(args) {
		checkHelp(constants.CmdCommitPush, []string{"--help"})

		return nil
	}

	cleanArgs, meta := ExtractCommitFlags(args)
	if len(cleanArgs) == 0 {
		return apperror.NewSimple("Usage: gitmap commit-push <message>", "E9000")
	}
	commitMessage := FormatCommitMessage("cp", cleanArgs)

	return executeCommitPush(commitMessage, meta)
}

// RunPullCommitPush pulls first, then stages, commits, and pushes.
func RunPullCommitPush(args []string) error {
	if IsCommitPushHelpArg(args) {
		checkHelp(constants.CmdPullCommitPush, []string{"--help"})

		return nil
	}

	cleanArgs, meta := ExtractCommitFlags(args)
	if len(cleanArgs) == 0 {
		return apperror.NewSimple("Usage: gitmap pull-commit-push <message>", "E9000")
	}
	commitMessage := FormatCommitMessage("pcp", cleanArgs)

	return executePullCommitPush(commitMessage, meta)
}

func executePullCommitPush(commitMessage string, meta CommitMetadata) error {
	// Dry-run must return before the pull: a preview must not move the
	// tree, and a failing or unreachable remote must not fail a preview.
	if meta.IsDryRun {
		return previewCommitPush(commitMessage, meta, true)
	}

	fmt.Println()
	printPaddedInfo("Pulling latest changes before commit...")

	pullArgs := []string{}
	if errPull := cmdpull.RunPull(pullArgs); errPull != nil {
		printPaddedError("Pull failed: %v", errPull)

		return apperror.WrapSimple(errPull, "pre-commit pull")
	}

	return executeCommitPush(commitMessage, meta)
}

// RunCommitPushBug commits with a "Bug: " prefix.
func RunCommitPushBug(args []string) error {
	if IsCommitPushHelpArg(args) {
		checkHelp(constants.CmdCommitPushBug, []string{"--help"})

		return nil
	}

	cleanArgs, meta := ExtractCommitFlags(args)
	if len(cleanArgs) == 0 {
		return apperror.NewSimple("Usage: gitmap commit-push-bug <message>", string(constants.ErrCodeGenericExecution))
	}
	commitMessage := FormatCommitMessage(string(constants.CommitVerbBugAlt), cleanArgs)

	return executeCommitPush(commitMessage, meta)
}

// RunCommitPushFeature commits with a "Feature: " prefix.
func RunCommitPushFeature(args []string) error {
	if IsCommitPushHelpArg(args) {
		checkHelp(constants.CmdCommitPushFeature, []string{"--help"})

		return nil
	}

	cleanArgs, meta := ExtractCommitFlags(args)
	if len(cleanArgs) == 0 {
		return apperror.NewSimple("Usage: gitmap commit-push-feature <message>", string(constants.ErrCodeGenericExecution))
	}
	commitMessage := FormatCommitMessage(string(constants.CommitVerbFeature), cleanArgs)

	return executeCommitPush(commitMessage, meta)
}

// RunCommitPushChore commits with a "Chore: " prefix.
func RunCommitPushChore(args []string) error {
	if IsCommitPushHelpArg(args) {
		checkHelp(constants.CmdCommitPushChore, []string{"--help"})

		return nil
	}

	cleanArgs, meta := ExtractCommitFlags(args)
	if len(cleanArgs) == 0 {
		return apperror.NewSimple("Usage: gitmap commit-push-chore <message>", string(constants.ErrCodeGenericExecution))
	}
	commitMessage := FormatCommitMessage(string(constants.CommitVerbChore), cleanArgs)

	return executeCommitPush(commitMessage, meta)
}

// RunCommitPushRelease commits with a "Release: " prefix.
func RunCommitPushRelease(args []string) error {
	if IsCommitPushHelpArg(args) {
		checkHelp(constants.CmdCommitPushRelease, []string{"--help"})

		return nil
	}

	cleanArgs, meta := ExtractCommitFlags(args)
	if len(cleanArgs) == 0 {
		return apperror.NewSimple("Usage: gitmap commit-push-release <message>", string(constants.ErrCodeGenericExecution))
	}
	commitMessage := FormatCommitMessage(string(constants.CommitVerbRelease), cleanArgs)

	return executeCommitPush(commitMessage, meta)
}

// RunCommitPushRefactor commits with a "REFACTOR: " prefix.
func RunCommitPushRefactor(args []string) error {
	if IsCommitPushHelpArg(args) {
		checkHelp(constants.CmdCommitPushRefactor, []string{"--help"})

		return nil
	}

	cleanArgs, meta := ExtractCommitFlags(args)
	if len(cleanArgs) == 0 {
		return apperror.NewSimple("Usage: gitmap commit-push-refactor <message>", string(constants.ErrCodeGenericExecution))
	}
	commitMessage := FormatCommitMessage(string(constants.CommitVerbRefactor), cleanArgs)

	return executeCommitPush(commitMessage, meta)
}

// executeCommitPush is the shared logic for all commit-push variants.
// When meta.TaskSlug is set (--task flag or GITMAP_TASK env), staging is scoped to
// that task's claimed files instead of staging everything.
func executeCommitPush(commitMessage string, meta CommitMetadata) *apperror.AppError {
	if !cmdpull.IsGitRepoCWD() {
		return handleNonGitRepoCommitPush()
	}

	// Dry-run early return: the flag was parsed into meta.IsDryRun but
	// never checked here, so --dry-run previously fell through to staging,
	// committing, and pushing. Preview and return before any mutation.
	if meta.IsDryRun {
		return previewCommitPush(commitMessage, meta, false)
	}

	hasChanges, errPrep := stageAndCheckChanges(meta.TaskSlug)
	if errPrep != nil {
		return errPrep
	}

	if !hasChanges {
		return handleCleanWorkingTreeCP()
	}

	return performCommitPush(commitMessage, meta)
}

func performCommitPush(commitMessage string, meta CommitMetadata) *apperror.AppError {
	printPaddedInfo("Committing: %s", commitMessage)
	if errCommit := performCommitStep(commitMessage); errCommit != nil {
		return errCommit
	}

	RecordCurrentCommitSafely(meta)

	printPaddedInfo("Pushing to remote...")
	if errPush := performPushStep(); errPush != nil {
		return errPush
	}
	renderSuccessSummary(commitMessage)

	return nil
}

func performCommitStep(commitMessage string) *apperror.AppError {
	if idErr := ensureGitIdentity(); idErr != nil {
		return idErr
	}

	startTime := time.Now()
	out, err := execGitPaddedFilteredCapture("commit", "-m", commitMessage)
	durMs := time.Since(startTime).Milliseconds()
	recordGitCommandDirect("git commit", "git commit -m "+commitMessage, ExtractGitExitCode(err), durMs)
	if err != nil {
		return apperror.WrapSimple(gitErrorWithOutput(err, out), "git commit failed:")
	}

	return nil
}

func performPushStep() *apperror.AppError {
	startTime := time.Now()
	out, err := execGitPaddedFilteredCapture("push")
	durMs := time.Since(startTime).Milliseconds()
	recordGitCommandDirect("git push", "git push", ExtractGitExitCode(err), durMs)
	if err != nil {
		return apperror.WrapSimple(gitErrorWithOutput(err, out), "git push failed:")
	}

	return nil
}

func renderSuccessSummary(commitMessage string) {
	branch := resolveCurrentBranchOrHead()
	sha := resolveShortHeadSha()

	renderCommitPushSummaryCard(branch, sha, commitMessage, "pushed")
}

func resolveCurrentBranchOrHead() string {
	branch, _ := getCurrentBranchName()
	if branch == "" {
		return "HEAD"
	}

	return branch
}
