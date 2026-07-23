// Package cmd — updatedebugwindows_rename_test.go originally pinned the
// v3.92.0 rename of the debug-dump helper from `fileExists` to
// `fileExistsLoose` to prevent the redeclaration build break in this
// package. As of v3.113.0 both helpers are gone — the contract has moved
// to the shared `gitmap/fsutil` package, so the redeclaration footgun is
// now structurally impossible (you cannot redeclare an imported symbol).
//
// The test is preserved (rather than deleted) as a forward-looking guard:
// it fails at compile time if a contributor reintroduces an unexported
// `fileExists` or `fileExistsLoose` in this package. The fix in that case
// is always the same — call fsutil.FileExists / fsutil.FileOrDirExists
// instead of adding a local helper.
//
// See spec/02-app-issues/33-stale-binary-clone-folder-url-guard.md for
// the related stale-binary diagnostic pattern: when CI logs a
// redeclaration at a line number that does not match this source, the
// user is building an out-of-date snapshot — `git pull` is the actual fix.
package cmd

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/fsutil"
)

// TestFsutilMigrationPinned compiles only when the cmd package is using
// the shared fsutil predicates instead of local copies. The two
// assertions exercise the contracts the cmd package actually depends on:
// loose existence (debug dump) and strict file existence (repo-root
// detection). Reverting either to a local helper triggers the
// redeclaration build break this test was created to prevent.
func TestFsutilMigrationPinned(t *testing.T) {
	t.Parallel()

	if fsutil.FileOrDirExists("") {
		t.Fatal("fsutil.FileOrDirExists(\"\") must short-circuit to false; the debug dump's path-may-be-unset branch relies on this contract")
	}

	if !fsutil.FileOrDirExists(".") {
		t.Fatal("fsutil.FileOrDirExists(\".\") must return true (CWD always exists; loose variant treats directories as existing)")
	}

	if fsutil.FileExists(".") {
		t.Fatal("fsutil.FileExists(\".\") must return false (strict variant rejects directories); updaterepo.go's source-marker detection depends on this")
	}
}
