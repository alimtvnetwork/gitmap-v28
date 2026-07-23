//go:build windows

package startup

// Registry backend for `gitmap startup-add` / `startup-remove` /
// `startup-list` on Windows. Writes ONE direct value per entry to:
//
//   <hive>\Software\Microsoft\Windows\CurrentVersion\Run
//     gitmap-<name> = "<exec>"
//
// Where <hive> is HKCU for the per-user `--backend=registry` (the
// long-standing default) and HKLM for `--backend=registry-hklm`
// (machine-wide, opt-in, requires admin). Both hives use the same
// relative path so the only thing that changes between the two
// code paths is the registry root key handed to the helpers.
//
// Ownership is tracked OUT-OF-BAND under a separate scope:
//
//   <hive>\Software\Gitmap\StartupRegistry\<name>
//     Exec      = "<exec>"
//     CreatedAt = "<RFC3339-UTC>"
//     Source    = "registry" | "registry-hklm"
//     WorkingDir = "<dir>" (only when --working-dir was passed)
//
// Why no sibling marker in the Run key: Windows treats EVERY value
// under Run as an autostart command and feeds it to the shell at
// login. A `gitmap-<name>.gitmap-managed = "true"` sibling value
// shows up in Task Manager's Startup tab and is dispatched as a
// command (the literal string "true" — silently fails, but
// pollutes the user's startup surface). Keeping the ownership
// marker in a SEPARATE scope under <hive>\Software\Gitmap means the
// Run key contains only real autostart commands, exactly like a
// hand-edited entry would.
//
// Build tag: only compiled on windows. The non-windows stub in
// winregistry_other.go provides the same symbols so cross-platform
// callers in winbackend.go compile everywhere.

import (
	"fmt"
	"time"

	"golang.org/x/sys/windows/registry"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// hiveLabel is the user-visible registry-hive prefix written into
// AddResult.Path / RemoveResult.Path / Entry.Path so a user can
// paste the value straight into regedit. Centralized so the
// rendering helper (`runValuePathFor`) and the list-filter
// discriminator stay in lockstep.
const (
	hiveLabelHKCU = `HKCU`
	hiveLabelHKLM = `HKLM`
)

// addWindowsRegistry implements the per-user (HKCU) Add path. Thin
// adapter over addWindowsRegistryAt so the public dispatcher in
// winbackend.go can keep using the old function name without
// caring which hive is targeted.
func addWindowsRegistry(clean string, opts AddOptions) (AddResult, error) {
	return addWindowsRegistryAt(registry.CURRENT_USER, hiveLabelHKCU,
		constants.StartupBackendRegistry, clean, opts)
}

// addWindowsRegistryAt is the hive-agnostic Add path. The two
// writes (Run value + tracking subkey) are NOT atomic at the
// Windows API level — there is no transactional registry write
// across keys. We tolerate this: the tracking subkey is written
// FIRST, so a crash between the two writes leaves an "owned but
// inactive" record that future Add re-runs can safely overwrite
// (because classifyRunValueAt will see managed=true).
func addWindowsRegistryAt(root registry.Key, hive, source string,
	clean string, opts AddOptions) (AddResult, error) {
	valueName := constants.StartupWinValuePrefix + clean
	exists, managed, err := classifyRunValueAt(root, valueName, clean)
	if err != nil {

		return AddResult{}, err
	}
	if exists && !managed {

		return AddResult{Status: AddRefused, Path: runValuePathFor(hive, valueName)}, nil
	}
	if exists && managed && !opts.Force {

		return AddResult{Status: AddExists, Path: runValuePathFor(hive, valueName)}, nil
	}
	if err := writeTrackingSubkeyAt(root, constants.RegGitmapRegistrySub, clean,
		opts.Exec, source, opts.WorkingDir); err != nil {

		return AddResult{}, err
	}
	if err := writeRunValueAt(root, valueName, opts.Exec); err != nil {

		return AddResult{}, err
	}
	if exists {

		return AddResult{Status: AddOverwritten, Path: runValuePathFor(hive, valueName)}, nil
	}

	return AddResult{Status: AddCreated, Path: runValuePathFor(hive, valueName)}, nil
}

// classifyRunValue is the HKCU adapter retained so existing call
// sites keep compiling. New code paths should call
// classifyRunValueAt directly so the hive is explicit.
func classifyRunValue(valueName, clean string) (bool, bool, error) {
	return classifyRunValueAt(registry.CURRENT_USER, valueName, clean)
}

// classifyRunValueAt returns (exists, managed, error) for a
// Run-key value under the given hive. "Managed" means the tracking
// subkey under <hive>\Software\Gitmap\StartupRegistry\<clean>
// exists. The Run key itself carries NO marker — keeping the
// autostart surface clean is the whole point of the direct-value
// model.
func classifyRunValueAt(root registry.Key, valueName, clean string) (bool, bool, error) {
	k, err := registry.OpenKey(root, constants.RegRunKeyPath, registry.QUERY_VALUE)
	if err != nil {
		if err == registry.ErrNotExist {

			return false, false, nil
		}

		return false, false, fmt.Errorf(constants.ErrStartupRegistryOpen, constants.RegRunKeyPath, err)
	}
	defer k.Close()

	if _, _, err := k.GetStringValue(valueName); err != nil {
		if err == registry.ErrNotExist {

			return false, false, nil
		}

		return false, false, fmt.Errorf(constants.ErrStartupRegistryRead, valueName, err)
	}
	hasTracking := trackingSubkeyExistsAt(root, constants.RegGitmapRegistrySub, clean)

	return true, hasTracking, nil
}

// trackingSubkeyExists is the HKCU-rooted convenience wrapper used
// by every legacy call site (winshortcut.go's StartupFolder
// classify, the registry remove path, etc.). New hive-aware code
// should call trackingSubkeyExistsAt directly.
func trackingSubkeyExists(parent, name string) bool {
	return trackingSubkeyExistsAt(registry.CURRENT_USER, parent, name)
}

// trackingSubkeyExistsAt checks for <root>\<parent>\<name>. Returns
// false on any open error — same conservative posture as
// classifyTarget on Linux: an unreadable key is "not ours" so we
// refuse to delete it.
func trackingSubkeyExistsAt(root registry.Key, parent, name string) bool {
	full := parent + `\` + name
	k, err := registry.OpenKey(root, full, registry.QUERY_VALUE)
	if err != nil {

		return false
	}
	k.Close()

	return true
}

// writeRunValue is the HKCU-rooted convenience wrapper retained for
// any older caller. The HKLM Add path goes through writeRunValueAt
// directly.
func writeRunValue(valueName, exec string) error {
	return writeRunValueAt(registry.CURRENT_USER, valueName, exec)
}

// writeRunValueAt writes the autostart command to the Run key
// under the given hive. ONE value per entry — no sibling marker.
// Ownership tracking lives entirely under <hive>\Software\Gitmap\
// StartupRegistry.
func writeRunValueAt(root registry.Key, valueName, exec string) error {
	k, _, err := registry.CreateKey(root,
		constants.RegRunKeyPath, registry.SET_VALUE)
	if err != nil {

		return fmt.Errorf(constants.ErrStartupRegistryOpen, constants.RegRunKeyPath, err)
	}
	defer k.Close()

	if err := k.SetStringValue(valueName, exec); err != nil {

		return fmt.Errorf(constants.ErrStartupRegistryWrite, valueName, err)
	}

	return nil
}

// writeTrackingSubkey is the HKCU-rooted adapter; existing callers
// (the .lnk Startup-folder backend) keep working untouched.
func writeTrackingSubkey(parent, name, exec, source, workingDir string) error {
	return writeTrackingSubkeyAt(registry.CURRENT_USER, parent, name, exec, source, workingDir)
}

// writeTrackingSubkeyAt creates <root>\<parent>\<name> with the
// standard metadata values (Exec / CreatedAt / Source) plus an
// optional WorkingDir value when the caller passed a non-empty
// working directory. CreatedAt is RFC3339 UTC for stable cross-
// locale parsing by future tooling. Used by every backend (HKCU
// registry, HKLM registry, .lnk Startup-folder) — Source labels
// the originating backend so a future audit pass can attribute
// each tracking record to the flag value the user typed.
func writeTrackingSubkeyAt(root registry.Key, parent, name, exec, source, workingDir string) error {
	full := parent + `\` + name
	k, _, err := registry.CreateKey(root, full, registry.SET_VALUE)
	if err != nil {

		return fmt.Errorf(constants.ErrStartupRegistryOpen, full, err)
	}
	defer k.Close()

	if err := k.SetStringValue(constants.RegTrackKeyExec, exec); err != nil {

		return fmt.Errorf(constants.ErrStartupRegistryWrite, constants.RegTrackKeyExec, err)
	}
	if err := k.SetStringValue(constants.RegTrackKeyCreatedAt,
		time.Now().UTC().Format(time.RFC3339)); err != nil {

		return fmt.Errorf(constants.ErrStartupRegistryWrite, constants.RegTrackKeyCreatedAt, err)
	}
	if err := k.SetStringValue(constants.RegTrackKeySource, source); err != nil {

		return fmt.Errorf(constants.ErrStartupRegistryWrite, constants.RegTrackKeySource, err)
	}
	if workingDir != "" {
		if err := k.SetStringValue(constants.RegTrackKeyWorkingDir, workingDir); err != nil {

			return fmt.Errorf(constants.ErrStartupRegistryWrite, constants.RegTrackKeyWorkingDir, err)
		}
	}

	return nil
}

// runValuePath formats a stable user-facing locator for a HKCU
// Run-key value. Kept as a thin alias around runValuePathFor so
// existing callers and tests keep compiling unchanged.
func runValuePath(valueName string) string {
	return runValuePathFor(hiveLabelHKCU, valueName)
}

// runValuePathFor formats `<hive>\<RunPath>\<valueName>`. Used for
// AddResult.Path / Entry.Path so `gitmap startup-list` shows the
// full registry path the user would see in regedit. The hive label
// is also the cross-platform discriminator the list-filter relies
// on (HKCU\ vs HKLM\) — see cmd/startuplistfilter.go.
func runValuePathFor(hive, valueName string) string {
	return hive + `\` + constants.RegRunKeyPath + `\` + valueName
}
