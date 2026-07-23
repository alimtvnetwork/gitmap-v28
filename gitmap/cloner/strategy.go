// Package cloner — branch selection strategy.
//
// pickCloneStrategy decides whether `git clone` should be invoked with an
// explicit `-b <branch>` flag or without it (letting the remote's default
// HEAD decide). The decision is driven by ScanRecord.BranchSource so that
// untrustworthy values (literal "HEAD", a detached SHA, or "unknown") never
// reach the git command line and produce "Remote branch not found" errors.
package cloner

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// cloneStrategy describes how a clone should be invoked.
type cloneStrategy struct {
	// useBranch is true when `-b <branch>` should be passed to git clone.
	useBranch bool
	// branch is the branch name to check out (only used when useBranch).
	branch string
	// reason is a short human-readable description of why this strategy
	// was chosen. It is propagated into CloneResult.Notes for diagnostics.
	reason string
}

// pickCloneStrategy decides whether to checkout the recorded branch, the
// remote-tracking branch, the repo default, or fall back to the remote's
// default HEAD. Decisions are made from BranchSource so that scans that
// captured a detached / unknown state never produce "branch not found"
// errors during clone.
func pickCloneStrategy(rec model.ScanRecord) cloneStrategy {
	branch := rec.Branch

	switch rec.BranchSource {
	case gitutil.BranchSourceHEAD:
		if branch == "" || branch == gitutil.BranchSourceHEAD {
			return cloneStrategy{
				reason: "branchSource=HEAD but branch empty; using remote default",
			}
		}

		return cloneStrategy{
			useBranch: true,
			branch:    branch,
			reason:    "branchSource=HEAD; checking out detected branch",
		}

	case gitutil.BranchSourceRemoteTracking:
		if branch == "" {
			return cloneStrategy{
				reason: "branchSource=remote-tracking but branch empty; using remote default",
			}
		}

		return cloneStrategy{
			useBranch: true,
			branch:    branch,
			reason:    "branchSource=remote-tracking; checking out tracked branch",
		}

	case gitutil.BranchSourceDefault:
		if branch == "" {
			return cloneStrategy{
				reason: "branchSource=default but branch empty; using remote default",
			}
		}

		return cloneStrategy{
			useBranch: true,
			branch:    branch,
			reason:    "branchSource=default; checking out repo default branch",
		}

	case gitutil.BranchSourceDetached:
		return cloneStrategy{
			reason: "branchSource=detached; using remote default HEAD",
		}

	case gitutil.BranchSourceUnknown, "":
		return cloneStrategy{
			reason: "branchSource=unknown; using remote default HEAD",
		}
	}

	// Unrecognized source label — be safe and skip -b.
	return cloneStrategy{
		reason: "branchSource=" + rec.BranchSource + " (unrecognized); using remote default HEAD",
	}
}

// applyDefaultBranchFallback rewrites records whose recorded
// (Branch, BranchSource) pair would otherwise leave the cloner with
// no branch to check out — i.e., where pickCloneStrategy returns
// useBranch=false. When `fallback` is non-empty, those records are
// rebuilt with Branch=fallback and BranchSource=BranchSourceDefault
// so the existing "trusted default" path takes over and emits
// `git clone -b <fallback> ...`. A breadcrumb is appended to Notes
// so audits and CloneResult.Notes can show that the value came from
// the CLI fallback, not from the original scan.
//
// Records that already have a usable branch are returned unchanged.
// The input slice is never mutated; a fresh slice is returned. When
// `fallback` is empty the input slice is returned as-is, preserving
// the legacy "remote default HEAD" behavior bit-for-bit.
func applyDefaultBranchFallback(records []model.ScanRecord, fallback string) []model.ScanRecord {
	if len(fallback) == 0 {

		return records
	}
	out := make([]model.ScanRecord, len(records))
	for i, rec := range records {
		if pickCloneStrategy(rec).useBranch {
			out[i] = rec

			continue
		}
		patched := rec
		patched.Branch = fallback
		patched.BranchSource = gitutil.BranchSourceDefault
		patched.Notes = appendFallbackNote(rec.Notes, fallback)
		out[i] = patched
	}

	return out
}

// appendFallbackNote appends a breadcrumb that the default-branch
// fallback was applied. Kept as its own helper so the formatting
// stays consistent across audit / clone-result surfaces and so the
// 15-line function budget for applyDefaultBranchFallback is preserved.
func appendFallbackNote(existing, fallback string) string {
	note := "default-branch fallback applied: " + fallback
	if len(existing) == 0 {

		return note
	}

	return existing + "; " + note
}
