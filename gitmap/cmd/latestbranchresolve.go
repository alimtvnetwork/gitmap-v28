// Package cmd — latest-branch resolve helpers.
package cmd

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/gitutil"
)

// latestBranchResult holds the resolved display data for the latest branch.
type latestBranchResult struct {
	branchNames    []string
	selectedRemote string
	shortSha       string
	commitDate     string
	latest         gitutil.RemoteBranchInfo
}

// resolveLatestResult picks the latest branch and resolves its name.
func resolveLatestResult(items []gitutil.RemoteBranchInfo, cfg latestBranchConfig) latestBranchResult {
	latest := items[0]
	selectedRemote := extractRemoteName(latest.RemoteRef)
	branchNames := resolveBranchNames(latest.Sha, selectedRemote, cfg.containsFallback)

	return latestBranchResult{
		branchNames:    branchNames,
		selectedRemote: selectedRemote,
		shortSha:       gitutil.TruncSha(latest.Sha),
		commitDate:     gitutil.FormatDisplayDate(latest.CommitDate),
		latest:         latest,
	}
}

// extractRemoteName extracts the remote name from a ref (e.g. "origin/main" → "origin").
func extractRemoteName(remoteRef string) string {
	idx := strings.Index(remoteRef, "/")
	if idx >= 0 {

		return remoteRef[:idx]
	}

	return remoteRef
}

// resolveBranchNames resolves human-readable branch names from a SHA.
func resolveBranchNames(sha, remote string, containsFallback bool) []string {
	names := gitutil.ResolvePointsAt(sha, remote)
	if len(names) > 0 {

		return names
	}
	if containsFallback {
		names = gitutil.ResolveContains(sha, remote)
		if len(names) > 0 {

			return names
		}
	}

	return []string{constants.LBUnknownBranch}
}
