package startup

// Add implements `gitmap startup-add` with a strict "managed-only,
// never escalate" contract that mirrors Remove. The flow validates
// the name, ensures the autostart dir exists, resolves the
// platform-specific filename, refuses to overwrite third-party
// files (even with --force), and writes the rendered body
// atomically. See per-OS render/dispatch helpers for format
// details (renderDesktop, renderPlist, addWindows*).

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// AddOptions captures every knob the runner exposes. Kept as a
// struct (not positional args) so future fields like Comment,
// Categories, OnlyShowIn don't keep growing the signature.
type AddOptions struct {
	// Name is the entry's logical identifier. The on-disk filename
	// becomes "gitmap-<Name>.desktop" (prefix added if missing) and
	// the listed Name= field in the .desktop body is the same value.
	Name string
	// Exec is the absolute command line to launch at login. Callers
	// MUST pre-quote any path containing spaces (the .desktop spec
	// allows quoted strings inside Exec=); we do not parse it.
	Exec string
	// DisplayName, when non-empty, overrides the value written to
	// the Name= field. Empty means "reuse Name verbatim".
	DisplayName string
	// Comment populates the Comment= field. Optional.
	Comment string
	// NoDisplay sets NoDisplay=true to hide the entry from desktop
	// app menus while still autostarting it — useful for headless
	// helpers users don't want cluttering their app drawer.
	NoDisplay bool
	// Force allows overwriting a previously gitmap-created entry
	// with the same name. Has NO effect on third-party files; those
	// always refuse.
	Force bool
	// WorkingDir, when non-empty, sets the process working directory
	// for the autostart entry. Rendered as `Path=<dir>` in .desktop
	// files (XDG spec field), `WorkingDirectory` in LaunchAgent
	// plists, and stored as the `WorkingDir` value of the gitmap
	// tracking subkey on Windows (HKCU\Software\Gitmap\Startup*\<name>).
	// Callers MUST pass an absolute path; relative paths are accepted
	// as-is and interpreted by the OS at login time. Empty means
	// "inherit whatever the login session provides".
	WorkingDir string
	// Backend selects the Windows autostart target (Registry vs
	// Startup-folder shortcut). Ignored on Linux/macOS, which each
	// have one canonical backend. Zero value (BackendUnspecified)
	// means "use the OS default" which is BackendRegistry on
	// Windows.
	Backend Backend
}

// AddStatus tags the four mutually-exclusive Add outcomes. Kept
// parallel to RemoveStatus so the CLI rendering layer can switch on
// either with the same shape.
type AddStatus int

const (
	// AddCreated = file did not exist; was written fresh.
	AddCreated AddStatus = iota
	// AddOverwritten = previously gitmap-managed file was replaced
	// because Force was set.
	AddOverwritten
	// AddRefused = a non-gitmap-managed file with the same name
	// already exists; we did NOT touch it.
	AddRefused
	// AddBadName = name failed validation (empty / separator / NUL).
	AddBadName
	// AddExists = a gitmap-managed file with the same name already
	// exists and Force was NOT set; nothing was written.
	AddExists
)

// AddResult mirrors RemoveResult. Path is the absolute target file
// for Created / Overwritten / Refused / Exists; empty for BadName.
type AddResult struct {
	Status AddStatus
	Path   string
}

// Add is the public entry point. Returns (result, nil) for every
// "soft" outcome; only real I/O failures produce a non-nil error.
//
// OS dispatch:
//
//   - linux/unix → writes a `.desktop` file with the
//     X-Gitmap-Managed=true marker into AutostartDir().
//   - darwin     → writes a LaunchAgent `.plist` with the
//     XGitmapManaged <true/> marker into ~/Library/LaunchAgents/.
//   - windows    → routes via opts.Backend (Registry by default,
//     or Startup-folder .lnk shortcut). Both backends share the
//     same managed-marker contract enforced by HKCU\Software\Gitmap
//     tracking subkeys + a sibling marker value next to Run-key
//     entries. See addWindows / winbackend.go for details.
//
// Both OS paths share the same five-status outcome model
// (Created/Overwritten/Refused/BadName/Exists) and the same
// "managed-only, never escalate" guard — Force only lifts the
// "already exists AND is ours" check; a third-party file is NEVER
// overwritten.
func Add(opts AddOptions) (AddResult, error) {
	clean := normalizeName(opts.Name)
	if !isValidName(clean) {
		return AddResult{Status: AddBadName}, nil
	}
	if runtime.GOOS == "windows" {

		return addWindows(clean, opts)
	}
	dir, err := AutostartDir()
	if err != nil {
		return AddResult{}, err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return AddResult{}, fmt.Errorf("create autostart dir %s: %w", dir, err)
	}
	full := joinPath(dir, platformFilename(clean))

	return writeManaged(full, clean, opts)
}

// addWindows lives in winbackend.go (kept with the other Windows
// dispatch helpers so the public Add API has only the OS routing
// switch in this file).

// platformFilename picks the OS-specific filename shape. macOS uses
// the reverse-DNS `gitmap.<name>.plist` convention; everything else
// uses the XDG `gitmap-<name>.desktop` convention. Centralized here
// so add.go's dispatch stays small and Remove's platformExt() and
// this function can never disagree on which extension Add wrote.
func platformFilename(clean string) string {
	if runtime.GOOS == "darwin" {

		return prefixedFilenamePlist(clean)
	}

	return prefixedFilename(clean)
}

// prefixedFilename ensures the on-disk name starts with the gitmap-
// prefix exactly once. Callers passing "gitmap-foo" or "foo" both
// land on "gitmap-foo.desktop". Linux/Unix only — the macOS analog
// is prefixedFilenamePlist (different prefix shape: `gitmap.` not
// `gitmap-`, per LaunchAgent reverse-DNS labeling convention).
func prefixedFilename(clean string) string {
	if strings.HasPrefix(clean, constants.StartupFilePrefix) {
		return clean + constants.StartupDesktopExt
	}

	return constants.StartupFilePrefix + clean + constants.StartupDesktopExt
}

// writeManaged performs the existence + ownership checks then writes
// (or refuses). Splits the four Add outcomes into a flat sequence so
// each branch is its own clear `return`.
func writeManaged(full, clean string, opts AddOptions) (AddResult, error) {
	exists, managed := classifyTarget(full)
	if exists && !managed {
		return AddResult{Status: AddRefused, Path: full}, nil
	}
	if exists && managed && !opts.Force {
		return AddResult{Status: AddExists, Path: full}, nil
	}
	body := renderForOS(clean, opts)
	if err := atomicWrite(full, body); err != nil {
		return AddResult{}, err
	}
	if exists {
		return AddResult{Status: AddOverwritten, Path: full}, nil
	}

	return AddResult{Status: AddCreated, Path: full}, nil
}

// renderForOS picks the per-OS body renderer. Same dispatch shape
// as platformFilename so the format and the on-disk extension can
// never drift apart (e.g., a future refactor that writes a .plist
// body into a .desktop file would have to change both functions).
func renderForOS(clean string, opts AddOptions) []byte {
	if runtime.GOOS == "darwin" {

		return renderPlist(clean, opts)
	}

	return renderDesktop(clean, opts)
}

// classifyTarget returns (exists, managed). A read error is treated
// as "exists, not managed" so we err on the side of refusing to
// overwrite — preferable to deleting a file we can't read.
func classifyTarget(full string) (bool, bool) {
	if _, err := os.Stat(full); os.IsNotExist(err) {
		return false, false
	}

	return true, isManagedFile(full)
}
