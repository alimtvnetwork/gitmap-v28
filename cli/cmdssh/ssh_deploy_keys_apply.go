// Package cmdssh — ssh_deploy_keys_apply.go applies unique SSH public keys to remote authorized_keys.
package cmdssh

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/secrets"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"golang.org/x/crypto/ssh"
)

func deployKeysToRemoteNode(c db.SSHConnection, uniqueKeys []string, isDryRun bool) (DeployKeysNodeResult, error) {
	header := fmt.Sprintf("[%s|%s]", c.Alias, c.IPAddress)
	res := DeployKeysNodeResult{Alias: c.Alias, IPAddress: c.IPAddress, KeysTotal: len(uniqueKeys)}
	if !checkRemoteNodeOnline(c.IPAddress, header) {
		res.ErrorMsg = "offline"
		return res, fmt.Errorf("node %s is offline", header)
	}
	res.IsOnline = true
	client, isConnected := connectSSHClient(c, header)
	if !isConnected {
		res.ErrorMsg = "authentication failed"
		logDeployKeysError(c, "authentication failed", len(uniqueKeys))
		return res, fmt.Errorf("auth failed for %s", header)
	}
	defer client.Close()
	nodeRes, err := syncAuthorizedKeysOnClient(client, res, uniqueKeys, c.OS, isDryRun)
	if err != nil {
		logDeployKeysError(c, err.Error(), len(uniqueKeys))
	}
	return nodeRes, err
}

func logDeployKeysError(c db.SSHConnection, errMsg string, keysTotal int) {
	store.LogInternalErrorRecord(store.InternalErrorRecord{
		ErrorCode:     "ERR_SSH_DEPLOY_KEYS",
		ErrorType:     "SSHDeploymentError",
		Command:       "gitmap ssh deploy keys",
		Message:       fmt.Sprintf("Mesh SSH key deployment failed on node %s (%s): %s", c.Alias, c.IPAddress, errMsg),
		Details:       errMsg,
		SourceFile:    "cli/cmdssh/ssh_deploy_keys_apply.go",
		ContextJson:   fmt.Sprintf(`{"node":"%s","ip":"%s","keys_total":%d}`, c.Alias, c.IPAddress, keysTotal),
		GitMapVersion: constants.Version,
	})
}

func syncAuthorizedKeysOnClient(client *ssh.Client, res DeployKeysNodeResult, uniqueKeys []string, osType string, isDryRun bool) (DeployKeysNodeResult, error) {
	existingAuth, err := readRemoteAuthorizedKeys(client, osType)
	if err != nil {
		res.ErrorMsg = err.Error()
		return res, err
	}
	existingSignatures := buildExistingSignaturesMap(existingAuth)
	missingKeys := filterMissingPublicKeys(uniqueKeys, existingSignatures)
	res.KeysAdded = len(missingKeys)
	if isDryRun || len(missingKeys) == 0 {
		return res, nil
	}
	return appendMissingKeysToRemote(client, res, missingKeys, osType)
}

func readRemoteAuthorizedKeys(client *ssh.Client, osType string) (string, error) {
	if isWindowsOS(osType) {
		cmd := `powershell -NoProfile -Command "Get-Content -Path (Join-Path $env:USERPROFILE '.ssh\authorized_keys') -ErrorAction SilentlyContinue; if (Test-Path (Join-Path $env:ProgramData 'ssh\administrators_authorized_keys')) { Get-Content -Path (Join-Path $env:ProgramData 'ssh\administrators_authorized_keys') -ErrorAction SilentlyContinue }"`
		return secrets.RunCommand(client, cmd, "")
	}
	prepCmd := "mkdir -p ~/.ssh && chmod 700 ~/.ssh && touch ~/.ssh/authorized_keys && chmod 600 ~/.ssh/authorized_keys && cat ~/.ssh/authorized_keys"
	return secrets.RunCommand(client, prepCmd, "")
}

func buildExistingSignaturesMap(authContent string) map[string]bool {
	out := make(map[string]bool)
	for _, line := range strings.Split(authContent, "\n") {
		trimmed := strings.TrimSpace(line)
		if isValidSSHPublicKeyLine(trimmed) {
			out[extractKeySignature(trimmed)] = true
		}
	}
	return out
}

func filterMissingPublicKeys(uniqueKeys []string, existing map[string]bool) []string {
	var missing []string
	for _, k := range uniqueKeys {
		if !existing[extractKeySignature(k)] {
			missing = append(missing, k)
		}
	}
	return missing
}

func appendMissingKeysToRemote(client *ssh.Client, res DeployKeysNodeResult, missing []string, osType string) (DeployKeysNodeResult, error) {
	if isWindowsOS(osType) {
		return appendMissingKeysWindows(client, res, missing)
	}
	return appendMissingKeysUnix(client, res, missing)
}

func appendMissingKeysUnix(client *ssh.Client, res DeployKeysNodeResult, missing []string) (DeployKeysNodeResult, error) {
	builder := strings.Builder{}
	for _, k := range missing {
		builder.WriteString(fmt.Sprintf("printf '%%s\\n' %q >> ~/.ssh/authorized_keys\n", k))
	}
	builder.WriteString("chmod 600 ~/.ssh/authorized_keys\n")
	_, err := secrets.RunCommand(client, builder.String(), "")
	if err != nil {
		res.ErrorMsg = err.Error()
		return res, err
	}
	return res, nil
}

func appendMissingKeysWindows(client *ssh.Client, res DeployKeysNodeResult, missing []string) (DeployKeysNodeResult, error) {
	for _, k := range missing {
		script := buildInjectAuthKeyScript(k, "windows")
		_, err := secrets.RunCommand(client, script, "")
		if err != nil {
			res.ErrorMsg = err.Error()
			return res, err
		}
	}
	return res, nil
}
