package cmdagy

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

func renderPromptsTable(prompts []AgyPromptEntry) {
	fmt.Printf("\n  %s● AGY Prompts (%d entries)%s\n\n", constants.ColorCyan, len(prompts), constants.ColorReset)
	cfg := termout.TableConfig{
		Columns: []termout.Column{
			{Title: "DATE", Align: termout.AlignLeft, MinWidth: 16},
			{Title: "WORKSPACE", Align: termout.AlignLeft, MinWidth: 20},
			{Title: "PROMPT PREVIEW", Align: termout.AlignLeft, MinWidth: 45},
		},
		Rows: buildPromptRows(prompts),
	}
	termout.PrintTable(cfg)
	termout.EnsureBottomPadding("\n")
}

func buildPromptRows(prompts []AgyPromptEntry) []termout.Row {
	var rows []termout.Row
	for _, p := range prompts {
		rows = append(rows, termout.Row{
			Cells: []string{
				p.CreatedAt.Format("2006-01-02 15:04"),
				truncatePromptString(p.Workspace, 20),
				truncatePromptString(p.Content, 45),
			},
		})
	}

	return rows
}

func truncatePromptString(s string, max int) string {
	clean := strings.ReplaceAll(strings.TrimSpace(s), "\n", " ")
	if len(clean) <= max {
		return clean
	}

	return clean[:max-3] + "..."
}
