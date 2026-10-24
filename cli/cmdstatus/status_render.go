package cmdstatus

import (
	"fmt"

	statusconstants "github.com/alimtvnetwork/gitmap-v28/cli/cmdstatus/constants"
	statusenums "github.com/alimtvnetwork/gitmap-v28/cli/cmdstatus/enums"
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
		{Title: statusenums.StatusTableColRepo.String(), Alignment: termout.AlignLeft},
		{Title: statusenums.StatusTableColBranch.String(), Alignment: termout.AlignLeft},
		{Title: statusenums.StatusTableColStatus.String(), Alignment: termout.AlignCenter},
		{Title: statusenums.StatusTableColSync.String(), Alignment: termout.AlignLeft},
		{Title: statusenums.StatusTableColStash.String(), Alignment: termout.AlignRight},
		{Title: statusenums.StatusTableColFiles.String(), Alignment: termout.AlignLeft},
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
	missingIcon := fmt.Sprintf("%s%s%s", constants.ColorRed, statusconstants.StatusNotFoundText, constants.ColorReset)
	dash := statusconstants.StatusTableEmptyDash

	return termout.Row{
		Cells: []string{repoColored, dash, missingIcon, dash, dash, dash},
	}
}

func createNormalStatusRow(r statusRow, pastelColor string) termout.Row {
	repoColored := fmt.Sprintf("%s%s%s", pastelColor, r.RepoName, constants.ColorReset)
	branchColored := fmt.Sprintf("%s%s%s", constants.ColorCyan, r.Branch, constants.ColorReset)

	return termout.Row{
		Cells: []string{repoColored, branchColored, r.StateIcon, r.SyncText, r.StashText, r.FilesText},
	}
}
