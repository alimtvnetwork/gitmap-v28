package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/setup"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

const fallbackGitSetupJSON = `{
  "diffTool": {
    "name": "vscode",
    "cmd": "code --wait --diff $LOCAL $REMOTE",
    "trustExitCode": true
  },
  "mergeTool": {
    "name": "vscode",
    "cmd": "code --wait --merge $REMOTE $LOCAL $BASE $MERGED",
    "trustExitCode": true
  },
  "aliases": {
    "co": "checkout",
    "br": "branch",
    "ci": "commit",
    "st": "status",
    "unstage": "reset HEAD --",
    "last": "log -1 HEAD",
    "lg": "log --oneline --graph --decorate --all",
    "df": "diff --stat",
    "pushf": "push --force-with-lease",
    "amend": "commit --amend --no-edit"
  },
  "credentialHelper": "manager",
  "core": {
    "autocrlf": "true",
    "longpaths": "true",
    "editor": "code --wait",
    "defaultBranch": "main",
    "safecrlf": "warn"
  }
}`

// mustLoadSetupConfig loads the resolved setup config or uses a fallback.
func mustLoadSetupConfig(configPath string) setup.GitSetupConfig {
	cfg, err := setup.LoadConfig(configPath)
	if err == nil {
		return cfg
	}

	fmt.Fprintf(os.Stderr, "  WARN  could not load %s, using embedded fallback config\n", filepath.Base(configPath))
	var fallbackCfg setup.GitSetupConfig
	_ = json.Unmarshal([]byte(fallbackGitSetupJSON), &fallbackCfg)
	return fallbackCfg
}

// resolveSetupConfigPath prefers the bundled config unless overridden.
func resolveSetupConfigPath(configPath string, hasConfig bool) string {
	if hasConfig {
		return configPath
	}

	return resolveDefaultSetupConfigPath(configPath, store.BinaryDataDir(), constants.RepoPath)
}

// resolveDefaultSetupConfigPath picks the best default setup config path.
func resolveDefaultSetupConfigPath(configPath, binaryDataDir, repoPath string) string {
	name := filepath.Base(configPath)
	repoConfigPath := resolveRepoSetupConfigPath(repoPath, name)

	return firstExistingPath(
		filepath.Join(binaryDataDir, name),
		repoConfigPath,
		filepath.Join(constants.GitMapSubdir, constants.DBDir, name),
		configPath,
	)
}

// resolveRepoSetupConfigPath returns the source-repo setup config path.
func resolveRepoSetupConfigPath(repoPath, name string) string {
	if len(repoPath) == 0 {
		return ""
	}

	return filepath.Join(repoPath, constants.GitMapSubdir, constants.DBDir, name)
}

// firstExistingPath returns the first existing path or the first candidate.
func firstExistingPath(paths ...string) string {
	for _, path := range paths {
		if len(path) == 0 {
			continue
		}

		_, err := os.Stat(path)
		if err == nil || !errors.Is(err, os.ErrNotExist) {
			return path
		}
	}

	return paths[0]
}
