package model

import "path/filepath"

// CleanRelativePath returns the OS-native, defensively-normalized form
// of a manifest-supplied RelativePath. Used at every call site that
// joins a parsed JSON/CSV manifest row's RelativePath onto a target
// directory so the resulting absolute path has a single canonical
// shape regardless of how the row was authored.
//
// Why this exists (the failure mode it prevents):
//
//	A manifest authored on macOS / Linux typically writes RelativePath
//	as `"acme/widget"`. When the same manifest is consumed on Windows
//	and joined onto `C:\code\` via `filepath.Join`, the result is
//	`C:\code\acme\widget` — but ONLY because Join happens to call
//	Clean internally. Sloppier inputs slip through:
//
//	  "acme//widget"   → Join leaves the doubled slash on Windows
//	                     when the second segment is treated as a
//	                     single token (rare, but observed with
//	                     hand-edited manifests).
//	  "./acme/widget"  → leading `.` survives into AbsolutePath that
//	                     then differs byte-for-byte from a re-scan
//	                     result, breaking dedup keys in projects.json
//	                     and the SQLite RepoTracking table.
//	  "acme/widget/"   → trailing separator survives, producing two
//	                     distinct strings ("…/widget" and "…/widget/")
//	                     for the same physical folder.
//
// Normalization steps (intentionally minimal — each one earns its
// keep against a real failure mode):
//
//  1. filepath.FromSlash — converts forward slashes to the OS-native
//     separator. No-op on Unix, the Windows-correctness step that
//     stops mixed-separator paths from leaking into AbsolutePath.
//  2. filepath.Clean — collapses `.`, `..`, doubled separators, and
//     trailing separators. Same canonicalization rule used by
//     canonicalizePMPath in gitmap/cmd/clonepmsync.go and by
//     vscodepm.normalizePath, so all three surfaces agree on "what
//     counts as the same relative path".
//
// Empty input is preserved (returned as ""), NOT normalized to ".".
// Callers treat empty RelativePath as "no manifest row" — see
// reclone_summary.go and reclone_confirm.go which both early-return
// on the empty case. Returning "." here would silently turn that
// signal into a clone-into-cwd, which is exactly the kind of invisible
// behavior change this helper exists to prevent.
//
// This helper is the SINGLE source of truth for manifest-side
// RelativePath normalization. Any new call site that joins a
// manifest-supplied RelativePath onto a target dir MUST route the
// RelativePath through CleanRelativePath first.
func CleanRelativePath(rel string) string {
	if rel == "" {
		return ""
	}

	return filepath.Clean(filepath.FromSlash(rel))
}
