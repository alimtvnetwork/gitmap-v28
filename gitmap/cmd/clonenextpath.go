package cmd

import (
	"path/filepath"
	"strings"
)

// samePath returns true when a and b refer to the same on-disk location.
//
// We deliberately avoid filepath.EvalSymlinks here — the caller is the
// `cn -f` pre-flatten check, where both inputs come from os.Getwd() and
// filepath.Join(parent, base). They're already absolute and clean, and
// resolving symlinks would just add a stat call that can fail mid-clone.
//
// On Windows, path comparison is case-insensitive (NTFS treats "C:\Foo"
// and "C:\foo" as the same dir). filepath.Clean is applied first so a
// trailing slash or "." segment doesn't produce a false negative. We
// don't gate on runtime.GOOS because the strings.EqualFold fallback is
// safe everywhere — Linux paths that differ only in case are pathological
// edge cases not worth optimizing against.
func samePath(a, b string) bool {
	aClean := filepath.Clean(a)
	bClean := filepath.Clean(b)
	if aClean == bClean {
		return true
	}

	return strings.EqualFold(aClean, bClean)
}
