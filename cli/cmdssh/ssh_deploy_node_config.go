// Package cmdssh — ssh_deploy_node_config.go handles deploying node configuration across the fleet.
package cmdssh

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

// RunSSHDeployNodeConfigCLI deploys local SSH node topology to remote nodes (--except id,ip,alias).
func RunSSHDeployNodeConfigCLI(args []string) error {
	return RunSSHDeployConfigSSHCLI(args)
}

func executeNodeConfigDeployWorker(c db.SSHConnection, remoteCmd string) (string, error) {
	header := fmt.Sprintf("[%s|%s]", c.Alias, c.IPAddress)
	if !checkRemoteNodeOnline(c.IPAddress, header) {
		return "", fmt.Errorf("node %s is offline", header)
	}
	client, isConnected := connectSSHClient(c, header)
	if !isConnected {
		return "", fmt.Errorf("auth failed for %s", header)
	}
	defer client.Close()
	return crypto.RunCommand(client, remoteCmd, "")
}

func parseNodeConfigDeployFlags(args []string) (string, bool, bool) {
	var exceptParts []string
	isDryRun, isJSON := false, false
	for i := 0; i < len(args); i++ {
		a := args[i]
		low := strings.ToLower(a)
		if low == "node-config" || low == "nc" || low == "deploy" || low == "ssh" || low == "all" {
			continue
		}
		if low == "--dry-run" || low == "-n" {
			isDryRun = true
			continue
		}
		if low == "--json" {
			isJSON = true
			continue
		}
		if isExceptOrExcepFlag(low) {
			for i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				exceptParts = append(exceptParts, args[i+1])
				i++
			}
			continue
		}
		if strings.HasPrefix(low, "--except=") || strings.HasPrefix(low, "--excep=") || strings.HasPrefix(low, "--accept=") || strings.HasPrefix(low, "--exclude=") {
			idx := strings.IndexByte(a, '=')
			exceptParts = append(exceptParts, a[idx+1:])
		}
	}
	return strings.Join(exceptParts, ","), isDryRun, isJSON
}
