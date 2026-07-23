package clonenext

// Local-state reader for `gitmap cn` update detection.
//
// Reads the origin remote URL and HEAD commit SHA from a repo path
// without shelling out to git. Both are derived from .git plumbing
// files so the reader works on cloned repos that don't have git
// installed at exec time (CI, sandboxes).

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LocalRepoState is the snapshot used to drive the "is an update
// available?" decision for one repo.
type LocalRepoState struct {
	OriginURL string // contents of .git/config [remote "origin"] url
	HeadSHA   string // resolved 40-char commit hash, or "" when unborn
}

// ErrNotAGitRepo is returned when path has no .git entry.
var ErrNotAGitRepo = errors.New("clonenext: path is not a git repository")

// ReadLocalRepoState parses .git/config for the origin URL and resolves
// HEAD to a commit SHA. Both fields may be empty individually (no
// origin set; unborn HEAD) without producing an error — only a missing
// .git entry is fatal.
func ReadLocalRepoState(repoPath string) (LocalRepoState, error) {
	gitDir, err := resolveGitDir(repoPath)
	if err != nil {
		return LocalRepoState{}, err
	}

	state := LocalRepoState{}
	state.OriginURL = readOriginURL(gitDir)
	state.HeadSHA = readHeadSHA(gitDir)

	return state, nil
}

// resolveGitDir returns the actual git directory for repoPath, handling
// both regular repos (.git is a directory) and worktrees (.git is a
// file containing "gitdir: <path>").
func resolveGitDir(repoPath string) (string, error) {
	dotGit := filepath.Join(repoPath, ".git")
	info, err := os.Stat(dotGit)
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrNotAGitRepo, repoPath)
	}
	if info.IsDir() {
		return dotGit, nil
	}

	return resolveWorktreeGitDir(dotGit)
}

// resolveWorktreeGitDir reads the "gitdir: <path>" pointer used by
// worktrees and submodules.
func resolveWorktreeGitDir(dotGitFile string) (string, error) {
	body, err := os.ReadFile(dotGitFile)
	if err != nil {
		return "", fmt.Errorf("read %s: %w", dotGitFile, err)
	}
	line := strings.TrimSpace(string(body))
	if !strings.HasPrefix(line, "gitdir:") {
		return "", fmt.Errorf("malformed .git pointer in %s", dotGitFile)
	}

	return strings.TrimSpace(strings.TrimPrefix(line, "gitdir:")), nil
}

// readOriginURL parses .git/config and returns the url under
// [remote "origin"]. Returns "" when not present (best-effort).
func readOriginURL(gitDir string) string {
	body, err := os.ReadFile(filepath.Join(gitDir, "config"))
	if err != nil {
		return ""
	}

	return extractOriginURL(string(body))
}

// extractOriginURL is the pure-string parser for .git/config; split out
// to keep readOriginURL flat and testable.
func extractOriginURL(config string) string {
	inOriginSection := false
	for _, raw := range strings.Split(config, "\n") {
		line := strings.TrimSpace(raw)
		if strings.HasPrefix(line, "[") {
			inOriginSection = line == `[remote "origin"]`

			continue
		}
		if !inOriginSection {
			continue
		}
		if strings.HasPrefix(line, "url") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1])
			}
		}
	}

	return ""
}

// readHeadSHA reads .git/HEAD and resolves a symbolic ref one level
// down to its commit SHA. Returns "" for unborn HEAD or unreadable
// files (callers treat empty SHA as "unknown", not an error).
func readHeadSHA(gitDir string) string {
	body, err := os.ReadFile(filepath.Join(gitDir, "HEAD"))
	if err != nil {
		return ""
	}
	head := strings.TrimSpace(string(body))
	if !strings.HasPrefix(head, "ref:") {
		return head
	}
	refPath := strings.TrimSpace(strings.TrimPrefix(head, "ref:"))
	refBytes, err := os.ReadFile(filepath.Join(gitDir, refPath))
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(refBytes))
}
