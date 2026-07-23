package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// PrintBinaryLocations prints the Active / Deployed / Config binary triplet
// to stdout. Called from bare `gitmap` (no args) and from the post-update
// readout. The output is suppressed when --no-banner is in os.Args or when
// the GITMAP_QUIET env var is set to "1".
//
// Definitions (see spec/01-app/89-deploy-layout-and-binary-readout.md):
//
//   - Active   = os.Executable() after filepath.EvalSymlinks. The file the
//     OS actually loaded for this process.
//   - Deployed = <powershell.json.deployPath>/gitmap-cli/<binaryName> if it
//     exists on disk; "(not found)" otherwise.
//   - Config   = literal path the config declares, whether or not the file
//     exists. Represents config intent.
func PrintBinaryLocations() {
	if isBannerSuppressed() {
		return
	}

	active := resolveActiveBinaryPath()
	deployed, configPath := resolveDeployedAndConfigPaths()

	fmt.Printf(constants.BinaryReadoutActive, displayPath(active))
	fmt.Printf(constants.BinaryReadoutDeployed, displayPath(deployed))
	fmt.Printf(constants.BinaryReadoutConfig, displayPath(configPath))
	fmt.Println()
}

// isBannerSuppressed reports whether --no-banner or GITMAP_QUIET=1 is set.
func isBannerSuppressed() bool {
	if os.Getenv(constants.EnvGitMapQuiet) == constants.EnvGitMapQuietTrue {
		return true
	}
	for _, arg := range os.Args[1:] {
		if arg == constants.FlagNoBanner {
			return true
		}
	}

	return false
}

// resolveActiveBinaryPath returns the symlink-resolved path of the running binary.
func resolveActiveBinaryPath() string {
	exe, err := os.Executable()
	if err != nil {
		return ""
	}
	resolved, evalErr := filepath.EvalSymlinks(exe)
	if evalErr != nil {
		return filepath.Clean(exe)
	}

	return filepath.Clean(resolved)
}

// resolveDeployedAndConfigPaths reads powershell.json from the source repo
// and returns (deployedPathIfExists, configIntentPath). The deployed path is
// empty when the file is missing on disk; the config path is always returned
// when the JSON declares a deployPath.
func resolveDeployedAndConfigPaths() (string, string) {
	if len(constants.RepoPath) == 0 {
		return "", ""
	}
	configFile := filepath.Join(constants.RepoPath, constants.GitMapSubdir, constants.PowershellConfigFile)
	data, err := os.ReadFile(configFile)
	if err != nil {
		return "", ""
	}

	deployRoot := extractJSONString(data, constants.JSONKeyDeployPath)
	if len(deployRoot) == 0 {
		return "", ""
	}
	binaryName := extractJSONString(data, constants.JSONKeyBinaryName)
	if len(binaryName) == 0 {
		binaryName = constants.DoctorDefaultBinary
	}

	configPath := filepath.Join(deployRoot, constants.GitMapCliSubdir, binaryName)
	deployed := configPath
	if _, statErr := os.Stat(deployed); statErr != nil {
		deployed = ""
	}

	return deployed, configPath
}

// displayPath returns the path or a "(not found)" placeholder when empty.
func displayPath(path string) string {
	if len(path) == 0 {
		return constants.BinaryReadoutMissing
	}

	return path
}
