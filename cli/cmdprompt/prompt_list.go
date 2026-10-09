package cmdprompt

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdprompttemplate"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
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
	cfg := termout.TableConfig{
		Columns: buildPromptTemplateColumns(),
		Rows:    buildPromptTemplateRows(templates),
	}
	termout.PrintTable(cfg)
	termout.EnsureBottomPadding("\n")
}

func buildPromptTemplateColumns() []termout.Column {
	return []termout.Column{
		{Title: "NAME", Align: termout.AlignLeft, MinWidth: 15},
		{Title: "DESCRIPTION", Align: termout.AlignLeft, MinWidth: 25},
		{Title: "PREVIEW", Align: termout.AlignLeft, MinWidth: 35},
		{Title: "INVOCATION EXAMPLE", Align: termout.AlignLeft, MinWidth: 35},
	}
}

func buildPromptTemplateRows(templates []cmdprompttemplate.PromptTemplate) []termout.Row {
	rows := make([]termout.Row, 0, len(templates))
	for _, tpl := range templates {
		rows = append(rows, buildSingleTemplateRow(tpl))
	}

	return rows
}

func buildSingleTemplateRow(tpl cmdprompttemplate.PromptTemplate) termout.Row {
	desc := resolveTemplateDesc(tpl.Description)
	preview := truncateTemplatePreview(tpl.Content, 35)
	example := fmt.Sprintf("gitmap agy prompt -n %s", tpl.Name)

	return termout.Row{
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
