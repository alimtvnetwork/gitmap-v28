package cmdpipeline

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

func isTagRef(ref string) bool {
	trimmed := strings.TrimSpace(ref)
	if len(trimmed) == 0 {
		return false
	}
	lower := strings.ToLower(trimmed)
	if strings.HasPrefix(lower, "refs/tags/") {
		return true
	}
	if strings.HasPrefix(lower, "v") && isSemverLike(trimmed) {
		return true
	}

	return false
}

func getBranchPrecedence(branch string) int {
	clean := strings.TrimSpace(branch)
	if len(clean) == 0 {
		return 0
	}
	if isTagRef(clean) {
		return 1
	}
	if strings.HasPrefix(strings.ToLower(clean), "release/") {
		return 2
	}
	active := gitutil.GetActiveBranch(".")
	isMatchingActive := len(active) > 0 && active != "-" && strings.EqualFold(clean, active)
	if isMatchingActive {
		return 4
	}

	return 3
}

func resolveActiveOrMainBranch() string {
	active := gitutil.GetActiveBranch(".")
	isValidActive := len(active) > 0 && active != "-" && active != "(detached)"
	if isValidActive {
		return active
	}

	return "main"
}

func resolveBestBranchForRuns(runs []ghRunItem) string {
	bestBranch := ""
	bestPrec := 0
	for _, r := range runs {
		prec := getBranchPrecedence(r.HeadBranch)
		if prec > bestPrec {
			bestPrec = prec
			bestBranch = r.HeadBranch
		}
	}
	if bestPrec <= 1 {
		return resolveActiveOrMainBranch()
	}

	return bestBranch
}

func resolveLatestBranchName(p *PipelineErrorLogsPayload, runs []ghRunItem) string {
	hasNonTagBranch := len(p.Branch) > 0 && !isTagRef(p.Branch)
	if hasNonTagBranch {
		return p.Branch
	}
	branch := resolveBestBranchForRuns(runs)
	if len(branch) > 0 {
		return branch
	}

	return resolveActiveOrMainBranch()
}
