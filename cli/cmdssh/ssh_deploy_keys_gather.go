// Package cmdssh — ssh_deploy_keys_gather.go gathers and deduplicates SSH public keys across nodes.
package cmdssh

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

func gatherLocalPublicKeys() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	sshDir := filepath.Join(home, ".ssh")
	matches, _ := filepath.Glob(filepath.Join(sshDir, "*.pub"))
	var out []string
	for _, pubPath := range matches {
		data, readErr := os.ReadFile(pubPath)
		if readErr != nil {
			continue
		}
		lines := strings.Split(string(data), "\n")
		out = append(out, lines...)
	}
	return out
}

func gatherRemotePublicKeys(c db.SSHConnection) ([]string, error) {
	header := fmt.Sprintf("[%s|%s]", c.Alias, c.IPAddress)
	if !checkRemoteNodeOnline(c.IPAddress, header) {
		return nil, fmt.Errorf("node %s is offline", header)
	}
	client, isConnected := connectSSHClient(c, header)
	if !isConnected {
		return nil, fmt.Errorf("auth failed for %s", header)
	}
	defer client.Close()
	cmd := "cat ~/.ssh/id_ed25519.pub ~/.ssh/id_rsa.pub ~/.ssh/id_ecdsa.pub ~/.ssh/*.pub 2>/dev/null"
	out, err := crypto.RunCommand(client, cmd, "")
	if err != nil {
		return nil, err
	}
	return strings.Split(out, "\n"), nil
}

func deduplicatePublicKeys(allKeys []string) []string {
	seen := make(map[string]bool)
	var unique []string
	for _, raw := range allKeys {
		trimmed := strings.TrimSpace(raw)
		if !isValidSSHPublicKeyLine(trimmed) {
			continue
		}
		keySignature := extractKeySignature(trimmed)
		if seen[keySignature] {
			continue
		}
		seen[keySignature] = true
		unique = append(unique, trimmed)
	}
	return unique
}

func isValidSSHPublicKeyLine(line string) bool {
	if line == "" || strings.HasPrefix(line, "#") {
		return false
	}
	return strings.HasPrefix(line, "ssh-") ||
		strings.HasPrefix(line, "ecdsa-") ||
		strings.HasPrefix(line, "sk-ssh-") ||
		strings.HasPrefix(line, "sk-ecdsa-")
}

func extractKeySignature(line string) string {
	parts := strings.Fields(line)
	if len(parts) >= 2 {
		return parts[0] + " " + parts[1]
	}
	return line
}
