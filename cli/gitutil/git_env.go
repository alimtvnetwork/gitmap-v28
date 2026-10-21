package gitutil

import (
	"os"
	"runtime"

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
		return append(env, "GCM_CREDENTIAL_STORE=cache")
	}

	if isHeadlessSSHSession() {
		return append(env, "GCM_CREDENTIAL_STORE=dpapi", "GCM_NO_PERSIST=1")
	}

	return append(env, "GCM_CREDENTIAL_STORE=wincredman")
}

func defaultSafeGitVars() []string {
	return []string{
		constants.EnvGitTerminalPromptZero,
		constants.EnvGitAskpassEmpty,
		constants.EnvSSHAskpassEmpty,
		constants.EnvGCMInteractiveNever,
		"GCM_NO_PERSIST=1",
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
