package clonenow

import (
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func runInteractiveGitClone(cmd *exec.Cmd) (string, bool) {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err.Error(), false
	}

	return "", true
}

func isSSHCloneURL(url string) bool {
	lower := strings.ToLower(strings.TrimSpace(url))
	return strings.HasPrefix(lower, constants.PrefixSSH) ||
		strings.HasPrefix(lower, constants.PrefixSSHScheme)
}
