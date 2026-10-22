package cmdstatus

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// RenderStatusTable renders the status summary table using centralized termout table engine.
func RenderStatusTable(c *statusTableContext) {
	renderStatusTableEngine(c)
}

func renderStatusTableEngine(c *statusTableContext) {
	if c == nil {
		return
	}
	cfg := termout.TableConfig{
		BorderStyle: termout.BorderRounded,
		HeaderColor: constants.ColorYellow,
		RowStripe:   true,
		Columns:     buildStatusColumns(),
		Rows:        buildStatusRows(c),
	}
	termout.PrintTable(cfg)
	printMissingRepoRemediation(c)
}

func buildStatusColumns() []termout.ColumnConfig {
	return []termout.ColumnConfig{
		{Title: "Repository", Alignment: termout.AlignLeft},
		{Title: "Branch", Alignment: termout.AlignLeft},
		{Title: "Status", Alignment: termout.AlignCenter},
		{Title: "Sync", Alignment: termout.AlignLeft},
		{Title: "Stash", Alignment: termout.AlignRight},
		{Title: "Files", Alignment: termout.AlignLeft},
	}
}

func buildStatusRows(c *statusTableContext) []termout.Row {
	var rows []termout.Row
	for i, r := range c.Rows {
		pastelColor := constants.ColorCycle[i%len(constants.ColorCycle)]
		rows = append(rows, convertStatusRow(r, pastelColor))
	}

	return rows
}

func convertStatusRow(r statusRow, pastelColor string) termout.Row {
	if r.Missing {
		return createMissingStatusRow(r, pastelColor)
	}

	return createNormalStatusRow(r, pastelColor)
}

func createMissingStatusRow(r statusRow, pastelColor string) termout.Row {
	repoColored := fmt.Sprintf("%s%s%s", pastelColor, r.RepoName, constants.ColorReset)
	missingIcon := fmt.Sprintf("%s✖ not found%s", constants.ColorRed, constants.ColorReset)

	return termout.Row{
		Cells: []string{repoColored, "-", missingIcon, "-", "-", "-"},
	}
}

func createNormalStatusRow(r statusRow, pastelColor string) termout.Row {
	repoColored := fmt.Sprintf("%s%s%s", pastelColor, r.RepoName, constants.ColorReset)
	branchColored := fmt.Sprintf("%s%s%s", constants.ColorCyan, r.Branch, constants.ColorReset)

	return termout.Row{
		Cells: []string{repoColored, branchColored, r.StateIcon, r.SyncText, r.StashText, r.FilesText},
	}
}
