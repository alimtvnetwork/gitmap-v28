package cmdssh

import (
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

// RunSSHCloneCLI executes the remote repository clone with self-healing authentication.
func RunSSHCloneCLI(args []string) error {
	opts := parseSSHCloneOptions(args)
	if opts.IsHelp {
		printSSHCloneUsage()
		return nil
	}

	repoURL, repoName, rawDest, err := resolveCloneTargetAndPath(opts)
	if err != nil {
		printSSHCloneUsage()
		return err
	}

	conns, err := loadTargetNodes(opts.Target)
	if err != nil {
		return err
	}

	filtered := filterSSHConns(conns, opts.Except)
	return executeCloneOnFleet(filtered, repoURL, repoName, rawDest)
}

func executeCloneOnFleet(conns []db.SSHConnection, repoURL, repoName, rawDest string) error {
	fmt.Printf("\n%s● Remote Cloning '%s' across %d node(s):%s\n\n",
		constants.ColorCyan, repoURL, len(conns), constants.ColorReset)

	var lastErr error
	for _, c := range conns {
		if err := cloneOnNode(c, repoURL, repoName, rawDest); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

func cloneOnNode(c db.SSHConnection, repoURL, repoName, rawDest string) error {
	header := fmt.Sprintf("[%s|%s]", c.Alias, c.IPAddress)
	if !checkRemoteNodeOnline(c.IPAddress, header) {
		return apperror.NewExecutionError("node is offline: " + header)
	}

	client, isConnected := connectSSHClient(c, header)
	if !isConnected {
		return apperror.NewExecutionError("ssh authentication failed for " + header)
	}
	defer client.Close()

	destPath := resolveRemoteDestPathForNode(repoName, rawDest, c.OS)
	return executeCloneWithHealing(client, c, repoURL, destPath, header)
}

func executeCloneWithHealing(client *ssh.Client, c db.SSHConnection, repoURL, destPath, header string) error {
	shell := resolveRemoteShell(c.OS)
	cloneCmd := buildRemoteCloneCmd(repoURL, destPath, isWindowsOS(c.OS))

	out, err := crypto.RunCommand(client, cloneCmd, shell)
	if err == nil {
		fmt.Printf("  %s %s✓ Successfully cloned to %s%s\n", header, constants.ColorGreen, destPath, constants.ColorReset)
		return nil
	}

	if !isAuthFailure(out, err) {
		printAppErrorWithStack(header, "Clone Error", apperror.WrapSimple(err, "remote clone"))
		return err
	}

	healErr := executeSelfHealingAuth(client, c, repoURL, header)
	if healErr != nil {
		printAppErrorWithStack(header, "Auth Heal Failed", apperror.WrapSimple(healErr, "self-healing auth"))
		return healErr
	}

	retryOut, retryErr := crypto.RunCommand(client, cloneCmd, shell)
	if retryErr != nil {
		printAppErrorWithStack(header, "Clone Retry Failed", apperror.WrapSimple(retryErr, "retry clone"))
		return retryErr
	}

	fmt.Printf("  %s %s✓ Successfully cloned to %s (after self-healing)%s\n%s\n",
		header, constants.ColorGreen, destPath, constants.ColorReset, strings.TrimSpace(retryOut))
	return nil
}

func buildRemoteCloneCmd(repoURL, destPath string, isWin bool) string {
	if isWin {
		return fmt.Sprintf(`powershell -NoProfile -Command "$dir = [System.IO.Path]::GetDirectoryName('%s'); if ($dir -and -not (Test-Path $dir)) { New-Item -ItemType Directory -Path $dir -Force | Out-Null }; git clone '%s' '%s'"`, destPath, repoURL, destPath)
	}
	return fmt.Sprintf(`mkdir -p "$(dirname '%s')" && git clone '%s' '%s'`, destPath, repoURL, destPath)
}
