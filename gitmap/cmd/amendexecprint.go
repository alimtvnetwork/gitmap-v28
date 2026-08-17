// Package cmd — amendexecprint.go handles output and display for amend operations.
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// buildEnvFilter constructs the env-filter script for filter-branch.
func buildEnvFilter(f amendFlags) string {
	var lines []string

	if f.name != "" {
		lines = append(lines, "export GIT_AUTHOR_NAME='"+f.name+"'")
		lines = append(lines, "export GIT_COMMITTER_NAME='"+f.name+"'")
	}

	if f.email != "" {
		lines = append(lines, "export GIT_AUTHOR_EMAIL='"+f.email+"'")
		lines = append(lines, "export GIT_COMMITTER_EMAIL='"+f.email+"'")
	}

	return strings.Join(lines, "\n")
}

// getGitConfig reads a git config value.
func getGitConfig(key string) string {
	out, err := exec.Command("git", "config", key).Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not read git %s: %v\n", key, err)
	}

	return strings.TrimSpace(string(out))
}

// buildAuthorString creates the --author flag value.
func buildAuthorString(f amendFlags) string {
	name := f.name
	email := f.email

	if name == "" {
		name = getGitConfig("user.name")
	}

	if email == "" {
		email = getGitConfig("user.email")
	}

	return name + " <" + email + ">"
}

// runForcePush executes git push --force-with-lease.
func runForcePush() {
	cmd := exec.Command("git", "push", "--force-with-lease")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrAmendForcePush, err)

		return
	}

	fmt.Print(constants.MsgAmendForcePush)
}

// printAmendHeader outputs the operation header.
func printAmendHeader(f amendFlags, commits []model.CommitEntry, branch, prevName, prevEmail string) {
	if f.commitHash == "" {
		fmt.Printf(constants.MsgAmendHeaderAll, len(commits), branch)
	} else {
		fmt.Printf(constants.MsgAmendHeader, len(commits), commits[0].SHA[:7], commits[len(commits)-1].SHA[:7], branch)
	}

	oldAuthor := prevName + " <" + prevEmail + ">"
	newAuthor := buildAuthorString(f)
	fmt.Printf(constants.MsgAmendAuthor, oldAuthor, newAuthor)
}

// printAmendProgress outputs per-commit progress lines.
func printAmendProgress(commits []model.CommitEntry) {
	for i, c := range commits {
		sha := c.SHA
		if len(sha) > 7 {
			sha = sha[:7]
		}

		fmt.Printf(constants.MsgAmendProgress, i+1, len(commits), sha, c.Message)
	}
}

// printAmendDryRun outputs dry-run preview.
func printAmendDryRun(commits []model.CommitEntry, f amendFlags, prevName, prevEmail string) {
	fmt.Printf(constants.MsgAmendDryHeader, len(commits))

	for i, c := range commits {
		sha := c.SHA
		if len(sha) > 7 {
			sha = sha[:7]
		}

		fmt.Printf(constants.MsgAmendDryLine, i+1, sha, c.Message, prevName, prevEmail)
	}

	fmt.Print(constants.MsgAmendDrySkip)
}
