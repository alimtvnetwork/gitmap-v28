package cmd

import (
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/pterm/pterm"
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
		return apperror.New("Usage: gitmap commit-push \"<commit message>\"", "E9000", nil)
	}
	commitMessage := strings.Join(args, " ")
	return executeCommitPush(commitMessage)
}

// runCommitPushPull pulls first, then stages, commits, and pushes.
func runCommitPushPull(args []string) error {
	if isCommitPushHelpArg(args) {
		checkHelp(constants.CmdCommitPushPull, []string{"--help"})
		return nil
	}
	if len(args) == 0 {
		return apperror.New("Usage: gitmap commit-push-pull \"<commit message>\"", "E9000", nil)
	}
	commitMessage := strings.Join(args, " ")

	pterm.Info.Println("Pulling latest changes first...")
	if err := execGitInheritCP("pull", "--rebase"); err != nil {
		pterm.Warning.Println("Pull failed — you may need to resolve conflicts manually.")
		pterm.Warning.Printf("Error: %v\n", err)
		return apperror.New("fatal error", "E9000", nil)
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
		return apperror.New("Usage: gitmap commit-push-bug \"<what was fixed>\"", "E9000", nil)
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
		return apperror.New("Usage: gitmap commit-push-feature \"<what feature was added>\"", "E9000", nil)
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
		return apperror.New("Usage: gitmap commit-push-release \"<what release changes>\"", "E9000", nil)
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
		return apperror.New("Usage: gitmap rm-git <last-4-digits-of-sha>", "E9000", nil)
	}
	shaFragment := args[0]
	if len(shaFragment) < 4 {
		return apperror.New("SHA fragment must be at least 4 characters.", "E9000", nil)
	}

	// Resolve the full SHA from the fragment
	fullSha, err := execGitOutputCP("rev-parse", shaFragment)
	if err != nil {
		// Try log search if rev-parse fails
		fullSha, err = execGitOutputCP("log", "--all", "--format=%H", "--grep="+shaFragment)
		if err != nil || fullSha == "" {
			return apperror.New("Could not resolve SHA fragment: "+shaFragment, "E9000", nil)
		}
		// Take first match
		fullSha = strings.Split(strings.TrimSpace(fullSha), "\n")[0]
	}
	fullSha = strings.TrimSpace(fullSha)

	pterm.Info.Printf("Resolved SHA: %s\n", fullSha)
	pterm.Warning.Println("This will rewrite history. Use with caution.")

	// Use git rebase to drop the commit
	if err := execGitInheritCP("rebase", "--onto", fullSha+"^", fullSha); err != nil {
		pterm.Error.Printf("Failed to remove commit %s: %v\n", fullSha, err)
		pterm.Info.Println("You may need to run: git rebase --abort")
		return apperror.New("fatal error", "E9000", nil)
	}
	pterm.Success.Printf("Commit %s removed successfully.\n", fullSha[:8])
	return nil
}

// executeCommitPush is the shared logic for all commit-push variants.
func executeCommitPush(commitMessage string) *apperror.AppError {
	pterm.Info.Println("Staging all changes...")
	if err := execGitInheritCP("add", "-A"); err != nil {
		return apperror.Wrap(err, "git add failed:", nil)
	}

	pterm.Info.Printf("Committing: %s\n", commitMessage)
	if err := execGitInheritCP("commit", "-m", commitMessage); err != nil {
		return apperror.Wrap(err, "git commit failed:", nil)
	}

	pterm.Info.Println("Pushing to remote...")
	if err := execGitInheritCP("push"); err != nil {
		return apperror.Wrap(err, "git push failed:", nil)
	}

	pterm.Success.Println("Done! Changes committed and pushed.")
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
