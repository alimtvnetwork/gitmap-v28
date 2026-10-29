package setup

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// PathResolver provides cross-platform discovery of IDE configuration files.
type PathResolver struct {
	BaseHomeDir    string
	BaseAppDataDir string
	TargetOS       string
}

// NewPathResolver creates a default platform-aware PathResolver.
func NewPathResolver() *PathResolver {
	return &PathResolver{}
}

func (r *PathResolver) resolveOS() string {
	if r.TargetOS != "" {
		return r.TargetOS
	}

	return runtime.GOOS
}

func (r *PathResolver) resolveHomeDir() (string, *apperror.AppError) {
	if r.BaseHomeDir != "" {
		return r.BaseHomeDir, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", apperror.WrapExecution(err, "resolveHomeDir failed")
	}

	return home, nil
}

func (r *PathResolver) resolveAppDataDir() (string, *apperror.AppError) {
	if r.BaseAppDataDir != "" {
		return r.BaseAppDataDir, nil
	}

	if appData := os.Getenv("APPDATA"); appData != "" {
		return appData, nil
	}

	return r.fallbackAppDataDir()
}

func (r *PathResolver) fallbackAppDataDir() (string, *apperror.AppError) {
	home, err := r.resolveHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, "AppData", "Roaming"), nil
}

// GetClineStorageDir locates the Cline global storage directory for the target OS.
func (r *PathResolver) GetClineStorageDir() (string, *apperror.AppError) {
	if r.resolveOS() == "windows" {
		return r.resolveWindowsClineStorage()
	}

	return r.resolveUnixClineStorage()
}

func (r *PathResolver) resolveWindowsClineStorage() (string, *apperror.AppError) {
	appData, err := r.resolveAppDataDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(appData, "Code", "User", "globalStorage", "saoudrizwan.claude-dev"), nil
}

func (r *PathResolver) resolveUnixClineStorage() (string, *apperror.AppError) {
	home, err := r.resolveHomeDir()
	if err != nil {
		return "", err
	}

	if r.resolveOS() == "darwin" {
		return filepath.Join(home, "Library", "Application Support", "Code", "User", "globalStorage", "saoudrizwan.claude-dev"), nil
	}

	return filepath.Join(home, ".config", "Code", "User", "globalStorage", "saoudrizwan.claude-dev"), nil
}

// GetClineSettingsPath returns the absolute path to cline_mcp_settings.json.
func (r *PathResolver) GetClineSettingsPath() (string, *apperror.AppError) {
	dir, err := r.GetClineStorageDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "settings", "cline_mcp_settings.json"), nil
}

// GetClineProviderSettingsPath returns the path to cline_provider_settings.json.
func (r *PathResolver) GetClineProviderSettingsPath() (string, *apperror.AppError) {
	dir, err := r.GetClineStorageDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "settings", "cline_provider_settings.json"), nil
}

// GetContinueConfigPath locates ~/.continue/config.json.
func (r *PathResolver) GetContinueConfigPath() (string, *apperror.AppError) {
	home, err := r.resolveHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".continue", "config.json"), nil
}

// GetZedSettingsPath locates Zed settings.json across OS platforms.
func (r *PathResolver) GetZedSettingsPath() (string, *apperror.AppError) {
	if r.resolveOS() == "windows" {
		return r.resolveWindowsZedPath()
	}

	return r.resolveUnixZedPath()
}

func (r *PathResolver) resolveWindowsZedPath() (string, *apperror.AppError) {
	appData, err := r.resolveAppDataDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(appData, "Zed", "settings.json"), nil
}

func (r *PathResolver) resolveUnixZedPath() (string, *apperror.AppError) {
	home, err := r.resolveHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".config", "zed", "settings.json"), nil
}

// GetAllTargetPaths returns resolved paths for all supported IDE targets.
func (r *PathResolver) GetAllTargetPaths() (TargetConfigPaths, *apperror.AppError) {
	paths := TargetConfigPaths{}
	if err := r.populateClinePaths(&paths); err != nil {
		return paths, err
	}

	if err := r.populateEditorPaths(&paths); err != nil {
		return paths, err
	}

	return paths, nil
}

func (r *PathResolver) populateClinePaths(paths *TargetConfigPaths) *apperror.AppError {
	settings, err := r.GetClineSettingsPath()
	if err != nil {
		return err
	}

	paths.ClineSettingsPath = settings

	return r.populateClineStorage(paths)
}

func (r *PathResolver) populateClineStorage(paths *TargetConfigPaths) *apperror.AppError {
	storage, err := r.GetClineStorageDir()
	if err != nil {
		return err
	}

	paths.ClineStorageDir = storage

	return nil
}

func (r *PathResolver) populateEditorPaths(paths *TargetConfigPaths) *apperror.AppError {
	contPath, err := r.GetContinueConfigPath()
	if err != nil {
		return err
	}

	paths.ContinuePath = contPath

	return r.populateZedPath(paths)
}

func (r *PathResolver) populateZedPath(paths *TargetConfigPaths) *apperror.AppError {
	zedPath, err := r.GetZedSettingsPath()
	if err != nil {
		return err
	}

	paths.ZedSettingsPath = zedPath

	return nil
}
