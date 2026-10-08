package cmdupdate

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/termtable"
)

func printFleetNoTargetsBanner(pkg string, excludedCount int) {
	fmt.Printf("\n%s[FLEET UPDATE]%s No matching active targets found for package '%s' (excluded: %d).\n\n",
		constants.ColorYellow, constants.ColorReset, pkg, excludedCount)
}

func renderFleetUpdateSummary(results []FleetUpdateNodeResult, pkg string, excludedCount int) {
	if len(results) == 0 {
		return
	}
	successCount, failCount, offlineCount := calculateFleetMetrics(results)
	fmt.Printf("\n%s================================================================================%s\n",
		constants.ColorCyan, constants.ColorReset)
	fmt.Printf(" %sSSH Fleet Update Summary [%s]:%s Total: %d | Succeeded: %d | Failed: %d | Offline: %d | Excluded: %d\n",
		constants.ColorBold, pkg, constants.ColorReset, len(results), successCount, failCount, offlineCount, excludedCount)
	fmt.Printf("%s--------------------------------------------------------------------------------%s\n",
		constants.ColorDim, constants.ColorReset)

	cfg := termtable.TableConfig{
		Columns: []termtable.Column{
			{Title: "ALIAS", Align: termtable.AlignLeft, MinWidth: 15},
			{Title: "IP", Align: termtable.AlignLeft, MinWidth: 16},
			{Title: "STATUS", Align: termtable.AlignLeft, MinWidth: 10},
			{Title: "DURATION", Align: termtable.AlignRight, MinWidth: 10},
			{Title: "DETAILS", Align: termtable.AlignLeft, MinWidth: 25},
		},
		Rows: buildFleetUpdateRows(results),
	}
	termtable.PrintTable(cfg)
	fmt.Printf("%s================================================================================%s\n\n",
		constants.ColorCyan, constants.ColorReset)
}

func calculateFleetMetrics(results []FleetUpdateNodeResult) (int, int, int) {
	succeeded := 0
	failed := 0
	offline := 0
	for _, r := range results {
		if r.IsSuccess {
			succeeded++
			continue
		}
		if r.IsOffline {
			offline++
			continue
		}
		failed++
	}
	return succeeded, failed, offline
}

func buildFleetUpdateRows(results []FleetUpdateNodeResult) []termtable.Row {
	rows := make([]termtable.Row, 0, len(results))
	for _, r := range results {
		statusStr := constants.ColorRed + "FAILED" + constants.ColorReset
		if r.IsSuccess {
			statusStr = constants.ColorGreen + "SUCCESS" + constants.ColorReset
		} else if r.IsOffline {
			statusStr = constants.ColorYellow + "OFFLINE" + constants.ColorReset
		}
		detail := sanitizeTableRowDetail(r.Details)
		rows = append(rows, termtable.Row{
			Cells: []string{
				r.Alias,
				r.IP,
				statusStr,
				fmt.Sprintf("%dms", r.DurationMs),
				detail,
			},
		})
	}
	return rows
}

func sanitizeTableRowDetail(raw string) string {
	clean := strings.ReplaceAll(raw, "\r", "")
	clean = strings.ReplaceAll(clean, "\n", " ")
	clean = strings.TrimSpace(clean)
	if len(clean) > 50 {
		return clean[:47] + "..."
	}
	if clean == "" {
		return "OK"
	}
	return clean
}
