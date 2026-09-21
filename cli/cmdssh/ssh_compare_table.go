package cmdssh

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func RunSSHCompareCLI(args []string) error {
	PrintArchitectureComparisonTable()

	return nil
}

func runSSHCompareCLI(args []string) error {
	return RunSSHCompareCLI(args)
}

// PrintArchitectureComparisonTable displays comparison matrix for ssh vs cluster vs sc.
func PrintArchitectureComparisonTable() {
	fmt.Printf("\n  %sGitMap Remote Architecture Comparison:%s\n\n", constants.ColorCyan, constants.ColorReset)
	printSSHCompareCards()
	printComparisonWorkflowGuidance()
}

func printSSHCompareCards() {
	printCompareCard("gitmap ssh (Direct Node Management)", constants.ColorGreen,
		"Direct ad-hoc node execution & setup",
		"gitmap ssh join <u@ip> [alias]",
		"gitmap ssh exec <command>",
		"gitmap ssh scan / status",
		"Ad-hoc command run, remote install/update, AGY/code open")
	printCompareCard("gitmap cluster (Multi-Node Orchestration)", constants.ColorCyan,
		"Multi-node cluster orchestration & recipes",
		"gitmap cluster node add <ip>",
		"gitmap cluster exec <target> <cmd>",
		"gitmap cluster node ls",
		"K8s bootstrap, cluster recipes, distributed scripts")
	printCompareCard("gitmap sc (Servers-Clients Fleet Daemon)", constants.ColorYellow,
		"Client-server daemon topology & sync",
		"gitmap sc join <server-url>",
		"gitmap sc exec <command>",
		"gitmap sc status",
		"Master-worker topology, continuous sync, live telemetry")
}

func printCompareCard(title, color, focus, joinCmd, execCmd, monitor, bestUsed string) {
	borderLen := 58 - len(title)
	hasPositiveLen := borderLen > 2
	if !hasPositiveLen {
		borderLen = 2
	}
	border := strings.Repeat("─", borderLen)
	fmt.Printf("  %s┌─ %s %s%s\n", color, title, border, constants.ColorReset)
	fmt.Printf("  │  Focus:     %s\n", focus)
	fmt.Printf("  │  Join:      %s%s%s\n", constants.ColorWhite, joinCmd, constants.ColorReset)
	fmt.Printf("  │  Exec:      %s%s%s\n", constants.ColorWhite, execCmd, constants.ColorReset)
	fmt.Printf("  │  Monitor:   %s\n", monitor)
	fmt.Printf("  │  Best Used: %s\n", bestUsed)
	fmt.Printf("  %s└──────────────────────────────────────────────────────────%s\n\n", color, constants.ColorReset)
}

func printComparisonWorkflowGuidance() {
	fmt.Printf("  %sWhen to use which remote command:%s\n", constants.ColorYellow, constants.ColorReset)
	fmt.Println("    • Use 'ssh'     for instant terminal commands, checking liveness, or installing gitmap/AGY on remote machines.")
	fmt.Println("    • Use 'cluster' for orchestrating multi-node infrastructure, running distributed recipes, and node control.")
	fmt.Println("    • Use 'sc'      for continuous daemon monitoring, client-server sync, and central fleet management.")
	fmt.Println()
}
