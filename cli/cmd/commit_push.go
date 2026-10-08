package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// isCommitPushHelpArg returns true if the first argument is a help trigger word.
func isCommitPushHelpArg(args []string) bool {
	if len(args) == 0 {
		return false
	}

	return args[0] == "help" || args[0] == "--help" || args[0] == "-h"
}

func printPaddedInfo(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("  %s%s%s %s\n", constants.ColorCyan, "INFO", constants.ColorReset, msg)
}

func printPaddedSuccess(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("  %s%s%s %s\n", constants.ColorGreen, "SUCCESS", constants.ColorReset, msg)
}

func printPaddedWarning(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("  %s%s%s %s\n", constants.ColorYellow, "WARNING", constants.ColorReset, msg)
}

func printPaddedError(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("  %s%s%s %s\n", constants.ColorRed, "ERROR", constants.ColorReset, msg)
}

func printCommandVersionFooter() {
	fmt.Println()
	fmt.Printf("  %sgitmap v%s%s\n", constants.ColorDim, constants.Version, constants.ColorReset)
}

// runCommitPush stages all changes, commits with the given message, and pushes.
func runCommitPush(args []string) error {
	if isCommitPushHelpArg(args) {
		checkHelp(constants.CmdCommitPush, []string{"--help"})

		return nil
	}

	if len(args) == 0 {
		return apperror.NewSimple("Usage: gitmap commit-push \"<commit message>\"", "E9000")
	}

	commitMessage := strings.Join(args, " ")

	return executeCommitPush(commitMessage)
}

// runPullCommitPush pulls first, then stages, commits, and pushes.
func runPullCommitPush(args []string) error {
	if isCommitPushHelpArg(args) {
		checkHelp(constants.CmdPullCommitPush, []string{"--help"})

		return nil
	}

	if len(args) == 0 {
		return apperror.NewSimple("Usage: gitmap pull-commit-push \"<commit message>\"", "E9000")
	}

	return executePullCommitPush(strings.Join(args, " "))
}

func executePullCommitPush(commitMessage string) error {
	printPaddedInfo("Pulling latest changes first...")

	if err := execGitInheritCP("pull", "--rebase"); err != nil {
		printPaddedWarning("Pull failed — you may need to resolve conflicts manually.")
		printPaddedWarning("Error: %v", err)

		return apperror.WrapSimple(err, "pull --rebase")
	}

	printPaddedSuccess("Pull complete.")

	return executeCommitPush(commitMessage)
}

// runCommitPushBug commits with a "Bug: " prefix.
func runCommitPushBug(args []string) error {
	if isCommitPushHelpArg(args) {
		checkHelp(constants.CmdCommitPushBug, []string{"--help"})

		return nil
	}

	if len(args) == 0 {
		return apperror.NewSimple("Usage: gitmap commit-push-bug \"<what was fixed>\"", "E9000")
	}

	commitMessage := "Bug: " + strings.Join(args, " ")

	return executeCommitPush(commitMessage)
}

// runCommitPushFeature commits with a "Feature: " prefix.
func runCommitPushFeature(args []string) error {
	if isCommitPushHelpArg(args) {
		checkHelp(constants.CmdCommitPushFeature, []string{"--help"})

		return nil
	}

	if len(args) == 0 {
		return apperror.NewSimple("Usage: gitmap commit-push-feature \"<what feature was added>\"", "E9000")
	}

	commitMessage := "Feature: " + strings.Join(args, " ")

	return executeCommitPush(commitMessage)
}

// runCommitPushChore commits with a "Chore: " prefix.
func runCommitPushChore(args []string) error {
	if isCommitPushHelpArg(args) {
		checkHelp(constants.CmdCommitPushChore, []string{"--help"})

		return nil
	}

	if len(args) == 0 {
		return apperror.NewSimple("Usage: gitmap commit-push-chore \"<what chore was done>\"", "E9000")
	}

	commitMessage := "Chore: " + strings.Join(args, " ")

	return executeCommitPush(commitMessage)
}

// runCommitPushRelease commits with a "Release: " prefix.
func runCommitPushRelease(args []string) error {
	if isCommitPushHelpArg(args) {
		checkHelp(constants.CmdCommitPushRelease, []string{"--help"})

		return nil
	}

	if len(args) == 0 {
		return apperror.NewSimple("Usage: gitmap commit-push-release \"<what release changes>\"", "E9000")
	}

	commitMessage := "Release: " + strings.Join(args, " ")

	return executeCommitPush(commitMessage)
}

// parseRewriteFlags extracts target SHA and push flag from args.
func parseRewriteFlags(args []string) (string, bool) {
	isPush := true
	targetSha := ""

	for _, arg := range args {
		targetSha, isPush = evaluateRewriteFlag(arg, targetSha, isPush)
	}

	return targetSha, isPush
}

func evaluateRewriteFlag(arg, targetSha string, isPush bool) (string, bool) {
	if arg == "--no-push" || arg == "--local" || arg == "-n" {
		return targetSha, false
	}

	if strings.HasPrefix(arg, "-") {
		return targetSha, isPush
	}

	if targetSha == "" {
		return arg, isPush
	}

	return targetSha, isPush
}

// runRmGit removes a commit by its SHA prefix and synchronizes remote tracking.
func runRmGit(args []string) error {
	if isCommitPushHelpArg(args) {
		checkHelp(constants.CmdRmGit, []string{"--help"})

		return nil
	}

	targetSha, isPush := parseRewriteFlags(args)
	if len(targetSha) < 4 {
		return apperror.NewSimple("Usage: gitmap rm-git <sha-fragment> [--no-push]", "E9000")
	}

	return executeRmGit(targetSha, isPush)
}

func executeRmGit(targetSha string, isPush bool) error {
	fmt.Println()
	fullSha, errResolve := resolveSHAFragment(targetSha)
	if errResolve != nil {
		return errResolve
	}

	printPaddedInfo("Resolved SHA: %s", fullSha)
	printPaddedWarning("This will rewrite history. Use with caution.")

	return finishRmGit(fullSha, isPush)
}

func finishRmGit(fullSha string, isPush bool) error {
	if errDrop := executeDropCommit(fullSha); errDrop != nil {
		return errDrop
	}

	printPaddedSuccess("Commit %s removed successfully.", fullSha[:8])

	if isPush {
		syncRemoteAfterRewrite()
	}

	printCommandVersionFooter()

	return nil
}

// runGitReset resets the current branch to a target SHA and synchronizes remote.
func runGitReset(args []string) error {
	if isCommitPushHelpArg(args) {
		checkHelp(constants.CmdGitReset, []string{"--help"})

		return nil
	}

	targetSha, isPush := parseRewriteFlags(args)
	if len(targetSha) < 4 {
		return apperror.NewSimple("Usage: gitmap git-reset <target-sha> [--no-push]", "E9000")
	}

	return executeGitReset(targetSha, isPush)
}

func executeGitReset(targetSha string, isPush bool) error {
	fmt.Println()
	fullSha, errResolve := resolveSHAFragment(targetSha)
	if errResolve != nil {
		return errResolve
	}

	subject := getCommitSubject(fullSha)
	preSha, _ := execGitOutputCP("rev-parse", "HEAD")

	printPaddedInfo("Resolved target SHA: %s", fullSha)
	printPaddedWarning("Resetting branch history to %s (%s).", fullSha[:8], subject)

	return finishGitReset(fullSha, subject, preSha, isPush)
}

func finishGitReset(fullSha, subject, preSha string, isPush bool) error {
	if errReset := execGitPadded("reset", "--hard", fullSha); errReset != nil {
		printPaddedError("Failed to reset branch: %v", errReset)

		return apperror.WrapSimple(errReset, "reset branch")
	}

	printPaddedSuccess("Branch reset to %s (%s).", fullSha[:8], subject)

	if isPush {
		syncRemoteAfterRewrite()
	}

	printUndoResetHint(preSha)
	printCommandVersionFooter()

	return nil
}

func printUndoResetHint(preSha string) {
	if len(preSha) >= 8 {
		printPaddedInfo("To undo this reset locally: git reset --hard %s", strings.TrimSpace(preSha)[:8])
	}
}

// RunGitReset is the exported entry point for git-reset.
func RunGitReset(args []string) error {
	return runGitReset(args)
}

// RunRmGit is the exported entry point for rm-git.
func RunRmGit(args []string) error {
	return runRmGit(args)
}

func executeDropCommit(fullSha string) error {
	headSha, errHead := execGitOutputCP("rev-parse", "HEAD")
	isHead := errHead == nil && strings.HasPrefix(strings.TrimSpace(headSha), fullSha)

	if isHead {
		return execGitPadded("reset", "--hard", "HEAD~1")
	}

	return rebaseDropCommit(fullSha)
}

func rebaseDropCommit(fullSha string) error {
	if errRebase := execGitPadded("rebase", "--onto", fullSha+"^", fullSha, "HEAD"); errRebase != nil {
		execGitOutputCP("rebase", "--abort")
		printPaddedError("Failed to rebase commit %s: %v", fullSha[:8], errRebase)

		return apperror.WrapSimple(errRebase, "rebase commit")
	}

	return nil
}

func syncRemoteAfterRewrite() {
	branch, errBranch := getCurrentBranchName()
	if errBranch != nil || !hasRemoteTracking(branch) {
		return
	}

	printPaddedInfo("Synchronizing remote branch 'origin/%s' (force-pushing with lease)...", branch)
	if errPush := forcePushRemote(branch); errPush != nil {
		printPaddedWarning("Remote push with lease failed: %v", errPush)
		printPaddedWarning("Run 'git push --force origin %s' if you wish to overwrite remote history.", branch)

		return
	}

	printPaddedSuccess("Remote 'origin/%s' synchronized successfully.", branch)
}

func getCurrentBranchName() (string, error) {
	out, err := execGitOutputCP("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(out), nil
}

func hasRemoteTracking(branch string) bool {
	remote, err := execGitOutputCP("config", "--get", fmt.Sprintf("branch.%s.remote", branch))
	if err != nil || strings.TrimSpace(remote) == "" {
		return false
	}

	return true
}

func forcePushRemote(branch string) error {
	return execGitPadded("push", "--force-with-lease", "origin", branch)
}

func getCommitSubject(sha string) string {
	out, err := execGitOutputCP("log", "-1", "--format=%s", sha)
	if err != nil {
		return "commit"
	}

	return strings.TrimSpace(out)
}

// hasStagedChangesCP checks whether there are any staged or unstaged changes in the working tree.
func hasStagedChangesCP() (bool, error) {
	out, err := execGitOutputCP("status", "--porcelain")
	if err != nil {
		return false, err
	}

	hasChanges := len(strings.TrimSpace(out)) > 0

	return hasChanges, nil
}

// countUnpushedCommitsCP returns the count of unpushed commits ahead of upstream remote tracking branch.
func countUnpushedCommitsCP() int {
	branch, errBranch := getCurrentBranchName()
	if errBranch != nil || !hasRemoteTracking(branch) {
		return 0
	}

	return parseUnpushedCount()
}

func parseUnpushedCount() int {
	out, errCount := execGitOutputCP("rev-list", "--count", "@{u}..HEAD")
	if errCount != nil {
		return 0
	}

	count, errParse := strconv.Atoi(strings.TrimSpace(out))
	if errParse != nil {
		return 0
	}

	return count
}

// executeCommitPush is the shared logic for all commit-push variants.
func executeCommitPush(commitMessage string) *apperror.AppError {
	if !isGitRepoCWD() {
		return handleNonGitRepoCommitPush()
	}

	hasChanges, errPrep := stageAndCheckChanges()
	if errPrep != nil {
		return errPrep
	}

	if !hasChanges {
		return handleCleanWorkingTreeCP()
	}

	return performCommitPush(commitMessage)
}

func stageAndCheckChanges() (bool, *apperror.AppError) {
	printPaddedInfo("Staging all changes...")

	if errAdd := executeStageAllStep(); errAdd != nil {
		return false, errAdd
	}

	hasChanges, errStatus := hasStagedChangesCP()
	if errStatus != nil {
		return false, apperror.WrapSimple(errStatus, "check git status failed:")
	}

	return hasChanges, nil
}

func executeStageAllStep() *apperror.AppError {
	startTime := time.Now()
	err := execGitPaddedFiltered("add", "-A")
	durMs := time.Since(startTime).Milliseconds()
	recordGitCommandDirect("git add", "git add -A", extractGitExitCode(err), durMs)
	if err != nil {
		return apperror.WrapSimple(err, "git add failed:")
	}

	return nil
}

func performCommitPush(commitMessage string) *apperror.AppError {
	printPaddedInfo("Committing: %s", commitMessage)
	if errCommit := performCommitStep(commitMessage); errCommit != nil {
		return errCommit
	}
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
	recordGitCommandDirect("git commit", "git commit -m "+commitMessage, extractGitExitCode(err), durMs)
	if err != nil {
		return apperror.WrapSimple(err, "git commit failed:")
	}

	return nil
}

func performPushStep() *apperror.AppError {
	startTime := time.Now()
	err := execGitPaddedFiltered("push")
	durMs := time.Since(startTime).Milliseconds()
	recordGitCommandDirect("git push", "git push", extractGitExitCode(err), durMs)
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

// execGitInheritCP runs a git command with inherited stdio.
func execGitInheritCP(gitArgs ...string) error {
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

func resolveSHAFragment(shaFragment string) (string, error) {
	fullSha, err := execGitOutputCP("rev-parse", shaFragment)
	if err == nil && fullSha != "" {
		return strings.TrimSpace(fullSha), nil
	}

	return searchSHAFragment(shaFragment)
}

func searchSHAFragment(shaFragment string) (string, error) {
	fullSha, err := execGitOutputCP("log", "--all", "--format=%H", "--grep="+shaFragment)
	if err != nil {
		return "", apperror.WrapSimple(err, "resolve SHA fragment: "+shaFragment)
	}

	if fullSha == "" {
		return "", apperror.NewSimple("Could not resolve SHA fragment: "+shaFragment, "E9000")
	}

	firstMatch := strings.Split(strings.TrimSpace(fullSha), "\n")[0]

	return strings.TrimSpace(firstMatch), nil
}

func extractGitExitCode(err error) int {
	if err == nil {
		return 0
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode()
	}

	return 1
}

func recordGitCommandDirect(cmdName, cmdLine string, exitCode int, durationMs int64) {
	store.RecordCommandSafely(cmdLine, cmdName, exitCode, durationMs)
}

func recordGitCommandHistory(gitArgs []string, exitCode int, durationMs int64) {
	if len(gitArgs) == 0 {
		recordGitCommandDirect("git", "git", exitCode, durationMs)
		return
	}
	cmdName := "git " + gitArgs[0]
	cmdLine := "git " + strings.Join(gitArgs, " ")
	recordGitCommandDirect(cmdName, cmdLine, exitCode, durationMs)
}
