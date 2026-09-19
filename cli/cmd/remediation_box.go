package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/render"
)

type RemediationItem struct {
	RepoPath      string                      `json:"repoPath"`
	RepoName      string                      `json:"repoName"`
	SummaryReason string                      `json:"summaryReason"`
	Recipes       []gitutil.RemediationRecipe `json:"recipes"`
	Files         []string                    `json:"files,omitempty"`
}

type RemediationBatchState struct {
	Items []RemediationItem `json:"items"`
}

func getRemediationStateFile() string {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, ".gitmap", "output")
	_ = os.MkdirAll(dir, 0755)

	return filepath.Join(dir, "last_remediation.json")
}

func SaveRemediationState(items []RemediationItem) error {
	state := RemediationBatchState{Items: items}
	b, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(getRemediationStateFile(), b, 0644)
}

func LoadRemediationState() []RemediationItem {
	b, err := os.ReadFile(getRemediationStateFile())
	if err != nil {
		return nil
	}

	var batch RemediationBatchState
	if err := json.Unmarshal(b, &batch); err == nil && len(batch.Items) > 0 {
		return batch.Items
	}

	var single RemediationItem
	if err := json.Unmarshal(b, &single); err == nil && single.RepoName != "" {
		return []RemediationItem{single}
	}

	return nil
}

func PrintRemediationBox(repoName, repoPath string, d gitutil.DirtyDiagnosis) {
	recipes := gitutil.GenerateRemediationRecipes(repoPath, d)
	if len(recipes) == 0 {
		return
	}

	item := RemediationItem{
		RepoPath:      repoPath,
		RepoName:      repoName,
		SummaryReason: d.SummaryReason,
		Recipes:       recipes,
		Files:         d.AllFiles,
	}

	PrintRemediationSummary([]RemediationItem{item})
}

func printRemediationStrategyBox() {
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffb86c"))
	cmdStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#8be9fd"))
	fmt.Println()
	fmt.Println(headerStyle.Render("  ── Pull Remediation Options ──"))
	fmt.Println(cmdStyle.Render("    [1] Stash & Re-apply : Stash local changes (-u), pull remote, then pop stash"))
	fmt.Println(cmdStyle.Render("    [2] Commit WIP       : Commit all modified and untracked files, then pull --rebase"))
	fmt.Println(cmdStyle.Render("    [3] Discard Local    : Permanently discard changes (reset --hard & clean -fd), then pull"))
}

func printPendingReposList(items []RemediationItem) {
	fmt.Println()
	fmt.Printf("  Pending Repositories (%d):\n", len(items))
	for i, item := range items {
		printPendingRepoEntry(i+1, item)
	}

	fmt.Println()
}

func printPendingRepoEntry(idx int, item RemediationItem) {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#50fa7b"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#6272a4"))
	cmdStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#8be9fd"))
	reason := item.SummaryReason
	if reason == "" {
		reason = "uncommitted changes"
	}

	fmt.Printf("    %2d. %s %s\n", idx, titleStyle.Render(item.RepoName), dimStyle.Render("("+reason+")"))
	fixCmd := fmt.Sprintf("gitmap fix %s 1  (or 2=wip, 3=discard)", item.RepoName)
	fmt.Printf("        ↳ Gitmap Fix:    %s\n", cmdStyle.Render(fixCmd))
}

func printRemediationCLIHelp() {
	cmdStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#8be9fd"))
	fmt.Println("  To fix:")
	fmt.Printf("    %-32s %s\n", cmdStyle.Render("gitmap fix"), "(interactive prompt walkthrough)")
	fmt.Printf("    %-32s %s\n", cmdStyle.Render("gitmap fix all"), "(apply stash to all repositories)")
	fmt.Printf("    %-32s %s\n", cmdStyle.Render("gitmap fix all [1|2|3]"), "(apply specific strategy to all)")
	fmt.Printf("    %-32s %s\n", cmdStyle.Render("gitmap fix --prompt"), "(step-by-step interactive prompt)")
	fmt.Printf("    %-32s %s\n\n", cmdStyle.Render("gitmap fix <repo> [1|2|3]"), "(target specific repository)")
}

func promptForRemediation(items []RemediationItem) {
	promptStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffb86c"))
	fmt.Printf("  %s ", promptStyle.Render("Remediate dirty repository(ies) now? [y/N]:"))
	reader := bufio.NewReader(os.Stdin)
	ans, _ := reader.ReadString('\n')
	ans = strings.TrimSpace(strings.ToLower(ans))
	if ans == "y" || ans == "yes" {
		_ = runInteractiveRemediation(items)

		return
	}

	printRemediationCLIHelp()
}

func PrintRemediationSummaryNoPrompt(items []RemediationItem) {
	if len(items) == 0 {
		return
	}

	_ = SaveRemediationState(items)
	printRemediationStrategyBox()
	printPendingReposList(items)
	printRemediationCLIHelp()
}

func PrintRemediationSummaryAutoFix(items []RemediationItem) {
	if len(items) == 0 {
		return
	}

	_ = SaveRemediationState(items)
	printRemediationStrategyBox()
	printPendingReposList(items)
	_ = runInteractiveRemediation(items)
}

func PrintRemediationSummary(items []RemediationItem) {
	if len(items) == 0 {
		return
	}

	_ = SaveRemediationState(items)
	printRemediationStrategyBox()
	printPendingReposList(items)
	if !render.StdoutIsTerminal() {
		printRemediationCLIHelp()

		return
	}

	promptForRemediation(items)
}
