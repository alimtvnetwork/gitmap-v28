package cmdssh

import (
	"fmt"

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
