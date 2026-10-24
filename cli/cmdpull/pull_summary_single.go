package cmdpull

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

// RenderSingleRepoPullSummary prints a concise summary of single-repo pull outcome.
func RenderSingleRepoPullSummary(state *PullRepoState, cwd string) {
	diag := gitutil.InspectDirtyState(cwd)
	fmt.Println(FormatSingleRepoSummaryLineWithDiag(state, cwd, diag))
	renderIncomingPullChanges(state, cwd)
	if diag.IsDirty {
		renderGitDirtyStatus(cwd)
	}
}

// FormatSingleRepoSummaryLine builds the single-repo summary line.
func FormatSingleRepoSummaryLine(state *PullRepoState, cwd string) string {
	diag := gitutil.InspectDirtyState(cwd)
	return FormatSingleRepoSummaryLineWithDiag(state, cwd, diag)
}

// FormatSingleRepoSummaryLineWithDiag builds the single-repo summary line using precomputed dirty diagnosis.
func FormatSingleRepoSummaryLineWithDiag(state *PullRepoState, cwd string, diag gitutil.DirtyDiagnosis) string {
	repo := resolveSummaryRepoName(state, cwd)
	branch := resolveSummaryBranch(state, cwd)
	shortSHA := resolveSummaryShortSHA(state, cwd)
	commitRange := resolveSummaryCommitRange(state, cwd, shortSHA)
	diffStr := ""
	if state != nil {
		diffStr = formatDiffStat(state.Changes)
	}

	isError, statusDesc := resolveSummaryStatusDesc(state, cwd, shortSHA, commitRange, diffStr)
	isDirty, cleanDirtyText := resolveCleanDirtySummaryWithDiag(state, diag)
	glyph := resolveSummaryGlyph(isError, isDirty)

	return fmt.Sprintf("  %s %s: %s [%s] (%s)", glyph, repo, statusDesc, branch, cleanDirtyText)
}

func resolveSummaryStatusDesc(state *PullRepoState, cwd, shortSHA, commitRange, diffStr string) (bool, string) {
	if state == nil {
		return false, fmt.Sprintf("up-to-date (%s)", shortSHA)
	}

	isFailed := state.Step == PullStepTypeError || state.Changes == "failed"
	if isFailed {
		return true, fmt.Sprintf("pull failed - %s", formatErrorSummary(state.ErrorMsg))
	}

	if state.IsAutoMerged {
		return false, fmt.Sprintf("auto-merge executed %s%s", commitRange, diffStr)
	}

	isUpToDate := state.Step == PullStepTypeUpToDate || state.Changes == "up-to-date"
	if !isUpToDate && state.OldSHA != "" && state.NewSHA != "" && state.OldSHA == state.NewSHA {
		isUpToDate = true
	}

	if isUpToDate {
		return false, fmt.Sprintf("up-to-date (%s)", shortSHA)
	}

	return false, fmt.Sprintf("updated %s%s", commitRange, diffStr)
}

func resolveSummaryGlyph(isError, isDirty bool) string {
	if isError {
		return constants.ColorRed + "✗" + constants.ColorReset
	}
	if isDirty {
		return constants.ColorYellow + "⚠" + constants.ColorReset
	}

	return constants.ColorGreen + "✓" + constants.ColorReset
}

func resolveCleanDirtySummaryWithDiag(state *PullRepoState, diag gitutil.DirtyDiagnosis) (bool, string) {
	isDirty := diag.IsDirty || (state != nil && state.IsDirty)
	if !isDirty {
		return false, "clean"
	}

	reason := diag.SummaryReason
	if reason == "" {
		reason = "uncommitted changes"
	}

	return true, fmt.Sprintf("dirty: %s", reason)
}

func resolveSummaryRepoName(state *PullRepoState, cwd string) string {
	if state != nil && state.RepoName != "" {
		return state.RepoName
	}
	if cwd == "" {
		return "repo"
	}

	base := filepath.Base(cwd)
	if base != "" && base != "." && base != "/" && base != "\\" {
		return base
	}

	return "repo"
}

func resolveSummaryBranch(state *PullRepoState, cwd string) string {
	if state != nil && state.Branch != "" {
		return state.Branch
	}
	if cwd == "" {
		return "main"
	}

	b := gitutil.GetActiveBranch(cwd)
	if b != "" {
		return b
	}

	return "main"
}

func resolveSummaryShortSHA(state *PullRepoState, cwd string) string {
	if state != nil && state.NewSHA != "" && state.NewSHA != "-" {
		return ShortenSHA(state.NewSHA)
	}
	if state != nil && state.OldSHA != "" && state.OldSHA != "-" {
		return ShortenSHA(state.OldSHA)
	}
	if cwd == "" {
		return "head"
	}

	sha := gitutil.GetLastCommitSHA(cwd)
	if sha != "" && sha != "-" {
		return ShortenSHA(sha)
	}

	return "head"
}

func resolveSummaryCommitRange(state *PullRepoState, cwd, shortSHA string) string {
	if state != nil && state.CommitRange != "" && state.CommitRange != "-" {
		return state.CommitRange
	}
	if state != nil && state.OldSHA != "" && state.NewSHA != "" && state.OldSHA != state.NewSHA {
		return FormatCommitRange(state.OldSHA, state.NewSHA)
	}

	return shortSHA
}

func formatDiffStat(changes string) string {
	changes = strings.TrimSpace(changes)
	isUnchanged := changes == "" || changes == "updated" || changes == "up-to-date" || changes == "synced" || changes == "failed"
	if isUnchanged {
		return ""
	}
	parenRes := extractParenDiffStat(changes)
	if parenRes != "" {
		return parenRes
	}
	if strings.Contains(changes, "file") {
		return fmt.Sprintf(" (%s)", changes)
	}

	return fmt.Sprintf(" (%s)", changes)
}

func extractParenDiffStat(changes string) string {
	openParen := strings.Index(changes, "(")
	closeParen := strings.Index(changes, ")")
	if openParen == -1 || closeParen <= openParen {
		return ""
	}
	diffParts := strings.TrimSpace(changes[:openParen])
	fileCount := strings.TrimSpace(changes[openParen+1 : closeParen])

	return formatParenDiffParts(diffParts, fileCount)
}

func formatParenDiffParts(diffParts, fileCount string) string {
	if diffParts == "" || fileCount == "" {
		return ""
	}
	if fileCount == "1" {
		return fmt.Sprintf(" (%s, 1 file)", diffParts)
	}

	return fmt.Sprintf(" (%s, %s files)", diffParts, fileCount)
}

func formatErrorSummary(errMsg string) string {
	errMsg = strings.TrimSpace(errMsg)
	if errMsg == "" {
		return "unknown error"
	}
	if strings.Contains(errMsg, "active unmerged conflict") || strings.Contains(errMsg, "merge conflict") {
		return "merge conflict detected"
	}
	lines := strings.Split(errMsg, "\n")
	first := strings.TrimSpace(lines[0])
	idx := strings.LastIndex(first, ": ")
	if idx == -1 {
		return first
	}

	rest := strings.TrimSpace(first[idx+2:])
	if rest != "" && !strings.HasPrefix(rest, "repo=") {
		return rest
	}

	return first
}
