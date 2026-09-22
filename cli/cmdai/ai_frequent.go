package cmdai

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/atotto/clipboard"
)

// RunAiFrequent displays the most frequently executed AI commands and optionally copies them.
func RunAiFrequent(limit int, isCopy bool) *apperror.AppError {
	frequent, err := store.GetFrequentAiCommands(limit)
	if err != nil {
		return apperror.WrapSimple(err, "ai.frequent")
	}

	if len(frequent) == 0 {
		renderEmptyFrequentNotice()
		return nil
	}

	renderFrequentTable(frequent)
	if isCopy {
		copyFrequentToClipboard(frequent)
	}

	return nil
}

func renderEmptyFrequentNotice() {
	fmt.Printf("\n%s[AI Split DB — Frequent Commands]%s\n", constants.ColorBold, constants.ColorReset)
	fmt.Println("  No AI commands recorded yet in ~/.gitmap/ai/instructions.db.")
	fmt.Println("  Tip: Run PowerShell or automation commands with '--ai' to track frequency.")
	fmt.Println()
}

func renderFrequentTable(items []store.AiFrequentCommand) {
	fmt.Printf("\n%s[AI Split DB — Top Frequent Commands]%s\n", constants.ColorBold, constants.ColorReset)
	fmt.Printf("  %-4s  %-6s  %-8s  %-35s  %s\n", "#", "RUNS", "SUCCESS", "COMMAND", "LAST RUN")
	fmt.Printf("  %-4s  %-6s  %-8s  %-35s  %s\n", "---", "----", "-------", "-----------------------------------", "-------------------")

	for i, item := range items {
		successRate := calculateSuccessRate(item.SuccessCount, item.RunCount)
		cmdDisplay := truncateCommand(item.CommandText, 35)
		fmt.Printf("  %-4d  %-6d  %-8s  %-35s  %s\n",
			i+1, item.RunCount, successRate, cmdDisplay, item.LastExecutedAt)
	}
	fmt.Println()
}

func calculateSuccessRate(success, total int) string {
	if total == 0 {
		return "0%"
	}
	pct := (success * 100) / total
	return fmt.Sprintf("%d%%", pct)
}

func truncateCommand(cmd string, maxLen int) string {
	clean := strings.TrimSpace(cmd)
	if len(clean) > maxLen {
		return clean[:maxLen-3] + "..."
	}
	return clean
}

func copyFrequentToClipboard(items []store.AiFrequentCommand) {
	var lines []string
	for _, item := range items {
		lines = append(lines, item.CommandText)
	}
	payload := strings.Join(lines, "\n")
	if err := clipboard.WriteAll(payload); err == nil {
		fmt.Printf("  %s📋 Copied %d frequent AI command(s) to clipboard ✅%s\n\n",
			constants.ColorGreen, len(items), constants.ColorReset)
	}
}
