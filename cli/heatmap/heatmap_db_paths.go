// Package heatmap — heatmap_db_paths.go resolves relative/absolute DB paths and repo-secrets URLs.
package heatmap

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/secrets"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// resolveDbPathDetails returns (absoluteDbPath, relativeDbPath) formatted with forward slashes.
func resolveDbPathDetails(dbPath string) (string, string) {
	repoRoot := findRepoRoot()
	absPath, err := filepath.Abs(dbPath)
	if err != nil {
		absPath = dbPath
	}
	absRoot, err := filepath.Abs(repoRoot)
	if err != nil {
		absRoot = repoRoot
	}
	relPath, relErr := filepath.Rel(absRoot, absPath)
	if relErr == nil {
		return filepath.ToSlash(absPath), filepath.ToSlash(relPath)
	}
	return filepath.ToSlash(absPath), filepath.ToSlash(dbPath)
}

func findRepoRoot() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "."
	}
	dir := cwd
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return cwd
}

// ResolveRepoSecretsURL finds the remote URL for repo-secrets via Split DB or git config.
func ResolveRepoSecretsURL() string {
	if url := querySpecialRepoSecretsURL(); len(url) > 0 {
		return url
	}
	return queryGitSecretsRemoteURL()
}

func querySpecialRepoSecretsURL() string {
	sdb, err := store.OpenSpecialReposSplitDB()
	if err != nil {
		return ""
	}
	defer sdb.Close()

	rec, recErr := sdb.GetSpecialRepo("rs")
	if recErr == nil && rec != nil && len(strings.TrimSpace(rec.RemoteURL)) > 0 {
		return strings.TrimSpace(rec.RemoteURL)
	}
	return ""
}

func queryGitSecretsRemoteURL() string {
	secretsRoot := secrets.ResolveRepoSecretsRoot()
	if len(secretsRoot) == 0 {
		return ""
	}
	cmd := exec.Command("git", "-C", secretsRoot, "config", "--get", "remote.origin.url")
	out, err := cmd.Output()
	if err == nil && len(strings.TrimSpace(string(out))) > 0 {
		return strings.TrimSpace(string(out))
	}
	return ""
}
