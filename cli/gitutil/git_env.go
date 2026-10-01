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

// BuildSafeGitEnvWithBase appends safe Git execution variables to an existing environment slice.
func BuildSafeGitEnvWithBase(base []string, extra ...string) []string {
	safe := []string{
		constants.EnvGitTerminalPromptZero,
		constants.EnvGitAskpassEmpty,
		constants.EnvSSHAskpassEmpty,
		constants.EnvGCMInteractiveNever,
		"GCM_NO_PERSIST=1",
	}
	if runtime.GOOS != constants.PlatformWindows {
		safe = append(safe, "GCM_CREDENTIAL_STORE=cache")
	}
	res := append(base, safe...)
	if len(extra) > 0 {
		res = append(res, extra...)
	}
	return res
}
