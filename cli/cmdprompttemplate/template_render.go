package cmdprompttemplate

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/termpad"
	"github.com/alimtvnetwork/gitmap-v28/cli/termtable"
)

// RenderTemplatesTable displays registered prompt templates in a formatted table.
func RenderTemplatesTable(list []PromptTemplate) {
	fmt.Printf("\n%s  ● Prompt Templates (%d registered)%s\n\n", constants.ColorCyan, len(list), constants.ColorReset)
	cfg := termtable.TableConfig{
		Columns: []termtable.Column{
			{Title: "NAME", Align: termtable.AlignLeft, MinWidth: 15},
			{Title: "CONTENT PREVIEW", Align: termtable.AlignLeft, MinWidth: 45},
			{Title: "UPDATED", Align: termtable.AlignLeft, MinWidth: 12},
		},
		Rows: buildTableRows(list),
	}
	termtable.PrintTable(cfg)
	termpad.EnsureBottomPadding("\n")
}

func buildTableRows(list []PromptTemplate) []termtable.Row {
	var rows []termtable.Row
	for _, item := range list {
		rows = append(rows, termtable.Row{
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
