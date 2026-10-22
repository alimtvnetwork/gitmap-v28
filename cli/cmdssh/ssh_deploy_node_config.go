package cmdssh

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/secrets"
	"golang.org/x/crypto/ssh"
	"time"
)

// RunSSHDeployNodeConfigCLI deploys local SSH node topology to remote nodes (--except id,ip,alias).
func RunSSHDeployNodeConfigCLI(args []string) error {
	return RunSSHDeployConfigSSHCLI(args)
}

func executeNodeConfigDeployWorker(c db.SSHConnection, payload []byte) (string, error) {
	header := fmt.Sprintf("[%s|%s]", c.Alias, c.IPAddress)
	isOnline, reason := CheckConnLiveness(context.Background(), c.IPAddress, 22, 1500*time.Millisecond)
	if !isOnline {
		return "", apperror.NewExecutionError(fmt.Sprintf("node %s is offline (%s)", header, reason))
	}
	client, err := connectSSHClientWithErr(c, header)
	if err != nil {
		return "", apperror.WrapSimple(err, "auth failed for "+header)
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
		return "", apperror.WrapSimple(streamErr, "failed to stream config payload to "+c.Alias)
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
