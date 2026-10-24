package gitutil

import (
	"os"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// BuildSafeGitEnv constructs a sanitized subprocess environment for automated Git commands.
// On Windows, GCM_CREDENTIAL_STORE=cache is excluded because Git for Windows lacks UNIX socket support.
func BuildSafeGitEnv(extra ...string) []string {
	return BuildSafeGitEnvWithBase(os.Environ(), extra...)
}

// AppendGitCredentialEnv appends safe credential store environment variables for Git subprocesses.
func AppendGitCredentialEnv(env []string) []string {
	if runtime.GOOS != constants.PlatformWindows {
		return append(env, constants.EnvGCMCredentialStoreCache)
	}

	if isHeadlessSSHSession() {
		return append(env, constants.EnvGCMCredentialStoreDPAPI, constants.EnvGCMNoPersistOne)
	}

	return append(env, constants.EnvGCMCredentialStoreWinCM)
}

func defaultSafeGitVars() []string {
	return []string{
		constants.EnvGitTerminalPromptZero,
		constants.EnvGitAskpassEmpty,
		constants.EnvSSHAskpassEmpty,
		constants.EnvGCMInteractiveNever,
		constants.EnvGCMNoPersistOne,
	}
}

// BuildSafeGitEnvWithBase appends safe Git execution variables to an existing environment slice.
func BuildSafeGitEnvWithBase(base []string, extra ...string) []string {
	safe := AppendGitCredentialEnv(defaultSafeGitVars())
	res := make([]string, 0, len(base)+len(safe)+len(extra))
	res = append(res, base...)
	res = append(res, safe...)
	res = append(res, extra...)

	return res
}

// HasGitSSHCommand checks whether an environment slice contains a GIT_SSH_COMMAND setting.
func HasGitSSHCommand(env []string) bool {
	prefix := constants.EnvGitSSHCommand + constants.EnvAssignmentSeparator
	for _, v := range env {
		if strings.HasPrefix(v, prefix) {
			return true
		}
	}

	return false
}

// ExtractGitSSHCommand extracts the GIT_SSH_COMMAND value from an environment slice if present.
func ExtractGitSSHCommand(env []string) (string, bool) {
	prefix := constants.EnvGitSSHCommand + constants.EnvAssignmentSeparator
	for _, v := range env {
		if strings.HasPrefix(v, prefix) {
			return strings.TrimPrefix(v, prefix), true
		}
	}

	return "", false
}
