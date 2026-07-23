package cmd

// In-memory filters for `gitmap startup-list --backend=... --name=...`.
// Filtering happens AFTER startup.List() returns so we don't have to
// fork the platform-specific List dispatchers — those still return
// the canonical full set; this layer just narrows it.
//
// Why post-filter (not push-down): startup.List on Windows already
// concatenates BOTH backends; the only stable cross-backend
// discriminator that survives the round-trip is the Path shape
// (registry entries are formatted as `HKCU\...`, .lnk entries are
// real filesystem paths). Push-down would require plumbing a
// filter struct through three OS-specific code paths for a savings
// of at most a few dozen registry reads — not worth the surface.

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/startup"
)

// filterStartupList applies --backend and --name to the full
// listing. An empty filter value means "no filter for this field".
// Returns a fresh slice so the caller's input is not mutated even
// when both filters are no-ops.
func filterStartupList(entries []startup.Entry, backend, name string) []startup.Entry {
	if backend == "" && name == "" {
		out := make([]startup.Entry, len(entries))
		copy(out, entries)

		return out
	}
	out := make([]startup.Entry, 0, len(entries))
	for _, e := range entries {
		if !matchesBackend(e, backend) {
			continue
		}
		if !matchesName(e, name) {
			continue
		}
		out = append(out, e)
	}

	return out
}

// matchesBackend returns true when the entry came from the
// requested backend. The discriminator is the Path shape:
//
//   - HKCU registry-backend entries set Path to `HKCU\...\<value>`
//     (see runValuePathFor in startup/winregistry_windows.go).
//   - HKLM registry-backend entries set Path to `HKLM\...\<value>`
//     (machine-wide, opt-in via --backend=registry-hklm).
//   - Startup-folder entries set Path to a real filesystem path
//     ending in .lnk.
//   - Linux/macOS entries don't belong to any Windows backend, so
//     a --backend filter on those OSes always matches zero
//     entries — exactly the behavior a user scripting cross-OS
//     would expect (no false positives).
//
// Empty backend = no filter.
func matchesBackend(e startup.Entry, backend string) bool {
	if backend == "" {
		return true
	}
	switch backend {
	case constants.StartupBackendRegistry:
		return strings.HasPrefix(e.Path, `HKCU\`)
	case constants.StartupBackendRegistryHKLM:
		return strings.HasPrefix(e.Path, `HKLM\`)
	case constants.StartupBackendStartupFolder:
		return strings.HasSuffix(strings.ToLower(e.Path), constants.StartupLnkExt)
	}

	return false
}

// matchesName returns true when the entry's logical name matches
// the requested filter. The on-disk Name carries OS-specific
// decoration (gitmap- / gitmap. prefix, optional .lnk suffix on
// Windows-folder backend) — strip those so the filter value is the
// same string the user passed to `startup-add --name`. Empty name
// = no filter.
func matchesName(e startup.Entry, name string) bool {
	if name == "" {
		return true
	}

	return logicalEntryName(e) == name
}

// logicalEntryName strips every prefix/suffix shape Add/List can
// produce so the returned string equals the --name value the user
// originally passed. Centralized here so a future Add format
// (e.g., a new OS) only adds one case to this function.
func logicalEntryName(e startup.Entry) string {
	n := e.Name
	n = strings.TrimSuffix(n, constants.StartupLnkExt)
	n = strings.TrimSuffix(n, constants.StartupDesktopExt)
	n = strings.TrimSuffix(n, constants.StartupPlistExt)
	n = strings.TrimPrefix(n, constants.StartupWinValuePrefix) // "gitmap-"
	n = strings.TrimPrefix(n, constants.StartupFilePrefix)     // "gitmap-"
	n = strings.TrimPrefix(n, constants.StartupPlistPrefix)    // "gitmap."

	return n
}
