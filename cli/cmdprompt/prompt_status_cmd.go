// Package cmdprompt — prompt_status_cmd.go displays prompt architect status table.
package cmdprompt

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func RunPromptStatus(targets []string) error {
	if len(targets) == 0 {
		targets = ResolvePromptTarget("").Data
	}

	layout := NewPromptStatusTableLayout()
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#bd93f9"))
	fmt.Println(titleStyle.Render("Prompt Architect Status:"))
	layout.PrintHeader()

	for _, t := range targets {
		meta, _ := ReadPromptArchitectMetadata(t)
		layout.PrintRow(t, meta)
	}

	fmt.Println("  --------------------------------------------------------------------------------")

	return nil
}
