package cmdclone

import (
	"os"
	"path/filepath"
	"strings"
)

const githubKnownHostsKey = "github.com ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOMqqnkVzrm0SdG6UOoqKLsabgH5C9okWi0dh2l9GKJl\n"

func ensureKnownHostsFile() {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}

	sshDir := filepath.Join(home, ".ssh")
	_ = os.MkdirAll(sshDir, 0700)

	knownHostsPath := filepath.Join(sshDir, "known_hosts")
	if hasHostInKnownHosts(knownHostsPath, "github.com") {
		return
	}

	f, err := os.OpenFile(knownHostsPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return
	}
	defer f.Close()

	_, _ = f.WriteString(githubKnownHostsKey)
}

func hasHostInKnownHosts(filePath, host string) bool {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return false
	}

	return strings.Contains(string(data), host)
}
