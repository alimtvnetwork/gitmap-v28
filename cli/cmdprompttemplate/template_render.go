package cmdprompttemplate

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// RenderTemplatesTable displays registered prompt templates in a formatted table.
func RenderTemplatesTable(list []PromptTemplate) {
	fmt.Printf("\n%s  ● Prompt Templates (%d registered)%s\n\n", constants.ColorCyan, len(list), constants.ColorReset)
	cfg := termout.TableConfig{
		Columns: []termout.Column{
			{Title: "NAME", Align: termout.AlignLeft, MinWidth: 15},
			{Title: "CONTENT PREVIEW", Align: termout.AlignLeft, MinWidth: 45},
			{Title: "UPDATED", Align: termout.AlignLeft, MinWidth: 12},
		},
		Rows: buildTableRows(list),
	}
	termout.PrintTable(cfg)
	termout.EnsureBottomPadding("\n")
}

func buildTableRows(list []PromptTemplate) []termout.Row {
	var rows []termout.Row
	for _, item := range list {
		rows = append(rows, termout.Row{
			Cells: []string{
				item.Name,
				truncateContent(item.Content, 45),
				item.UpdatedAt.Format("2006-01-02"),
			},
		})
	}

	return rows
}

func truncateContent(s string, max int) string {
	clean := strings.ReplaceAll(strings.TrimSpace(s), "\n", " ")
	if len(clean) <= max {
		return clean
	}

	return clean[:max-3] + "..."
}
