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
