package cmdupdate

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"golang.org/x/crypto/ssh"
)

func dialFleetSSH(target FleetTarget) (*ssh.Client, error) {
	conn := db.SSHConnection{
		Alias:             target.Alias,
		IPAddress:         target.IP,
		Username:          target.Username,
		EncryptedPassword: target.Password,
		KeyPath:           target.KeyPath,
		OS:                target.OS,
	}
	client, isOk := cmdssh.ConnectSSHClient(conn, fmt.Sprintf("[%s|%s]", target.Alias, target.IP))
	if isOk && client != nil {
		return client, nil
	}
	if c, ok := tryDialFleetPassword(target); ok {
		return c, nil
	}
	if c, ok := tryDialFleetKey(target); ok {
		return c, nil
	}
	if c, ok := tryDialFleetCandidateKeys(target); ok {
		return c, nil
	}
	return nil, fmt.Errorf("ssh dial failed for %s@%s", target.Username, target.IP)
}

func tryDialFleetPassword(target FleetTarget) (*ssh.Client, bool) {
	if target.Password == "" {
		return nil, false
	}
	plain, err := crypto.DecryptStoredPassword(target.Password)
	if err != nil || plain == "" {
		plain = target.Password
	}
	c, connErr := crypto.ConnectWithPassword(target.IP, target.Username, plain)
	return c, connErr == nil
}

func tryDialFleetKey(target FleetTarget) (*ssh.Client, bool) {
	if target.KeyPath == "" {
		return nil, false
	}
	c, err := crypto.ConnectWithKey(target.IP, target.Username, target.KeyPath)
	return c, err == nil
}

func tryDialFleetCandidateKeys(target FleetTarget) (*ssh.Client, bool) {
	home, _ := os.UserHomeDir()
	candidateKeys := []string{
		fmt.Sprintf("%s/.ssh/id_ed25519", home),
		fmt.Sprintf("%s/.ssh/id_rsa", home),
	}
	for _, k := range candidateKeys {
		c, err := crypto.ConnectWithKey(target.IP, target.Username, k)
		if err == nil {
			return c, true
		}
	}
	return nil, false
}
