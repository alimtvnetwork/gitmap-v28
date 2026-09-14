package osuser

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// InstallSSHKey appends an SSH public key to authorized_keys with deduplication and secure permissions.
func InstallSSHKey(opts SSHKeyOptions) *apperror.AppError {
	appErr := validateSSHKeyOptions(opts)
	if appErr != nil {
		return appErr
	}
	keyText, appErr := resolveKeyText(opts)
	if appErr != nil {
		return appErr
	}
	return appendKeyToHome(opts, keyText)
}

func validateSSHKeyOptions(opts SSHKeyOptions) *apperror.AppError {
	hasKey := len(strings.TrimSpace(opts.PublicKey)) > 0
	hasFile := len(strings.TrimSpace(opts.KeyFilePath)) > 0
	if !hasKey && !hasFile {
		return apperror.NewValidationError("either public key text or key file path is required")
	}
	return nil
}

func resolveKeyText(opts SSHKeyOptions) (string, *apperror.AppError) {
	hasText := len(strings.TrimSpace(opts.PublicKey)) > 0
	if hasText {
		return strings.TrimSpace(opts.PublicKey), nil
	}
	return readKeyFromFile(opts.KeyFilePath)
}

func readKeyFromFile(filePath string) (string, *apperror.AppError) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", apperror.Wrap(err, "osuser.readKeyFromFile", map[string]any{"path": filePath})
	}
	return strings.TrimSpace(string(data)), nil
}

func appendKeyToHome(opts SSHKeyOptions, keyText string) *apperror.AppError {
	homeDir, appErr := determineHomeDir(opts)
	if appErr != nil {
		return appErr
	}
	if opts.IsDryRun {
		return nil
	}
	return syncAuthorizedKeys(homeDir, keyText)
}

func determineHomeDir(opts SSHKeyOptions) (string, *apperror.AppError) {
	hasHome := len(opts.TargetHome) > 0
	if hasHome {
		return opts.TargetHome, nil
	}
	return resolveUserOrCurrentHome(opts.Username)
}

func resolveUserOrCurrentHome(username string) (string, *apperror.AppError) {
	hasUser := len(username) > 0
	if hasUser && runtime.GOOS == "linux" {
		return filepath.Join("/home", username), nil
	}
	u, err := user.Current()
	if err != nil {
		return "", apperror.Wrap(err, "osuser.resolveUserOrCurrentHome", nil)
	}
	return u.HomeDir, nil
}

func syncAuthorizedKeys(homeDir, keyText string) *apperror.AppError {
	sshDir := filepath.Join(homeDir, ".ssh")
	appErr := ensureSSHDirectory(sshDir)
	if appErr != nil {
		return appErr
	}
	authKeysPath := filepath.Join(sshDir, "authorized_keys")
	return writeDedupAuthorizedKeys(authKeysPath, keyText)
}

func ensureSSHDirectory(sshDir string) *apperror.AppError {
	err := os.MkdirAll(sshDir, 0700)
	if err != nil {
		return apperror.Wrap(err, "osuser.ensureSSHDirectory", map[string]any{"path": sshDir})
	}
	if chmodErr := os.Chmod(sshDir, 0700); chmodErr != nil {
		return apperror.Wrap(chmodErr, "osuser.ensureSSHDirectory.chmod", map[string]any{"path": sshDir})
	}
	return nil
}

func writeDedupAuthorizedKeys(authKeysPath, keyText string) *apperror.AppError {
	existingLines := loadExistingKeyLines(authKeysPath)
	hasKey := hasExistingKey(existingLines, keyText)
	if hasKey {
		return enforceKeyFilePermissions(authKeysPath)
	}
	existingLines = append(existingLines, keyText)
	return persistKeyFile(authKeysPath, existingLines)
}

func loadExistingKeyLines(authKeysPath string) []string {
	data, err := os.ReadFile(authKeysPath)
	if err != nil {
		return make([]string, 0)
	}
	return parseTrimmedLines(string(data))
}

func parseTrimmedLines(content string) []string {
	lines := strings.Split(content, "\n")
	resultLines := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		hasContent := len(trimmed) > 0
		if hasContent {
			resultLines = append(resultLines, trimmed)
		}
	}
	return resultLines
}

func hasExistingKey(lines []string, targetKey string) bool {
	trimmedTarget := strings.TrimSpace(targetKey)
	for _, line := range lines {
		isMatch := strings.TrimSpace(line) == trimmedTarget
		if isMatch {
			return true
		}
	}
	return false
}

func persistKeyFile(authKeysPath string, lines []string) *apperror.AppError {
	content := strings.Join(lines, "\n") + "\n"
	err := os.WriteFile(authKeysPath, []byte(content), 0600)
	if err != nil {
		return apperror.Wrap(err, "osuser.persistKeyFile", map[string]any{"path": authKeysPath})
	}
	return enforceKeyFilePermissions(authKeysPath)
}

func enforceKeyFilePermissions(authKeysPath string) *apperror.AppError {
	err := os.Chmod(authKeysPath, 0600)
	if err != nil {
		return apperror.Wrap(err, "osuser.enforceKeyFilePermissions", map[string]any{"path": authKeysPath})
	}
	return nil
}
