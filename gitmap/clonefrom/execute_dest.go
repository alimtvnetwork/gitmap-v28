package clonefrom

// Destination-side helpers used by the executor. Split out of
// execute.go to keep that file under the project's 200-line cap
// after adding parent-dir creation. All three functions are pure
// FS / path utilities with NO knowledge of git, results, or rows
// beyond what executeRow hands them.

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// resolveDest computes (rowDest, absoluteDest). rowDest is what
// goes into the Result + report; absDest drives the skip-check
// and parent-mkdir against the actual filesystem.
func resolveDest(r Row, cwd string) (string, string) {
	dest := r.Dest
	if len(dest) == 0 {
		dest = DeriveDest(r.URL)
	}
	absDest := dest
	if !filepath.IsAbs(absDest) {
		absDest = filepath.Join(cwd, dest)
	}

	return dest, absDest
}

// prepareDestParent ensures the parent dir of the resolved dest
// exists so nested dest paths (e.g. `org-a/repo-1`) work even on
// a fresh checkout where `org-a/` doesn't yet exist. MkdirAll is
// idempotent: pre-existing parent → no-op, no race against other
// concurrent clones into siblings. On failure we log to stderr in
// the project's Code Red format AND return a row Detail so the
// per-row line + CSV report carry the same diagnosis.
func prepareDestParent(absDest string) (string, bool) {
	parent := filepath.Dir(absDest)
	if err := os.MkdirAll(parent, constants.DirPermission); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrCloneFromMkdirParent, parent, err)
		return fmt.Sprintf(constants.MsgCloneFromMkdirParentFailFmt, err), false
	}

	return "", true
}

// shouldSkip returns true when the dest is a non-empty directory.
// Errors reading the dir (permission denied) → false (let git try
// and fail with a clearer message than we could craft).
func shouldSkip(absDest string) bool {
	info, err := os.Stat(absDest)
	if err != nil || !info.IsDir() {
		return false
	}
	entries, err := os.ReadDir(absDest)
	if err != nil {
		return false
	}

	return len(entries) > 0
}
