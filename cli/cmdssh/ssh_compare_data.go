package cmdssh

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

func compareTableColumns() []termout.Column {
	return []termout.Column{
		{Title: "SUBSYSTEM", Align: termout.AlignLeft},
		{Title: "PRIMARY FOCUS", Align: termout.AlignLeft},
		{Title: "JOIN COMMAND", Align: termout.AlignLeft},
		{Title: "EXEC COMMAND", Align: termout.AlignLeft},
		{Title: "MONITORING", Align: termout.AlignLeft},
		{Title: "BEST USED WHEN", Align: termout.AlignLeft},
	}
}

func compareTableRows() []termout.Row {
	return []termout.Row{
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
