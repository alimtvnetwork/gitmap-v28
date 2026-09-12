// Package fsutil provides small, dependency-free filesystem predicates
// that historically existed as duplicated unexported helpers across
// many gitmap packages (cmd, release, lockfile, localdirs, vscodepm,
// detector). Centralizing them here removes the redeclaration footgun
// that bit the cmd package twice (see updatedebugwindows_rename_test.go
// for the v3.92.0 regression context) and gives every caller one
// well-documented contract per variant.
//
// Three variants are exported because callers genuinely need different
// semantics — collapsing them would silently change behavior in the
// repo-detection and debug-dump code paths:
//
//   - FileExists:      strict file (rejects directories, rejects empty)
//   - FileOrDirExists: loose existence check (accepts dirs, rejects empty)
//   - DirExists:       strict directory (rejects files, rejects empty)
//
// All three short-circuit on the empty string so callers don't need to
// guard their inputs — this matches the previous fileExistsLoose
// contract that the cmd package relied on for "path may be unset"
// debug-dump branches.
package fsutil

import "os"

// FileExists reports whether path resolves to a regular file (or any
// non-directory entry). Returns false for directories, missing paths,
// stat errors, and the empty string. This is the strict variant that
// updaterepo.go's repo-root detection depends on — relaxing it would
// cause directories named like sentinel files (e.g. a folder called
// "constants.go") to be treated as the file marker.
func FileExists(path string) bool {
	if len(path) == 0 {
		return false
	}

	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !info.IsDir()
}

// FileOrDirExists reports whether path resolves to ANY filesystem
// entry (file, directory, symlink target, etc.). Returns false on
// missing paths, stat errors, and the empty string. This is the loose
// variant the debug-dump code uses because the path it inspects may be
// unset, may point at a deploy directory, or may point at the new
// binary — all three cases are "exists" for diagnostic purposes.
func FileOrDirExists(path string) bool {
	if len(path) == 0 {
		return false
	}

	_, err := os.Stat(path)

	return err == nil
}

// DirExists reports whether path resolves to a directory. Returns
// false for files, missing paths, stat errors, and the empty string.
// Pairs with FileExists as the strict directory variant — callers that
// want "exists at all" should use FileOrDirExists instead.
func DirExists(path string) bool {
	if len(path) == 0 {
		return false
	}

	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}
