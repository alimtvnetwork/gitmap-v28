//go:build !windows

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// setEnvPersistent sets an environment variable in the shell profile.
func setEnvPersistent(name, value string, _ bool, shell string) error {
	profilePath := resolveShellProfile(shell)
	exportLine := fmt.Sprintf(constants.EnvExportFmt, name, value)
	return appendToProfile(profilePath, name, exportLine)
}

// deleteEnvPersistent removes a variable from the shell profile.
func deleteEnvPersistent(name string, _ bool, shell string) error {
	profilePath := resolveShellProfile(shell)
	return removeFromProfile(profilePath, name)
}

// addPathPersistent adds a directory to PATH in the shell profile.
func addPathPersistent(dir string, _ bool, shell string) error {
	profilePath := resolveShellProfile(shell)
	exportLine := fmt.Sprintf(constants.EnvPathExportFmt, dir)
	marker := constants.EnvManagedComment + " path:" + dir

	return appendLineToProfile(profilePath, exportLine, marker)
}

// removePathPersistent removes a directory from PATH in the shell profile.
func removePathPersistent(dir string, _ bool, shell string) error {
	profilePath := resolveShellProfile(shell)
	marker := constants.EnvManagedComment + " path:" + dir

	return removeLineFromProfile(profilePath, marker)
}

// resolveShellProfile returns the profile path for the given or detected shell.
func resolveShellProfile(shell string) string {
	if shell != "" {
		return profilePathForShell(shell)
	}

	return detectShellProfile()
}

// profilePathForShell maps a shell name to its profile file path.
func profilePathForShell(shell string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not determine home directory: %v\n", err)

		return ""
	}

	if shell == constants.ShellZsh {
		return filepath.Join(home, constants.EnvProfileZshRC)
	}

	return filepath.Join(home, constants.EnvProfileBashRC)
}

// detectShellProfile returns the path to the active shell profile.
func detectShellProfile() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not determine home directory: %v\n", err)

		return ""
	}
	shell := os.Getenv("SHELL")

	if strings.Contains(shell, constants.ShellZsh) {
		return filepath.Join(home, constants.EnvProfileZshRC)
	}

	return filepath.Join(home, constants.EnvProfileBashRC)
}

// appendToProfile adds or updates an export line in the profile.
func appendToProfile(profilePath, name, exportLine string) error {
	marker := constants.EnvManagedComment + " " + name
	content := readProfileContent(profilePath)
	updatedContent := replaceOrAppendLine(content, marker, exportLine+" "+marker)

	return writeProfileContent(profilePath, updatedContent)
}

// removeFromProfile removes a managed variable line from the profile.
func removeFromProfile(profilePath, name string) error {
	marker := constants.EnvManagedComment + " " + name

	return removeLineFromProfile(profilePath, marker)
}

// appendLineToProfile adds a line with a marker to the profile.
func appendLineToProfile(profilePath, line, marker string) error {
	content := readProfileContent(profilePath)
	updatedContent := replaceOrAppendLine(content, marker, line+" "+marker)

	return writeProfileContent(profilePath, updatedContent)
}

// removeLineFromProfile removes a line matching a marker.
func removeLineFromProfile(profilePath, marker string) error {
	content := readProfileContent(profilePath)
	lines := strings.Split(content, "\n")
	filtered := make([]string, 0, len(lines))

	for _, line := range lines {
		if strings.Contains(line, marker) {
			continue
		}

		filtered = append(filtered, line)
	}

	return writeProfileContent(profilePath, strings.Join(filtered, "\n"))
}

// readProfileContent reads a shell profile file.
func readProfileContent(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}

	return string(data)
}

// writeProfileContent writes content to a shell profile file.
func writeProfileContent(path, content string) error {
	err := os.WriteFile(path, []byte(content), constants.FilePermission)
	if err != nil {
		return apperror.NewSimple(constants.ErrEnvProfileWrite, "E9000")
	}

	return nil
}

// replaceOrAppendLine replaces an existing marked line or appends new.
func replaceOrAppendLine(content, marker, newLine string) string {
	lines := strings.Split(content, "\n")

	for idx, line := range lines {
		if strings.Contains(line, marker) {
			lines[idx] = newLine

			return strings.Join(lines, "\n")
		}
	}

	return appendNewLine(content, newLine)
}

// appendNewLine appends a new line with trailing newline guaranteed.
func appendNewLine(content, newLine string) string {
	if content != "" && !strings.HasSuffix(content, "\n") {
		content += "\n"
	}

	return content + newLine + "\n"
}
