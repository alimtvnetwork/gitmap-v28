package osuser

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// CreateRootUser provisions a new user on Windows or Linux with administrative privileges.
func CreateRootUser(opts UserCreateOptions) *apperror.AppError {
	appErr := validateCreateOptions(opts)
	if appErr != nil {
		return appErr
	}
	return dispatchPlatformCreate(opts)
}

func validateCreateOptions(opts UserCreateOptions) *apperror.AppError {
	hasUsername := len(strings.TrimSpace(opts.Username)) > 0
	if !hasUsername {
		return apperror.NewValidationError("username is required for user creation")
	}
	return nil
}

func dispatchPlatformCreate(opts UserCreateOptions) *apperror.AppError {
	switch runtime.GOOS {
	case "windows":
		return createWindowsUser(opts)
	case "linux":
		return createLinuxUser(opts)
	default:
		return apperror.NewExecutionError(fmt.Sprintf("unsupported os for user creation: %s", runtime.GOOS))
	}
}

func createWindowsUser(opts UserCreateOptions) *apperror.AppError {
	if opts.IsDryRun {
		return nil
	}
	appErr := executeWindowsNetUser(opts.Username, opts.Password)
	if appErr != nil {
		return appErr
	}
	return configureWindowsAdminIfRequested(opts)
}

func executeWindowsNetUser(username, password string) *apperror.AppError {
	args := buildWindowsNetUserArgs(username, password)
	cmd := exec.Command("net", args...)
	out, err := defaultOSCommandRunner(cmd)
	if err != nil {
		return apperror.Wrap(err, "osuser.executeWindowsNetUser", map[string]any{"output": string(out)})
	}
	return nil
}

func buildWindowsNetUserArgs(username, password string) []string {
	hasPassword := len(password) > 0
	if hasPassword {
		return []string{"user", username, password, "/ADD"}
	}
	return []string{"user", username, "/ADD"}
}

func configureWindowsAdminIfRequested(opts UserCreateOptions) *apperror.AppError {
	if !opts.IsSudoer {
		return nil
	}
	cmd := exec.Command("net", "localgroup", "Administrators", opts.Username, "/ADD")
	out, err := defaultOSCommandRunner(cmd)
	if err != nil {
		return apperror.Wrap(err, "osuser.configureWindowsAdminIfRequested", map[string]any{"output": string(out)})
	}
	return nil
}

func createLinuxUser(opts UserCreateOptions) *apperror.AppError {
	appErr := executeLinuxUserAdd(opts)
	if appErr != nil {
		return appErr
	}
	return configureLinuxUserPostAdd(opts)
}

func executeLinuxUserAdd(opts UserCreateOptions) *apperror.AppError {
	if opts.IsDryRun {
		return nil
	}
	args := buildLinuxUserAddArgs(opts)
	cmd := exec.Command("useradd", args...)
	out, err := defaultOSCommandRunner(cmd)
	if err != nil {
		return apperror.Wrap(err, "osuser.executeLinuxUserAdd", map[string]any{"output": string(out)})
	}
	return nil
}

func buildLinuxUserAddArgs(opts UserCreateOptions) []string {
	args := []string{"-m"}
	hasHome := len(opts.HomeDir) > 0
	if hasHome {
		args = append(args, "-d", opts.HomeDir)
	}
	shell := resolveLinuxShell(opts.Shell)
	args = append(args, "-s", shell, opts.Username)
	return args
}

func resolveLinuxShell(shell string) string {
	hasShell := len(shell) > 0
	if hasShell {
		return shell
	}
	return selectAvailableShell()
}

func selectAvailableShell() string {
	candidates := []string{"/bin/zsh", "/usr/bin/zsh", "/bin/bash", "/bin/sh"}
	for _, cand := range candidates {
		if _, err := os.Stat(cand); err == nil {
			return cand
		}
	}
	return DefaultLinuxShell
}

func configureLinuxUserPostAdd(opts UserCreateOptions) *apperror.AppError {
	appErr := applyLinuxPassword(opts)
	if appErr != nil {
		return appErr
	}
	appErr = applyLinuxSudoers(opts)
	if appErr != nil {
		return appErr
	}
	return applyLinuxZshIfRequested(opts)
}

func applyLinuxPassword(opts UserCreateOptions) *apperror.AppError {
	hasPassword := len(opts.Password) > 0
	if !hasPassword || opts.IsDryRun {
		return nil
	}
	return pipePasswordToChpasswd(opts.Username, opts.Password)
}

func pipePasswordToChpasswd(username, password string) *apperror.AppError {
	cmd := exec.Command("chpasswd")
	cmd.Stdin = strings.NewReader(fmt.Sprintf("%s:%s\n", username, password))
	out, err := defaultOSCommandRunner(cmd)
	if err != nil {
		return apperror.Wrap(err, "osuser.pipePasswordToChpasswd", map[string]any{"output": string(out)})
	}
	return nil
}

func applyLinuxSudoers(opts UserCreateOptions) *apperror.AppError {
	if !opts.IsSudoer {
		return nil
	}
	appErr := ConfigureSudoersEntry(opts.Username, opts.IsDryRun)
	if appErr != nil {
		return appErr
	}
	return AddUserToSudoGroup(opts.Username, opts.IsDryRun)
}

func applyLinuxZshIfRequested(opts UserCreateOptions) *apperror.AppError {
	if !opts.IsConfigureZsh {
		return nil
	}
	return configureZshEnvironment(opts)
}

func configureZshEnvironment(opts UserCreateOptions) *apperror.AppError {
	homeDir := resolveUserHome(opts.Username, opts.HomeDir)
	appErr := createWorkspaceDirectories(homeDir, opts.IsDryRun)
	if appErr != nil {
		return appErr
	}
	return setupUserZshrc(homeDir, opts)
}

func resolveUserHome(username, homeDir string) string {
	hasHome := len(homeDir) > 0
	if hasHome {
		return homeDir
	}
	if runtime.GOOS == constants.OSWindows {
		return filepath.Join(os.Getenv("SystemDrive")+"\\Users", username)
	}
	return filepath.Join("/home", username)
}

func createWorkspaceDirectories(homeDir string, isDryRun bool) *apperror.AppError {
	if isDryRun {
		return nil
	}
	appErr := makeWorkspaceFolders(homeDir)
	if appErr != nil {
		return appErr
	}
	return enforceDirectoryMode(filepath.Join(homeDir, ".ssh"), 0700)
}

func makeWorkspaceFolders(homeDir string) *apperror.AppError {
	dirs := []string{"scripts", "gitlab", "github", ".ssh"}
	for _, dir := range dirs {
		path := filepath.Join(homeDir, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return apperror.Wrap(err, "osuser.makeWorkspaceFolders", map[string]any{"path": path})
		}
	}
	return nil
}

func enforceDirectoryMode(dirPath string, mode os.FileMode) *apperror.AppError {
	err := os.Chmod(dirPath, mode)
	if err != nil {
		return apperror.Wrap(err, "osuser.enforceDirectoryMode", map[string]any{"path": dirPath})
	}
	return nil
}

func setupUserZshrc(homeDir string, opts UserCreateOptions) *apperror.AppError {
	if opts.IsDryRun {
		return nil
	}
	theme := resolveZshTheme(opts.Theme)
	zshrcPath := filepath.Join(homeDir, ".zshrc")
	content := fmt.Sprintf("export ZSH=\"$HOME/.oh-my-zsh\"\nZSH_THEME=\"%s\"\nplugins=(git)\n", theme)
	err := os.WriteFile(zshrcPath, []byte(content), 0644)
	if err != nil {
		return apperror.Wrap(err, "osuser.setupUserZshrc", map[string]any{"path": zshrcPath})
	}
	return nil
}

func resolveZshTheme(theme string) string {
	hasTheme := len(theme) > 0
	if hasTheme {
		return theme
	}
	return DefaultZshTheme
}
