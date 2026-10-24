package startup

// Windows backend selection + dispatch. Lives in its own file (not
// inline in startup.go) so the cross-OS startup.go stays focused on
// Linux/.desktop logic and the Windows surface area is easy to find
// and code-review independently.
//
// Two backends are supported and a single Backend value flows
// through Add → writeManagedWindows. List enumerates BOTH backends
// unconditionally so users do not have to remember which one they
// chose at Add time. Remove also looks in BOTH backends; the
// tracking subkey under HKCU\Software\Gitmap pins which backend the
// entry actually lives in so we never delete a Run value that
// belongs to the .lnk path's tracking record (or vice versa).

import (
	"fmt"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)
import "github.com/alimtvnetwork/gitmap-v28/cli/startup/enums"

// BackendType enumerates the Windows-specific add targets. Linux/macOS
// have one canonical backend each (XDG .desktop / LaunchAgents
// .plist) so they ignore this value entirely. The zero value
// (BackendUnspecified) means "let the dispatcher pick the OS default"
// which is `BackendRegistry` on Windows.
// BackendType aliases the strongly-typed enum from the enums package.
type BackendType = enums.BackendType

const (
	BackendUnspecified   = enums.BackendUnspecified
	BackendRegistry      = enums.BackendRegistry
	BackendStartupFolder = enums.BackendStartupFolder
	BackendRegistryHKLM  = enums.BackendRegistryHKLM
)

// ParseBackend aliases the enums package parser so existing callers keep working.
func ParseBackend(s string) (BackendType, *apperror.AppError) {
	return enums.ParseBackend(s)
}

// resolveBackendForAdd picks the concrete backend to use when the
// caller passed BackendUnspecified. Registry is the Windows default
// (matches what most provisioning scripts use and what `regedit`
// users expect to see after running `gitmap startup-add`).
func resolveBackendForAdd(b BackendType) BackendType {
	if b != BackendUnspecified {
		return b
	}

	if runtime.GOOS == "windows" {
		return BackendRegistry
	}

	return BackendUnspecified
}

// addWindows is the Windows branch of Add. Routes through the
// chosen backend (Registry or Startup-folder). Lives here (not in
// add.go) so all Windows dispatch code lives together. Both
// backends return AddResult with Path set to a backend-appropriate
// locator: registry value path for the Run-key entry, filesystem
// path for the .lnk Startup-folder shortcut.
func addWindows(clean string, opts AddOptions) (AddResult, error) {
	backend := resolveBackendForAdd(opts.Backend)
	switch backend {
	case BackendRegistry:

		return addWindowsRegistry(clean, opts)
	case BackendStartupFolder:

		return addWindowsStartupFolder(clean, opts)
	case BackendRegistryHKLM:

		return addWindowsRegistryHKLM(clean, opts)
	case BackendUnspecified:

		return AddResult{}, apperror.NewValidationError(fmt.Sprintf(constants.ErrStartupAddBadBackend, backend.String()))
	default:

		return AddResult{}, apperror.NewValidationError(fmt.Sprintf(constants.ErrStartupAddBadBackend, backend.String()))
	}
}

// removeWindows routes the deletion to the requested backend. When
// opts.Backend is BackendUnspecified, falls back to the legacy
// multi-backend probe: HKCU registry → HKLM registry → startup-
// folder, first non-NoOp wins. When the user passes --backend
// explicitly, only that backend is touched and a missing entry
// returns NoOp without silently checking the other backends
// (important: a user removing a registry entry must not have a
// same-named .lnk deleted as a "courtesy"). Returns NoOp only when
// the chosen backend (or all three, in fallback mode) report no
// entry by that name.
func removeWindows(clean string, opts RemoveOptions) (RemoveResult, error) {
	switch opts.Backend {
	case BackendRegistry:

		return removeWindowsRegistry(clean, opts)
	case BackendStartupFolder:

		return removeWindowsStartupFolder(clean, opts)
	case BackendRegistryHKLM:

		return removeWindowsRegistryHKLM(clean, opts)
	case BackendUnspecified:
		// Fall through to the multi-backend probe below. Listed
		// explicitly so the exhaustive linter is satisfied and a
		// future Backend addition forces a deliberate decision
		// here instead of silently landing in the fallback.
	}

	res, err := removeWindowsRegistry(clean, opts)
	if err != nil || res.Status != RemoveNoOp {
		return res, err
	}

	res, err = removeWindowsRegistryHKLM(clean, opts)
	if err != nil || res.Status != RemoveNoOp {
		return res, err
	}

	return removeWindowsStartupFolder(clean, opts)
}

// listWindows enumerates ALL three backends and concatenates the
// results. Order is HKCU-registry → HKLM-registry → startup-folder
// so the CLI can group by backend in the rendered output. Returns
// nil (NOT an error) when every backend is empty — fresh accounts
// shouldn't see a scary error from `gitmap startup-list`.
func listWindows() ([]Entry, error) {
	reg, err := listWindowsRegistry()
	if err != nil {
		return nil, err
	}

	regHKLM, err := listWindowsRegistryHKLM()
	if err != nil {
		return nil, err
	}

	folder, err := listWindowsStartupFolder()
	if err != nil {
		return nil, err
	}

	out := make([]Entry, 0, len(reg)+len(regHKLM)+len(folder))
	out = append(out, reg...)
	out = append(out, regHKLM...)
	out = append(out, folder...)

	return out, nil
}
