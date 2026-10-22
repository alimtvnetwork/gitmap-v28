package gitutil

import (
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// buildGitDirArgs constructs arguments prepended with the git directory flag when repoPath is provided.
func buildGitDirArgs(repoPath string, args ...string) []string {
	hasRepoPath := len(repoPath) > 0
	if !hasRepoPath {
		return args
	}

	full := make([]string, 0, len(args)+2)
	full = append(full, constants.GitDirFlag, repoPath)
	full = append(full, args...)

	return full
}

// createGitCmd prepares an exec.Cmd configured with directory arguments and safe Git execution environment.
func createGitCmd(repoPath string, args ...string) *exec.Cmd {
	full := buildGitDirArgs(repoPath, args...)
	cmd := exec.Command(constants.GitBin, full...)
	cmd.Env = BuildSafeGitEnv(constants.EnvGitSSHCommandBatchYes)

	return cmd
}

// RunGitInDir executes a git command with directory arguments and returns stdout output.
func RunGitInDir(repoPath string, args ...string) result.ByteSliceResult {
	cmd := createGitCmd(repoPath, args...)
	out, err := cmd.Output()
	if err != nil {
		return result.Fail[[]byte](apperror.WrapSimple(err, "RunGitInDir"))
	}

	return result.Ok(out)
}

// RunGitCombinedInDir executes a git command with directory arguments and returns combined stdout and stderr.
func RunGitCombinedInDir(repoPath string, args ...string) result.ByteSliceResult {
	cmd := createGitCmd(repoPath, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		// Preserve the combined output alongside the error: callers such as
		// the pull remediation flow inspect git's own text to classify failures.
		return result.ByteSliceResult{Value: out, Data: out, Err: apperror.WrapSimple(err, "RunGitCombinedInDir")}
	}

	return result.Ok(out)
}

// GetRepoConfig retrieves a git configuration value for key in repoPath.
func GetRepoConfig(repoPath, key string) result.StringResult {
	res := RunGitInDir(repoPath, constants.GitConfigCmd, constants.GitGetFlag, key)
	if res.IsFailure() {
		return result.Fail[string](res.AppError())
	}

	return result.Ok(strings.TrimSpace(string(res.Data)))
}

// GetRemoteOriginURL retrieves the origin remote URL or an empty string on error.
func GetRemoteOriginURL(repoPath string) string {
	res := GetRepoConfig(repoPath, constants.GitRemoteOrigin)

	return res.UnwrapOr("")
}

// ReadRemoteOriginURL retrieves the configured remote.origin.url for repoPath.
func ReadRemoteOriginURL(repoPath string) result.StringResult {
	return GetRepoConfig(repoPath, constants.GitRemoteOrigin)
}

// SetRemoteOriginURL updates the remote origin URL for repoPath.
func SetRemoteOriginURL(repoPath, url string) error {
	res := RunGitCombinedInDir(repoPath, "remote", "set-url", constants.GitOrigin, url)
	if res.IsFailure() {
		return res.AppError()
	}

	return nil
}

// GetRepoCurrentBranch returns the active branch name of the repository at repoPath.
func GetRepoCurrentBranch(repoPath string) result.StringResult {
	res := RunGitInDir(repoPath, constants.GitRevParse, constants.GitAbbrevRef, constants.GitHEAD)
	if res.IsFailure() {
		return result.Fail[string](res.AppError())
	}

	return result.Ok(strings.TrimSpace(string(res.Data)))
}

// IsRepoClean reports whether the working tree at repoPath has no uncommitted changes.
func IsRepoClean(repoPath string) result.BoolResult {
	res := RunGitInDir(repoPath, "status", "--porcelain")
	if res.IsFailure() {
		return result.Fail[bool](res.AppError())
	}

	isClean := len(strings.TrimSpace(string(res.Data))) == 0

	return result.Ok(isClean)
}
