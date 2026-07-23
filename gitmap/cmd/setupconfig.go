package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/setup"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// mustLoadSetupConfig loads the resolved setup config or exits.
func mustLoadSetupConfig(configPath string) setup.GitSetupConfig {
	cfg, err := setup.LoadConfig(configPath)
	if err == nil {
		return cfg
	}

	fmt.Fprintf(os.Stderr, constants.ErrSetupLoadFailed, configPath, err)
	os.Exit(1)

	return setup.GitSetupConfig{}
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
