//go:build windows

package startup

// Registry-backend Remove + List for Windows. Pairs with
// winregistry_windows.go (which has Add). Split for the per-file
// budget and so each operation is independently scannable.
//
// Ownership model: the Run key holds ONLY the direct autostart
// value (`gitmap-<name> = "<exec>"`); the tracking subkey under
// <hive>\Software\Gitmap\StartupRegistry\<name> is the SOLE marker.
// Remove deletes the Run value only when the tracking subkey
// confirms gitmap ownership; List enumerates the tracking scope
// first, then verifies each entry has a live Run value.
//
// Both HKCU (default `--backend=registry`) and HKLM (opt-in
// `--backend=registry-hklm`) flow through the *At helpers; the
// thin HKCU-only wrappers exist so the public surface stays
// backwards compatible with code that doesn't care about the hive.

import (
	"fmt"
	"strings"

	"golang.org/x/sys/windows/registry"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// removeWindowsRegistry is the HKCU adapter retained for the
// public dispatcher in winbackend.go. The HKLM Remove path lives
// in winregistry_hklm_windows.go and reuses the same *At helpers.
func removeWindowsRegistry(clean string, opts RemoveOptions) (RemoveResult, error) {
	return removeWindowsRegistryAt(registry.CURRENT_USER, hiveLabelHKCU, clean, opts)
}

// removeWindowsRegistryAt deletes a managed Run-key value plus the
// tracking subkey under the given hive. Returns RemoveNoOp if the
// value does not exist, RemoveRefused if it exists but the
// tracking subkey is missing. Mirrors the Linux RemoveStatus
// contract one-for-one.
func removeWindowsRegistryAt(root registry.Key, hive, clean string,
	opts RemoveOptions) (RemoveResult, error) {
	valueName := constants.StartupWinValuePrefix + clean
	exists, managed, err := classifyRunValueAt(root, valueName, clean)
	if err != nil {

		return RemoveResult{}, err
	}
	if !exists {

		return RemoveResult{Status: RemoveNoOp, DryRun: opts.DryRun}, nil
	}
	if !managed {

		return RemoveResult{Status: RemoveRefused,
			Path: runValuePathFor(hive, valueName), DryRun: opts.DryRun}, nil
	}
	if opts.DryRun {

		return RemoveResult{Status: RemoveDeleted,
			Path: runValuePathFor(hive, valueName), DryRun: true}, nil
	}
	if err := deleteRunValueAt(root, valueName); err != nil {

		return RemoveResult{}, err
	}
	if err := deleteTrackingSubkeyAt(root, constants.RegGitmapRegistrySub, clean); err != nil {

		return RemoveResult{}, err
	}

	return RemoveResult{Status: RemoveDeleted, Path: runValuePathFor(hive, valueName)}, nil
}

// deleteRunValueAt removes the single autostart value under the
// given hive. Missing value is not an error — a previous partial
// Remove (crash after value delete, before subkey delete) would
// re-enter here with the value already gone, and the subkey-delete
// pass below should still run.
func deleteRunValueAt(root registry.Key, valueName string) error {
	k, err := registry.OpenKey(root,
		constants.RegRunKeyPath, registry.SET_VALUE)
	if err != nil {

		return fmt.Errorf(constants.ErrStartupRegistryOpen, constants.RegRunKeyPath, err)
	}
	defer k.Close()

	if err := k.DeleteValue(valueName); err != nil && err != registry.ErrNotExist {

		return fmt.Errorf(constants.ErrStartupRegistryWrite, valueName, err)
	}

	return nil
}

// deleteTrackingSubkey is the HKCU adapter retained for the .lnk
// Startup-folder backend (which never targets HKLM).
func deleteTrackingSubkey(parent, name string) error {
	return deleteTrackingSubkeyAt(registry.CURRENT_USER, parent, name)
}

// deleteTrackingSubkeyAt removes <root>\<parent>\<name>. Missing
// key is not an error — Add could have failed mid-way and left
// only the Run value, in which case the missing tracking subkey
// is expected.
func deleteTrackingSubkeyAt(root registry.Key, parent, name string) error {
	full := parent + `\` + name
	if err := registry.DeleteKey(root, full); err != nil &&
		err != registry.ErrNotExist {

		return fmt.Errorf(constants.ErrStartupRegistryWrite, full, err)
	}

	return nil
}

// listWindowsRegistry enumerates Run-key values whose name starts
// with the gitmap- prefix AND whose tracking subkey under
// HKCU\Software\Gitmap\StartupRegistry confirms gitmap ownership.
// No sibling-marker filter needed — the Run key now contains only
// real autostart values, never a `.gitmap-managed` companion.
// Returns nil (NOT an error) when the Run key itself does not
// exist — fresh accounts shouldn't see a scary error from
// `gitmap startup-list`.
func listWindowsRegistry() ([]Entry, error) {
	return listWindowsRegistryAt(registry.CURRENT_USER, hiveLabelHKCU)
}

// listWindowsRegistryAt enumerates Run-key values under the given
// hive whose name starts with the gitmap- prefix AND whose
// tracking subkey under <hive>\Software\Gitmap\StartupRegistry
// confirms gitmap ownership. Returns nil (NOT an error) when the
// Run key itself does not exist — fresh accounts shouldn't see a
// scary error from `gitmap startup-list`.
func listWindowsRegistryAt(root registry.Key, hive string) ([]Entry, error) {
	k, err := registry.OpenKey(root,
		constants.RegRunKeyPath, registry.QUERY_VALUE|registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		if err == registry.ErrNotExist {

			return nil, nil
		}

		return nil, fmt.Errorf(constants.ErrStartupRegistryOpen, constants.RegRunKeyPath, err)
	}
	defer k.Close()

	names, err := k.ReadValueNames(-1)
	if err != nil {

		return nil, fmt.Errorf(constants.ErrStartupRegistryRead, constants.RegRunKeyPath, err)
	}

	return collectRegistryManagedAt(root, hive, k, names), nil
}

// collectRegistryManagedAt is the per-entry filter. Skips:
//   - values without the gitmap- name prefix
//   - values whose tracking subkey under
//     <hive>\Software\Gitmap\StartupRegistry is missing
//     (third-party values that happen to share our prefix —
//     refuse to claim them)
func collectRegistryManagedAt(root registry.Key, hive string,
	k registry.Key, names []string) []Entry {
	var out []Entry
	for _, name := range names {
		if !strings.HasPrefix(name, constants.StartupWinValuePrefix) {
			continue
		}
		entry, ok := readManagedRegistryValueAt(root, hive, k, name)
		if !ok {
			continue
		}
		out = append(out, entry)
	}

	return out
}

// readManagedRegistryValueAt reads the command value and verifies
// the tracking subkey exists in the same hive; returns ok=false if
// either is missing.
func readManagedRegistryValueAt(root registry.Key, hive string,
	k registry.Key, valueName string) (Entry, bool) {
	exec, _, err := k.GetStringValue(valueName)
	if err != nil {

		return Entry{}, false
	}
	clean := strings.TrimPrefix(valueName, constants.StartupWinValuePrefix)
	if !trackingSubkeyExistsAt(root, constants.RegGitmapRegistrySub, clean) {

		return Entry{}, false
	}

	return Entry{
		Name: valueName,
		Path: runValuePathFor(hive, valueName),
		Exec: exec,
	}, true
}
