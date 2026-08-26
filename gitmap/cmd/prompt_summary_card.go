// Package cmd — prompt_summary_card.go renders installation summary.
package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func RenderPromptInstallSummary(results []model.PromptInstallResult) {
	var successCount, failCount int
	for _, r := range results {
		if r.IsSuccess {
			successCount++
		} else {
			failCount++
		}
	}

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#50fa7b")).
		Padding(0, 1)

	text := fmt.Sprintf("Prompt Architect Summary\nSuccess: %d  |  Failed: %d", successCount, failCount)
	fmt.Println(boxStyle.Render(text))
}
