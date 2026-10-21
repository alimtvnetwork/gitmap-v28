package cmdpull

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func renderFailedGroup(w io.Writer, colWidth int, failed []*PullRepoState, collisions map[string]bool) {
	deduped := DeduplicateRepoStates(failed)
	fmt.Fprintf(w, "\n  %s%s Failed Repositories (%d):%s\n", constants.ColorRed, constants.ColorBold, len(deduped), constants.ColorReset)
	for _, s := range deduped {
		renderSingleFailedItem(w, colWidth, s, collisions)
	}
	if len(deduped) == 0 {
		return
	}

	renderUnifiedResolutionFooters(w, deduped)
}

func renderUnifiedResolutionFooters(w io.Writer, deduped []*PullRepoState) {
	if hasAnyCredOrAuthFailure(deduped) {
		fmt.Fprintf(w, "\n  %s🔐 Unified Credential & SSH Resolution:%s To fix credential store and convert all failed repos to SSH at once:\n",
			constants.ColorYellow, constants.ColorReset)
		fmt.Fprintf(w, "     %s%sgitmap fix-credential && gitmap fix --all%s  (or: %sgitmap pa --ssh%s)\n",
			constants.ColorYellow, constants.ColorBold, constants.ColorReset,
			constants.ColorDim, constants.ColorReset)
	}
	fmt.Fprintf(w, "\n  %s💡 Unified Resolution:%s To resolve all %d failed/missing repositories together:\n",
		constants.ColorYellow, constants.ColorReset, len(deduped))
	fmt.Fprintf(w, "     %s%sgitmap fix --all%s  (or: %sgitmap pf --all%s)\n\n",
		constants.ColorYellow, constants.ColorBold, constants.ColorReset,
		constants.ColorDim, constants.ColorReset)
}

func hasAnyCredOrAuthFailure(deduped []*PullRepoState) bool {
	for _, s := range deduped {
		if isStateCredOrAuthFailure(s) {
			return true
		}
	}

	return false
}

func isStateCredOrAuthFailure(s *PullRepoState) bool {
	if s == nil {
		return false
	}

	return isWincredmanFailure(s.ErrorMsg) || isCacheFailure(s.ErrorMsg) || isAuthFailure(s.ErrorMsg)
}

func renderSingleFailedItem(w io.Writer, colWidth int, s *PullRepoState, collisions map[string]bool) {
	displayName := resolveRepoDisplayName(s, collisions)
	fmt.Fprintln(w, FormatConciseActiveResultLine(colWidth, displayName, "failed"))
	errDetails := ResolvePullErrorDetails(s)
	if errDetails == "" {
		errDetails = "pull execution failed"
	}
	fmt.Fprintf(w, "    %sReason: %s%s%s\n", treeBranch, constants.ColorDim, errDetails, constants.ColorReset)
	renderStructuredSolutionsSubtree(w, s)
	repoTarget := resolveFallbackRepoTarget(s)
	fmt.Fprintf(w, "    %sDiagnostic: To inspect stack trace: %sgitmap pull-error %s%s (or: %sgitmap pe%s)\n",
		treeTerminal, constants.ColorYellow, repoTarget, constants.ColorReset, constants.ColorYellow, constants.ColorReset)
}

func resolveFallbackRepoTarget(s *PullRepoState) string {
	if s.RepoName != "" {
		return s.RepoName
	}
	if s.RepoPath != "" {
		return filepath.Base(s.RepoPath)
	}

	return "all"
}

func renderStructuredSolutionsSubtree(w io.Writer, s *PullRepoState) {
	renderStructuredSolutionsSubtreeWithTerminal(w, s, false)
}

func resolveDefaultOptions(s *PullRepoState) []RemediationOption {
	structured := ResolveStructuredRemediation(s)
	if len(structured.Options) > 0 {
		return structured.Options
	}
	repoTarget := resolveFallbackRepoTarget(s)

	return []RemediationOption{
		{OptionNumber: 1, Title: "Auto-Fix", Command: "gitmap fix " + repoTarget},
		{OptionNumber: 2, Title: "Inspect Status", Command: "gitmap status " + repoTarget},
	}
}

func resolveTreePrefixes(isTerminal bool) (string, string) {
	if isTerminal {
		return treeTerminal, treeIndent
	}

	return treeBranch, treeContinuation
}

func renderStructuredSolutionsSubtreeWithTerminal(w io.Writer, s *PullRepoState, isTerminal bool) {
	opts := resolveDefaultOptions(s)
	headConn, parentPrefix := resolveTreePrefixes(isTerminal)
	fmt.Fprintf(w, "    %sSolutions:\n", headConn)
	for i, opt := range opts {
		renderSolutionOptionBranchWithParent(w, opt, i, len(opts), parentPrefix)
	}
}

func renderSolutionOptionBranch(w io.Writer, opt RemediationOption, idx, total int) {
	renderSolutionOptionBranchWithParent(w, opt, idx, total, treeContinuation)
}

func renderSolutionOptionBranchWithParent(w io.Writer, opt RemediationOption, idx, total int, parentPrefix string) {
	isLast := idx == total-1
	tier1Conn := treeBranch
	tier2Prefix := treeContinuation
	if isLast {
		tier1Conn = treeTerminal
		tier2Prefix = treeIndent
	}
	optNum := opt.OptionNumber
	if optNum <= 0 {
		optNum = idx + 1
	}
	fmt.Fprintf(w, "    %s%sOption %d (%s):\n", parentPrefix, tier1Conn, optNum, opt.Title)
	fmt.Fprintf(w, "    %s%s%s%s%s%s\n",
		parentPrefix, tier2Prefix, treeTerminal,
		constants.ColorYellow, opt.Command, constants.ColorReset)
}
