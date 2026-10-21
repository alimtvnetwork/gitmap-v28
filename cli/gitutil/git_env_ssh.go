package gitutil

import (
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/cloneconcurrency"
)

// isHeadlessSSHSession checks whether execution is running inside a headless SSH session.
func isHeadlessSSHSession() bool {
	if os.Getenv("GITMAP_HEADLESS_SSH") == "1" {
		return true
	}

	if cloneconcurrency.IsSSHSession() {
		return true
	}

	return hasSSHSessionEnv()
}

// hasSSHSessionEnv inspects SSH environment variables indicating a remote SSH session.
func hasSSHSessionEnv() bool {
	if os.Getenv("SSH_CONNECTION") != "" || os.Getenv("SSH_CLIENT") != "" {
		return true
	}

	return os.Getenv("SSH_TTY") != ""
}

// ResolveWindowsCredentialStore returns dpapi in headless SSH sessions and wincredman in interactive sessions.
func ResolveWindowsCredentialStore() string {
	if isHeadlessSSHSession() {
		return "dpapi"
	}

	return "wincredman"
}
