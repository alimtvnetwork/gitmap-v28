//go:build windows

package glyphs

import (
	"os"
)

// init swaps the platform-neutral stub for the Windows-aware detector.
func init() {
	isLegacyWindowsHost = legacyWindowsHost
}

// legacyWindowsHost reports true when gitmap is running under the
// classic ConsoleHost (powershell.exe 5.1 or cmd.exe), where the
// default font (Consolas / Lucida Console / Courier New) does not
// include emoji glyphs. Modern hosts (Windows Terminal, VS Code,
// ConEmu, WezTerm, Ghostty) advertise themselves via env vars.
func legacyWindowsHost() bool {
	if isModernWindowsHost() {
		return false
	}

	return true
}

func isModernWindowsHost() bool {
	if os.Getenv("WT_SESSION") != "" {
		return true
	}
	if os.Getenv("TERM_PROGRAM") == "vscode" || os.Getenv("VSCODE_PID") != "" {
		return true
	}
	if os.Getenv("ConEmuANSI") == "ON" || os.Getenv("ALACRITTY_LOG") != "" {
		return true
	}
	if os.Getenv("WEZTERM_PANE") != "" || os.Getenv("GHOSTTY_RESOURCES_DIR") != "" {
		return true
	}
	term := os.Getenv("TERM")

	return term != "" && term != "dumb"
}
