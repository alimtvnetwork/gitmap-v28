package cmdagy

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func renderActivePromptsTerminal(convs []AgyActiveConversation) {
	printActiveHeader(len(convs))
	isEmpty := len(convs) == 0
	if isEmpty {
		fmt.Printf("  %s✔ All Antigravity conversations are currently idle (0 active prompts running).%s\n\n",
			constants.ColorGreen, constants.ColorReset)
		return
	}
	printActiveTable(convs)
}

func printActiveHeader(count int) {
	fmt.Printf("\n%s● Antigravity Active Running Prompts (%d conversation(s)):%s\n\n",
		constants.ColorCyan, count, constants.ColorReset)
}

func printActiveTable(convs []AgyActiveConversation) {
	fmt.Printf("  %-3s %-24s %-32s %-8s %-12s %s\n",
		"#", "PROJECT / TITLE", "WORKSPACE / CWD", "STEPS", "ELAPSED", "CONVERSATION ID")
	fmt.Printf("  %s\n", strings.Repeat("─", 108))
	for idx, c := range convs {
		renderActiveConvRow(idx+1, c)
	}
	fmt.Printf("  %s\n\n", strings.Repeat("─", 108))
}

func renderActiveConvRow(seq int, c AgyActiveConversation) {
	title := truncateString(c.Title, 22)
	ws := truncateString(c.WorkspacePath, 30)
	elapsed := formatElapsedDuration(c.ElapsedSeconds)
	fmt.Printf("  %-3d %s%-24s%s %-32s %-8d %-12s %s%s%s\n",
		seq, constants.ColorYellow, title, constants.ColorReset,
		ws, c.StepCount, elapsed, constants.ColorDim, c.ConversationID, constants.ColorReset)
}

func formatElapsedDuration(secs int) string {
	if secs < 60 {
		return fmt.Sprintf("%ds", secs)
	}
	return fmt.Sprintf("%dm %ds", secs/60, secs%60)
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func renderActivePromptsJSON(convs []AgyActiveConversation) error {
	data, err := json.MarshalIndent(convs, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func hasArgFlag(args []string, flagName string) bool {
	for _, a := range args {
		if a == flagName {
			return true
		}
	}
	return false
}
