// Package release handles version parsing, release workflows,
// GitHub integration, and release metadata management.
package release

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/verbose"
)

// CreateBranch creates a release branch from the given source ref.
func CreateBranch(branchName, sourceRef string) error {
	if verbose.IsEnabled() {
		verbose.Get().Log("git: creating branch %s from %s", branchName, sourceRef)
	}

	args := []string{constants.GitCheckout, constants.GitBranchFlag, branchName}
	if len(sourceRef) > 0 {
		args = append(args, sourceRef)
	}

	return runGitCmd(args...)
}

// CreateTag creates an annotated git tag.
func CreateTag(tag, message string) error {
	if verbose.IsEnabled() {
		verbose.Get().Log("git: creating tag %s", tag)
	}

	return runGitCmd(constants.GitTag, constants.GitTagAnnotateFlag, tag, constants.GitTagMessageFlag, message)
}

// PushBranchAndTag pushes the branch and tag to origin.
func PushBranchAndTag(branchName, tag string) error {
	if verbose.IsEnabled() {
		verbose.Get().Log("git: pushing branch %s to origin", branchName)
	}

	err := runGitCmd(constants.GitPush, constants.GitOrigin, branchName)
	if err != nil {
		return fmt.Errorf("push branch: %w", err)
	}

	if verbose.IsEnabled() {
		verbose.Get().Log("git: pushing tag %s to origin", tag)
	}

	err = runGitCmd(constants.GitPush, constants.GitOrigin, tag)
	if err != nil {
		return fmt.Errorf("push tag: %w", err)
	}

	return nil
}

// CheckoutBranch checks out an existing branch.
func CheckoutBranch(branch string) error {
	return runGitCmd(constants.GitCheckout, branch)
}

// FetchBranch fetches the latest of a remote branch.
func FetchBranch(branch string) error {
	return runGitCmd(constants.GitFetch, constants.GitOrigin, branch)
}

// ResolveSourceRef returns the ref to use as the release base.
func ResolveSourceRef(commit, branch string) (string, string, error) {
	if len(commit) > 0 {
		return resolveFromCommit(commit)
	}
	if len(branch) > 0 {
		return resolveFromBranch(branch)
	}

	return resolveFromHead()
}

// resolveFromCommit validates and returns the commit ref.
func resolveFromCommit(commit string) (string, string, error) {
	if CommitExists(commit) {
		if verbose.IsEnabled() {
			verbose.Get().Log("source: using commit %s", commit)
		}
		return commit, constants.GitCommitPrefix + commit, nil
	}

	return "", "", fmt.Errorf("commit %s not found", commit)
}

// resolveFromBranch fetches and returns the branch tip.
func resolveFromBranch(branch string) (string, string, error) {
	err := FetchBranch(branch)
	if err != nil {
		return "", "", fmt.Errorf("branch %s not found: %w", branch, err)
	}

	if verbose.IsEnabled() {
		verbose.Get().Log("source: using branch %s (origin/%s)", branch, branch)
	}

	return constants.GitOriginPrefix + branch, branch, nil
}

// resolveFromHead returns HEAD as the source ref.
func resolveFromHead() (string, string, error) {
	branchName, err := CurrentBranchName()
	if err != nil {
		if verbose.IsEnabled() {
			verbose.Get().Log("source: using detached HEAD")
		}
		return constants.GitHEAD, constants.GitHEAD, nil
	}

	if verbose.IsEnabled() {
		verbose.Get().Log("source: using HEAD on branch %s", branchName)
	}

	return constants.GitHEAD, branchName, nil
}

// runGitCmd executes a git command, forwards stdout, and pipes stderr
// through filteredStderrWriter so cosmetic git warnings (see
// constants.GitStderrNoisePatterns) never reach the user's terminal.
func runGitCmd(args ...string) error {
	cmd := exec.Command(constants.GitBin, args...)
	cmd.Stdout = os.Stdout
	stderr := newFilteredStderr(os.Stderr)
	cmd.Stderr = stderr

	runErr := cmd.Run()
	if flushErr := stderr.Flush(); flushErr != nil && runErr == nil {
		return flushErr
	}

	return runErr
}
