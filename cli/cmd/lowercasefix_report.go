// Package cmd — lowercasefix_report.go formats status logs and summaries for lowercase fix.
package cmd

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func renderRenameHeader(opts LowerCaseFixOptions, isGit bool, matched int, root string) {
	modeStr := "Local Filesystem (2-step rename)"
	if isGit {
		modeStr = "Git Repository (2-step git mv)"
	}
	filterStr := strings.Join(opts.Patterns, ", ")
	if opts.IsReadmeOnly {
		filterStr = "Root README only"
	}
	fmt.Printf("\n%s⚡ GitMap Lowercase File Renamer%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  ● Mode:        %s\n", modeStr)
	fmt.Printf("  ● Working Dir: %s\n", root)
	fmt.Printf("  ● Filter:      %s\n", filterStr)
	fmt.Printf("  ● Found:       %d uppercase file(s)\n\n", matched)
}

func renderFileStepLog(idx, total int, p RenamePair, isSuccess bool, err error) {
	status := constants.ColorGreen + "OK" + constants.ColorReset
	if !isSuccess {
		status = constants.ColorRed + fmt.Sprintf("FAIL (%v)", err) + constants.ColorReset
	}
	prefix := "git mv"
	if !p.IsGitTracked {
		prefix = "fs rename"
	}
	fmt.Printf("  [%d/%d] %s → %s\n", idx, total, p.OldBase, p.NewBase)
	fmt.Printf("        Step 1: %s %s → %s.tmp-lcf (%s)\n", prefix, p.OldBase, p.OldBase, status)
	fmt.Printf("        Step 2: %s %s.tmp-lcf → %s (%s)\n", prefix, p.OldBase, p.NewBase, status)
}

func renderDryRunNotice(count int) {
	fmt.Printf("\n%sℹ [dry-run] %d file(s) would be renamed to lowercase.%s\n\n",
		constants.ColorYellow, count, constants.ColorReset)
}

func renderRenameSummary(s RenameSummary) {
	fmt.Printf("\n%s════════════════════════════════════════════════════════════════%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("%s✓ Lowercase Rename Summary:%s\n", constants.ColorGreen, constants.ColorReset)
	fmt.Printf("  ● Total Files Scanned: %d\n", s.TotalScanned)
	fmt.Printf("  ● Files Matched:       %d\n", s.TotalMatched)
	fmt.Printf("  ● Files Renamed:       %d\n", s.TotalRenamed)
	if s.IsGitRepo && s.CommitSHA != "" {
		fmt.Printf("  ● Git Status:          Committed (%s)\n", s.CommitSHA[:minLen(8, len(s.CommitSHA))])
	} else if s.IsGitRepo {
		fmt.Printf("  ● Git Status:          Staged (no-commit mode)\n")
	} else {
		fmt.Printf("  ● Filesystem Status:   Renamed on disk\n")
	}
	fmt.Printf("%s════════════════════════════════════════════════════════════════%s\n\n", constants.ColorCyan, constants.ColorReset)
}

func minLen(a, b int) int {
	if a < b {
		return a
	}
	return b
}
