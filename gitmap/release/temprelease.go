package release

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// TempReleaseCommit holds commit info for temp-release branch creation.
type TempReleaseCommit struct {
	SHA     string
	Short   string
	Message string
}

// ListRecentCommits returns the last N commits from HEAD (oldest first).
func ListRecentCommits(count int) ([]TempReleaseCommit, error) {
	cmd := exec.Command(constants.GitBin, constants.GitLog,
		fmt.Sprintf("-%d", count),
		"--format=%H|%s",
		"--reverse")

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list recent commits: %w", err)
	}

	return parseTempReleaseCommitLines(strings.TrimSpace(string(out))), nil
}

// parseTempReleaseCommitLines parses git log output into commit structs.
func parseTempReleaseCommitLines(output string) []TempReleaseCommit {
	if len(output) == 0 {
		return nil
	}

	lines := strings.Split(output, "\n")
	commits := make([]TempReleaseCommit, 0, len(lines))

	for _, line := range lines {
		c := parseOneCommitLine(line)
		if len(c.SHA) > 0 {
			commits = append(commits, c)
		}
	}

	return commits
}

// parseOneCommitLine extracts SHA and message from a single log line.
func parseOneCommitLine(line string) TempReleaseCommit {
	parts := strings.SplitN(line, "|", 2)
	if len(parts) < 2 {
		return TempReleaseCommit{}
	}

	sha := strings.TrimSpace(parts[0])
	short := sha
	if len(short) > constants.ShaDisplayLength {
		short = short[:constants.ShaDisplayLength]
	}

	return TempReleaseCommit{
		SHA:     sha,
		Short:   short,
		Message: strings.TrimSpace(parts[1]),
	}
}

// CreateBranchFromSHA creates a branch at a specific commit without checkout.
func CreateBranchFromSHA(branchName, sha string) error {
	cmd := exec.Command(constants.GitBin, constants.GitBranch, branchName, sha)
	cmd.Stderr = nil

	return cmd.Run()
}

// PushBranches pushes multiple branches to origin in a single command.
func PushBranches(branches []string) error {
	args := append([]string{constants.GitPush, constants.GitOrigin}, branches...)

	return runGitCmd(args...)
}

// DeleteLocalBranch deletes a local branch.
func DeleteLocalBranch(branch string) error {
	cmd := exec.Command(constants.GitBin, constants.GitBranch, "-D", branch)

	return cmd.Run()
}

// DeleteRemoteBranches deletes multiple branches from the remote.
func DeleteRemoteBranches(branches []string) error {
	args := append([]string{constants.GitPush, constants.GitOrigin, "--delete"}, branches...)

	return runGitCmd(args...)
}

// ListTempReleaseBranches returns all local branches matching temp-release/*.
func ListTempReleaseBranches() ([]string, error) {
	cmd := exec.Command(constants.GitBin, constants.GitBranch,
		constants.GitBranchListFlag, constants.TempReleaseBranchPrefix+"*")

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return parseBranchOutput(strings.TrimSpace(string(out))), nil
}

// parseBranchOutput extracts branch names from git branch --list output.
func parseBranchOutput(output string) []string {
	if len(output) == 0 {
		return nil
	}

	lines := strings.Split(output, "\n")
	var branches []string

	for _, line := range lines {
		branch := strings.TrimSpace(line)
		branch = strings.TrimPrefix(branch, "* ")

		if len(branch) > 0 {
			branches = append(branches, branch)
		}
	}

	return branches
}
