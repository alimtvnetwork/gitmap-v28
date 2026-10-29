package cmdcommit

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpull"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/timeutil"
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
	cleanArgs, errPrompt := PromptCommitMessageIfMissing(cleanArgs, "commit-push", "E9000")
	if errPrompt != nil {
		return errPrompt
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
	cleanArgs, errPrompt := PromptCommitMessageIfMissing(cleanArgs, "pull-commit-push", "E9000")
	if errPrompt != nil {
		return errPrompt
	}
	commitMessage := FormatCommitMessage("pcp", cleanArgs)

	return executePullCommitPush(commitMessage, meta)
}

func executePullCommitPush(commitMessage string, meta CommitMetadata) error {
	if meta.IsDryRun {
		return previewCommitPush(commitMessage, meta, true)
	}
	if errTime := verifyTimeRestriction(meta); errTime != nil {
		return errTime
	}
	defer cleanupTimeSimulatedEnv(meta.IsAutoWindow)

	return proceedPullCommitPush(commitMessage, meta)
}

func proceedPullCommitPush(commitMessage string, meta CommitMetadata) error {
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
	cleanArgs, errPrompt := PromptCommitMessageIfMissing(cleanArgs, "commit-push-bug", string(constants.ErrCodeGenericExecution))
	if errPrompt != nil {
		return errPrompt
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
	cleanArgs, errPrompt := PromptCommitMessageIfMissing(cleanArgs, "commit-push-feature", string(constants.ErrCodeGenericExecution))
	if errPrompt != nil {
		return errPrompt
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
	cleanArgs, errPrompt := PromptCommitMessageIfMissing(cleanArgs, "commit-push-chore", string(constants.ErrCodeGenericExecution))
	if errPrompt != nil {
		return errPrompt
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
	cleanArgs, errPrompt := PromptCommitMessageIfMissing(cleanArgs, "commit-push-release", string(constants.ErrCodeGenericExecution))
	if errPrompt != nil {
		return errPrompt
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
	cleanArgs, errPrompt := PromptCommitMessageIfMissing(cleanArgs, "commit-push-refactor", string(constants.ErrCodeGenericExecution))
	if errPrompt != nil {
		return errPrompt
	}
	commitMessage := FormatCommitMessage(string(constants.CommitVerbRefactor), cleanArgs)

	return executeCommitPush(commitMessage, meta)
}

// executeCommitPush is the shared logic for all commit-push variants.
func executeCommitPush(commitMessage string, meta CommitMetadata) *apperror.AppError {
	if !cmdpull.IsGitRepoCWD() {
		return handleNonGitRepoCommitPush()
	}
	if meta.IsDryRun {
		return previewCommitPush(commitMessage, meta, false)
	}
	if errTime := verifyTimeRestriction(meta); errTime != nil {
		return errTime
	}
	defer cleanupTimeSimulatedEnv(meta.IsAutoWindow)

	return proceedCommitPush(commitMessage, meta)
}

func proceedCommitPush(commitMessage string, meta CommitMetadata) *apperror.AppError {
	hasChanges, errPrep := stageAndCheckChanges(meta.TaskSlug)
	if errPrep != nil {
		return errPrep
	}
	if !hasChanges {
		return handleCleanWorkingTreeCP()
	}

	return performCommitPush(commitMessage, meta)
}

func verifyTimeRestriction(meta CommitMetadata) *apperror.AppError {
	if !isTimeRestricted(meta) {
		return nil
	}
	loc, errLoc := timeutil.ResolveLocation(meta.Timezone)
	if errLoc != nil {
		return apperror.NewSimple("Invalid timezone "+meta.Timezone+": "+errLoc.Error(), "E9001")
	}

	return evaluateTimeWindow(meta, loc)
}

func isTimeRestricted(meta CommitMetadata) bool {
	return meta.IsTimeEnforced || meta.IsAutoWindow
}

func evaluateTimeWindow(meta CommitMetadata, loc *time.Location) *apperror.AppError {
	now := time.Now().In(loc)
	printPaddedInfo("GitMap restricted status: %s (%s)", now.Format("02-Jan-06 03:04 PM"), loc.String())
	if timeutil.IsCommitWindowAllowed(now, loc) {
		return nil
	}
	if meta.IsAutoWindow {
		return applySimulatedMorningWindow(now, loc)
	}

	return newWindowClosedError(now, loc)
}

func applySimulatedMorningWindow(now time.Time, loc *time.Location) *apperror.AppError {
	morning := timeutil.SimulateMorningWindow(now, loc)
	stamp := morning.Format(time.RFC3339)
	os.Setenv("GIT_AUTHOR_DATE", stamp)
	os.Setenv("GIT_COMMITTER_DATE", stamp)
	printPaddedInfo("Time window closed; auto-window simulated compliant morning timestamp: %s", stamp)

	return nil
}

func newWindowClosedError(now time.Time, loc *time.Location) *apperror.AppError {
	nextSlot := timeutil.CalculateNextWindow(now, loc)
	msg := fmt.Sprintf("Commit rejected: outside allowed window for timezone %s. Next window opens at %s", loc.String(), nextSlot.Format(time.RFC3339))

	return apperror.NewSimple(msg, "E9002")
}

func cleanupTimeSimulatedEnv(isAutoWindow bool) {
	if isAutoWindow {
		os.Unsetenv("GIT_AUTHOR_DATE")
		os.Unsetenv("GIT_COMMITTER_DATE")
	}
}

func performCommitPush(commitMessage string, meta CommitMetadata) *apperror.AppError {
	if errCommit := commitAndRecord(commitMessage, meta); errCommit != nil {
		return errCommit
	}

	if meta.IsCommitOnly {
		printPaddedInfo("Commit-only: restrictions kept, push skipped")

		return nil
	}

	return pushAndSummarize(commitMessage, meta.Author)
}

func commitAndRecord(commitMessage string, meta CommitMetadata) *apperror.AppError {
	printPaddedInfo("Committing: %s", commitMessage)
	if errCommit := performCommitStep(commitMessage, meta.Author); errCommit != nil {
		return errCommit
	}
	RecordCurrentCommitSafely(meta)

	return nil
}

func pushAndSummarize(commitMessage, author string) *apperror.AppError {
	printPaddedInfo("Pushing to remote...")
	if errPush := performPushStep(); errPush != nil {
		return errPush
	}
	renderSuccessSummary(commitMessage, author)

	return nil
}

func performCommitStep(commitMessage, author string) *apperror.AppError {
	if idErr := ensureGitIdentity(); idErr != nil {
		return idErr
	}
	commitArgs := buildCommitArgs(commitMessage, author)

	return execAndRecordCommit(commitArgs)
}

func execAndRecordCommit(commitArgs []string) *apperror.AppError {
	startTime := time.Now()
	out, err := execGitPaddedFilteredCapture(commitArgs...)
	durMs := time.Since(startTime).Milliseconds()
	recordGitCommandDirect("git commit", "git "+strings.Join(commitArgs, " "), ExtractGitExitCode(err), durMs)
	if err != nil {
		return apperror.WrapSimple(gitErrorWithOutput(err, out), "git commit failed:")
	}

	return nil
}

func buildCommitArgs(commitMessage, author string) []string {
	if strings.TrimSpace(author) != "" {
		return []string{"commit", "--author=" + author, "-m", commitMessage}
	}

	return []string{"commit", "-m", commitMessage}
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

func renderSuccessSummary(commitMessage, author string) {
	branch := resolveCurrentBranchOrHead()
	sha := resolveShortHeadSha()

	renderCommitPushSummaryCard(branch, sha, commitMessage, "pushed", author)
}

func resolveCurrentBranchOrHead() string {
	branch, _ := getCurrentBranchName()
	if branch == "" {
		return "HEAD"
	}

	return branch
}
