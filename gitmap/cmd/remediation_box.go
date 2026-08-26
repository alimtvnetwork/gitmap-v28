// Package cmd — remediation_box.go renders styled remediation suggestions.
package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/gitutil"
)

func PrintRemediationBox(repoName, repoPath string, d gitutil.DirtyDiagnosis) {
	recipes := gitutil.GenerateRemediationRecipes(repoPath, d)
	if len(recipes) == 0 {
		return
	}

	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffb86c"))
	cmdStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#8be9fd"))

	fmt.Println()
	fmt.Printf("  %s %s (%s)\n", headerStyle.Render("Remediation Options for:"), repoName, d.SummaryReason)
	for _, rec := range recipes {
		fmt.Printf("    • %s: %s\n", rec.Title, rec.Description)
		fmt.Printf("      %s\n", cmdStyle.Render(rec.Command))
	}
	fmt.Println()
}
