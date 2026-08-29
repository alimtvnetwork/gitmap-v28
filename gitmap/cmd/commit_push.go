package cmd

import (
	"os"
	"os/exec"
	"strings"

	"github.com/pterm/pterm"

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

	pterm.Info.Println("Pulling latest changes first...")
	if err := execGitInheritCP("pull", "--rebase"); err != nil {
		pterm.Warning.Println("Pull failed — you may need to resolve conflicts manually.")
		pterm.Warning.Printf("Error: %v\n", err)
		return apperror.NewSimple("fatal error", "E9000")
	}
	pterm.Success.Println("Pull complete.")
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

// runRmGit removes a commit by its last 4-digit SHA prefix using rebase --onto.
func runRmGit(args []string) error {
	if isCommitPushHelpArg(args) {
		checkHelp(constants.CmdRmGit, []string{"--help"})
		return nil
	}
	if len(args) == 0 {
		return apperror.NewSimple("Usage: gitmap rm-git <last-4-digits-of-sha>", "E9000")
	}
	shaFragment := args[0]
	if len(shaFragment) < 4 {
		return apperror.NewSimple("SHA fragment must be at least 4 characters.", "E9000")
	}

	fullSha, errResolve := resolveSHAFragment(shaFragment)
	if errResolve != nil {
		return errResolve
	}

	pterm.Info.Printf("Resolved SHA: %s\n", fullSha)
	pterm.Warning.Println("This will rewrite history. Use with caution.")

	// Use git reset --hard to drop the commit
	if err := execGitInheritCP("reset", "--hard", fullSha+"^"); err != nil {
		pterm.Error.Printf("Failed to remove commit %s: %v\n", fullSha, err)
		pterm.Info.Println("You may need to check your git reflog to recover.")
		return apperror.NewSimple("fatal error", "E9000")
	}
	pterm.Success.Printf("Commit %s removed successfully.\n", fullSha[:8])
	return nil
}

// executeCommitPush is the shared logic for all commit-push variants.
func executeCommitPush(commitMessage string) *apperror.AppError {
	pterm.Info.Println("Staging all changes...")
	if err := execGitInheritCP("add", "-A"); err != nil {
		return apperror.WrapSimple(err, "git add failed:")
	}

	pterm.Info.Printf("Committing: %s\n", commitMessage)
	if err := execGitInheritCP("commit", "-m", commitMessage); err != nil {
		return apperror.WrapSimple(err, "git commit failed:")
	}

	pterm.Info.Println("Pushing to remote...")
	if err := execGitInheritCP("push"); err != nil {
		return apperror.WrapSimple(err, "git push failed:")
	}

	pterm.Success.Println("!Done Changes committed and pushed.")
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
