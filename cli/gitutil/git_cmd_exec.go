package gitutil

import (
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
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
func RunGitInDir(repoPath string, args ...string) ([]byte, error) {
	cmd := createGitCmd(repoPath, args...)

	return cmd.Output()
}

// RunGitCombinedInDir executes a git command with directory arguments and returns combined stdout and stderr.
func RunGitCombinedInDir(repoPath string, args ...string) ([]byte, error) {
	cmd := createGitCmd(repoPath, args...)

	return cmd.CombinedOutput()
}

// GetRepoConfig retrieves a git configuration value for key in repoPath.
func GetRepoConfig(repoPath, key string) (string, error) {
	out, err := RunGitInDir(repoPath, constants.GitConfigCmd, constants.GitGetFlag, key)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

// GetRemoteOriginURL retrieves the origin remote URL or an empty string on error.
func GetRemoteOriginURL(repoPath string) string {
	val, err := GetRepoConfig(repoPath, constants.GitRemoteOrigin)
	if err != nil {
		return ""
	}

	return val
}

// ReadRemoteOriginURL retrieves the configured remote.origin.url for repoPath.
func ReadRemoteOriginURL(repoPath string) (string, error) {
	return GetRepoConfig(repoPath, constants.GitRemoteOrigin)
}

// SetRemoteOriginURL updates the remote origin URL for repoPath.
func SetRemoteOriginURL(repoPath, url string) error {
	_, err := RunGitCombinedInDir(repoPath, "remote", "set-url", constants.GitOrigin, url)

	return err
}

// GetRepoCurrentBranch returns the active branch name of the repository at repoPath.
func GetRepoCurrentBranch(repoPath string) (string, error) {
	out, err := RunGitInDir(repoPath, constants.GitRevParse, constants.GitAbbrevRef, constants.GitHEAD)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

// IsRepoClean reports whether the working tree at repoPath has no uncommitted changes.
func IsRepoClean(repoPath string) (bool, error) {
	out, err := RunGitInDir(repoPath, "status", "--porcelain")
	if err != nil {
		return false, err
	}

	isClean := len(strings.TrimSpace(string(out))) == 0

	return isClean, nil
}
