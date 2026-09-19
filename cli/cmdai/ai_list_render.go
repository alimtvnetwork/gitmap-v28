package cmdai

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/termpad"
	"github.com/alimtvnetwork/gitmap-v28/cli/termtable"
)

func renderScriptsTable(scripts []ScriptMetadata) {
	fmt.Printf("\n  %s● AI Scripts Catalog (%d scripts)%s\n\n", constants.ColorCyan, len(scripts), constants.ColorReset)
	cfg := termtable.TableConfig{
		Columns: scriptTableColumns(),
		Rows:    buildScriptRows(scripts),
	}
	termtable.PrintTable(cfg)
	termpad.EnsureBottomPadding("\n")
}

func scriptTableColumns() []termtable.Column {
	return []termtable.Column{
		{Title: "#", Align: termtable.AlignLeft, MinWidth: 4},
		{Title: "NAME", Align: termtable.AlignLeft, MinWidth: 24},
		{Title: "CATEGORY", Align: termtable.AlignLeft, MinWidth: 12},
		{Title: "AUTOFIX", Align: termtable.AlignLeft, MinWidth: 8},
		{Title: "DESCRIPTION", Align: termtable.AlignLeft, MinWidth: 42},
	}
}

func buildScriptRows(scripts []ScriptMetadata) []termtable.Row {
	var rows []termtable.Row
	for _, s := range scripts {
		rows = append(rows, buildScriptRow(s))
	}

	return rows
}

func buildScriptRow(s ScriptMetadata) termtable.Row {
	fixLabel := formatAutofixBadge(s.HasFixMode)
	desc := truncateDescription(s.Description, 50)

	return termtable.Row{
		Cells: []string{
			s.Number,
			s.Slug,
			string(s.Category),
			fixLabel,
			desc,
		},
	}
}

func formatAutofixBadge(hasFix bool) string {
	if hasFix {
		return "yes"
	}

	return "-"
}

func truncateDescription(text string, maxLen int) string {
	clean := strings.ReplaceAll(strings.TrimSpace(text), "\n", " ")
	isShort := len(clean) <= maxLen
	if isShort {
		return clean
	}

	return clean[:maxLen-3] + "..."
}

func renderScriptsJson(scripts []ScriptMetadata) *apperror.AppError {
	data, err := json.MarshalIndent(scripts, "", "  ")
	hasErr := err != nil
	if hasErr {
		return apperror.WrapSimple(err, "marshal_json")
	}

	_, writeErr := fmt.Fprintln(os.Stdout, string(data))
	hasWriteErr := writeErr != nil
	if hasWriteErr {
		return apperror.WrapSimple(writeErr, "write_json")
	}

	return nil
}
