// Package cmd — clonenextfolderdispatch.go implements the v3.117.0
// folder-arg forms of `gitmap cn`:
//
//	gitmap cn vX <folder>     — explicit version, explicit folder
//	gitmap cn v+1 <folder>    — version-bump shortcut
//	gitmap cn <folder>        — folder only (defaults to v++)
//
// All three chdir into the resolved folder, run the existing in-place
// `runCloneNext` pipeline, then chdir back via the shared
// `performCrossDirCloneNext` helper in clonenextcrossdir.go.
//
// Disambiguation rules and the full test matrix are documented in
// spec/01-app/111-cn-folder-arg.md. The two interceptor functions
// here run BEFORE tryCrossDirCloneNext in runCloneNext so the
// path-shaped tokens win over the release-alias fallback (the alias
// resolver matches bare names like "gitmap" which would otherwise
// shadow a same-named local folder).
package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/fsutil"
)

// errCNFolderNotDir is returned when the resolved path exists but is
// not a directory. Sentinel so the dispatcher can distinguish
// "user typo" (not found) from "wrong target" (file, not dir) in the
// error message without parsing the underlying os error.
var errCNFolderNotDir = errors.New("resolved path is not a directory")

// tryFolderArgCloneNext intercepts the three new v3.117.0 forms and
// dispatches to the cross-dir helper. Returns true when the call was
// handled (success OR clean exit), false to fall through to the
// existing tryCrossDirCloneNext / in-place flow.
//
// Order of the four classification branches matters — see the
// truth table in spec/01-app/111-cn-folder-arg.md §Disambiguation.
func tryFolderArgCloneNext(args []string) bool {
	positional := extractPositionalArgs(args)

	switch len(positional) {
	case 1:
		return tryFolderArgSinglePositional(positional[0], args)
	case 2:
		return tryFolderArgTwoPositional(positional[0], positional[1], args)
	}

	return false
}

// tryFolderArgSinglePositional handles `gitmap cn <folder>`. Returns
// true only when the token is folder-shaped AND resolves to an
// existing directory — bare alias names (no path-sep, no on-disk
// match) fall through so the in-place / alias forms keep working.
func tryFolderArgSinglePositional(token string, originalArgs []string) bool {
	if looksLikeVersion(token) {
		return false
	}
	isNonFolderShaped := !isFolderShaped(token)
	if isNonFolderShaped {
		return false
	}

	resolved, err := resolveCloneNextFolder(token)
	if err != nil && hasFolderHint(token) == true {
		return false
	}
	if err != nil {
		return false
	}

	performCrossDirCloneNext(resolved, filepath.Base(resolved), constants.CloneNextDefaultVersionArg, originalArgs)

	return true
}

// tryFolderArgTwoPositional handles `gitmap cn vX <folder>`. The
// reverse order (`cn <folder> vX`) is intentionally NOT handled
// here — that's the existing release-alias cross-dir form
// (tryCrossDirCloneNext). The two paths are mutually exclusive
// because of the version-position swap.
func tryFolderArgTwoPositional(first, second string, originalArgs []string) bool {
	firstIsVersion := looksLikeVersion(first)
	secondIsVersion := looksLikeVersion(second)

	if firstIsVersion && secondIsVersion {
		return false
	}

	if !firstIsVersion {
		// Either folder+version (existing alias form) or
		// alias+alias / folder+folder (existing form's problem).
		// Fall through.
		return false
	}

	// firstIsVersion && !secondIsVersion → NEW form.
	isNonFolderShaped := !isFolderShaped(second)
	if isNonFolderShaped {
		// "cn v+1 gitmap" with no folder hint and no on-disk match
		// is almost certainly a typo — the user meant a folder but
		// gave a bare alias. Refuse with the canonical message
		// rather than silently swallowing it.
		return false
	}

	resolved, err := resolveCloneNextFolder(second)
	if err != nil {
		return false
	}

	performCrossDirCloneNext(resolved, filepath.Base(resolved), first, originalArgs)

	return true
}

// resolveCloneNextFolder expands ~, joins to cwd if relative, then
// stats. Returns the absolute path on success; errCNFolderNotDir if
// the path exists but is a file; the underlying stat error otherwise.
func resolveCloneNextFolder(token string) (string, error) {
	expanded := expandTilde(token)
	expanded = resolveEndpointString(expanded)

	expanded, errExpand := ensureAbsolutePath(expanded)
	if errExpand != nil {
		return "", errExpand
	}

	abs, err := filepath.Abs(expanded)
	if err != nil {
		return "", err
	}

	if fsutil.DirExists(abs) == false && fsutil.FileExists(abs) == true {
		return "", errCNFolderNotDir
	}
	if fsutil.DirExists(abs) == false {
		return "", os.ErrNotExist
	}

	return abs, nil
}

// isFolderShaped reports whether token has the syntactic shape of a
// path (separator or tilde) OR resolves to an existing directory.
// The two-pronged check lets bare names like "macro-ahk-v11" still
// be recognized when the user is in the parent folder and the dir
// exists, without false-positively claiming any random alias name.
func isFolderShaped(token string) bool {
	if hasFolderHint(token) {
		return true
	}

	resolved, err := resolveCloneNextFolder(token)

	return err == nil && len(resolved) > 0
}

// hasFolderHint returns true for tokens with an unambiguous path
// signal — a separator or a leading tilde. Used to escalate
// "folder not found" to a hard error instead of silently falling
// through to the alias resolver.
func hasFolderHint(token string) bool {
	if strings.HasPrefix(token, "~") {
		return true
	}
	if strings.ContainsAny(token, `/\`) {
		return true
	}

	return false
}

func ensureAbsolutePath(expanded string) (string, error) {
	if filepath.IsAbs(expanded) == true {
		return expanded, nil
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(cwd, expanded), nil
}
