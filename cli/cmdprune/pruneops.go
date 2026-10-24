package cmdprune

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdprune/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/release"
)

// staleBranch holds a release branch and its matching tag.
type staleBranch struct {
	name string
	tag  string
}

var (
	listBranchesFn       = listReleaseBranches
	tagExistsFn          = release.TagExistsLocally
	deleteLocalBranchFn  = deleteLocalBranch
	deleteRemoteBranchFn = deleteRemoteBranch
)

// listReleaseBranches returns all local branches matching release/*.
func listReleaseBranches() []string {
	pattern := constants.ReleaseBranchPrefix + constants.GitWildcard
	cmd := exec.Command(constants.GitBin,
		constants.GitBranch, constants.GitBranchListFlag, pattern)
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf(constants.ErrPruneListBranches, err)

		return nil
	}

	return parseGitBranchOutput(out)
}

// parseGitBranchOutput splits and trims git branch output lines.
func parseGitBranchOutput(out []byte) []string {
	raw := strings.TrimSpace(string(out))
	if len(raw) == 0 {
		return nil
	}

	return extractNonEmptyBranches(strings.Split(raw, constants.Newline))
}

func extractNonEmptyBranches(lines []string) []string {
	var branches []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) > 0 {
			branches = append(branches, trimmed)
		}
	}

	return branches
}

// filterStaleBranches returns branches whose corresponding tag exists.
func filterStaleBranches(branches []string) []staleBranch {
	var stale []staleBranch
	for _, branch := range branches {
		tag := branchToTag(branch)
		if tagExistsFn(tag) {
			stale = append(stale, staleBranch{name: branch, tag: tag})
		}
	}

	return stale
}

// branchToTag converts release/v2.20.0 to v2.20.0.
func branchToTag(branch string) string {
	return strings.TrimPrefix(branch, constants.ReleaseBranchPrefix)
}

// printStaleBranches displays the stale branch list.
func printStaleBranches(stale []staleBranch) {
	fmt.Printf(constants.MsgPruneStaleHeader, len(stale))
	for _, sb := range stale {
		fmt.Printf(constants.MsgPruneStaleItem, sb.name, sb.tag)
	}
}

// deleteLocalBranch force-deletes a local branch.
func deleteLocalBranch(name string) error {
	cmd := exec.Command(constants.GitBin,
		constants.GitBranch, constants.GitBranchDeleteFlag, name)

	return cmd.Run()
}

// deleteRemoteBranch attempts to delete a remote branch.
func deleteRemoteBranch(name string) {
	cmd := exec.Command(constants.GitBin,
		constants.GitPush, constants.GitOrigin,
		constants.GitPushDeleteFlag, name)
	if err := cmd.Run(); err != nil {
		fmt.Printf(constants.MsgPruneRemoteWarn, name, err)

		return
	}

	fmt.Printf(constants.MsgPruneRemoteDelete, name)
}
