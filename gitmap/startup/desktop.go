package startup

// Internal helpers for parsing and filtering autostart files. Split
// from startup.go to keep both files under the per-file budget and so
// the parser is independently testable without going through the
// filesystem-coupled List API.
//
// This file owns the Linux/Unix `.desktop` path. The macOS `.plist`
// path lives in plist.go. collectManaged is the OS-aware dispatcher
// — it stays here (not in startup.go) so the per-OS reader functions
// it calls are private to the package.

import (
	"bufio"
	"os"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// collectManaged scans `files` (a single ReadDir result) and returns
// only the ones that BOTH match the gitmap filename prefix AND carry
// the gitmap marker inside the file. The two-gate check is
// deliberate: filename alone is spoofable; marker alone would force
// us to read every file in the directory (slow on systems with many
// startup entries).
//
// Dispatches to the per-OS reader. Anything other than darwin uses
// the .desktop reader; darwin uses the .plist reader. Windows never
// reaches here because AutostartDir errors first.
func collectManaged(dir string, files []os.DirEntry) []Entry {
	if runtime.GOOS == "darwin" {

		return collectManagedPlist(dir, files)
	}

	return collectManagedDesktop(dir, files)
}

// collectManagedDesktop is the Linux/Unix reader.
func collectManagedDesktop(dir string, files []os.DirEntry) []Entry {
	var out []Entry
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		name := f.Name()
		if !looksLikeOursDesktop(name) {
			continue
		}
		entry, ok := readManagedDesktop(dir, name)
		if !ok {
			continue
		}
		out = append(out, entry)
	}

	return out
}

// looksLikeOursDesktop is the cheap pre-filter for Linux: filename
// must end in `.desktop` AND start with the gitmap- prefix. Files
// that fail either check are skipped without being opened.
func looksLikeOursDesktop(filename string) bool {
	if !strings.HasSuffix(filename, constants.StartupDesktopExt) {
		return false
	}

	return strings.HasPrefix(filename, constants.StartupFilePrefix)
}

// readManagedDesktop opens the file and parses just enough to decide
// whether it's gitmap-managed AND to surface its Exec line. Returns
// ok=false on any I/O error or when the marker is absent — both
// outcomes mean "not ours, skip it" from the caller's perspective.
func readManagedDesktop(dir, filename string) (Entry, bool) {
	full := joinPath(dir, filename)
	f, err := os.Open(full)
	if err != nil {
		return Entry{}, false
	}
	defer f.Close()

	managed, exec := parseDesktopFields(bufio.NewScanner(f))
	if !managed {
		return Entry{}, false
	}

	return Entry{
		Name: strings.TrimSuffix(filename, constants.StartupDesktopExt),
		Path: full,
		Exec: exec,
	}, true
}

// parseDesktopFields walks the .desktop file line-by-line looking for
// the marker key and the Exec= line. Returns (managed, exec). We do
// NOT use a full INI parser because .desktop files are UTF-8 with a
// strict key=value-per-line grammar that bufio.Scanner handles
// correctly, and the dependency budget for a single CLI subcommand
// shouldn't include a parser library.
func parseDesktopFields(sc *bufio.Scanner) (bool, string) {
	managed := false
	exec := ""
	wantedMarker := constants.StartupMarkerKey + "=" + constants.StartupMarkerVal
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == wantedMarker {
			managed = true
			continue
		}
		if strings.HasPrefix(line, "Exec=") {
			exec = strings.TrimPrefix(line, "Exec=")
		}
	}

	return managed, exec
}

// joinPath is a tiny wrapper kept here (rather than importing
// filepath.Join into every helper) so the hot paths in this file
// stay free of stdlib boilerplate. Linux-only context: forward
// slashes are always correct.
func joinPath(dir, name string) string {
	if strings.HasSuffix(dir, "/") {
		return dir + name
	}

	return dir + "/" + name
}
