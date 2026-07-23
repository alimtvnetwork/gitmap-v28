//go:build !windows

package startup

// Non-Windows stubs for the Windows-only registry backend. These
// exist so cross-platform call sites (winbackend.go's switch,
// startup.go's List dispatcher) compile on Linux/macOS without
// build tags. The functions are unreachable at runtime on those
// OSes because the `runtime.GOOS == "windows"` guards in Add/
// Remove/List short-circuit before reaching them — the stubs
// return a clear "wrong OS" error in case some future caller
// bypasses the guard, so the failure mode is loud rather than
// silently writing nothing.

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// addWindowsRegistry stub: unreachable on non-Windows. Returns the
// standard unsupported-OS error so the failure mode is identical
// to AutostartDir() refusing on Windows in older versions of this
// package — callers see one consistent message.
func addWindowsRegistry(_ string, _ AddOptions) (AddResult, error) {
	return AddResult{}, fmt.Errorf(constants.ErrStartupUnsupportedOS)
}

// removeWindowsRegistry stub: unreachable on non-Windows.
func removeWindowsRegistry(_ string, _ RemoveOptions) (RemoveResult, error) {
	return RemoveResult{}, fmt.Errorf(constants.ErrStartupUnsupportedOS)
}

// listWindowsRegistry stub: returns nil entries (no error) so a
// hypothetical cross-OS lister that always calls every backend
// gets an empty slice on non-Windows rather than a hard error.
func listWindowsRegistry() ([]Entry, error) {
	return nil, nil
}

// addWindowsRegistryHKLM stub: machine-wide HKLM backend is
// Windows-only. Returns the standard unsupported-OS error so
// cross-platform callers get one consistent message.
func addWindowsRegistryHKLM(_ string, _ AddOptions) (AddResult, error) {
	return AddResult{}, fmt.Errorf(constants.ErrStartupUnsupportedOS)
}

// removeWindowsRegistryHKLM stub: unreachable on non-Windows.
func removeWindowsRegistryHKLM(_ string, _ RemoveOptions) (RemoveResult, error) {
	return RemoveResult{}, fmt.Errorf(constants.ErrStartupUnsupportedOS)
}

// listWindowsRegistryHKLM stub: returns nil entries (no error)
// matching the HKCU listWindowsRegistry stub above so the
// cross-OS list dispatcher gets identical empty-slice behavior
// on Linux/macOS for both registry backends.
func listWindowsRegistryHKLM() ([]Entry, error) {
	return nil, nil
}

// trackingSubkeyExists stub: false on non-Windows so the .lnk
// renderer's classify checks degrade gracefully if any code path
// reaches them off-Windows.
func trackingSubkeyExists(_, _ string) bool {
	return false
}

// writeTrackingSubkey stub: unreachable on non-Windows but kept
// symmetric with the windows file so .lnk code can call it. The
// trailing workingDir parameter mirrors the Windows signature.
func writeTrackingSubkey(_, _, _, _, _ string) error {
	return fmt.Errorf(constants.ErrStartupUnsupportedOS)
}

// deleteTrackingSubkey stub.
func deleteTrackingSubkey(_, _ string) error {
	return fmt.Errorf(constants.ErrStartupUnsupportedOS)
}
