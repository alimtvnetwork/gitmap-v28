package cmdcd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// RenderCDSelectionTable outputs an interactive repository selection table strictly to os.Stderr.
func RenderCDSelectionTable(title string, records []model.ScanRecord) {
	if len(records) == 0 {
		return
	}

	cfg := buildCDTableConfig(records)
	rendered := termout.RenderTable(cfg)
	fmt.Fprintf(os.Stderr, "\n  %s\n%s\n", title, rendered)
}

func buildCDTableConfig(records []model.ScanRecord) termout.TableConfig {
	return termout.TableConfig{
		Columns:     buildCDTableColumns(),
		Rows:        buildCDTableRows(records),
		HeaderColor: termout.RoleSGR(termout.RolePrimary),
		HasBorders:  true,
		BorderStyle: termout.BorderRounded,
	}
}

func buildCDTableColumns() []termout.ColumnConfig {
	return []termout.ColumnConfig{
		{Title: "#", MinWidth: 4, MaxWidth: 4, Alignment: termout.AlignRight},
		{Title: "Repository", MinWidth: 16, MaxWidth: 24, Alignment: termout.AlignLeft},
		{Title: "Branch", MinWidth: 10, MaxWidth: 16, Alignment: termout.AlignLeft},
		{Title: "Path", MinWidth: 28, MaxWidth: 56, Alignment: termout.AlignLeft},
	}
}

func buildCDTableRows(records []model.ScanRecord) []termout.Row {
	var rows []termout.Row
	for i, r := range records {
		rows = append(rows, formatCDRecordRow(i+1, r))
	}

	return rows
}

func formatCDRecordRow(idx int, r model.ScanRecord) termout.Row {
	branch := resolveCDBranchText(r.Branch)
	path := termout.TruncateMiddleWithEllipsis(r.AbsolutePath, 56)

	return termout.Row{
		Cells: []string{
			fmt.Sprintf("%d", idx),
			r.RepoName,
			branch,
			path,
		},
		Color: termout.RoleSGR(termout.RoleFG),
	}
}

func resolveCDBranchText(branch string) string {
	if branch == "" {
		return "-"
	}

	return branch
}
