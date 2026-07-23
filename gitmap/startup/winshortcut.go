package startup

// .lnk Startup folder backend. Uses the in-process Shell Link
// writer in winshortcut_writer.go + winshortcut_linkinfo.go — no
// PowerShell shellout. Pure-Go, microsecond write, no powershell.exe
// dependency, unit-testable on Linux CI.
//
// Marker contract: the .lnk filename uses the `gitmap-` prefix
// AND a tracking subkey under HKCU\Software\Gitmap\StartupFolder\
// <name> records the entry as ours. Both must agree before Remove
// will delete the .lnk. Same belt-and-suspenders rule as the
// Registry backend.

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// addWindowsStartupFolder writes a .lnk via PowerShell, then
// records ownership in the tracking subkey. Same five-status
// outcome model as every other backend.
func addWindowsStartupFolder(clean string, opts AddOptions) (AddResult, error) {
	if runtime.GOOS != "windows" {

		return AddResult{}, fmt.Errorf(constants.ErrStartupUnsupportedOS)
	}
	dir, err := startupFolderDir()
	if err != nil {

		return AddResult{}, err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {

		return AddResult{}, fmt.Errorf("create startup folder %s: %w", dir, err)
	}
	full := filepath.Join(dir, constants.StartupWinValuePrefix+clean+constants.StartupLnkExt)

	return writeStartupShortcut(full, clean, opts)
}

// writeStartupShortcut performs the existence + ownership checks
// then either writes the .lnk (and tracking subkey) or refuses.
// Mirrors the structure of writeManaged on Linux so the five-status
// flow is visually parallel.
func writeStartupShortcut(full, clean string, opts AddOptions) (AddResult, error) {
	exists := fileExists(full)
	managed := exists && trackingSubkeyExists(constants.RegGitmapStartupFolder, clean)
	if exists && !managed {

		return AddResult{Status: AddRefused, Path: full}, nil
	}
	if exists && managed && !opts.Force {

		return AddResult{Status: AddExists, Path: full}, nil
	}
	if err := writeShortcutFile(full, opts.Exec); err != nil {

		return AddResult{}, fmt.Errorf(constants.ErrStartupShortcutCreate, full, err)
	}
	if err := writeTrackingSubkey(constants.RegGitmapStartupFolder, clean,
		opts.Exec, constants.StartupBackendStartupFolder, opts.WorkingDir); err != nil {

		return AddResult{}, err
	}
	if exists {

		return AddResult{Status: AddOverwritten, Path: full}, nil
	}

	return AddResult{Status: AddCreated, Path: full}, nil
}

// removeWindowsStartupFolder deletes a managed .lnk and its
// tracking subkey. Returns NoOp / Refused / Deleted with the same
// semantics as the Linux Remove path.
func removeWindowsStartupFolder(clean string, opts RemoveOptions) (RemoveResult, error) {
	if runtime.GOOS != "windows" {

		return RemoveResult{}, fmt.Errorf(constants.ErrStartupUnsupportedOS)
	}
	dir, err := startupFolderDir()
	if err != nil {

		return RemoveResult{}, err
	}
	full := filepath.Join(dir, constants.StartupWinValuePrefix+clean+constants.StartupLnkExt)
	if !fileExists(full) {

		return RemoveResult{Status: RemoveNoOp, DryRun: opts.DryRun}, nil
	}
	if !trackingSubkeyExists(constants.RegGitmapStartupFolder, clean) {

		return RemoveResult{Status: RemoveRefused, Path: full, DryRun: opts.DryRun}, nil
	}
	if opts.DryRun {

		return RemoveResult{Status: RemoveDeleted, Path: full, DryRun: true}, nil
	}
	if err := os.Remove(full); err != nil {

		return RemoveResult{}, fmt.Errorf("delete %s: %w", full, err)
	}
	if err := deleteTrackingSubkey(constants.RegGitmapStartupFolder, clean); err != nil {

		return RemoveResult{}, err
	}

	return RemoveResult{Status: RemoveDeleted, Path: full}, nil
}

// listWindowsStartupFolder enumerates managed .lnk files. Two-gate
// filter: filename prefix scan, then tracking-subkey re-check.
// Missing folder is "zero entries", not an error.
func listWindowsStartupFolder() ([]Entry, error) {
	if runtime.GOOS != "windows" {

		return nil, nil
	}
	dir, err := startupFolderDir()
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

	return collectStartupFolderManaged(dir, files), nil
}

// collectStartupFolderManaged is the per-file filter. Same shape
// as collectManagedDesktop / collectManagedPlist; differs only in
// the ownership re-check (subkey-based instead of file-marker).
func collectStartupFolderManaged(dir string, files []os.DirEntry) []Entry {
	var out []Entry
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		name := f.Name()
		if !looksLikeOursLnk(name) {
			continue
		}
		clean := strings.TrimSuffix(strings.TrimPrefix(name,
			constants.StartupWinValuePrefix), constants.StartupLnkExt)
		if !trackingSubkeyExists(constants.RegGitmapStartupFolder, clean) {
			continue
		}
		out = append(out, Entry{
			Name: strings.TrimSuffix(name, constants.StartupLnkExt),
			Path: filepath.Join(dir, name),
			Exec: "", // .lnk Exec lives inside the binary; not parsed here
		})
	}

	return out
}

// looksLikeOursLnk is the cheap pre-filter: filename ends in .lnk
// AND starts with the gitmap- prefix.
func looksLikeOursLnk(filename string) bool {
	if !strings.HasSuffix(filename, constants.StartupLnkExt) {
		return false
	}

	return strings.HasPrefix(filename, constants.StartupWinValuePrefix)
}

// startupFolderDir + fileExists helpers live in winshortcut_helpers.go.
// Writer (writeShortcutFile, buildShortcutBytes) lives in
// winshortcut_writer.go + winshortcut_linkinfo.go. Legacy
// PowerShell helper lives in winshortcut_ps.go — retained as a
// future fallback for non-trivial target shapes.
