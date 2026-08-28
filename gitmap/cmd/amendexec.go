// Package cmd — amendexec.go handles git operations for the amend command.
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func getAmendLogArgs(f amendFlags) []string {
	switch {
	case f.commitHash == "":
		return []string{"log", "--format=%H %s", "--reverse"}
	case f.commitHash == constants.GitHEAD:
		return []string{"log", "--format=%H %s", "-1"}
	default:
		return []string{"log", "--format=%H %s", "--reverse", f.commitHash + "^.." + constants.GitHEAD}
	}
}

// listCommitsForAmend returns commits that will be rewritten.
func listCommitsForAmend(f amendFlags) []model.CommitEntry {
	args := getAmendLogArgs(f)
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
		if e, ok := parseSingleCommitLine(line); ok {
			entries = append(entries, e)
		}
	}
	return entries
}

func parseSingleCommitLine(line string) (model.CommitEntry, bool) {
	if line == "" {
		return model.CommitEntry{}, false
	}
	parts := strings.SplitN(line, " ", 2)
	msg := ""
	if len(parts) > 1 {
		msg = parts[1]
	}
	return model.CommitEntry{SHA: parts[0], Message: msg}, true
}

func getGitAuthorField(sha, format string) string {
	out, err := exec.Command("git", "log", "-1", format, sha).Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not read author info for %s: %v\n", sha, err)
		return ""
	}
	return strings.TrimSpace(string(out))
}

// detectPreviousAuthor reads the author of the first commit in the range.
func detectPreviousAuthor(commits []model.CommitEntry) (string, string) {
	if len(commits) == 0 {
		return "", ""
	}
	sha := commits[0].SHA
	return getGitAuthorField(sha, "--format=%an"), getGitAuthorField(sha, "--format=%ae")
}

// getCurrentBranch returns the current Git branch name.
func getCurrentBranch() string {
	out, err := exec.Command("git", "rev-parse", "--abbrev-ref", constants.GitHEAD).Output()
	if err != nil {
		return constants.DefaultBranch
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

func buildFilterBranchArgs(f amendFlags) []string {
	envFilter := buildEnvFilter(f)
	if f.commitHash == "" {
		return []string{"filter-branch", "-f", "--env-filter", envFilter, "--", constants.GitHEAD}
	}
	return []string{"filter-branch", "-f", "--env-filter", envFilter, f.commitHash + "^.." + constants.GitHEAD}
}

// runFilterBranch executes the git filter-branch command.
func runFilterBranch(f amendFlags, commits []model.CommitEntry) error {
	if f.commitHash == constants.GitHEAD {
		runAmendHead(f)
		return nil
	}
	args := buildFilterBranchArgs(f)
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrAmendFilter, err)
		os.Exit(1)
	}
	return nil
}

// runAmendHead uses git commit --amend for single HEAD commit.
func runAmendHead(f amendFlags) error {
	author := buildAuthorString(f)
	args := []string{"commit", "--amend", "--no-edit", "--author", author}
	cmd := exec.Command("git", args...)
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrAmendCommitAmend, err)
		os.Exit(1)
	}
	return nil
}
