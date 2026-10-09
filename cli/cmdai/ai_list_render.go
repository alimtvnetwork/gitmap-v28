package cmdai

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

func renderScriptsTable(scripts []ScriptMetadata) {
	fmt.Printf("\n  %s● AI Scripts Catalog (%d scripts)%s\n\n", constants.ColorCyan, len(scripts), constants.ColorReset)
	cfg := termout.TableConfig{
		Columns: scriptTableColumns(),
		Rows:    buildScriptRows(scripts),
	}
	termout.PrintTable(cfg)
	termout.EnsureBottomPadding("\n")
}

func scriptTableColumns() []termout.Column {
	return []termout.Column{
		{Title: "#", Align: termout.AlignLeft, MinWidth: 4},
		{Title: "NAME", Align: termout.AlignLeft, MinWidth: 24},
		{Title: "CATEGORY", Align: termout.AlignLeft, MinWidth: 12},
		{Title: "AUTOFIX", Align: termout.AlignLeft, MinWidth: 8},
		{Title: "DESCRIPTION", Align: termout.AlignLeft, MinWidth: 42},
	}
}

func buildScriptRows(scripts []ScriptMetadata) []termout.Row {
	var rows []termout.Row
	for _, s := range scripts {
		rows = append(rows, buildScriptRow(s))
	}

	return rows
}

func buildScriptRow(s ScriptMetadata) termout.Row {
	fixLabel := formatAutofixBadge(s.HasFixMode)
	desc := truncateDescription(s.Description, 50)

	return termout.Row{
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
