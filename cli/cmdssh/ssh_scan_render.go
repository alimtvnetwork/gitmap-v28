package cmdssh

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

func renderScanTable(results []SSHHealthResult) {
	fmt.Printf("\n%s SSH Fleet Liveness & Reachability Scan:%s\n\n", constants.ColorCyan, constants.ColorReset)

	cfg := termout.TableConfig{
		Columns: scanTableColumns(),
		Rows:    buildScanTableRows(results),
	}
	termout.PrintTable(cfg)

	onlineCount := countOnlineResults(results)
	fmt.Printf("\n  Summary: %d/%d nodes online\n\n", onlineCount, len(results))
}

func scanTableColumns() []termout.Column {
	return []termout.Column{
		{Title: "ALIAS", Align: termout.AlignLeft},
		{Title: "IP", Align: termout.AlignLeft},
		{Title: "USER", Align: termout.AlignLeft},
		{Title: "PORT", Align: termout.AlignRight},
		{Title: "STATUS", Align: termout.AlignLeft},
		{Title: "LATENCY", Align: termout.AlignRight},
		{Title: "DETAILS", Align: termout.AlignLeft},
	}
}

func buildScanTableRows(results []SSHHealthResult) []termout.Row {
	var rows []termout.Row
	for _, r := range results {
		color := constants.ColorGreen
		if !r.IsOnline {
			color = constants.ColorRed
		}

		rows = append(rows, termout.Row{
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
