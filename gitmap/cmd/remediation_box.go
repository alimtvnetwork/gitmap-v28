package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/gitutil"
	"github.com/charmbracelet/lipgloss"
)

type RemediationState struct {
	RepoPath string                      `json:"repoPath"`
	RepoName string                      `json:"repoName"`
	Recipes  []gitutil.RemediationRecipe `json:"recipes"`
}

func getRemediationStateFile() string {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, ".gitmap", "output")
	os.MkdirAll(dir, 0755)
	return filepath.Join(dir, "last_remediation.json")
}

func PrintRemediationBox(repoName, repoPath string, d gitutil.DirtyDiagnosis) {
	recipes := gitutil.GenerateRemediationRecipes(repoPath, d)
	if len(recipes) == 0 {
		return
	}

	state := RemediationState{
		RepoPath: repoPath,
		RepoName: repoName,
		Recipes:  recipes,
	}
	b, _ := json.MarshalIndent(state, "", "  ")
	os.WriteFile(getRemediationStateFile(), b, 0644)

	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffb86c"))
	cmdStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#8be9fd"))

	fmt.Println()
	fmt.Printf("  %s %s (%s)\n", headerStyle.Render("Remediation Options for:"), repoName, d.SummaryReason)
	for i, rec := range recipes {
		// Rewrite the output to recommend the gitmap fix command
		alias := ""
		if i == 0 {
			alias = "stash"
		} else if i == 1 {
			alias = "wip"
		} else if i == 2 {
			alias = "discard"
		}

		fmt.Printf("    • %s: %s\n", rec.Title, rec.Description)
		fmt.Printf("      %s\n", cmdStyle.Render(fmt.Sprintf("Run: gitmap fix %d  (or 'gitmap %s')", i+1, alias)))
	}
	fmt.Println()
}
