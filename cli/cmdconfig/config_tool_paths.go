package cmdconfig

import (
	"os"
	"path/filepath"
	"runtime"

	cfgconst "github.com/alimtvnetwork/gitmap-v28/cli/cmdconfig/constants"
	cfgenum "github.com/alimtvnetwork/gitmap-v28/cli/cmdconfig/enums"
)

// resolveWindowsConfigPath resolves path within APPDATA on Windows.
func resolveWindowsConfigPath(subpath string) string {
	appData := os.Getenv(cfgconst.EnvAppData)
	if appData != "" {
		return filepath.Join(appData, filepath.FromSlash(subpath))
	}

	return ""
}

// resolveUnixConfigPath resolves configuration directory on Darwin or Linux.
func resolveUnixConfigPath(home, macSub, linuxSub string) string {
	if runtime.GOOS == cfgconst.OSDarwin {
		return filepath.Join(home, cfgconst.PathMacLibrary, cfgconst.PathMacAppSupport, filepath.FromSlash(macSub))
	}

	return filepath.Join(home, cfgconst.PathLinuxConfig, filepath.FromSlash(linuxSub))
}

// resolveOSConfigPath resolves cross-platform configuration directory.
func resolveOSConfigPath(winSub, macSub, linuxSub string) string {
	if runtime.GOOS == cfgconst.OSWindows {
		return resolveWindowsConfigPath(winSub)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return resolveUnixConfigPath(home, macSub, linuxSub)
}

// resolveVSCodeConfigDir returns the user settings directory for VS Code.
func resolveVSCodeConfigDir() string {
	return resolveOSConfigPath(cfgconst.PathVSCodeUser, cfgconst.PathVSCodeUser, cfgconst.PathVSCodeUser)
}

// resolveQBittorrentConfigDir returns the configuration directory for qBittorrent.
func resolveQBittorrentConfigDir() string {
	return resolveOSConfigPath(cfgconst.PathQBittorrent, cfgconst.PathQBittorrent, cfgconst.PathQBittorrent)
}

// resolveUTorrentConfigDir returns the configuration directory for uTorrent.
func resolveUTorrentConfigDir() string {
	return resolveOSConfigPath(cfgconst.PathUTorrent, cfgconst.PathUTorrent, cfgconst.PathUTorrentLower)
}

// dispatchToolDir dispatches directory resolution for a known tool type.
func dispatchToolDir(toolType cfgenum.ConfigToolType) string {
	switch toolType {
	case cfgenum.ConfigToolVSCode:
		return resolveVSCodeConfigDir()
	case cfgenum.ConfigToolQTorrent:
		return resolveQBittorrentConfigDir()
	case cfgenum.ConfigToolUTorrent:
		return resolveUTorrentConfigDir()
	}

	return ""
}

// resolveToolConfigDir returns the configuration directory for the given canonical tool.
func resolveToolConfigDir(tool string) string {
	toolType, err := cfgenum.ParseConfigToolType(tool)
	if err != nil {
		return ""
	}

	return dispatchToolDir(toolType)
}

// resolveQBittorrentFileName returns .ini on Windows and .conf elsewhere.
func resolveQBittorrentFileName() string {
	if runtime.GOOS == cfgconst.OSWindows {
		return cfgconst.FileQBittorrentIni
	}

	return cfgconst.FileQBittorrentConf
}

// resolveDefaultConfigFileNameForOS returns primary config file name for the OS.
func resolveDefaultConfigFileNameForOS(tool string) string {
	toolType, err := cfgenum.ParseConfigToolType(tool)
	if err != nil || (!toolType.IsQTorrent() && !toolType.IsUTorrent()) {
		return cfgconst.FileSettingsJSON
	}

	if toolType.IsQTorrent() {
		return resolveQBittorrentFileName()
	}

	return cfgconst.FileSettingsDat
}
