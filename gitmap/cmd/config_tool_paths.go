package cmd

import (
	"os"
	"path/filepath"
	"runtime"
)

// resolveWindowsConfigPath resolves path within APPDATA on Windows.
func resolveWindowsConfigPath(subpath string) string {
	appData := os.Getenv("APPDATA")
	if appData != "" {
		return filepath.Join(appData, filepath.FromSlash(subpath))
	}

	return ""
}

// resolveOSConfigPath resolves cross-platform configuration directory.
func resolveOSConfigPath(winSub, macSub, linuxSub string) string {
	if runtime.GOOS == "windows" {
		return resolveWindowsConfigPath(winSub)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	if runtime.GOOS == "darwin" {
		return filepath.Join(home, "Library", "Application Support", filepath.FromSlash(macSub))
	}

	return filepath.Join(home, ".config", filepath.FromSlash(linuxSub))
}

// resolveVSCodeConfigDir returns the user settings directory for VS Code.
func resolveVSCodeConfigDir() string {
	return resolveOSConfigPath("Code/User", "Code/User", "Code/User")
}

// resolveQBittorrentConfigDir returns the configuration directory for qBittorrent.
func resolveQBittorrentConfigDir() string {
	return resolveOSConfigPath("qBittorrent", "qBittorrent", "qBittorrent")
}

// resolveUTorrentConfigDir returns the configuration directory for uTorrent.
func resolveUTorrentConfigDir() string {
	return resolveOSConfigPath("uTorrent", "uTorrent", "utorrent")
}

// resolveToolConfigDir returns the configuration directory for the given canonical tool.
func resolveToolConfigDir(tool string) string {
	switch tool {
	case "vscode":

		return resolveVSCodeConfigDir()
	case "qtorrent":

		return resolveQBittorrentConfigDir()
	case "utorrent":

		return resolveUTorrentConfigDir()
	default:

		return ""
	}
}

// resolveQBittorrentFileName returns .ini on Windows and .conf elsewhere.
func resolveQBittorrentFileName() string {
	if runtime.GOOS == "windows" {
		return "qBittorrent.ini"
	}

	return "qBittorrent.conf"
}

// resolveDefaultConfigFileNameForOS returns primary config file name for the OS.
func resolveDefaultConfigFileNameForOS(tool string) string {
	if tool == "qtorrent" {
		return resolveQBittorrentFileName()
	}

	if tool == "utorrent" {
		return "settings.dat"
	}

	return "settings.json"
}
