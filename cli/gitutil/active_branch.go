// Package gitutil — active_branch.go resolves active branch name or detached state.
package gitutil

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GetActiveBranch returns the current active branch name or "(detached)" if in detached HEAD.
func GetActiveBranch(repoPath string) string {
	cleanPath := strings.TrimSpace(repoPath)
	if cleanPath == "" {
		return "-"
	}

	branch := resolveActiveBranchFromHEAD(cleanPath)
	if branch != "" {
		return branch
	}

	return resolveActiveBranchViaGit(cleanPath)
}

func resolveActiveBranchFromHEAD(cleanPath string) string {
	headBytes, err := os.ReadFile(filepath.Join(cleanPath, ".git", "HEAD"))
	if err != nil {
		return ""
	}

	content := strings.TrimSpace(string(headBytes))
	if strings.HasPrefix(content, "ref: refs/heads/") {
		return strings.TrimPrefix(content, "ref: refs/heads/")
	}
	if len(content) == 40 {
		return "(detached)"
	}

	return ""
}

func resolveActiveBranchViaGit(cleanPath string) string {
	cmd := exec.Command("git", "-C", cleanPath, "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "-"
	}

	branch := strings.TrimSpace(string(out))
	if branch == "HEAD" {
		return "(detached)"
	}

	return branch
}
