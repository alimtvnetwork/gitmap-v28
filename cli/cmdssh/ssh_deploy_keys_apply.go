// Package cmdssh — ssh_deploy_keys_apply.go applies unique SSH public keys to remote authorized_keys.
package cmdssh

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
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
		return res, fmt.Errorf("auth failed for %s", header)
	}
	defer client.Close()
	return syncAuthorizedKeysOnClient(client, res, uniqueKeys, isDryRun)
}

func syncAuthorizedKeysOnClient(client *ssh.Client, res DeployKeysNodeResult, uniqueKeys []string, isDryRun bool) (DeployKeysNodeResult, error) {
	prepCmd := "mkdir -p ~/.ssh && chmod 700 ~/.ssh && touch ~/.ssh/authorized_keys && chmod 600 ~/.ssh/authorized_keys && cat ~/.ssh/authorized_keys"
	existingAuth, err := crypto.RunCommand(client, prepCmd, "")
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
	return appendMissingKeysToRemote(client, res, missingKeys)
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

func appendMissingKeysToRemote(client *ssh.Client, res DeployKeysNodeResult, missing []string) (DeployKeysNodeResult, error) {
	builder := strings.Builder{}
	for _, k := range missing {
		builder.WriteString(fmt.Sprintf("printf '%%s\\n' %q >> ~/.ssh/authorized_keys\n", k))
	}
	builder.WriteString("chmod 600 ~/.ssh/authorized_keys\n")
	_, err := crypto.RunCommand(client, builder.String(), "")
	if err != nil {
		res.ErrorMsg = err.Error()
		return res, err
	}
	return res, nil
}
