package osclean

import (
	"fmt"
	"strings"
)

// RenderTerminalSummaryTable outputs an aligned terminal summary table for terminal cleaning.
func RenderTerminalSummaryTable(summary TerminalCleanSummary) {
	fmt.Println()
	fmt.Printf("  %-30s %-12s %8s %12s %16s\n", "Shell / Category", "Status", "Files", "Freed", "Reseeded")
	fmt.Println("  " + strings.Repeat("─", 82))
	for _, s := range summary.Shells {
		renderTerminalTableRow(s, summary.IsDryRun)
	}
	fmt.Println("  " + strings.Repeat("─", 82))
	renderTerminalTableTotals(summary)
	printTerminalSessionTip()
}

func renderTerminalTableRow(s ShellCleanStats, isDryRun bool) {
	status := "Cleaned"
	if isDryRun {
		status = "Reclaimable"
	}
	sizeStr := FormatCleanSize(s.BytesFreed)
	reseedStr := "None"
	if s.IsReseeded {
		if isDryRun {
			reseedStr = fmt.Sprintf("%d cmds (plan)", s.ReseedCount)
		} else {
			reseedStr = fmt.Sprintf("✔ %d cmds", s.ReseedCount)
		}
	}
	fmt.Printf("  • %-28s %-12s %8d %12s %16s\n",
		s.Label, status, s.FilesCleared, sizeStr, reseedStr)
}

func renderTerminalTableTotals(summary TerminalCleanSummary) {
	mode := "Cleaned"
	if summary.IsDryRun {
		mode = "Dry-Run Reclaimable"
	}
	totalSize := FormatCleanSize(summary.TotalBytesFreed)
	fmt.Printf("  ✔ Total %s: %s across %d history file(s), reseeded %d suggestions (%dms)\n\n",
		mode, totalSize, summary.TotalFilesCleared, summary.TotalReseedCount, summary.DurationMs)
}

func printTerminalSessionTip() {
	fmt.Println("  💡 Active Session Tips:")
	fmt.Println("     • PowerShell: Run '[Microsoft.PowerShell.PSConsoleReadLine]::ClearHistory()' or restart session")
	fmt.Println("     • Bash:       Run 'history -c' to clear active in-memory buffer")
	fmt.Println("     • Zsh:        Run 'history -c' or restart session to load fresh suggestions")
	fmt.Println()
}

// RenderVerboseTerminalNotes displays individual cleared paths, notes, and warnings.
func RenderVerboseTerminalNotes(shells []ShellCleanStats) {
	for _, s := range shells {
		if len(s.ClearedPaths) > 0 {
			fmt.Printf("    [%s] Cleared paths:\n", s.Shell)
			for _, p := range s.ClearedPaths {
				fmt.Printf("      - %s\n", p)
			}
		}
		for _, n := range s.Notes {
			fmt.Printf("      note: %s\n", n)
		}
		for _, e := range s.Errors {
			fmt.Printf("      warning: %s\n", e)
		}
	}
	fmt.Println()
}
