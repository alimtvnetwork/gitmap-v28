// Package cmdpullerror provides rendering utilities for pull error inspection.
package cmdpullerror

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func renderPullErrorHelp() {
	fmt.Printf("\n%s%sGitMap Pull Errors Inspector (gitmap pull-error)%s\n\n",
		constants.ColorCyan, constants.ColorBold, constants.ColorReset)
	printHelpUsage()
	printHelpFlags()
}

func printHelpUsage() {
	fmt.Println("  Usage:")
	fmt.Println("    gitmap pull-error [target] [flags]")
	fmt.Println("    gitmap pull-errors [target] [flags]")
	fmt.Println("    gitmap pulle [target] [flags]")
	fmt.Println("    gitmap pull-e [target] [flags]")
	fmt.Println()
	fmt.Println("  Arguments:")
	fmt.Println("    target                        Repository slug, folder name, or 'all' (default: current repo or all)")
	fmt.Println()
}

func printHelpFlags() {
	fmt.Println("  Flags:")
	fmt.Println("    --json, -j                    Output structured JSON for AI agent ingestion")
	fmt.Println("    --ssh                         Query pull errors across SSH fleet nodes (PAS Formula)")
	fmt.Println("    --limit, -l <n>               Maximum number of records to return (default: 50)")
	fmt.Println("    --help, -h                    Show this help message")
	fmt.Println()
}

func renderNoPullErrors(target string) {
	fmt.Printf("\n  %s✓%s %sNo pull errors recorded for %s.%s\n\n",
		constants.ColorGreen, constants.ColorReset, constants.ColorBold, target, constants.ColorReset)
}

func renderPullErrorsList(records []store.PullErrorRecord, target string) {
	fmt.Printf("\n%s%sPull Error Diagnostics (%d recorded for '%s'):%s\n",
		constants.ColorRed, constants.ColorBold, len(records), target, constants.ColorReset)

	for _, rec := range records {
		renderSinglePullErrorCard(rec)
	}
}

func renderSinglePullErrorCard(rec store.PullErrorRecord) {
	fmt.Printf("\n  %s%s┌── Pull Error Diagnostic ───────────────────────────────────────%s\n",
		constants.ColorRed, constants.ColorBold, constants.ColorReset)
	renderCardHeader(rec)
	renderCardErrorBody(rec)
	renderCardTraceAndRemedy(rec)
	fmt.Printf("  %s└─────────────────────────────────────────────────────────────────%s\n",
		constants.ColorDim, constants.ColorReset)
}

func renderCardHeader(rec store.PullErrorRecord) {
	repoDisplay := rec.RepoSlug
	if rec.RepoPath != "" && rec.RepoPath != rec.RepoSlug {
		repoDisplay = fmt.Sprintf("%s (%s)", rec.RepoSlug, rec.RepoPath)
	}
	fmt.Printf("  │ %sRepo Name:%s          %s\n", constants.ColorCyan, constants.ColorReset, repoDisplay)
	if rec.NodeID != "" {
		fmt.Printf("  │ %sNode ID:%s            %s\n", constants.ColorCyan, constants.ColorReset, rec.NodeID)
	}
	if rec.NodeVersion != "" {
		fmt.Printf("  │ %sGitMap Version:%s     %s\n", constants.ColorCyan, constants.ColorReset, rec.NodeVersion)
	}
	ts := rec.CreatedAt.UTC().Format("2006-01-02 15:04:05 UTC")
	fmt.Printf("  │ %sTimestamp:%s          %s\n", constants.ColorCyan, constants.ColorReset, ts)
}

func renderCardErrorBody(rec store.PullErrorRecord) {
	if rec.ErrorType != "" {
		fmt.Printf("  │ %sClassification:%s     %s%s%s\n",
			constants.ColorYellow, constants.ColorReset, constants.ColorBold, rec.ErrorType, constants.ColorReset)
	}
	if rec.ErrorText != "" {
		fmt.Printf("  │ %sRoot Cause:%s         %s\n", constants.ColorYellow, constants.ColorReset, rec.ErrorText)
	}
}

func renderCardTraceAndRemedy(rec store.PullErrorRecord) {
	if rec.StackTrace != "" {
		fmt.Printf("  │ %sStack Trace / Error Lines:%s\n%s\n",
			constants.ColorRed, constants.ColorReset, indentLines(rec.StackTrace, "  │   "))
	}
	if rec.RemediationCmd != "" {
		fmt.Printf("  │ %sActionable Remediation Command:%s\n  │   %s%s%s\n",
			constants.ColorGreen, constants.ColorReset, constants.ColorBold, rec.RemediationCmd, constants.ColorReset)
	}
}

func indentLines(s, prefix string) string {
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		lines[i] = prefix + l
	}
	return strings.Join(lines, "\n")
}
