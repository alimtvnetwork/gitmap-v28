//go:build windows

package startup

// HKLM (machine-wide) registry backend for `gitmap startup-add` /
// `startup-remove` / `startup-list`. Targets the SAME relative
// Run-key path as the per-user HKCU backend but rooted under
// HKEY_LOCAL_MACHINE so the autostart value fires for every
// interactive user on the machine. Writes require administrator
// privileges (UAC elevation); the elevation check runs UP-FRONT
// in addWindowsRegistryHKLM / removeWindowsRegistryHKLM so the
// user sees a friendly, actionable error instead of a raw Win32
// "Access is denied" mid-write.
//
// All the heavy lifting (classify / write / delete / list) is
// shared with the HKCU backend via the *At helpers in
// winregistry_windows.go and winregistry_remove_windows.go — this
// file is intentionally thin so the HKCU and HKLM code paths can
// never disagree on the marker contract.

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// addWindowsRegistryHKLM is the HKLM Add path. Validates elevation
// FIRST, then dispatches to the hive-agnostic addWindowsRegistryAt
// with registry.LOCAL_MACHINE as the root and the registry-hklm
// Source label so the tracking record records which backend the
// user typed.
func addWindowsRegistryHKLM(clean string, opts AddOptions) (AddResult, error) {
	if err := requireWindowsAdminForHKLM(); err != nil {

		return AddResult{}, err
	}

	return addWindowsRegistryAt(registry.LOCAL_MACHINE, hiveLabelHKLM,
		constants.StartupBackendRegistryHKLM, clean, opts)
}

// removeWindowsRegistryHKLM is the HKLM Remove path. The admin
// check is gated on opts.DryRun: a dry-run is read-only and works
// without elevation (so a non-admin user can preview what `sudo`
// would do); a live remove requires elevation up-front.
func removeWindowsRegistryHKLM(clean string, opts RemoveOptions) (RemoveResult, error) {
	if !opts.DryRun {
		if err := requireWindowsAdminForHKLM(); err != nil {

			return RemoveResult{}, err
		}
	}

	return removeWindowsRegistryAt(registry.LOCAL_MACHINE, hiveLabelHKLM, clean, opts)
}

// listWindowsRegistryHKLM enumerates HKLM Run-key entries managed
// by gitmap. Reading HKLM\Software\Microsoft\Windows\CurrentVersion
// \Run does NOT require admin (any user can query it), so no
// elevation check is needed here — `gitmap startup-list
// --backend=registry-hklm` works for every user.
func listWindowsRegistryHKLM() ([]Entry, error) {
	return listWindowsRegistryAt(registry.LOCAL_MACHINE, hiveLabelHKLM)
}

// requireWindowsAdminForHKLM probes the current process token for
// elevation via TokenElevation and returns ErrStartupHKLMNotAdmin
// when the process is not elevated. We probe the token (not the
// HKLM ACL) so the failure mode is uniform across systems with
// custom registry ACLs and so the user sees the actionable
// "re-run from an elevated shell" message before any mutation
// attempt — important because a partial write (Run value set,
// tracking subkey not) would leave a third-party-looking entry
// on disk.
func requireWindowsAdminForHKLM() error {
	if isProcessElevated() {

		return nil
	}

	return fmt.Errorf(constants.ErrStartupHKLMNotAdmin)
}

// isProcessElevated returns true when the current process token
// has TokenElevation.TokenIsElevated set. Returns false on any
// API failure so the caller surfaces the friendly admin error
// instead of an opaque internal one — the absolute worst case is
// that an admin user is told to re-elevate, which is corrected by
// re-running the same command (idempotent) from any shell.
func isProcessElevated() bool {
	var token windows.Token
	if err := windows.OpenProcessToken(windows.CurrentProcess(),
		windows.TOKEN_QUERY, &token); err != nil {

		return false
	}
	defer token.Close()

	var elevation uint32
	var returned uint32
	err := windows.GetTokenInformation(token, windows.TokenElevation,
		(*byte)(unsafe.Pointer(&elevation)), uint32(unsafe.Sizeof(elevation)),
		&returned)
	if err != nil {

		return false
	}

	return elevation != 0
}
