package cmdssh

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

func tryDialWithSpecifiedKey(c db.SSHConnection) (*ssh.Client, bool) {
	if c.KeyPath == "" {
		return nil, false
	}
	if fi, err := os.Stat(c.KeyPath); err != nil || fi.IsDir() {
		return nil, false
	}
	client, err := crypto.ConnectWithKey(c.IPAddress, c.Username, c.KeyPath)
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
	client, err := crypto.ConnectWithPassword(c.IPAddress, c.Username, plain)
	if err != nil {
		return nil, false
	}
	maybePersistAuthenticatedPassword(c, plain)
	return client, true
}

func dialNodeWithFallback(c db.SSHConnection, header string) (*ssh.Client, error) {
	if client, isKeyOk := tryDialWithSpecifiedKey(c); isKeyOk {
		return client, nil
	}
	if client, isDefaultKeyOk := tryDialWithDefaultKeys(c); isDefaultKeyOk {
		return client, nil
	}
	if client, isPassOk := tryDialWithPasswordCandidate(c); isPassOk {
		return client, nil
	}
	return nil, fmt.Errorf("ssh dial failed for %s@%s: unable to authenticate with key or password", c.Username, c.IPAddress)
}

// DialSSHConnectionWithFallback connects to an SSH node using key, default keys, or vaulted/fallback password.
func DialSSHConnectionWithFallback(c db.SSHConnection) (*ssh.Client, error) {
	return dialNodeWithFallback(c, "["+c.Alias+"|"+c.IPAddress+"]")
}
