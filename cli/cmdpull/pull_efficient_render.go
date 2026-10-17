package cmdpull

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

const (
	minConciseRepoColWidth = 26
	conciseReservedColGap  = 30

	treeBranch       = "├── "
	treeTerminal     = "└── "
	treeContinuation = "│   "
	treeIndent       = "    "
)

// renderConciseActiveResults prints bullet list with aligned columns to stdout.
func renderConciseActiveResults(states []*PullRepoState, allRecords ...[]model.ScanRecord) {
	RenderConciseActiveResultsTo(os.Stdout, states, allRecords...)
}

// RenderConciseActiveResultsTo renders concise repo states grouped into categories.
func RenderConciseActiveResultsTo(w io.Writer, states []*PullRepoState, allRecords ...[]model.ScanRecord) {
	deduped := DeduplicateRepoStates(states)
	colWidth := ResolveConciseRepoColWidth(deduped, allRecords...)
	cat := categorizeRepoStates(deduped)
	if len(cat.updated) == 0 && len(cat.dirty) == 0 && len(cat.failed) == 0 {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "  (all repositories are up-to-date)")
		return
	}
	collisions := detectRepoCollisions(deduped, allRecords...)
	renderCategorizedGroups(w, colWidth, cat, collisions)
}

type categorizedStates struct {
	updated  []*PullRepoState
	upToDate []*PullRepoState
	dirty    []*PullRepoState
	failed   []*PullRepoState
}

func categorizeRepoStates(states []*PullRepoState) categorizedStates {
	var cat categorizedStates
	for _, s := range states {
		if s == nil {
			continue
		}
		switch {
		case s.IsDirty || s.Changes == "dirty":
			cat.dirty = append(cat.dirty, s)
		case isFailedRepoState(s):
			cat.failed = append(cat.failed, s)
		case isUpdatedRepoState(s):
			cat.updated = append(cat.updated, s)
		case !s.IsDirty && s.ErrorMsg == "" && (s.Step == PullStepTypeUpToDate || s.Changes == "up-to-date" || ResolveRepoStatusLabel(s.Changes) == "up-to-date"):
			cat.upToDate = append(cat.upToDate, s)
		}
	}
	return cat
}

func renderCategorizedGroups(w io.Writer, colWidth int, cat categorizedStates, collisions map[string]bool) {
	if len(cat.updated) > 0 {
		renderUpdatedGroup(w, colWidth, cat.updated, collisions)
	}
	if len(cat.upToDate) > 0 {
		renderUpToDateGroup(w, cat.upToDate, collisions)
	}
	if len(cat.dirty) > 0 {
		renderDirtyGroup(w, colWidth, cat.dirty, collisions)
	}
	if len(cat.failed) > 0 {
		renderFailedGroup(w, colWidth, cat.failed, collisions)
	}
}

func renderUpToDateGroup(w io.Writer, upToDate []*PullRepoState, collisions map[string]bool) {
	fmt.Fprintf(w, "\n  %s%s✓ Up-to-Date Repositories (%d):%s\n", constants.ColorGreen, constants.ColorBold, len(upToDate), constants.ColorReset)
	names := make([]string, 0, len(upToDate))
	for _, s := range upToDate {
		if s == nil {
			continue
		}
		displayName := resolveRepoDisplayName(s, collisions)
		names = append(names, displayName)
	}
	termWidth := detectTerminalWidth()
	wrapped := FormatWrappedInactiveList(names, "    ", termWidth-2)
	fmt.Fprintf(w, "%s%s%s\n", constants.ColorDim, wrapped, constants.ColorReset)
}

func isFailedRepoState(s *PullRepoState) bool {
	if s == nil {
		return false
	}
	return s.ErrorMsg != "" || s.Step == PullStepTypeError || s.Step == PullStepTypeConflict || s.Changes == "failed"
}

// DeduplicateRepoStates ensures each repository appears at most once in output states.
// When collisions occur, error/failure states take precedence over up-to-date states.
func DeduplicateRepoStates(states []*PullRepoState) []*PullRepoState {
	if len(states) <= 1 {
		return states
	}
	seen := make(map[string]int, len(states))
	unique := make([]*PullRepoState, 0, len(states))
	for _, s := range states {
		unique = appendUniqueRepoState(unique, s, seen)
	}
	return unique
}

func appendUniqueRepoState(unique []*PullRepoState, s *PullRepoState, seen map[string]int) []*PullRepoState {
	if s == nil {
		return unique
	}
	key := resolveRepoStateKey(s)
	if idx, exists := seen[key]; exists {
		mergeRepoStateIfPrioritized(unique, idx, s)
		return unique
	}
	seen[key] = len(unique)
	return append(unique, s)
}

func resolveRepoStateKey(s *PullRepoState) string {
	key := CanonicalRepoPathKey(s.RepoPath)
	if key != "" {
		return key
	}
	return strings.ToLower(strings.TrimSpace(s.RepoName))
}

func isActionableRepoState(s *PullRepoState) bool {
	if s == nil {
		return false
	}
	return s.IsDirty || s.ErrorMsg != "" || s.Changes == "failed" || s.Step == PullStepTypeError || s.Step == PullStepTypeConflict
}

func mergeRepoStateIfPrioritized(unique []*PullRepoState, existingIdx int, incoming *PullRepoState) {
	existing := unique[existingIdx]
	hasIncomingAction := isActionableRepoState(incoming)
	hasExistingAction := isActionableRepoState(existing)
	if hasIncomingAction && !hasExistingAction {
		unique[existingIdx] = incoming
	}
}

func isUpdatedRepoState(s *PullRepoState) bool {
	if s.Step == PullStepTypeFastForward || s.Step == PullStepTypeMerging {
		return true
	}
	status := ResolveRepoStatusLabel(s.Changes)
	return status != "up-to-date" && status != "dirty" && status != "failed"
}

func renderUpdatedGroup(w io.Writer, colWidth int, updated []*PullRepoState, collisions map[string]bool) {
	fmt.Fprintf(w, "\n  %s%s Updated Repositories (%d):%s\n", constants.ColorGreen, constants.ColorBold, len(updated), constants.ColorReset)
	for _, s := range updated {
		statusLabel := resolveUpdatedStatusLabel(s)
		displayName := resolveRepoDisplayName(s, collisions)
		fmt.Fprintln(w, FormatConciseActiveResultLine(colWidth, displayName, statusLabel))
		if s.CommitRange != "" {
			fmt.Fprintf(w, "        %s↳ Commits: %s%s\n", constants.ColorDim, s.CommitRange, constants.ColorReset)
		}
	}
}

func resolveUpdatedStatusLabel(s *PullRepoState) string {
	statusLabel := ResolveRepoStatusLabel(s.Changes)
	if statusLabel != "up-to-date" {
		return statusLabel
	}
	if s.Step == PullStepTypeFastForward {
		return "fast-forward"
	}
	if s.Step == PullStepTypeMerging {
		return "merged"
	}
	return statusLabel
}

func renderDirtyGroup(w io.Writer, colWidth int, dirty []*PullRepoState, collisions map[string]bool) {
	fmt.Fprintf(w, "\n  %s%s Dirty Repositories (%d):%s\n", constants.ColorYellow, constants.ColorBold, len(dirty), constants.ColorReset)
	for _, s := range dirty {
		renderSingleDirtyItem(w, colWidth, s, collisions)
	}
}

func renderSingleDirtyItem(w io.Writer, colWidth int, s *PullRepoState, collisions map[string]bool) {
	displayName := resolveRepoDisplayName(s, collisions)
	fmt.Fprintln(w, FormatConciseActiveResultLine(colWidth, displayName, "dirty"))
	maybeRenderDirtyFiles(w, s.RepoPath)
	renderDirtyReasonsAndOptions(w, s)
}

func maybeRenderDirtyFiles(w io.Writer, repoPath string) {
	if repoPath != "" {
		diag := gitutil.InspectDirtyState(repoPath)
		renderItemizedDirtyFiles(w, diag)
	}
}

func renderDirtyReasonsAndOptions(w io.Writer, s *PullRepoState) {
	errDetails := ResolvePullErrorDetails(s)
	if errDetails != "" {
		fmt.Fprintf(w, "    %sReason: %s%s%s\n", treeBranch, constants.ColorDim, errDetails, constants.ColorReset)
	}
	renderStructuredSolutionsSubtree(w, s)
}

func renderItemizedDirtyFiles(w io.Writer, diag gitutil.DirtyDiagnosis) {
	for i, f := range diag.ModifiedFiles {
		if i >= 5 {
			fmt.Fprintf(w, "        %s... and %d more modified files%s\n", constants.ColorDim, len(diag.ModifiedFiles)-5, constants.ColorReset)
			break
		}
		cleanFile := strings.ReplaceAll(filepath.ToSlash(f), "\\", "/")
		fmt.Fprintf(w, "        %smodified: %s%s\n", constants.ColorDim, cleanFile, constants.ColorReset)
	}
	for i, f := range diag.UntrackedFiles {
		if i >= 5 {
			fmt.Fprintf(w, "        %s... and %d more untracked files%s\n", constants.ColorDim, len(diag.UntrackedFiles)-5, constants.ColorReset)
			break
		}
		cleanFile := strings.ReplaceAll(filepath.ToSlash(f), "\\", "/")
		fmt.Fprintf(w, "        %suntracked: %s%s\n", constants.ColorDim, cleanFile, constants.ColorReset)
	}
}
