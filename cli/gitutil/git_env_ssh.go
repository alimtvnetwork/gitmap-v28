package gitutil

import (
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/cloneconcurrency"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// isHeadlessSSHSession checks whether execution is running inside a headless SSH session.
func isHeadlessSSHSession() bool {
	if os.Getenv(constants.EnvGitMapHeadlessSSH) == constants.EnvValueOne {
		return true
	}

	if cloneconcurrency.IsSSHSession() {
		return true
	}

	return hasSSHSessionEnv()
}

// hasSSHSessionEnv inspects SSH environment variables indicating a remote SSH session.
func hasSSHSessionEnv() bool {
	if os.Getenv(constants.EnvSSHConnection) != "" || os.Getenv(constants.EnvSSHClient) != "" {
		return true
	}

	return os.Getenv(constants.EnvSSHTTY) != ""
}

// ResolveWindowsCredentialStore returns dpapi in headless SSH sessions and wincredman in interactive sessions.
func ResolveWindowsCredentialStore() constants.CredentialStoreType {
	if isHeadlessSSHSession() {
		return constants.CredentialStoreDPAPI
	}

	return constants.CredentialStoreWinCredMan
}
