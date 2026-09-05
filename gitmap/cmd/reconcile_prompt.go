package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/gitutil"
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
	files := item.Files
	isMissingFiles := len(files) <= 0
	hasRepoPath := len(item.RepoPath) > 0
	if isMissingFiles && hasRepoPath {
		diag := gitutil.InspectDirtyState(item.RepoPath)

		return diag.AllFiles
	}

	return files
}

func printOverflowNote(totalCount, maxShowCount int) {
	if totalCount <= maxShowCount {
		return
	}
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#6272a4"))
	remainingCount := totalCount - maxShowCount
	fmt.Printf("    %s\n", dimStyle.Render(fmt.Sprintf("... and %d more files", remainingCount)))
}

func renderDirtyFileList(files []string) {
	fileStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#f1fa8c"))
	maxShowCount := 10
	displayCount := len(files)
	if displayCount > maxShowCount {
		displayCount = maxShowCount
	}
	for i := 0; i < displayCount; i++ {
		fmt.Printf("    • %s\n", fileStyle.Render(files[i]))
	}
	printOverflowNote(len(files), maxShowCount)
}

func printRepoDirtyFiles(item *RemediationItem) {
	files := resolveDirtyFiles(item)
	if len(files) <= 0 {
		return
	}
	renderDirtyFileList(files)
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
	promptStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#8be9fd"))
	fmt.Printf("\n[%d/%d] %s (%s)\n", idx, total, titleStyle.Render(item.RepoName), item.SummaryReason)
	printRepoDirtyFiles(item)
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

func runInteractiveReconciliation(items []RemediationItem) error {
	return runInteractiveRemediation(items)
}
