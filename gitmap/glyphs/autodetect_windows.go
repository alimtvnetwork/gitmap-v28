//go:build windows

package glyphs

import "os"

// init swaps the platform-neutral stub for the Windows-aware detector.
func init() {
	isLegacyWindowsHost = legacyWindowsHost
}

// legacyWindowsHost reports true when gitmap is running under the
// classic ConsoleHost (powershell.exe 5.1 or cmd.exe), where the
// default font (Consolas / Lucida Console / Courier New) does not
// include emoji glyphs. Modern hosts (Windows Terminal, VS Code,
// ConEmu) advertise themselves via env vars.
func legacyWindowsHost() bool {
	if os.Getenv("WT_SESSION") != "" {
		return false
	}
	if os.Getenv("TERM_PROGRAM") == "vscode" {
		return false
	}
	if os.Getenv("ConEmuANSI") == "ON" {
		return false
	}
	if os.Getenv("ALACRITTY_LOG") != "" {
		return false
	}

	return true
}
