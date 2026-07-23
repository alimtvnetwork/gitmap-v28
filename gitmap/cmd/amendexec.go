// Package cmd — amendexec.go handles git operations for the amend command.
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// listCommitsForAmend returns commits that will be rewritten.
func listCommitsForAmend(f amendFlags) []model.CommitEntry {
	var args []string

	switch {
	case f.commitHash == "":
		args = []string{"log", "--format=%H %s", "--reverse"}
	case f.commitHash == "HEAD":
		args = []string{"log", "--format=%H %s", "-1"}
	default:
		args = []string{"log", "--format=%H %s", "--reverse", f.commitHash + "^..HEAD"}
	}

	out, err := exec.Command("git", args...).Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrAmendListCommits, err)

		return nil
	}

	return parseCommitLines(string(out))
}

// parseCommitLines splits git log output into CommitEntry slices.
func parseCommitLines(output string) []model.CommitEntry {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	var entries []model.CommitEntry

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, " ", 2)
		msg := ""
		if len(parts) > 1 {
			msg = parts[1]
		}

		entries = append(entries, model.CommitEntry{
			SHA:     parts[0],
			Message: msg,
		})
	}

	return entries
}

// detectPreviousAuthor reads the author of the first commit in the range.
func detectPreviousAuthor(commits []model.CommitEntry) (string, string) {
	if len(commits) == 0 {
		return "", ""
	}

	sha := commits[0].SHA

	nameOut, nameErr := exec.Command("git", "log", "-1", "--format=%an", sha).Output()
	if nameErr != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not read author name for %s: %v\n", sha, nameErr)
	}

	emailOut, emailErr := exec.Command("git", "log", "-1", "--format=%ae", sha).Output()
	if emailErr != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not read author email for %s: %v\n", sha, emailErr)
	}

	return strings.TrimSpace(string(nameOut)), strings.TrimSpace(string(emailOut))
}

// getCurrentBranch returns the current Git branch name.
func getCurrentBranch() string {
	out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return "main"
	}

	return strings.TrimSpace(string(out))
}

// switchBranch checks out the specified branch.
func switchBranch(branch string) {
	fmt.Printf(constants.MsgAmendCheckout, branch)

	cmd := exec.Command("git", "checkout", branch)
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrAmendCheckout, branch, err)
		os.Exit(1)
	}
}

// runFilterBranch executes the git filter-branch command.
func runFilterBranch(f amendFlags, commits []model.CommitEntry) {
	if f.commitHash == "HEAD" {
		runAmendHead(f)

		return
	}

	envFilter := buildEnvFilter(f)
	var args []string

	if f.commitHash == "" {
		args = []string{"filter-branch", "-f", "--env-filter", envFilter, "--", "HEAD"}
	} else {
		args = []string{"filter-branch", "-f", "--env-filter", envFilter, f.commitHash + "^..HEAD"}
	}

	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrAmendFilter, err)
		os.Exit(1)
	}
}

// runAmendHead uses git commit --amend for single HEAD commit.
func runAmendHead(f amendFlags) {
	author := buildAuthorString(f)
	args := []string{"commit", "--amend", "--no-edit", "--author", author}

	cmd := exec.Command("git", args...)
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrAmendCommitAmend, err)
		os.Exit(1)
	}
}
