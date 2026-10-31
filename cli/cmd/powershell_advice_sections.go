package cmd

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// Card-section builders split out of powershell_advice.go: the file-size
// gate caps Go files at 300 lines, and the advice file exceeded it once
// the dead helper was removed. Logic here is unchanged, only relocated.

func getDirectoryReplacements(target string) []ReplacementOption {
	return []ReplacementOption{
		{Command: "gitmap dir " + target + "   OR   gitmap ls " + target, Description: "List Directory"},
		{Command: "gitmap agent subtask ls", Description: "Inspect Subtasks"},
		{Command: "gitmap cat <filepath>   OR   gitmap heatmap <file>", Description: "Read File / Spec"},
		{Command: fmt.Sprintf("gitmap ps \"Get-ChildItem -Path '%s'\"", target), Description: "PowerShell Query"},
	}
}

func getDocumentReplacements(target string) []ReplacementOption {
	return []ReplacementOption{
		{Command: "gitmap heatmap " + target + "   OR   gitmap cat " + target, Description: "Read File / Spec"},
		{Command: "gitmap cat " + target + " -L 1:50", Description: "Read Line Range"},
		{Command: "gitmap dir <path>   OR   gitmap ls <path>", Description: "List Directory"},
	}
}

func getDefaultReplacements(target string) []ReplacementOption {
	return []ReplacementOption{
		{Command: "gitmap help", Description: "All Commands"},
		{Command: fmt.Sprintf("gitmap ps \"%s\"", target), Description: "Quote Expression"},
	}
}

func getReplacementsForTarget(kind PowerShellTargetKind, target string) []ReplacementOption {
	if kind == TargetKindDirectory {
		return getDirectoryReplacements(target)
	}

	if kind == TargetKindDocument {
		return getDocumentReplacements(target)
	}

	return getDefaultReplacements(target)
}

func formatReplacementRow(r ReplacementOption) string {
	cmdText := termout.SafeTokyoStyle(termout.TokyoAccent).Render(r.Command)

	return fmt.Sprintf("• %-18s %s", r.Description+":", cmdText)
}

func buildNativeReplacements(opts PowerShellAdviceOptions) string {
	header := termout.SafeTokyoStyle(termout.TokyoOK).Bold(true).Render("💡 Recommended Native Replacements:")
	replacements := getReplacementsForTarget(opts.TargetKind, opts.RawTarget)
	rows := make([]string, 0, len(replacements))
	for _, r := range replacements {
		rows = append(rows, formatReplacementRow(r))
	}

	return fmt.Sprintf("%s\n%s", header, strings.Join(rows, "\n"))
}

func build3StepLearningActions() string {
	header := termout.SafeTokyoStyle(termout.TokyoPrimary).Bold(true).Render("🚀 3-Step Immediate Learning Action:")
	s1 := "Step 1: Identify your real intent (list directory vs execute script)."
	s2 := "Step 2: Use the specialized native GitMap command shown above."
	s3 := "Step 3: If running PowerShell, enclose full expression in quotes (\"...\")."

	return fmt.Sprintf("%s\n%s\n%s\n%s", header, s1, s2, s3)
}

func renderCardDivider() string {
	return termout.SafeTokyoStyle(termout.TokyoMuted).Render(strings.Repeat("─", 74))
}
