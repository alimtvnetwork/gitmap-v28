package cmd

import (
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

func applyCloneAssumeYesEnv(isAssumeYes bool) {
	if !isAssumeYes {
		return
	}
	cmd := withSSHAcceptNew(os.Getenv(constants.EnvGitSSHCommand))
	if err := os.Setenv(constants.EnvGitSSHCommand, cmd); err != nil {
		fmtCloneEnvError(err)
	}
}

func cloneEnvWithSSHAcceptNew() []string {
	cmd := withSSHAcceptNew(os.Getenv(constants.EnvGitSSHCommand))
	return envWithOverride(constants.EnvGitSSHCommand, cmd)
}

func envWithOverride(key, value string) []string {
	prefix := key + constants.EnvAssignmentSeparator
	out := make([]string, 0, len(os.Environ())+1)
	for _, entry := range os.Environ() {
		if strings.HasPrefix(entry, prefix) {
			continue
		}
		out = append(out, entry)
	}
	return append(out, prefix+value)
}

func withSSHAcceptNew(existing string) string {
	trimmed := strings.TrimSpace(existing)
	if trimmed == "" {
		return constants.SSHBin + " " + constants.SSHOptionFlag + " " +
			constants.SSHStrictHostKeyAcceptNew
	}
	if strings.Contains(trimmed, constants.SSHStrictHostKeyChecking) {
		return trimmed
	}

	return trimmed + " " + constants.SSHOptionFlag + " " +
		constants.SSHStrictHostKeyAcceptNew
}

func isSSHCloneURL(url string) bool {
	lower := strings.ToLower(strings.TrimSpace(url))
	return strings.HasPrefix(lower, constants.PrefixSSH) ||
		strings.HasPrefix(lower, constants.PrefixSSHScheme)
}
