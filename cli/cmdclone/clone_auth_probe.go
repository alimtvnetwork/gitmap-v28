// Package cmdclone — clone_auth_probe.go probes SSH access before falling back to tokens.
package cmdclone

import (
	"os"
	"os/exec"
	"sync"
)

var sshProbeCache sync.Map

// ProbeSSHAccess tests if the target repository is accessible over SSH.
// If accessible, returns the SSH URL and true; otherwise returns empty and false.
func ProbeSSHAccess(repoURL string) (string, bool) {
	sshURL, ok := ConvertURLToSSH(repoURL)
	if !ok {
		return "", false
	}

	if val, isFound := sshProbeCache.Load(sshURL); isFound {
		isSuccess, _ := val.(bool)

		return sshURL, isSuccess
	}

	isAccessible := executeSSHProbe(sshURL)
	sshProbeCache.Store(sshURL, isAccessible)

	return sshURL, isAccessible
}

func executeSSHProbe(sshURL string) bool {
	cmd := exec.Command("git", "ls-remote", "--exit-code", "-h", sshURL, "HEAD")
	cmd.Env = append(os.Environ(),
		"GIT_SSH_COMMAND=ssh -o BatchMode=yes -o ConnectTimeout=5 -o StrictHostKeyChecking=accept-new",
		"GIT_TERMINAL_PROMPT=0",
	)

	return cmd.Run() == nil
}
