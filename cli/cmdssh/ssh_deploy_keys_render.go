// Package cmdssh — ssh_deploy_keys_render.go formats terminal output for mesh key deployment.
package cmdssh

import (
	"encoding/json"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func renderDeployKeysSummary(summary DeployKeysSummary, except string, isJSON bool) error {
	if isJSON {
		b, _ := json.MarshalIndent(summary, "", "  ")
		fmt.Println(string(b))
		return nil
	}
	printDeployKeysHeader(summary, except)
	printDeployKeysRows(summary.NodeResults)
	printDeployKeysFooter(summary)
	return nil
}

func printDeployKeysHeader(s DeployKeysSummary, except string) {
	fmt.Printf("\n  %s🔑 Mesh SSH Public Key Deployment%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    • Unique Keys Identified: %s%d%s (local: %d, remote: %d)\n",
		constants.ColorBold, s.UniqueKeysIdentified, constants.ColorReset, s.LocalKeysCollected, s.RemoteKeysCollected)
	fmt.Printf("    • Target Fleet Nodes:     %d (succeeded: %d, except=%q)\n\n",
		s.NodesTargeted, s.NodesSucceeded, except)
	fmt.Printf("    %-4s %-20s %-16s %-12s %s\n", "#", "ALIAS", "IP ADDRESS", "STATUS", "KEYS ADDED")
	fmt.Printf("    %s\n", constants.TermTableRule)
}

func printDeployKeysRows(results []DeployKeysNodeResult) {
	for idx, r := range results {
		statusStr := formatNodeStatus(r)
		addedStr := fmt.Sprintf("+%d key(s)", r.KeysAdded)
		if r.KeysAdded == 0 && r.IsOnline {
			addedStr = "already synced"
		}
		fmt.Printf("    %-4d %-20s %-16s %-12s %s\n", idx+1, r.Alias, r.IPAddress, statusStr, addedStr)
	}
}

func formatNodeStatus(r DeployKeysNodeResult) string {
	if !r.IsOnline {
		return constants.ColorRed + "offline" + constants.ColorReset
	}
	if r.ErrorMsg != "" {
		return constants.ColorRed + "failed" + constants.ColorReset
	}
	return constants.ColorGreen + "synced" + constants.ColorReset
}

func printDeployKeysFooter(s DeployKeysSummary) {
	fmt.Println()
	if s.IsDryRun {
		fmt.Printf("  %s[dry-run] No remote authorized_keys were modified.%s\n\n", constants.ColorYellow, constants.ColorReset)
		return
	}
	fmt.Printf("  %s✓ Mesh public key synchronization complete! All nodes now trust cluster keys passwordlessly.%s\n\n",
		constants.ColorGreen, constants.ColorReset)
}
