// Package startup manages user-scoped autostart entries created by
// gitmap on the host OS:
//
//   - Linux/Unix: XDG `.desktop` files in `$XDG_CONFIG_HOME/autostart/`
//     (or `$HOME/.config/autostart/`) carrying the `X-Gitmap-Managed=true`
//     key.
//   - macOS: LaunchAgent `.plist` files in `~/Library/LaunchAgents/`
//     carrying a top-level `<key>XGitmapManaged</key><true/>` marker.
//
// On both OSes, List enumerates ONLY entries that satisfy BOTH the
// filename prefix gate AND the in-file marker; Remove deletes one
// named entry only after re-confirming it carries the marker. A
// request to remove a third-party entry becomes a refused no-op,
// never a deletion.
//
// Windows is intentionally NOT covered by this package — Windows uses
// Registry `Run` keys / Startup folder shortcuts handled by separate
// code. The directory resolver returns an error on Windows so the
// CLI prints the "unsupported OS" message instead of touching a
// non-existent directory.
//
// macOS LaunchAgent lifecycle (launchctl load/unload) is NOT
// triggered here — list/remove operate on the .plist file ONLY. A
// removed plist takes effect at the next login or after a manual
// `launchctl unload`. This is intentional: invoking launchctl
// requires a running user GUI session and would make automated
// uninstall scripts brittle on CI / SSH sessions.
package startup

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// Entry is one gitmap-managed autostart record. Path is the absolute
// file path; Name is the basename WITHOUT the platform-specific
// extension (the form users pass to `startup-remove`); Exec is the
// command line surfaced so `startup-list` shows what would actually
// run at login. On macOS, Exec is the space-joined ProgramArguments
// (or the Program string if ProgramArguments is absent).
type Entry struct {
	Name string
	Path string
	Exec string
}

// AutostartDir returns the absolute path to the user's autostart
// directory for the current OS.
//
//   - $GITMAP_AUTOSTART_DIR (if set and non-empty): returned verbatim
//     on EVERY OS, bypassing the per-OS default below. Lets tests run
//     cross-platform without OS-specific skips and lets power users
//     override the install location (Nix profiles, portable installs).
//   - Linux/Unix: honors $XDG_CONFIG_HOME, falls back to
//     $HOME/.config/autostart per freedesktop.org base-dir spec.
//   - macOS: $HOME/Library/LaunchAgents (per Apple's LaunchAgents
//     documentation; the user-domain location requires no sudo).
//
// Returns an error on Windows so callers print the platform-specific
// "unsupported OS" message instead of touching a non-existent dir.
func AutostartDir() (string, error) {
	if override := strings.TrimSpace(os.Getenv(constants.EnvStartupAutostartDir)); len(override) > 0 {

		return override, nil
	}
	if runtime.GOOS == "windows" {

		return "", fmt.Errorf(constants.ErrStartupUnsupportedOS)
	}
	if runtime.GOOS == "darwin" {

		return darwinLaunchAgentsDir()
	}
	if xdg := strings.TrimSpace(os.Getenv("XDG_CONFIG_HOME")); len(xdg) > 0 {

		return filepath.Join(xdg, "autostart"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {

		return "", err
	}

	return filepath.Join(home, ".config", "autostart"), nil
}

// darwinLaunchAgentsDir resolves $HOME/Library/LaunchAgents. Kept
// separate so the macOS path-shape choice is in one place and tests
// can override $HOME to point at a temp dir.
func darwinLaunchAgentsDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {

		return "", err
	}

	return filepath.Join(home, "Library", "LaunchAgents"), nil
}

// List returns every gitmap-managed entry. A MISSING directory /
// registry key is treated as "zero entries", NOT an error — fresh
// accounts that have never had any autostart entry shouldn't see a
// scary error from `gitmap startup-list`. On Windows, enumerates
// BOTH the Registry Run-key and Startup-folder backends.
func List() ([]Entry, error) {
	if runtime.GOOS == "windows" {

		return listWindows()
	}
	dir, err := AutostartDir()
	if err != nil {

		return nil, err
	}
	files, err := os.ReadDir(dir)
	if os.IsNotExist(err) {

		return nil, nil
	}
	if err != nil {

		return nil, fmt.Errorf(constants.ErrStartupReadDir, dir, err)
	}

	return collectManaged(dir, files), nil
}
