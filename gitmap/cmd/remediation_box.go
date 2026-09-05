package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/render"
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

func FindRemediationItem(items []RemediationItem, query string) *RemediationItem {
	cleanQuery := strings.TrimSpace(query)
	for i := range items {
		if strings.EqualFold(items[i].RepoName, cleanQuery) {
			return &items[i]
		}
		if strings.EqualFold(filepath.Base(items[i].RepoPath), cleanQuery) {
			return &items[i]
		}
	}
	for i := range items {
		if strings.Contains(strings.ToLower(items[i].RepoName), strings.ToLower(cleanQuery)) {
			return &items[i]
		}
	}
	if num, err := strconv.Atoi(cleanQuery); err == nil && num > 0 && num <= len(items) {
		return &items[num-1]
	}
	return nil
}

func RemoveRemediationItem(repoName string) {
	items := LoadRemediationState()
	var remaining []RemediationItem
	for _, item := range items {
		if !strings.EqualFold(item.RepoName, repoName) {
			remaining = append(remaining, item)
		}
	}
	if len(remaining) == 0 {
		_ = os.Remove(getRemediationStateFile())
		return
	}
	_ = SaveRemediationState(remaining)
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
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#50fa7b"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#6272a4"))
	fmt.Println()
	fmt.Printf("  Pending Repositories (%d):\n", len(items))
	for i, item := range items {
		reason := item.SummaryReason
		if reason == "" {
			reason = "uncommitted changes"
		}
		fmt.Printf("    %2d. %s %s\n", i+1, titleStyle.Render(item.RepoName), dimStyle.Render("("+reason+")"))
	}
	fmt.Println()
}

func printRemediationCLIHelp() {
	cmdStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#8be9fd"))
	fmt.Println("  To fix later:")
	fmt.Printf("    %s\n", cmdStyle.Render("gitmap fix               (interactive walkthrough)"))
	fmt.Printf("    %s\n", cmdStyle.Render("gitmap fix <repo> [1|2|3] (target specific repo)"))
	fmt.Printf("    %s\n\n", cmdStyle.Render("gitmap fix --all stash    (apply stash to all)"))
}

func promptForRemediation(items []RemediationItem) {
	promptStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffb86c"))
	fmt.Printf("  %s ", promptStyle.Render("Remediate these repositories now? [Y/n]:"))
	reader := bufio.NewReader(os.Stdin)
	ans, _ := reader.ReadString('\n')
	ans = strings.TrimSpace(strings.ToLower(ans))
	if ans == "" || ans == "y" || ans == "yes" {
		_ = runInteractiveRemediation(items)

		return
	}
	printRemediationCLIHelp()
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
