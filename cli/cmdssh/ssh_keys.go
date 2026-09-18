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

func connectWithDefaultKey(ip, user, header string) (*ssh.Client, bool) {
	keyPath := findDefaultUserSSHKey()
	if keyPath == "" {
		return nil, false
	}

	client, err := crypto.ConnectWithKey(ip, user, keyPath)
	if err != nil {
		fmt.Printf("%s Default key (%s) failed: %v\n", header, filepath.Base(keyPath), err)

		return nil, false
	}

	return client, true
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
