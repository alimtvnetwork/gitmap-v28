package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func runInteractiveReconciliation(items []RemediationItem) error {
	reader := bufio.NewReader(os.Stdin)
	var applyAllAction string
	for i := 0; i < len(items); i++ {
		item := &items[i]
		if applyAllAction != "" {
			executeAllAction(item, applyAllAction)
			continue
		}
		action, shouldQuit := promptSingleRepo(reader, i+1, len(items), item)
		if shouldQuit {
			fmt.Println("Reconciliation paused.")
			break
		}
		applyAllAction = handlePromptAction(item, action)
	}
	return nil
}

func executeAllAction(item *RemediationItem, action string) {
	idx := parseRecipeIndex(action, item.Recipes)
	if idx >= 0 && idx < len(item.Recipes) {
		_ = executeFixRecipe(item, item.Recipes[idx])
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

func promptSingleRepo(reader *bufio.Reader, idx, total int, item *RemediationItem) (string, bool) {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#50fa7b"))
	promptStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#8be9fd"))
	fmt.Printf("\n[%d/%d] %s (%s)\n", idx, total, titleStyle.Render(item.RepoName), item.SummaryReason)
	fmt.Printf("  %s ", promptStyle.Render("Pick [1=stash, 2=wip, 3=discard, s=skip, a=all-stash, q=quit]:"))

	input, err := reader.ReadString('\n')
	if err != nil {
		return "", true
	}
	return resolvePromptChoice(strings.TrimSpace(strings.ToLower(input)))
}

func resolvePromptChoice(choice string) (string, bool) {
	switch choice {
	case "1", "stash":
		return "stash", false
	case "2", "wip":
		return "wip", false
	case "3", "discard":
		return "discard", false
	case "s", "skip":
		return "skip", false
	case "a", "all", "all-stash":
		return "all-stash", false
	case "q", "quit", "exit":
		return "", true
	default:
		return "stash", false
	}
}
