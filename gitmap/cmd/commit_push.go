package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
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

	commitMessage := strings.Join(args, " ")
	printPaddedInfo("Pulling latest changes first...")

	if err := execGitInheritCP("pull", "--rebase"); err != nil {
		printPaddedWarning("Pull failed — you may need to resolve conflicts manually.")
		printPaddedWarning("Error: %v", err)
		return apperror.NewSimple("fatal error", "E9000")
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
		if arg == "--no-push" || arg == "--local" || arg == "-n" {
			isPush = false
			continue
		}

		if targetSha == "" && !strings.HasPrefix(arg, "-") {
			targetSha = arg
		}
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

	fmt.Println()
	fullSha, errResolve := resolveSHAFragment(targetSha)
	if errResolve != nil {
		return errResolve
	}

	printPaddedInfo("Resolved SHA: %s", fullSha)
	printPaddedWarning("This will rewrite history. Use with caution.")

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

	fmt.Println()
	fullSha, errResolve := resolveSHAFragment(targetSha)
	if errResolve != nil {
		return errResolve
	}

	subject := getCommitSubject(fullSha)
	preSha, _ := execGitOutputCP("rev-parse", "HEAD")

	printPaddedInfo("Resolved target SHA: %s", fullSha)
	printPaddedWarning("Resetting branch history to %s (%s).", fullSha[:8], subject)

	if errReset := execGitPadded("reset", "--hard", fullSha); errReset != nil {
		printPaddedError("Failed to reset branch: %v", errReset)
		return apperror.NewSimple("reset failed", "E9000")
	}

	printPaddedSuccess("Branch reset to %s (%s).", fullSha[:8], subject)

	if isPush {
		syncRemoteAfterRewrite()
	}

	if len(preSha) >= 8 {
		printPaddedInfo("To undo this reset locally: git reset --hard %s", strings.TrimSpace(preSha)[:8])
	}

	printCommandVersionFooter()

	return nil
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

	if errRebase := execGitPadded("rebase", "--onto", fullSha+"^", fullSha, "HEAD"); errRebase != nil {
		execGitOutputCP("rebase", "--abort")
		printPaddedError("Failed to rebase commit %s: %v", fullSha[:8], errRebase)
		return apperror.NewSimple("rebase failed", "E9000")
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

// executeCommitPush is the shared logic for all commit-push variants.
func executeCommitPush(commitMessage string) *apperror.AppError {
	printPaddedInfo("Staging all changes...")
	if err := execGitInheritCP("add", "-A"); err != nil {
		return apperror.WrapSimple(err, "git add failed:")
	}

	printPaddedInfo("Committing: %s", commitMessage)
	if err := execGitInheritCP("commit", "-m", commitMessage); err != nil {
		return apperror.WrapSimple(err, "git commit failed:")
	}

	printPaddedInfo("Pushing to remote...")
	if err := execGitInheritCP("push"); err != nil {
		return apperror.WrapSimple(err, "git push failed:")
	}

	printPaddedSuccess("Changes committed and pushed.")

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

	scanner := bufio.NewScanner(stdoutPipe)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("  %s\n", line)
	}

	return cmd.Wait()
}

func resolveSHAFragment(shaFragment string) (string, error) {
	fullSha, err := execGitOutputCP("rev-parse", shaFragment)
	if err == nil && fullSha != "" {
		return strings.TrimSpace(fullSha), nil
	}

	fullSha, err = execGitOutputCP("log", "--all", "--format=%H", "--grep="+shaFragment)
	if err != nil || fullSha == "" {
		return "", apperror.NewSimple("Could not resolve SHA fragment: "+shaFragment, "E9000")
	}

	firstMatch := strings.Split(strings.TrimSpace(fullSha), "\n")[0]

	return strings.TrimSpace(firstMatch), nil
}
