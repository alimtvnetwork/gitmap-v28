// Package gitutil — branch checkout helpers.
package gitutil

import (
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// CheckoutBranch runs `git -C <repoPath> checkout <branch>` and returns
// the combined stdout+stderr so callers can surface the user-friendly
// "Switched to branch 'foo'" / "Already on 'foo'" lines git prints to
// stderr. We deliberately use CombinedOutput (not Output) because git's
// progress and status messages all go to stderr.
//
// The branch name is taken verbatim — strip any "origin/" prefix at the
// caller (see cmd/latestbranchswitch.go) so this helper stays a thin
// pass-through with one job.
func CheckoutBranch(repoPath, branch string) (string, error) {
	cmd := exec.Command(constants.GitBin, constants.GitCheckout, branch)
	cmd.Dir = repoPath
	out, err := cmd.CombinedOutput()

	return strings.TrimRight(string(out), "\n"), err
}

// StripRemotePrefix already lives in latestbranch.go (same package);
// callers in this file reuse it directly.

// ResolveDefaultBranchName returns the repo's default branch name by
// asking `git symbolic-ref refs/remotes/origin/HEAD` first, then
// falling back to constants.DefaultBranch. The returned name is
// already stripped of the "origin/" prefix and ready for `git checkout`.
//
// This is the engine behind `gitmap branch default` / `gitmap b def`.
// Kept separate from DetectBranch (which has different semantics —
// it returns the *current* branch with a "default" sentinel only as
// a last-resort fallback).
func ResolveDefaultBranchName(repoPath string) string {
	out, err := runGit(repoPath, constants.GitSymbolicRef,
		constants.GitRefsRemotesOriginHEAD)
	if err == nil {
		name := strings.TrimSpace(out)
		if len(name) > 0 {
			return StripRemotePrefix(strings.TrimPrefix(name, "refs/remotes/"))
		}
	}

	return constants.DefaultBranch
}
