// Package cmd — lowercasefix_report.go formats status logs and summaries for lowercase fix.
package cmd

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func renderRenameHeader(opts LowerCaseFixOptions, isGit bool, matched, scanned int, root string) {
	modeStr := resolveModeStr(isGit)
	filterStr := resolveFilterStr(opts)
	fmt.Printf("\n%s⚡ GitMap Lowercase File Renamer%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  ● Mode:        %s\n", modeStr)
	fmt.Printf("  ● Working Dir: %s\n", root)
	fmt.Printf("  ● Filter:      %s\n", filterStr)
	fmt.Printf("  ● Scanned:     %d files in directory\n", scanned)
	fmt.Printf("  ● Matched:     %d uppercase file(s)\n\n", matched)
}

func resolveModeStr(isGit bool) string {
	if isGit {
		return "Git Repository (2-step git mv)"
	}
	return "Local Filesystem (2-step rename)"
}

func resolveFilterStr(opts LowerCaseFixOptions) string {
	if opts.IsReadmeOnly {
		return "Root README only"
	}
	return strings.Join(opts.Patterns, ", ")
}

func renderFileStepLog(idx, total int, p RenamePair, isSuccess bool, err error) {
	status := resolveStepStatus(isSuccess, err)
	prefix := resolveRenamePrefix(p.IsGitTracked)
	matched := resolveMatchedInfo(p.MatchedBy)
	fmt.Printf("  [%d/%d] %s → %s%s\n", idx, total, p.OldBase, p.NewBase, matched)
	fmt.Printf("        Step 1: %s %s → %s.tmp-lcf (%s)\n", prefix, p.OldBase, p.OldBase, status)
	fmt.Printf("        Step 2: %s %s.tmp-lcf → %s (%s)\n", prefix, p.OldBase, p.NewBase, status)
}

func resolveStepStatus(isSuccess bool, err error) string {
	if isSuccess {
		return constants.ColorGreen + "OK" + constants.ColorReset
	}
	return constants.ColorRed + fmt.Sprintf("FAIL (%v)", err) + constants.ColorReset
}

func resolveMatchedInfo(matchedBy string) string {
	if matchedBy == "" {
		return ""
	}
	return " (" + matchedBy + ")"
}

func renderDryRunNotice(count int) {
	fmt.Printf("\n%sℹ [dry-run] %d file(s) would be safely renamed using 2-step move. No files modified.%s\n\n",
		constants.ColorYellow, count, constants.ColorReset)
}

func renderRenameSummary(s RenameSummary) {
	renderSummaryBoxHeader(s)
	renderSummaryDetails(s)
	renderManipulationSteps(s)
	fmt.Printf("%s════════════════════════════════════════════════════════════════%s\n\n", constants.ColorCyan, constants.ColorReset)
}

func renderSummaryBoxHeader(_ RenameSummary) {
	fmt.Printf("\n%s════════════════════════════════════════════════════════════════%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("%s✓ Lowercase Rename Summary:%s\n", constants.ColorGreen, constants.ColorReset)
}

func renderSummaryDetails(s RenameSummary) {
	fmt.Printf("  ● Total Files Scanned: %d\n", s.TotalScanned)
	fmt.Printf("  ● Files Matched:       %d\n", s.TotalMatched)
	fmt.Printf("  ● Files Renamed:       %d\n", s.TotalRenamed)
	renderGitOrFSStatus(s)
}

func renderGitOrFSStatus(s RenameSummary) {
	if s.IsGitRepo && s.CommitSHA != "" {
		pushInfo := ""
		if s.IsPushed {
			pushInfo = " & Pushed to remote"
		}
		fmt.Printf("  ● Git Status:          Committed (%s)%s\n", s.CommitSHA[:minLen(8, len(s.CommitSHA))], pushInfo)
		return
	}
	if s.IsGitRepo {
		fmt.Printf("  ● Git Status:          Staged in index (no-commit mode)\n")
		return
	}
	fmt.Printf("  ● Filesystem Status:   Renamed on disk (non-git directory)\n")
}

func renderManipulationSteps(s RenameSummary) {
	fmt.Printf("\n  ● Manipulation Steps Performed:\n")
	if s.IsGitRepo {
		renderGitManipulationSteps(s)
		return
	}
	renderFSManipulationSteps()
}

func renderGitManipulationSteps(s RenameSummary) {
	fmt.Printf("    1. Step 1 (Safe Temp Move):  git mv <file> <file>.tmp-lcf\n")
	fmt.Printf("       Avoids silent case-collision / no-op on case-insensitive filesystems (Windows/macOS)\n")
	fmt.Printf("    2. Step 2 (Target Rename):   git mv <file>.tmp-lcf <file_lowercase>\n")
	fmt.Printf("       Registers true case rename in Git index tree\n")
	fmt.Printf("    3. Step 3 (Index Sync):      git add -A\n")
	if s.CommitSHA != "" {
		fmt.Printf("    4. Step 4 (Atomic Commit):   git commit -m \"chore: rename ...\"\n")
		if s.IsPushed {
			fmt.Printf("    5. Step 5 (Remote Push):     git push origin <branch>\n")
		} else {
			fmt.Printf("    5. Step 5 (Skip Push):       Push skipped (--no-push specified or push failed)\n")
		}
		return
	}
	fmt.Printf("    4. Step 4 (Skip Commit):     Staged without committing (--no-commit specified)\n")
}

func renderFSManipulationSteps() {
	fmt.Printf("    1. Step 1 (Safe Temp Move):  os.Rename(<file>, <file>.tmp-lcf)\n")
	fmt.Printf("       Avoids filesystem collision on case-insensitive disks\n")
	fmt.Printf("    2. Step 2 (Target Rename):   os.Rename(<file>.tmp-lcf, <file_lowercase>)\n")
	fmt.Printf("    3. Status:                   Files renamed on disk\n")
}

func minLen(a, b int) int {
	if a < b {
		return a
	}
	return b
}
