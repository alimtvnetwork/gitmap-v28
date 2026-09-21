package cmdprompt

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdprompttemplate"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/termpad"
	"github.com/alimtvnetwork/gitmap-v28/cli/termtable"
)

// RunPromptList displays all available prompt templates in a table.
func RunPromptList(args []string) error {
	templates, err := cmdprompttemplate.LoadTemplates()
	if err != nil {
		return apperror.WrapSimple(err, "load prompt templates")
	}

	RenderPromptTemplatesTable(templates)

	return nil
}

// RenderPromptTemplatesTable formats and prints prompt templates.
func RenderPromptTemplatesTable(templates []cmdprompttemplate.PromptTemplate) {
	fmt.Printf("\n%s  ● Available Prompt Templates (%d registered)%s\n\n",
		constants.ColorCyan, len(templates), constants.ColorReset)
	cfg := termtable.TableConfig{
		Columns: buildPromptTemplateColumns(),
		Rows:    buildPromptTemplateRows(templates),
	}
	termtable.PrintTable(cfg)
	termpad.EnsureBottomPadding("\n")
}

func buildPromptTemplateColumns() []termtable.Column {
	return []termtable.Column{
		{Title: "NAME", Align: termtable.AlignLeft, MinWidth: 15},
		{Title: "DESCRIPTION", Align: termtable.AlignLeft, MinWidth: 25},
		{Title: "PREVIEW", Align: termtable.AlignLeft, MinWidth: 35},
		{Title: "INVOCATION EXAMPLE", Align: termtable.AlignLeft, MinWidth: 35},
	}
}

func buildPromptTemplateRows(templates []cmdprompttemplate.PromptTemplate) []termtable.Row {
	rows := make([]termtable.Row, 0, len(templates))
	for _, tpl := range templates {
		rows = append(rows, buildSingleTemplateRow(tpl))
	}

	return rows
}

func buildSingleTemplateRow(tpl cmdprompttemplate.PromptTemplate) termtable.Row {
	desc := resolveTemplateDesc(tpl.Description)
	preview := truncateTemplatePreview(tpl.Content, 35)
	example := fmt.Sprintf("gitmap agy prompt -n %s", tpl.Name)

	return termtable.Row{
		Cells: []string{tpl.Name, desc, preview, example},
	}
}

func resolveTemplateDesc(desc string) string {
	hasDesc := len(strings.TrimSpace(desc)) > 0
	if hasDesc {
		return desc
	}

	return "(no description)"
}

func truncateTemplatePreview(content string, maxLen int) string {
	clean := strings.ReplaceAll(strings.TrimSpace(content), "\n", " ")
	hasSmall := len(clean) <= maxLen
	if hasSmall {
		return clean
	}

	return clean[:maxLen-3] + "..."
}
