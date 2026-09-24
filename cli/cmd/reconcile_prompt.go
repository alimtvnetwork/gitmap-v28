package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

var promptChoiceMap = map[string]string{
	"1":         "stash",
	"stash":     "stash",
	"2":         "wip",
	"wip":       "wip",
	"3":         "discard",
	"discard":   "discard",
	"s":         "skip",
	"skip":      "skip",
	"a":         "all-stash",
	"all":       "all-stash",
	"all-stash": "all-stash",
}

func executeAllAction(item *RemediationItem, action string) {
	idx := parseRecipeIndex(action, item.Recipes)
	isInvalidIndex := idx < 0 || idx >= len(item.Recipes)
	if isInvalidIndex {
		return
	}

	err := executeFixRecipe(item, item.Recipes[idx])
	if err != nil {
		fmt.Printf("Warning: fix action failed on %s: %v\n", item.RepoName, err)
	}
}

func handlePromptAction(item *RemediationItem, action string) string {
	if action == "skip" {
		return ""
	}

	if strings.HasPrefix(action, "all-") {
		allAction := strings.TrimPrefix(action, "all-")
		executeAllAction(item, allAction)

		return allAction
	}

	executeAllAction(item, action)

	return ""
}

func resolveDirtyFiles(item *RemediationItem) []string {
	if len(item.RepoPath) == 0 {
		return item.Files
	}

	diag := gitutil.InspectDirtyState(item.RepoPath)
	if len(diag.AllFiles) > 0 {
		item.Files = diag.AllFiles
		return diag.AllFiles
	}

	return item.Files
}

func printOverflowNote(totalCount, maxShowCount int) {
	if totalCount <= maxShowCount {
		return
	}

	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#8be9fd"))
	remainingCount := totalCount - maxShowCount
	fmt.Printf("      %s\n", dimStyle.Render(fmt.Sprintf("... and %d more files", remainingCount)))
}

func formatDirtyFileEntry(raw string) string {
	if strings.HasPrefix(raw, "staged: ") {
		tag := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#50fa7b")).Render("[staged]   ")
		return tag + " " + strings.TrimPrefix(raw, "staged: ")
	}
	if strings.HasPrefix(raw, "modified: ") {
		tag := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#f1fa8c")).Render("[modified] ")
		return tag + " " + strings.TrimPrefix(raw, "modified: ")
	}
	if strings.HasPrefix(raw, "untracked: ") {
		tag := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8be9fd")).Render("[untracked]")
		return tag + " " + strings.TrimPrefix(raw, "untracked: ")
	}
	if strings.HasPrefix(raw, "deleted: ") {
		tag := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ff5555")).Render("[deleted]  ")
		return tag + " " + strings.TrimPrefix(raw, "deleted: ")
	}

	return raw
}

func renderDirtyFileList(files []string) {
	fmt.Printf("    Pending Changes (%d files):\n", len(files))
	maxShowCount := 10
	displayCount := len(files)
	if displayCount > maxShowCount {
		displayCount = maxShowCount
	}

	for i := 0; i < displayCount; i++ {
		fmt.Printf("      • %s\n", formatDirtyFileEntry(files[i]))
	}

	printOverflowNote(len(files), maxShowCount)
}

func printRepoDirtyFiles(item *RemediationItem) {
	files := resolveDirtyFiles(item)
	if len(files) == 0 {
		return
	}

	renderDirtyFileList(files)
}

func printPromptOptions() {
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#8be9fd"))
	fmt.Println("    Remediation Options:")
	fmt.Printf("      [1] %-18s %s\n", "Stash & Re-apply", dimStyle.Render("Stash local changes (-u), pull remote, then pop stash"))
	fmt.Printf("      [2] %-18s %s\n", "Commit WIP", dimStyle.Render("Commit all modified/untracked files, then pull --rebase"))
	fmt.Printf("      [3] %-18s %s\n", "Discard Local", dimStyle.Render("Permanently discard changes (reset --hard & clean -fd), pull"))
	fmt.Printf("      [s] %-18s %s\n", "Skip", dimStyle.Render("Skip this repository for now"))
	fmt.Printf("      [a] %-18s %s\n", "Apply to All", dimStyle.Render("Apply stash to this and all remaining repositories"))
	fmt.Printf("      [q] %-18s %s\n", "Quit", dimStyle.Render("Exit interactive prompt"))
}

func resolvePromptChoice(choice string) (string, bool) {
	isQuit := choice == "q" || choice == "quit" || choice == "exit"
	if isQuit {
		return "", true
	}

	action, hasAction := promptChoiceMap[choice]
	if hasAction {
		return action, false
	}

	return "stash", false
}

func promptSingleRepo(reader *bufio.Reader, idx, total int, item *RemediationItem) (string, bool) {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#50fa7b"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#8be9fd"))
	promptStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8be9fd"))
	fmt.Printf("\n[%d/%d] %s %s\n", idx, total, titleStyle.Render(item.RepoName), dimStyle.Render("("+item.SummaryReason+")"))
	if len(item.RepoPath) > 0 {
		fmt.Printf("    Path: %s\n", dimStyle.Render(item.RepoPath))
	}
	printRepoDirtyFiles(item)
	printPromptOptions()
	fmt.Printf("  %s ", promptStyle.Render("Pick [1=stash, 2=wip, 3=discard, s=skip, a=all-stash, q=quit]:"))

	input, err := reader.ReadString('\n')
	if err != nil {
		return "", true
	}

	return resolvePromptChoice(strings.TrimSpace(strings.ToLower(input)))
}

func processPromptStep(reader *bufio.Reader, idx, total int, item *RemediationItem, applyAll string) (string, bool) {
	if len(applyAll) > 0 {
		executeAllAction(item, applyAll)

		return applyAll, false
	}

	action, shouldQuit := promptSingleRepo(reader, idx, total, item)
	if shouldQuit {
		return "", true
	}

	return handlePromptAction(item, action), false
}

func runInteractiveRemediation(items []RemediationItem) error {
	reader := bufio.NewReader(os.Stdin)
	var applyAllAction string
	for i := 0; i < len(items); i++ {
		nextAction, shouldQuit := processPromptStep(reader, i+1, len(items), &items[i], applyAllAction)
		if shouldQuit {
			return nil
		}

		applyAllAction = nextAction
	}

	return nil
}

//nolint:unused
func runInteractiveReconciliation(items []RemediationItem) error {
	return runInteractiveRemediation(items)
}
