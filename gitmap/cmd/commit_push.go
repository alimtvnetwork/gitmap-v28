package cmd

import (
	"os"
	"os/exec"
	"strings"

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
		pterm.Error.Println("Usage: gitmap commit-push \"<commit message>\"")
		os.Exit(1)
	}
	commitMessage := strings.Join(args, " ")
	executeCommitPush(commitMessage)
	return nil
}

// runCommitPushPull pulls first, then stages, commits, and pushes.
func runCommitPushPull(args []string) error {
	if isCommitPushHelpArg(args) {
		checkHelp(constants.CmdCommitPushPull, []string{"--help"})
		return nil
	}
	if len(args) == 0 {
		pterm.Error.Println("Usage: gitmap commit-push-pull \"<commit message>\"")
		os.Exit(1)
	}
	commitMessage := strings.Join(args, " ")

	pterm.Info.Println("Pulling latest changes first...")
	if err := execGitInheritCP("pull", "--rebase"); err != nil {
		pterm.Warning.Println("Pull failed — you may need to resolve conflicts manually.")
		pterm.Warning.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	pterm.Success.Println("Pull complete.")
	executeCommitPush(commitMessage)
	return nil
}

// runCommitPushBug commits with a "Bug: " prefix.
func runCommitPushBug(args []string) error {
	if isCommitPushHelpArg(args) {
		checkHelp(constants.CmdCommitPushBug, []string{"--help"})
		return nil
	}
	if len(args) == 0 {
		pterm.Error.Println("Usage: gitmap commit-push-bug \"<what was fixed>\"")
		os.Exit(1)
	}
	commitMessage := "Bug: " + strings.Join(args, " ")
	executeCommitPush(commitMessage)
	return nil
}

// runCommitPushFeature commits with a "Feature: " prefix.
func runCommitPushFeature(args []string) error {
	if isCommitPushHelpArg(args) {
		checkHelp(constants.CmdCommitPushFeature, []string{"--help"})
		return nil
	}
	if len(args) == 0 {
		pterm.Error.Println("Usage: gitmap commit-push-feature \"<what feature was added>\"")
		os.Exit(1)
	}
	commitMessage := "Feature: " + strings.Join(args, " ")
	executeCommitPush(commitMessage)
	return nil
}

// runCommitPushRelease commits with a "Release: " prefix.
func runCommitPushRelease(args []string) error {
	if isCommitPushHelpArg(args) {
		checkHelp(constants.CmdCommitPushRelease, []string{"--help"})
		return nil
	}
	if len(args) == 0 {
		pterm.Error.Println("Usage: gitmap commit-push-release \"<what release changes>\"")
		os.Exit(1)
	}
	commitMessage := "Release: " + strings.Join(args, " ")
	executeCommitPush(commitMessage)
	return nil
}

// runRmGit removes a commit by its last 4-digit SHA prefix using rebase --onto.
func runRmGit(args []string) error {
	if isCommitPushHelpArg(args) {
		checkHelp(constants.CmdRmGit, []string{"--help"})
		return nil
	}
	if len(args) == 0 {
		pterm.Error.Println("Usage: gitmap rm-git <last-4-digits-of-sha>")
		os.Exit(1)
	}
	shaFragment := args[0]
	if len(shaFragment) < 4 {
		pterm.Error.Println("SHA fragment must be at least 4 characters.")
		os.Exit(1)
	}

	// Resolve the full SHA from the fragment
	fullSha, err := execGitOutputCP("rev-parse", shaFragment)
	if err != nil {
		// Try log search if rev-parse fails
		fullSha, err = execGitOutputCP("log", "--all", "--format=%H", "--grep="+shaFragment)
		if err != nil || fullSha == "" {
			pterm.Error.Printf("Could not resolve SHA fragment: %s\n", shaFragment)
			os.Exit(1)
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
		os.Exit(1)
	}
	pterm.Success.Printf("Commit %s removed successfully.\n", fullSha[:8])
	return nil
}

// executeCommitPush is the shared logic for all commit-push variants.
func executeCommitPush(commitMessage string) {
	pterm.Info.Println("Staging all changes...")
	if err := execGitInheritCP("add", "-A"); err != nil {
		pterm.Error.Printf("git add failed: %v\n", err)
		os.Exit(1)
	}

	pterm.Info.Printf("Committing: %s\n", commitMessage)
	if err := execGitInheritCP("commit", "-m", commitMessage); err != nil {
		pterm.Error.Printf("git commit failed: %v\n", err)
		os.Exit(1)
	}

	pterm.Info.Println("Pushing to remote...")
	if err := execGitInheritCP("push"); err != nil {
		pterm.Error.Printf("git push failed: %v\n", err)
		os.Exit(1)
	}

	pterm.Success.Println("Done! Changes committed and pushed.")
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
