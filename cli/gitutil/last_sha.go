// Package gitutil — last_sha.go resolves short commit SHA of HEAD.
package gitutil

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GetLastCommitSHA returns the 7-character short commit hash for HEAD in repoPath.
// Returns "-" immediately if repoPath is empty, unresolvable, or not a git directory.
func GetLastCommitSHA(repoPath string) string {
	cleanPath := strings.TrimSpace(repoPath)
	if cleanPath == "" {
		return "-"
	}

	fi, err := os.Stat(cleanPath)
	isDir := err == nil && fi.IsDir()
	if !isDir {
		return "-"
	}

	fastSHA := resolveShortSHAFromDisk(cleanPath)
	if fastSHA != "" {
		return fastSHA
	}

	return resolveShortSHAViaGit(cleanPath)
}

func resolveShortSHAFromDisk(cleanPath string) string {
	headBytes, readErr := os.ReadFile(filepath.Join(cleanPath, ".git", "HEAD"))
	if readErr != nil {
		return ""
	}

	headStr := strings.TrimSpace(string(headBytes))
	if len(headStr) >= 7 && !strings.HasPrefix(headStr, "ref:") {
		return headStr[:7]
	}
	if strings.HasPrefix(headStr, "ref: ") {
		return resolveBranchRefSHA(cleanPath, strings.TrimPrefix(headStr, "ref: "))
	}

	return ""
}

func resolveBranchRefSHA(cleanPath, refRel string) string {
	refPath := filepath.Join(cleanPath, ".git", filepath.FromSlash(strings.TrimSpace(refRel)))
	refBytes, refErr := os.ReadFile(refPath)
	if refErr != nil {
		return ""
	}

	refStr := strings.TrimSpace(string(refBytes))
	if len(refStr) >= 7 {
		return refStr[:7]
	}

	return ""
}

func resolveShortSHAViaGit(cleanPath string) string {
	cmd := exec.Command("git", "-C", cleanPath, "rev-parse", "--short=7", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "-"
	}

	return strings.TrimSpace(string(out))
}
