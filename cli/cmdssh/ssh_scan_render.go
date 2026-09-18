package cmdssh

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/termtable"
)

func renderScanTable(results []SSHHealthResult) {
	fmt.Printf("\n%s SSH Fleet Liveness & Reachability Scan:%s\n\n", constants.ColorCyan, constants.ColorReset)

	cfg := termtable.TableConfig{
		Columns: scanTableColumns(),
		Rows:    buildScanTableRows(results),
	}
	termtable.PrintTable(cfg)

	onlineCount := countOnlineResults(results)
	fmt.Printf("\n  Summary: %d/%d nodes online\n\n", onlineCount, len(results))
}

func scanTableColumns() []termtable.Column {
	return []termtable.Column{
		{Title: "ALIAS", Align: termtable.AlignLeft},
		{Title: "IP", Align: termtable.AlignLeft},
		{Title: "USER", Align: termtable.AlignLeft},
		{Title: "PORT", Align: termtable.AlignRight},
		{Title: "STATUS", Align: termtable.AlignLeft},
		{Title: "LATENCY", Align: termtable.AlignRight},
		{Title: "DETAILS", Align: termtable.AlignLeft},
	}
}

func buildScanTableRows(results []SSHHealthResult) []termtable.Row {
	var rows []termtable.Row
	for _, r := range results {
		color := constants.ColorGreen
		if !r.IsOnline {
			color = constants.ColorRed
		}

		rows = append(rows, termtable.Row{
			Cells: []string{
				r.Alias,
				r.IP,
				r.User,
				fmt.Sprintf("%d", r.Port),
				r.Status,
				fmt.Sprintf("%v", r.Latency),
				r.Details,
			},
			Color: color,
		})
	}

	return rows
}

func countOnlineResults(results []SSHHealthResult) int {
	count := 0
	for _, r := range results {
		if r.IsOnline {
			count++
		}
	}

	return count
}
