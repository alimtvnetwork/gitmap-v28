package cmdagy

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func renderQueuesTerminal(summaries []AgyWorkspaceQueueSummary) {
	totalPrompts := countTotalQueuedPrompts(summaries)
	printQueuesHeader(len(summaries), totalPrompts)
	isEmpty := len(summaries) == 0
	if isEmpty {
		fmt.Printf("  %s✔ All Antigravity prompt queues are empty (0 queued prompts across workspaces).%s\n\n",
			constants.ColorGreen, constants.ColorReset)
		return
	}
	for idx, s := range summaries {
		renderSingleWorkspaceQueueSection(idx+1, s)
	}
}

func countTotalQueuedPrompts(summaries []AgyWorkspaceQueueSummary) int {
	total := 0
	for _, s := range summaries {
		total += s.TotalQueued
	}
	return total
}

func printQueuesHeader(wsCount, promptCount int) {
	fmt.Printf("\n%s● Antigravity Workspace Prompt Queues (%d prompt(s) across %d workspace(s)):%s\n\n",
		constants.ColorCyan, promptCount, wsCount, constants.ColorReset)
}

func renderSingleWorkspaceQueueSection(seq int, s AgyWorkspaceQueueSummary) {
	fmt.Printf("  %d. %s%s%s (%s)\n",
		seq, constants.ColorWhite, s.ProjectName, constants.ColorReset, s.Workspace)
	for _, q := range s.QueuedItems {
		renderQueueEntryLine(q)
	}
	fmt.Println()
}

func renderQueueEntryLine(q AgyPromptQueueEntry) {
	statusBadge := formatQueueStatusBadge(q.Status)
	preview := sanitizePromptSnippet(q.Prompt, 55)
	fmt.Printf("     [Queue #%d] %s • %s\n", q.ID, statusBadge, preview)
	if q.CreatedAt != "" {
		fmt.Printf("                 %sCreated: %s%s\n", constants.ColorDim, q.CreatedAt, constants.ColorReset)
	}
}

func formatQueueStatusBadge(status string) string {
	switch status {
	case "queued":
		return constants.ColorYellow + "queued" + constants.ColorReset
	case "requeued_with_check_prefix":
		return constants.ColorCyan + "requeued_check" + constants.ColorReset
	case "dispatched":
		return constants.ColorGreen + "dispatched" + constants.ColorReset
	default:
		return status
	}
}

func sanitizePromptSnippet(prompt string, maxLen int) string {
	lines := strings.Split(strings.TrimSpace(prompt), "\n")
	first := ""
	if len(lines) > 0 {
		first = strings.TrimSpace(lines[0])
	}
	if len(first) > maxLen {
		return first[:maxLen-3] + "..."
	}
	return first
}

func renderQueuesJSON(summaries []AgyWorkspaceQueueSummary) error {
	data, err := json.MarshalIndent(summaries, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
