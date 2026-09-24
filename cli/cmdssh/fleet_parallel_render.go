package cmdssh

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/termtable"
)

// PrintFleetStart displays immediate announcement when execution begins on a node.
func PrintFleetStart(alias, ip, task string) {
	fmt.Printf("%s[FLEET START]%s Executing on [%s] (IP: %s) [%s]...\n",
		constants.ColorCyan, constants.ColorReset, alias, ip, task)
}

// PrintFleetDone displays immediate completion output for a single node.
func PrintFleetDone(res FleetNodeResult) {
	if res.Success {
		fmt.Printf("%s[FLEET DONE]%s  [%s|%s]: %sSUCCESS%s (took %dms)\n",
			constants.ColorGreen, constants.ColorReset,
			res.Alias, res.IP, constants.ColorGreen, constants.ColorReset, res.DurationMs)
		if trimmed := strings.TrimSpace(res.Output); trimmed != "" {
			fmt.Printf("  %s\n", trimmed)
		}
		return
	}
	fmt.Printf("%s[FLEET FAIL]%s  [%s|%s]: %sFAILED%s - %v (took %dms)\n",
		constants.ColorRed, constants.ColorReset,
		res.Alias, res.IP, constants.ColorRed, constants.ColorReset, res.Error, res.DurationMs)
}

func printFleetNoTargets(task string) {
	fmt.Printf("%s[FLEET]%s No matching target nodes found for task '%s'.\n\n",
		constants.ColorYellow, constants.ColorReset, task)
}

// PrintFleetSummary prints an aggregated table of all fleet execution results.
func PrintFleetSummary(taskName string, results []FleetNodeResult) {
	if len(results) == 0 {
		return
	}
	fmt.Printf("\n%s================================================================================%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf(" %sSSH Fleet Execution Summary:%s %s\n", constants.ColorBold, constants.ColorReset, taskName)
	fmt.Printf(" Total: %d | Succeeded: %d | Failed: %d\n",
		len(results), countFleetSuccess(results), countFleetFailure(results))
	fmt.Printf("%s--------------------------------------------------------------------------------%s\n", constants.ColorDim, constants.ColorReset)

	cfg := termtable.TableConfig{
		Columns: []termtable.Column{
			{Title: "ALIAS", Align: termtable.AlignLeft, MinWidth: 15},
			{Title: "IP", Align: termtable.AlignLeft, MinWidth: 16},
			{Title: "STATUS", Align: termtable.AlignLeft, MinWidth: 10},
			{Title: "DURATION", Align: termtable.AlignRight, MinWidth: 10},
			{Title: "DETAILS", Align: termtable.AlignLeft, MinWidth: 25},
		},
		Rows: buildFleetSummaryRows(results),
	}
	termtable.PrintTable(cfg)
	fmt.Printf("%s================================================================================%s\n\n", constants.ColorCyan, constants.ColorReset)
}

func countFleetSuccess(results []FleetNodeResult) int {
	count := 0
	for _, r := range results {
		if r.Success {
			count++
		}
	}
	return count
}

func countFleetFailure(results []FleetNodeResult) int {
	count := 0
	for _, r := range results {
		if !r.Success {
			count++
		}
	}
	return count
}

func buildFleetSummaryRows(results []FleetNodeResult) []termtable.Row {
	rows := make([]termtable.Row, 0, len(results))
	for _, r := range results {
		statusStr := constants.ColorGreen + "SUCCESS" + constants.ColorReset
		details := "OK"
		if !r.Success {
			statusStr = constants.ColorRed + "FAILED" + constants.ColorReset
			if r.Error != nil {
				details = r.Error.Error()
			}
		}
		rows = append(rows, termtable.Row{
			Cells: []string{
				r.Alias,
				r.IP,
				statusStr,
				fmt.Sprintf("%dms", r.DurationMs),
				details,
			},
		})
	}
	return rows
}
