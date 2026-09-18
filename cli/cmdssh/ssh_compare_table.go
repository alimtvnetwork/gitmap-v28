package cmdssh

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/termtable"
)

func runSSHCompareCLI(args []string) error {
	PrintArchitectureComparisonTable()
	return nil
}

// PrintArchitectureComparisonTable displays comparison matrix for ssh vs cluster vs sc.
func PrintArchitectureComparisonTable() {
	fmt.Printf("\n%s GitMap Remote Subsystems Architecture Comparison:%s\n\n",
		constants.ColorCyan, constants.ColorReset)

	cfg := termtable.TableConfig{
		Columns: compareTableColumns(),
		Rows:    compareTableRows(),
	}
	termtable.PrintTable(cfg)
	printComparisonWorkflowGuidance()
}

func compareTableColumns() []termtable.Column {
	return []termtable.Column{
		{Title: "SUBSYSTEM", Align: termtable.AlignLeft},
		{Title: "PRIMARY FOCUS", Align: termtable.AlignLeft},
		{Title: "JOIN COMMAND", Align: termtable.AlignLeft},
		{Title: "EXEC COMMAND", Align: termtable.AlignLeft},
		{Title: "MONITORING", Align: termtable.AlignLeft},
		{Title: "BEST USED WHEN", Align: termtable.AlignLeft},
	}
}

func compareTableRows() []termtable.Row {
	return []termtable.Row{
		{
			Cells: []string{
				"gitmap ssh",
				"Direct node-level management",
				"gitmap ssh join <u@ip> [alias]",
				"gitmap ssh exec <command>",
				"gitmap ssh scan",
				"Ad-hoc command run, remote install/update, AGY/code open",
			},
			Color: constants.ColorGreen,
		},
		{
			Cells: []string{
				"gitmap cluster",
				"Multi-node cluster orchestration",
				"gitmap cluster node add <ip>",
				"gitmap cluster exec <target> <cmd>",
				"gitmap cluster node ls",
				"K8s bootstrap, cluster recipes, distributed scripts",
			},
			Color: constants.ColorCyan,
		},
		{
			Cells: []string{
				"gitmap sc",
				"Servers-Clients fleet daemon",
				"gitmap sc join <server-url>",
				"gitmap sc exec <command>",
				"gitmap sc status",
				"Master-worker topology, continuous sync, live telemetry",
			},
			Color: constants.ColorYellow,
		},
	}
}

func printComparisonWorkflowGuidance() {
	fmt.Printf("\n  %sWhen to use which remote command:%s\n", constants.ColorYellow, constants.ColorReset)
	fmt.Println("    • Use 'ssh'     for instant terminal commands, checking liveness, or installing gitmap/AGY on remote machines.")
	fmt.Println("    • Use 'cluster' for orchestrating multi-node infrastructure, running distributed recipes, and node control.")
	fmt.Println("    • Use 'sc'      for continuous daemon monitoring, client-server sync, and central fleet management.\n")
}
