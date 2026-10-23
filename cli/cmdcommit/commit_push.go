package cmdcommit

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
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
	commitMessage := strings.Join(cleanArgs, " ")

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
	commitMessage := strings.Join(cleanArgs, " ")

	return executePullCommitPush(commitMessage, meta)
}

func executePullCommitPush(commitMessage string, meta CommitMetadata) error {
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
		return apperror.NewSimple("Usage: gitmap commit-push-bug <message>", "E9000")
	}
	commitMessage := "Bug: " + strings.Join(cleanArgs, " ")

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
		return apperror.NewSimple("Usage: gitmap commit-push-feature <message>", "E9000")
	}
	commitMessage := "Feature: " + strings.Join(cleanArgs, " ")

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
		return apperror.NewSimple("Usage: gitmap commit-push-chore <message>", "E9000")
	}
	commitMessage := "Chore: " + strings.Join(cleanArgs, " ")

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
		return apperror.NewSimple("Usage: gitmap commit-push-release <message>", "E9000")
	}
	commitMessage := "Release: " + strings.Join(cleanArgs, " ")

	return executeCommitPush(commitMessage, meta)
}

// executeCommitPush is the shared logic for all commit-push variants.
// When meta.TaskSlug is set (--task flag or GITMAP_TASK env), staging is scoped to
// that task's claimed files instead of staging everything.
func executeCommitPush(commitMessage string, meta CommitMetadata) *apperror.AppError {
	if !cmdpull.IsGitRepoCWD() {
		return handleNonGitRepoCommitPush()
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
	startTime := time.Now()
	err := execGitPaddedFiltered("commit", "-m", commitMessage)
	durMs := time.Since(startTime).Milliseconds()
	recordGitCommandDirect("git commit", "git commit -m "+commitMessage, ExtractGitExitCode(err), durMs)
	if err != nil {
		return apperror.WrapSimple(err, "git commit failed:")
	}

	return nil
}

func performPushStep() *apperror.AppError {
	startTime := time.Now()
	err := execGitPaddedFiltered("push")
	durMs := time.Since(startTime).Milliseconds()
	recordGitCommandDirect("git push", "git push", ExtractGitExitCode(err), durMs)
	if err != nil {
		return apperror.WrapSimple(err, "git push failed:")
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

func resolveShortHeadSha() string {
	headSha, _ := execGitOutputCP("rev-parse", "--short", "HEAD")
	sha := strings.TrimSpace(headSha)
	if sha == "" {
		return "-"
	}

	return sha
}

func handleCleanWorkingTreeCP() *apperror.AppError {
	printPaddedInfo("Working tree clean, nothing to commit.")

	unpushed := countUnpushedCommitsCP()
	if unpushed == 0 {
		printPaddedSuccess("Everything is up to date.")

		return nil
	}

	return pushUnpushedCommitsCP(unpushed)
}

func pushUnpushedCommitsCP(unpushed int) *apperror.AppError {
	printPaddedInfo("Pushing %d unpushed commit(s) to remote...", unpushed)

	if errPush := performPushStep(); errPush != nil {
		return errPush
	}

	printPaddedSuccess("Pushed %d commit(s) to remote.", unpushed)

	return nil
}

// ExecGitInheritCP runs a git command with inherited stdio.
func ExecGitInheritCP(gitArgs ...string) error {
	cmd := exec.Command("git", gitArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// execGitOutputCP runs a git command and captures its stdout.
func execGitOutputCP(gitArgs ...string) (string, error) {
	cmd := exec.Command("git", gitArgs...)
	out, err := cmd.Output()

	return string(out), err
}

// execGitPadded runs a git command and indents all output lines by 2 spaces.
func execGitPadded(gitArgs ...string) error {
	cmd := exec.Command("git", gitArgs...)
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return cmd.Run()
	}

	cmd.Stderr = cmd.Stdout
	if err := cmd.Start(); err != nil {
		return err
	}

	scanPaddedLines(stdoutPipe)

	return cmd.Wait()
}

func scanPaddedLines(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Printf("  %s\n", scanner.Text())
	}
}
