//go:build windows

package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// setEnvPersistent sets an environment variable on Windows via setx.
func setEnvPersistent(name, value string, system bool, _ string) {
	args := buildSetxArgs(name, value, system)
	runSetx(args)
}

// deleteEnvPersistent removes an environment variable on Windows.
func deleteEnvPersistent(name string, system bool, _ string) {
	args := buildSetxArgs(name, "", system)
	runSetx(args)
}

// buildSetxArgs builds setx command arguments.
func buildSetxArgs(name, value string, system bool) []string {
	args := []string{name, value}

	if system {
		args = append(args, "/M")
	}

	return args
}

// runSetx executes the setx command.
func runSetx(args []string) {
	cmd := exec.Command("setx", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrEnvProfileWrite, "system registry", err)
		os.Exit(1)
	}
}

// addPathPersistent adds a directory to PATH on Windows via setx.
func addPathPersistent(dir string, system bool, _ string) {
	currentPath := os.Getenv("PATH")
	newPath := currentPath + ";" + dir

	setEnvPersistent("PATH", newPath, system, "")
}

// removePathPersistent removes a directory from PATH on Windows.
func removePathPersistent(dir string, system bool, _ string) {
	currentPath := os.Getenv("PATH")
	parts := strings.Split(currentPath, ";")
	filtered := filterPathParts(parts, dir)
	newPath := strings.Join(filtered, ";")

	setEnvPersistent("PATH", newPath, system, "")
}

// filterPathParts removes matching entries from PATH parts.
func filterPathParts(parts []string, dir string) []string {
	filtered := make([]string, 0, len(parts))

	for _, part := range parts {
		if strings.EqualFold(strings.TrimSpace(part), dir) {
			continue
		}

		filtered = append(filtered, part)
	}

	return filtered
}
