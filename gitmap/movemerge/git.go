package movemerge

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// runGit runs `git -C dir <args...>` and returns stdout/stderr combined.
func runGit(dir string, args ...string) (string, error) {
	full := append([]string{constants.GitDirFlag, dir}, args...)
	cmd := exec.Command(constants.GitBin, full...)
	out, err := cmd.CombinedOutput()

	return strings.TrimSpace(string(out)), err
}

// GetOriginURL reads `git config --get remote.origin.url` for dir.
func GetOriginURL(dir string) (string, error) {
	return runGit(dir, constants.GitConfigCmd, constants.GitGetFlag, constants.GitRemoteOrigin)
}

// CloneURL clones url into dir, optionally pinning a branch.
func CloneURL(url, branch, dir string) error {
	args := []string{constants.GitClone}
	if branch != "" {
		args = append(args, constants.GitBranchFlag, branch)
	}
	args = append(args, url, dir)
	cmd := exec.Command(constants.GitBin, args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintln(os.Stderr, string(out))

		return fmt.Errorf("git clone %s: %w", url, err)
	}

	return nil
}

// PullFFOnly runs `git pull --ff-only` inside dir.
func PullFFOnly(dir string) error {
	out, err := runGit(dir, constants.GitPull, constants.GitFFOnlyFlag)
	if err != nil {
		return fmt.Errorf("git pull --ff-only in %s: %w (%s)", dir, err, out)
	}

	return nil
}

// AddCommitPush stages, commits, and pushes the working folder.
// Returns the commit SHA on success.
func AddCommitPush(dir, msg string, push bool) (string, error) {
	if _, err := runGit(dir, constants.GitAddCmd, constants.GitAddAllArg); err != nil {
		return "", fmt.Errorf("git add -A in %s: %w", dir, err)
	}
	if _, err := runGit(dir, constants.GitCommitCmd, constants.GitMessageArg, msg); err != nil {
		// Empty commits aren't an error here — nothing changed.
		return "", nil
	}
	sha, _ := runGit(dir, constants.GitRevParse, constants.GitHEAD)
	if !push {
		return sha, nil
	}
	if _, err := runGit(dir, constants.GitPush); err != nil {
		return sha, fmt.Errorf("git push in %s: %w", dir, err)
	}

	return sha, nil
}
