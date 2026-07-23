// Package cmd — updatedebugwindows_source_test.go is a source-level
// invariant guard for `gitmap/cmd/updatedebugwindows.go`. It exists
// alongside the symbol-pin test (`updatedebugwindows_rename_test.go`)
// because the two failure modes are different:
//
//   - The rename-pin test catches a contributor who reverts the v3.113.0
//     fsutil migration by re-adding a LOCAL helper. It fails at compile
//     time because the new local helper would shadow the imported symbol
//     and the test references `fsutil.FileOrDirExists` directly.
//
//   - This file catches a contributor who adds a NEW local helper named
//     `fileExists` or `fileExistsLoose` to updatedebugwindows.go that
//     does NOT shadow anything (e.g. the cmd package no longer imports
//     fsutil for some unrelated reason). The compile-time guard would
//     not fire in that case, but the redeclaration footgun returns the
//     moment any sibling cmd file re-adds its own `fileExists`. By
//     parsing the file with go/parser we assert the invariant directly:
//     "this file declares zero functions named fileExists* — period."
//
// The two tests are intentionally cheap (parse one file, walk the AST
// once). Splitting them keeps each failure message specific to the
// regression it catches.
package cmd

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

const (
	// updatedebugwindowsPath is relative to this test file, which lives
	// in the same package directory as the source file under test. Using
	// a relative path avoids hardcoding the module root and keeps the
	// test working from any go-test invocation cwd.
	updatedebugwindowsPath = "updatedebugwindows.go"

	// fsutilImportPath is the canonical import path of the shared
	// existence-predicate package. Asserting the import is present
	// (rather than just the call site) catches a contributor who
	// removes the import and re-inlines the helper in the same change.
	fsutilImportPath = "github.com/alimtvnetwork/gitmap-v27/gitmap/fsutil"

	// expectedLooseCallSite is the FQN the file must contain at least
	// once. The string is asserted via the AST (selector expression),
	// not a substring scan, so comments mentioning it don't satisfy the
	// invariant — a real call site is required.
	expectedLooseCallSitePkg  = "fsutil"
	expectedLooseCallSiteName = "FileOrDirExists"
)

// TestUpdateDebugWindowsHasFsutilLooseCall asserts the file contains at
// least one real call to `fsutil.FileOrDirExists`. The historical name
// was `fileExistsLoose`; after the v3.113.0 migration the contract lives
// in fsutil and the call site is the only acceptable form. A call site
// (selector expression) is required — a comment mention does not count.
func TestUpdateDebugWindowsHasFsutilLooseCall(t *testing.T) {
	t.Parallel()

	file := parseUpdateDebugWindows(t)

	if !hasFsutilImport(file) {
		t.Fatalf("%s must import %q (the v3.113.0 migration target). Reverting to a local helper reintroduces the redeclaration footgun documented in updatedebugwindows_rename_test.go", updatedebugwindowsPath, fsutilImportPath)
	}

	if !hasSelectorCall(file, expectedLooseCallSitePkg, expectedLooseCallSiteName) {
		t.Fatalf("%s must contain at least one call to %s.%s — this is the post-v3.113.0 replacement for the local fileExistsLoose helper. If the call was removed, the debug-dump's path-exists branch is silently broken", updatedebugwindowsPath, expectedLooseCallSitePkg, expectedLooseCallSiteName)
	}
}

// TestUpdateDebugWindowsHasNoLocalFileExistsDecl asserts the file
// declares zero functions named `fileExists` or `fileExistsLoose`.
// Either name would re-trigger the redeclaration build break the moment
// a sibling cmd file (like updaterepo.go) declares the same name. The
// AST walk ignores comments and string literals, so the v3.113.0
// migration note that mentions `fileExistsLoose` in a doc comment does
// not false-positive.
func TestUpdateDebugWindowsHasNoLocalFileExistsDecl(t *testing.T) {
	t.Parallel()

	file := parseUpdateDebugWindows(t)

	forbidden := []string{"fileExists", "fileExistsLoose"}
	for _, name := range forbidden {
		if hasFuncDecl(file, name) {
			t.Fatalf("%s declares a local func %q — this is forbidden post-v3.113.0. Use fsutil.FileExists (strict) or fsutil.FileOrDirExists (loose) instead. Reintroducing the local helper recreates the redeclaration footgun the fsutil migration was designed to eliminate", updatedebugwindowsPath, name)
		}
	}
}

// parseUpdateDebugWindows parses the source file under test into an AST.
// Errors fail the test immediately because every assertion in this file
// depends on a successful parse — there is no useful partial result.
func parseUpdateDebugWindows(t *testing.T) *ast.File {
	t.Helper()

	fset := token.NewFileSet()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("runtime.Caller failed; cannot resolve absolute path to %s", updatedebugwindowsPath)
	}
	absPath := filepath.Join(filepath.Dir(thisFile), updatedebugwindowsPath)
	file, err := parser.ParseFile(fset, absPath, nil, parser.ParseComments)
	if err != nil {
		t.Fatalf("parse %s: %v", absPath, err)
	}

	return file
}

// hasFsutilImport reports whether the parsed file imports the canonical
// fsutil package path. Strips the surrounding quotes from the import
// literal before comparing because go/ast keeps them on the Path.Value.
func hasFsutilImport(file *ast.File) bool {
	for _, imp := range file.Imports {
		if strings.Trim(imp.Path.Value, "\"") == fsutilImportPath {
			return true
		}
	}

	return false
}

// hasSelectorCall reports whether the file contains at least one
// expression of the form `pkg.name(...)` — used to assert a real call
// site exists, not just a comment or import. Walks the full AST because
// the call may be nested arbitrarily deep inside other expressions.
func hasSelectorCall(file *ast.File, pkg, name string) bool {
	found := false
	ast.Inspect(file, func(n ast.Node) bool {
		if found {
			return false
		}

		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		ident, ok := sel.X.(*ast.Ident)
		if !ok {
			return true
		}

		if ident.Name == pkg && sel.Sel.Name == name {
			found = true

			return false
		}

		return true
	})

	return found
}

// hasFuncDecl reports whether the file declares a top-level function
// with the given name. Method declarations (with a receiver) are
// excluded because the redeclaration footgun only applies to
// package-level free functions sharing a namespace.
func hasFuncDecl(file *ast.File, name string) bool {
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		if fn.Recv != nil {
			continue
		}

		if fn.Name.Name == name {
			return true
		}
	}

	return false
}
