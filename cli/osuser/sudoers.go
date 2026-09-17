package osuser

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// ConfigureSudoersEntry writes a passwordless sudoers configuration file and validates it.
func ConfigureSudoersEntry(username string, isDryRun bool) *apperror.AppError {
	if isDryRun {
		return nil
	}
	filePath := filepath.Join(SudoersDirPath, username)
	content := fmt.Sprintf("%s ALL=(ALL) NOPASSWD:ALL\n", username)
	appErr := writeSudoersFile(filePath, content)
	if appErr != nil {
		return appErr
	}
	return validateAndEnforceSudoers(filePath)
}

func writeSudoersFile(filePath, content string) *apperror.AppError {
	err := os.WriteFile(filePath, []byte(content), 0440)
	if err != nil {
		return apperror.Wrap(err, "osuser.writeSudoersFile", map[string]any{"path": filePath})
	}
	return nil
}

func validateAndEnforceSudoers(filePath string) *apperror.AppError {
	appErr := validateVisudoSyntax(filePath)
	if appErr != nil {
		_ = os.Remove(filePath)
		return appErr
	}
	return nil
}

func validateVisudoSyntax(filePath string) *apperror.AppError {
	visudoPath, err := exec.LookPath("visudo")
	if err != nil {
		return nil
	}
	cmd := exec.Command(visudoPath, "-cf", filePath)
	out, err := defaultOSCommandRunner(cmd)
	if err != nil {
		return apperror.Wrap(err, "osuser.validateVisudoSyntax", map[string]any{"output": string(out)})
	}
	return nil
}

// AddUserToSudoGroup adds a Linux user to the sudo group.
func AddUserToSudoGroup(username string, isDryRun bool) *apperror.AppError {
	if isDryRun {
		return nil
	}
	cmd := exec.Command("usermod", "-aG", "sudo", username)
	out, err := defaultOSCommandRunner(cmd)
	if err != nil {
		return apperror.Wrap(err, "osuser.AddUserToSudoGroup", map[string]any{"output": string(out)})
	}
	return nil
}

// CleanSudoers purges user entries from sudoers drop-in and main sudoers file.
func CleanSudoers(username string, isDryRun bool) *apperror.AppError {
	if isDryRun {
		return nil
	}
	removeSudoersDropIn(username)
	return removeUserFromMainSudoers(username)
}

func removeSudoersDropIn(username string) {
	filePath := filepath.Join(SudoersDirPath, username)
	if _, err := os.Stat(filePath); err == nil {
		_ = os.Remove(filePath)
	}
}

func removeUserFromMainSudoers(username string) *apperror.AppError {
	if _, err := os.Stat(SudoersFilePath); err != nil {
		return nil
	}
	lines, appErr := filterSudoersLines(SudoersFilePath, username)
	if appErr != nil {
		return appErr
	}
	return writeFilteredSudoers(SudoersFilePath, lines)
}

func filterSudoersLines(filePath, username string) ([]string, *apperror.AppError) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, apperror.Wrap(err, "osuser.filterSudoersLines", map[string]any{"path": filePath})
	}
	return filterRetainedSudoersLines(string(data), username), nil
}

func filterRetainedSudoersLines(content, username string) []string {
	prefix := username + " "
	rawLines := strings.Split(content, "\n")
	retained := make([]string, 0, len(rawLines))
	for _, line := range rawLines {
		trimmed := strings.TrimSpace(line)
		isMatch := strings.HasPrefix(trimmed, prefix)
		if !isMatch {
			retained = append(retained, line)
		}
	}
	return retained
}

func writeFilteredSudoers(filePath string, lines []string) *apperror.AppError {
	output := strings.Join(lines, "\n")
	err := os.WriteFile(filePath, []byte(output), 0440)
	if err != nil {
		return apperror.Wrap(err, "osuser.writeFilteredSudoers", map[string]any{"path": filePath})
	}
	return nil
}
