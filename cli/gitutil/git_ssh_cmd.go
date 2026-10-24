package gitutil

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
)

// quoteKeyPath normalizes slashes and safely encloses a path in escaped quotes.
func quoteKeyPath(rawPath string) string {
	trimmed := strings.Trim(strings.TrimSpace(rawPath), "\"")
	cleanPath := fsutil.NormalizeToForwardSlashes(trimmed)

	return "\"" + cleanPath + "\""
}

const SSHDeterministicTemplate = "%s %s %s %s %s %s %s %s %s %s %s %s %s"

// BuildDeterministicSSHCommand synthesizes an isolated OpenSSH invocation string.
func BuildDeterministicSSHCommand(privateKeyPath string) string {
	quotedKey := quoteKeyPath(privateKeyPath)
	c := constants.SSHOptionFlag

	return fmt.Sprintf(
		SSHDeterministicTemplate,
		constants.SSHBin, constants.SSHFlagConfigFile, constants.DevNullPath,
		c, constants.SSHIdentitiesOnlyYes, constants.SSHFlagIdentity, quotedKey,
		c, constants.SSHStrictHostKeyAcceptNew, c, constants.SSHBatchModeYes,
		c, constants.SSHConnectTimeout15,
	)
}

// appendAuthorEnv appends Git author and committer overrides when provided.
func appendAuthorEnv(env []string, ctx ResolvedGitContext) []string {
	if len(ctx.AuthorName) > 0 {
		env = append(env, constants.EnvGitAuthorName+"="+ctx.AuthorName)
		env = append(env, constants.EnvGitCommitterName+"="+ctx.AuthorName)
	}

	if len(ctx.AuthorEmail) > 0 {
		env = append(env, constants.EnvGitAuthorEmail+"="+ctx.AuthorEmail)
		env = append(env, constants.EnvGitCommitterEmail+"="+ctx.AuthorEmail)
	}

	return env
}

// BuildSafeGitEnvForRepo constructs a subprocess environment with resolved SSH credentials.
func BuildSafeGitEnvForRepo(ctx ResolvedGitContext, extra ...string) []string {
	combinedExtra := make([]string, 0, len(extra)+5)
	hasSSH := ctx.HasInjectedSSH && len(ctx.GitSSHCommand) > 0
	if hasSSH {
		sshVar := constants.EnvGitSSHCommand + constants.EnvAssignmentSeparator + ctx.GitSSHCommand
		combinedExtra = append(combinedExtra, sshVar)
	}

	combinedExtra = appendAuthorEnv(combinedExtra, ctx)
	combinedExtra = append(combinedExtra, extra...)

	return BuildSafeGitEnv(combinedExtra...)
}

// CreateGitCmdInDirWithContext creates an exec.Cmd configured with directory and resolved context.
func CreateGitCmdInDirWithContext(repoPath string, ctx ResolvedGitContext, args ...string) *exec.Cmd {
	full := buildGitDirArgs(repoPath, args...)
	cmd := exec.Command(constants.GitBin, full...)
	cmd.Env = BuildSafeGitEnvForRepo(ctx)

	return cmd
}
