// Package cmd — `gitmap lb --switch` checkout post-step.
package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/gitutil"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
)

// maybeSwitchToLatest is the post-output hook for `gitmap lb --switch`
// (and its `-s` short form). It runs AFTER dispatchLatestOutput so the
// user still sees the normal latest-branch report, then we attempt the
// checkout and surface git's own status line.
//
// We intentionally only act when shouldSwitch is true so the default
// `gitmap lb` invocation stays a pure read-only query — no surprise
// working-tree mutations.
//
// Exits non-zero if checkout fails so shell pipelines (`gitmap lb -s &&
// make build`) short-circuit cleanly. Any "would lose your changes"
// errors from git are passed through verbatim — we don't try to be
// clever about stashing.
func maybeSwitchToLatest(result latestBranchResult, cfg latestBranchConfig) {
	if !cfg.shouldSwitch {
		return
	}
	target := pickSwitchTarget(result)
	if target == "" {
		cliexit.HandleError(apperror.NewSimple(constants.ErrLatestBranchSwitchNoTarget, "E9000"), 1)
	}
	fmt.Printf(constants.MsgLatestBranchSwitching, target)
	out, err := gitutil.CheckoutBranch(".", target)
	if len(out) > 0 {
		fmt.Println(out)
	}
	if err != nil {
		cliexit.HandleError(apperror.NewSimple(constants.ErrLatestBranchSwitchFailed, "E9000"), 1)
	}
}

// pickSwitchTarget chooses the branch name to feed to `git checkout`.
// Preference order:
//
//  1. The first resolved human-readable branch name (already stripped
//     of any "origin/" prefix by gitutil.ResolvePointsAt).
//  2. The remote ref itself with the remote prefix stripped — this is
//     the fallback when --contains-fallback is off and points-at
//     returned nothing useful (e.g. detached or orphan refs).
//
// Returns "" when neither path produces a usable name; the caller
// surfaces that as a clear "no target" error.
func pickSwitchTarget(result latestBranchResult) string {
	for _, name := range result.branchNames {
		if name != "" && name != constants.LBUnknownBranch {
			return name
		}
	}
	if result.latest.RemoteRef != "" {
		return gitutil.StripRemotePrefix(result.latest.RemoteRef)
	}

	return ""
}
