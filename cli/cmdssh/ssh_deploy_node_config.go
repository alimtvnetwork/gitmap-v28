package cmdssh

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"
	"golang.org/x/crypto/ssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/secrets"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

// RunSSHDeployNodeConfigCLI deploys local SSH node topology to remote nodes (--except id,ip,alias).
func RunSSHDeployNodeConfigCLI(args []string) error {
	return RunSSHDeployConfigSSHCLI(args)
}

func executeNodeConfigDeployWorker(c db.SSHConnection, payload []byte) (string, error) {
	header := fmt.Sprintf("[%s|%s]", c.Alias, c.IPAddress)
	isOnline, reason := CheckConnLiveness(context.Background(), c.IPAddress, 22, 1500*time.Millisecond)
	if !isOnline {
		return "", fmt.Errorf("node %s is offline (%s)", header, reason)
	}
	client, err := connectSSHClientWithErr(c, header)
	if err != nil {
		return "", fmt.Errorf("auth failed for %s: %w", header, err)
	}
	defer client.Close()
	return deployConfigPayloadToNode(client, c, payload)
}

func deployConfigPayloadToNode(client *ssh.Client, c db.SSHConnection, payload []byte) (string, error) {
	isWin := isWindowsOS(c.OS)
	remotePath := ".gitmap/fleet_config_deploy.json"
	b64 := base64.StdEncoding.EncodeToString(payload)
	if len(b64) < 1500 {
		cmd := fmt.Sprintf("gitmap ssh nodes import-json --base64 \"%s\"", b64)
		return secrets.RunCommand(client, cmd, "")
	}

	streamErr := streamDirectToRemote(client, remotePath, payload, isWin)
	if streamErr != nil {
		return "", fmt.Errorf("failed to stream config payload to %s: %w", c.Alias, streamErr)
	}
	defer cleanupRemoteTempFile(client, remotePath, isWin)

	importCmd := fmt.Sprintf("gitmap ssh nodes import-json %s", remotePath)
	return secrets.RunCommand(client, importCmd, "")
}
func cleanupRemoteTempFile(client *ssh.Client, remotePath string, isWin bool) {
	if isWin {
		_, _ = secrets.RunCommand(client, fmt.Sprintf("powershell.exe -NoProfile -Command \"Remove-Item -Force '%s' 2>$null\"", remotePath), "")
		return
	}
	_, _ = secrets.RunCommand(client, fmt.Sprintf("rm -f '%s'", remotePath), "")
}