// Package release — assets.go orchestrates cross-compilation for Go projects.
package release

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/verbose"
)

// BuildTarget represents a single GOOS/GOARCH pair for cross-compilation.
type BuildTarget struct {
	GOOS   string `json:"goos"`
	GOARCH string `json:"goarch"`
}

// CrossCompileResult holds the outcome of a cross-compile step.
type CrossCompileResult struct {
	Target  BuildTarget
	Output  string
	Success bool
	Error   string
}

// DetectGoProject checks if the current directory contains a buildable Go project.
func DetectGoProject() bool {
	_, err := os.Stat(constants.IndicatorGoMod)

	return err == nil
}

// ReadModuleName reads the module name from go.mod.
func ReadModuleName() (string, error) {
	f, err := os.Open(constants.IndicatorGoMod)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				return parts[1], nil
			}
		}
	}

	return "", fmt.Errorf("no module directive found in go.mod")
}

// BinaryName extracts the short name from a Go module path.
// "github.com/alimtvnetwork/gitmap-v27/gitmap" → "gitmap"
func BinaryName(moduleName string) string {
	parts := strings.Split(moduleName, "/")

	return parts[len(parts)-1]
}

// FindMainPackages locates buildable main package directories.
// Checks root main.go first, then cmd/ subdirectories.
func FindMainPackages() []string {
	if fileExists(constants.GoMainFile) {
		return []string{"."}
	}

	cmdDir := constants.GoCmdDir
	names, err := readDirNames(cmdDir)
	if err != nil {
		return nil
	}

	var packages []string

	for _, name := range names {
		entryPath := filepath.Join(cmdDir, name)
		info, statErr := os.Stat(entryPath)
		if statErr != nil || !info.IsDir() {
			continue
		}

		mainPath := filepath.Join(entryPath, constants.GoMainFile)
		if fileExists(mainPath) {
			packages = append(packages, entryPath)
		}
	}

	return packages
}

// CrossCompile builds binaries for all targets and packages.
// Returns the list of successfully built binary paths.
func CrossCompile(version string, targets []BuildTarget, packages []string, stagingDir string) []CrossCompileResult {
	var results []CrossCompileResult

	binName := resolveBinName()

	for _, pkg := range packages {
		pkgSuffix := ""
		if pkg != "." {
			pkgSuffix = "-" + filepath.Base(pkg)
		}

		for _, t := range targets {
			if verbose.IsEnabled() {
				verbose.Get().Log("build: %s/%s → %s", t.GOOS, t.GOARCH,
					filepath.Join(stagingDir, formatOutputName(binName+pkgSuffix, version, t)))
			}

			result := buildSingleTarget(binName+pkgSuffix, version, t, pkg, stagingDir)
			results = append(results, result)

			if result.Success {
				if verbose.IsEnabled() {
					info, statErr := os.Stat(result.Output)
					if statErr == nil {
						verbose.Get().Log("build: %s/%s complete (%d bytes)", t.GOOS, t.GOARCH, info.Size())
					}
				}

				fmt.Printf(constants.MsgAssetBuilt, filepath.Base(result.Output), t.GOOS, t.GOARCH)
			} else {
				if verbose.IsEnabled() {
					verbose.Get().Log("build: %s/%s failed: %s", t.GOOS, t.GOARCH, result.Error)
				}

				fmt.Fprintf(os.Stderr, constants.ErrAssetBuildFailed, t.GOOS, t.GOARCH, result.Error)
			}
		}
	}

	return results
}

// resolveBinName reads go.mod for the binary name.
func resolveBinName() string {
	mod, err := ReadModuleName()
	if err != nil {
		return "app"
	}

	return BinaryName(mod)
}

// CollectSuccessfulBuilds filters cross-compile results to only successful outputs.
func CollectSuccessfulBuilds(results []CrossCompileResult) []string {
	var paths []string

	for _, r := range results {
		if r.Success {
			paths = append(paths, r.Output)
		}
	}

	return paths
}

// EnsureStagingDir creates the release-assets staging directory.
func EnsureStagingDir() (string, error) {
	dir := constants.AssetsStagingDir
	err := os.MkdirAll(dir, constants.DirPermission)
	if err != nil {
		return "", fmt.Errorf("create staging dir: %w", err)
	}

	if verbose.IsEnabled() {
		verbose.Get().Log("staging: created directory %s", dir)
	}

	return dir, nil
}

// CleanupStagingDir removes the staging directory after upload.
func CleanupStagingDir() {
	if verbose.IsEnabled() {
		verbose.Get().Log("staging: removing directory %s", constants.AssetsStagingDir)
	}

	if err := os.RemoveAll(constants.AssetsStagingDir); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not remove staging directory %s: %v\n", constants.AssetsStagingDir, err)
	}
}
