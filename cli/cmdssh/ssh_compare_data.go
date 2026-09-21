package cmdssh

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/termtable"
)

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
