package cmdssh

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
)

func findDefaultUserSSHKey() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	candidates := []string{
		filepath.Join(home, ".ssh", "id_ed25519"),
		filepath.Join(home, ".ssh", "id_rsa"),
		filepath.Join(home, ".ssh", "id_ecdsa"),
	}

	for _, p := range candidates {
		if fi, statErr := os.Stat(p); statErr == nil && !fi.IsDir() {
			return p
		}
	}

	return ""
}

func findAllUserSSHKeys() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	candidates := []string{
		filepath.Join(home, ".ssh", "id_ed25519"),
		filepath.Join(home, ".ssh", "id_rsa"),
		filepath.Join(home, ".ssh", "id_ecdsa"),
		filepath.Join(home, ".ssh", "gitmap_id_rsa"),
		filepath.Join(home, ".ssh", "gitmap_id_ed25519"),
	}

	var found []string
	for _, p := range candidates {
		if fi, statErr := os.Stat(p); statErr == nil && fi.IsDir() == false {
			found = append(found, p)
		}
	}

	return found
}

func connectWithDefaultKey(ip, user, header string) (*ssh.Client, bool) {
	keys := findAllUserSSHKeys()
	if len(keys) == 0 {
		return nil, false
	}

	for _, keyPath := range keys {
		client, err := crypto.ConnectWithKey(ip, user, keyPath)
		if err == nil {
			return client, true
		}
	}

	return nil, false
}

func formatMissingAuthAdvice(alias, ip, user string) string {
	home, _ := os.UserHomeDir()
	sampleKey := filepath.Join(home, ".ssh", "id_rsa")

	return fmt.Sprintf(
		"No credentials configured for machine '%s' (%s@%s).\n"+
			"  To fix, run one of the following:\n"+
			"    1. gitmap ssh-join add-with-pass %s@%s <password> %s\n"+
			"    2. gitmap ssh join %s@%s %s --auth\n"+
			"    3. Configure standard key at: %s",
		alias, user, ip, user, ip, alias, user, ip, alias, sampleKey,
	)
}
