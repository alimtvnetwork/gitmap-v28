// Package cmdssh — ssh_deploy_node_config_render.go renders deployment feedback and advice.
package cmdssh

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

func isExceptOrExcepFlag(low string) bool {
	return low == "--except" || low == "--excep" || low == "--accept" || low == "--exclude" || low == "-e"
}

func renderNodeConfigDeploySummary(targets []db.SSHConnection, except string, totalConfigNodes int, isDryRun, isJSON bool) error {
	var names []string
	for _, t := range targets {
		names = append(names, fmt.Sprintf("%s(%s)", t.Alias, t.IPAddress))
	}
	if isJSON {
		payload := map[string]any{
			"command":             "ssh deploy node-config",
			"config_nodes_synced": totalConfigNodes,
			"target_nodes":        names,
			"target_count":        len(targets),
			"except":              except,
			"dry_run":             isDryRun,
		}
		b, _ := json.MarshalIndent(payload, "", "  ")
		fmt.Println(string(b))
		return nil
	}
	fmt.Printf("✓ Deployed SSH node-config (%d node definitions) to %d target machine(s) [except=%q]: %s\n",
		totalConfigNodes, len(targets), except, strings.Join(names, ", "))
	fmt.Printf("  %sTip: Next, run 'gitmap ssh deploy keys all' to sync public keys across all nodes.%s\n",
		constants.ColorCyan, constants.ColorReset)
	return nil
}

func printDeployConfigSSHHelp() {
	fmt.Printf("\n  %s🚀 GitMap Deploy Config SSH (Fleet Configuration Deployment)%s\n\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println("    gitmap deploy config ssh [all|<target>] [flags]")
	fmt.Println("    gitmap deploy config ssh --file <json-file> [all|<target>] [flags]")
	fmt.Println("    gitmap ssh deploy node-config [all|<target>] [flags] (alias: nc)")
	fmt.Println()
	fmt.Println("  Description:")
	fmt.Println("    Deploys SSH node topology, IP addresses, aliases, and encrypted credentials")
	fmt.Println("    from local CLI configuration (Split-DB) or an exported JSON file to remote fleet nodes.")
	fmt.Println("    Target machines automatically update matching records without duplicate inserts.")
	fmt.Println()
	fmt.Println("  Flags:")
	fmt.Println("    -f, --file <path>      Deploy from specified JSON file instead of local CLI config")
	fmt.Println("    -e, --except <tokens>  Exclude nodes by ID, worker ID, IP, or alias")
	fmt.Println("    -n, --dry-run          Preview deployment actions without executing remote changes")
	fmt.Println("    -j, --json             Output machine-readable JSON telemetry")
	fmt.Println("    -h, --help             Show this help menu")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println("    gitmap deploy config ssh")
	fmt.Println("    gitmap deploy config ssh w3")
	fmt.Println("    gitmap deploy config ssh --file D:\\work\\repo-secrets\\01-gitmap\\gitmap-ssh-nodes.json")
	fmt.Println("    gitmap deploy config ssh all --except w1")
	fmt.Println()
}
