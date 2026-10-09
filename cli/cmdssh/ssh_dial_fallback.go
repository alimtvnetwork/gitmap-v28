package cmdssh

import (
	"context"
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/cli/secrets"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

func tryDialWithSpecifiedKey(c db.SSHConnection) (*ssh.Client, bool) {
	if c.KeyPath == "" {
		return nil, false
	}
	if fi, err := os.Stat(c.KeyPath); err != nil || fi.IsDir() {
		return nil, false
	}
	client, err := secrets.ConnectWithKey(c.IPAddress, c.Username, c.KeyPath)
	if err == nil {
		return client, true
	}
	return nil, false
}

func tryDialWithDefaultKeys(c db.SSHConnection) (*ssh.Client, bool) {
	keys := findAllUserSSHKeys()
	for _, keyPath := range keys {
		client, err := cryptoConnectWithKeyFn(c.IPAddress, c.Username, keyPath)
		if err == nil {
			return client, true
		}
	}
	return nil, false
}

func tryDecryptCandidate(encrypted string) string {
	if encrypted == "" {
		return ""
	}
	plain, err := decryptPasswordCandidate(encrypted)
	if err == nil && plain != "" {
		return plain
	}
	return ""
}

func resolveCandidatePassword(c db.SSHConnection) string {
	if plain := tryDecryptCandidate(c.EncryptedPassword); plain != "" {
		return plain
	}
	encFromDB := queryHostPasswordFromDB(c.Alias, c.IPAddress)
	if plain := tryDecryptCandidate(encFromDB); plain != "" {
		return plain
	}
	return ResolveFallbackCredentials(c.Username, c.OS)
}

func maybePersistAuthenticatedPassword(c db.SSHConnection, plain string) {
	if c.EncryptedPassword == "" {
		saveSSHPasswordOnAuth(c.Alias, c.IPAddress, plain)
	}
}

func tryDialWithPasswordCandidate(c db.SSHConnection) (*ssh.Client, bool) {
	plain := resolveCandidatePassword(c)
	if plain == "" {
		return nil, false
	}
	client, err := secrets.ConnectWithPassword(c.IPAddress, c.Username, plain)
	if err != nil {
		return nil, false
	}
	maybePersistAuthenticatedPassword(c, plain)
	return client, true
}

func diagnoseDialFailure(c db.SSHConnection) error {
	isOnline, reason := CheckConnLiveness(context.Background(), c.IPAddress, 22, 1500*time.Millisecond)
	if !isOnline {
		return fmt.Errorf("ssh dial failed for %s@%s: network unreachable (%s) - verify hostname/IP, DNS, and port 22 firewall", c.Username, c.IPAddress, reason)
	}
	plain := resolveCandidatePassword(c)
	if plain == "" {
		return fmt.Errorf("ssh dial failed for %s@%s: no valid key found and no password in vault (enroll with 'gitmap sj add-with-pass %s@%s <password>' or 'gitmap ssh deploy keys')", c.Username, c.IPAddress, c.Username, c.IPAddress)
	}
	return fmt.Errorf("ssh dial failed for %s@%s: port 22 reachable, but authentication failed (credentials rejected by remote host - check password with 'gitmap ssh pass show %s')", c.Username, c.IPAddress, c.Alias)
}

func dialNodeWithFallback(c db.SSHConnection, header string) (*ssh.Client, error) {
	isOnline, reason := CheckConnLiveness(context.Background(), c.IPAddress, 22, 0)
	if !isOnline {
		return nil, fmt.Errorf("ssh dial failed for %s@%s: network unreachable (%s) - verify hostname/IP, DNS, and port 22 firewall", c.Username, c.IPAddress, reason)
	}
	if client, isKeyOk := tryDialWithSpecifiedKey(c); isKeyOk {
		return client, nil
	}
	if client, isDefaultKeyOk := tryDialWithDefaultKeys(c); isDefaultKeyOk {
		return client, nil
	}
	if client, isPassOk := tryDialWithPasswordCandidate(c); isPassOk {
		return client, nil
	}
	return nil, diagnoseDialFailure(c)
}

// DialSSHConnectionWithFallback connects to an SSH node using key, default keys, or vaulted/fallback password.
func DialSSHConnectionWithFallback(c db.SSHConnection) (*ssh.Client, error) {
	return dialNodeWithFallback(c, "["+c.Alias+"|"+c.IPAddress+"]")
}
