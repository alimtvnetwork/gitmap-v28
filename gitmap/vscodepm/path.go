// Package vscodepm syncs gitmap repos into the alefragnani.project-manager
// VS Code extension's projects.json file.
//
// Path resolution always discovers the VS Code USER-DATA root first, then
// appends the relative tail. The full path is never hardcoded; if the user
// has a portable VS Code install or non-standard APPDATA the resolver picks
// it up automatically.
//
// See spec/01-vscode-project-manager-sync/README.md
package vscodepm

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// ErrUserDataMissing is returned when the VS Code user-data root cannot be
// located (env vars unset OR the directory does not exist).
var ErrUserDataMissing = errors.New("vscode user-data root not found")

// ErrExtensionMissing is returned when projects.json's parent directory
// (the project-manager extension storage dir) does not exist — usually
// because the extension is not installed.
var ErrExtensionMissing = errors.New("project-manager extension storage dir not found")

// UserDataRoot returns the active VS Code user-data directory for the
// current OS, or ErrUserDataMissing if it cannot be located on disk.
func UserDataRoot() (string, error) {
	candidate := userDataCandidate()
	if candidate == "" {
		return "", ErrUserDataMissing
	}

	if !dirExists(candidate) {
		return candidate, ErrUserDataMissing
	}

	return candidate, nil
}

// ProjectsJSONPath returns the absolute path to projects.json by joining
// the discovered user-data root with the extension-relative tail. Returns
// ErrUserDataMissing when the root is not present, or ErrExtensionMissing
// when the extension storage dir is absent.
func ProjectsJSONPath() (string, error) {
	root, err := UserDataRoot()
	if err != nil {
		return root, err
	}

	extDir := filepath.Join(root,
		constants.VSCodePMUserDir,
		constants.VSCodePMGlobalStorageDir,
		constants.VSCodePMExtensionDir)

	if !dirExists(extDir) {
		return filepath.Join(extDir, constants.VSCodePMProjectsFile), ErrExtensionMissing
	}

	return filepath.Join(extDir, constants.VSCodePMProjectsFile), nil
}

// userDataCandidate returns the OS-specific user-data root candidate path.
// Empty string when neither the primary nor the fallback env var is set.
func userDataCandidate() string {
	switch runtime.GOOS {
	case "windows":
		return windowsUserDataCandidate()
	case "darwin":
		return darwinUserDataCandidate()
	default:
		return linuxUserDataCandidate()
	}
}

// windowsUserDataCandidate uses %APPDATA%\Code, falling back to
// %USERPROFILE%\AppData\Roaming\Code.
func windowsUserDataCandidate() string {
	if appData := os.Getenv(constants.VSCodeEnvAppData); appData != "" {
		return filepath.Join(appData, constants.VSCodeUserDataRootDirName)
	}

	if userProfile := os.Getenv(constants.VSCodeEnvUserProfile); userProfile != "" {
		return filepath.Join(userProfile,
			filepath.FromSlash(constants.VSCodeUserProfileAppDataRel))
	}

	return ""
}

// darwinUserDataCandidate uses $HOME/Library/Application Support/Code.
func darwinUserDataCandidate() string {
	home := os.Getenv(constants.VSCodeEnvHome)
	if home == "" {
		return ""
	}

	return filepath.Join(home, filepath.FromSlash(constants.VSCodeUserDataMacRel))
}

// linuxUserDataCandidate uses $XDG_CONFIG_HOME/Code, falling back to
// $HOME/.config/Code.
func linuxUserDataCandidate() string {
	if xdg := os.Getenv(constants.VSCodeEnvXDGConfigHome); xdg != "" {
		return filepath.Join(xdg, constants.VSCodeUserDataRootDirName)
	}

	if home := os.Getenv(constants.VSCodeEnvHome); home != "" {
		return filepath.Join(home, filepath.FromSlash(constants.VSCodeUserDataLinuxFallback))
	}

	return ""
}

// dirExists reports whether path exists and is a directory.
func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}
