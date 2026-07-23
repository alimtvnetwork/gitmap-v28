package startup

// Path + filesystem helpers for the .lnk Startup folder backend.
// Split from winshortcut.go to keep both files under the per-file
// budget. Both helpers are cross-platform Go (no Windows APIs);
// the runtime guard for non-Windows callers lives in
// addWindowsStartupFolder, not here.

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// startupFolderDir returns %APPDATA%\Microsoft\Windows\Start Menu\
// Programs\Startup. Honors $APPDATA so test fixtures can redirect
// to a temp dir without touching the real roaming profile.
func startupFolderDir() (string, error) {
	roaming := os.Getenv("APPDATA")
	if len(roaming) == 0 {

		return "", fmt.Errorf("APPDATA env var is empty")
	}

	return filepath.Join(roaming, constants.StartupFolderRelative), nil
}

// fileExists is a tiny os.Stat wrapper. Treats permission errors
// as "exists" — conservative posture: better to refuse a write
// than to silently overwrite a file we can't read.
func fileExists(p string) bool {
	if _, err := os.Stat(p); err == nil || !os.IsNotExist(err) {
		return true
	}

	return false
}
