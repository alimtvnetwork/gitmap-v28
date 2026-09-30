// Package cmdssh — ssh_deploy_keys_render.go formats terminal output for mesh key deployment.
package cmdssh

import (
	"encoding/json"
	"fmt"
	"strings"

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
		statusCell := formatNodeStatusCell(r, 12)
		actionText := formatNodeActionText(r)
		fmt.Printf("    %-4d %-20s %-16s %s%s\n", idx+1, r.Alias, r.IPAddress, statusCell, actionText)
	}
}

func formatNodeStatusCell(r DeployKeysNodeResult, width int) string {
	raw, color := "synced", constants.ColorGreen
	if !r.IsOnline {
		raw, color = "offline", constants.ColorRed
	} else if r.ErrorMsg != "" {
		raw, color = "failed", constants.ColorRed
	}
	padding := width - len(raw)
	if padding < 1 {
		padding = 1
	}
	return fmt.Sprintf("%s%s%s%s", color, raw, constants.ColorReset, strings.Repeat(" ", padding))
}

func formatNodeActionText(r DeployKeysNodeResult) string {
	if !r.IsOnline {
		return "+0 key(s)"
	}
	if r.ErrorMsg != "" {
		return fmt.Sprintf("error: %s", r.ErrorMsg)
	}
	if r.KeysAdded == 0 {
		return "already synced"
	}
	return fmt.Sprintf("+%d key(s)", r.KeysAdded)
}

func printDeployKeysFooter(s DeployKeysSummary) {
	fmt.Println()
	if s.IsDryRun {
		fmt.Printf("  %s[dry-run] No remote authorized_keys were modified.%s\n", constants.ColorYellow, constants.ColorReset)
		printDeployKeysSuggestions(s)
		return
	}
	if s.NodesSucceeded == 0 && s.NodesTargeted > 0 {
		fmt.Printf("  %s✗ Mesh public key deployment failed: 0/%d nodes updated.%s\n",
			constants.ColorRed, s.NodesTargeted, constants.ColorReset)
		printDeployKeysSuggestions(s)
		return
	}
	if s.NodesSucceeded < s.NodesTargeted {
		fmt.Printf("  %s⚠ Mesh public key deployment partially complete (%d/%d succeeded).%s\n",
			constants.ColorYellow, s.NodesSucceeded, s.NodesTargeted, constants.ColorReset)
		printDeployKeysSuggestions(s)
		return
	}
	fmt.Printf("  %s✓ Mesh public key synchronization complete! All nodes now trust cluster keys passwordlessly.%s\n",
		constants.ColorGreen, constants.ColorReset)
	printDeployKeysSuggestions(s)
}

func printDeployKeysSuggestions(s DeployKeysSummary) {
	fmt.Printf("\n  %s💡 SSH Key Deploy & Optimization Suggestions:%s\n", constants.ColorCyan, constants.ColorReset)
	if s.NodesSucceeded < s.NodesTargeted {
		fmt.Println("    • Repair auth on failed node:     gitmap ssh copy-id <alias|user@ip>")
		fmt.Println("    • Skip offline nodes during sync: gitmap ssh deploy-keys --except <offline-alias>")
		fmt.Println("    • Inspect SSH diagnostic logs:    gitmap ssh error-logs")
	}
	fmt.Println("    • Add external public key locally: gitmap ssh key add <pubkey-or-file>")
	fmt.Println("    • Sync ~/.ssh/config to all nodes: gitmap ssh deploy-config")
	fmt.Println("    • Verify passwordless mesh exec:   gitmap ssh exec all \"hostname\"")
	fmt.Println()
}
