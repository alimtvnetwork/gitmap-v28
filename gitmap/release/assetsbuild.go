// Package release — assetsbuild.go contains low-level build helpers
// for cross-compilation targets.
package release

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// buildSingleTarget compiles one GOOS/GOARCH combination.
func buildSingleTarget(binName, version string, target BuildTarget, pkgDir, stagingDir string) CrossCompileResult {
	outputName := formatOutputName(binName, version, target)
	outputPath := filepath.Join(stagingDir, outputName)

	ldflags := fmt.Sprintf("-s -w -X main.version=%s", version)

	cmd := exec.Command("go", "build",
		"-ldflags", ldflags,
		"-o", outputPath,
		"./"+pkgDir,
	)

	cmd.Env = buildEnv(target)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return CrossCompileResult{
			Target:  target,
			Output:  outputPath,
			Success: false,
			Error:   strings.TrimSpace(string(out)),
		}
	}

	return CrossCompileResult{
		Target:  target,
		Output:  outputPath,
		Success: true,
	}
}

// formatOutputName creates the binary filename with platform suffix.
func formatOutputName(binName, version string, target BuildTarget) string {
	name := fmt.Sprintf("%s_%s_%s_%s", binName, version, target.GOOS, target.GOARCH)
	if target.GOOS == "windows" {
		name += ".exe"
	}

	return name
}

// buildEnv returns the os.Environ() with CGO_ENABLED=0, GOOS, GOARCH set.
func buildEnv(target BuildTarget) []string {
	env := os.Environ()

	env = setEnv(env, "CGO_ENABLED", "0")
	env = setEnv(env, "GOOS", target.GOOS)
	env = setEnv(env, "GOARCH", target.GOARCH)

	return env
}

// setEnv sets or replaces an environment variable in a slice.
func setEnv(env []string, key, value string) []string {
	prefix := key + "="

	for i, e := range env {
		if strings.HasPrefix(e, prefix) {
			env[i] = prefix + value

			return env
		}
	}

	return append(env, prefix+value)
}

// fileExists checks if a file exists and is not a directory.
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !info.IsDir()
}
