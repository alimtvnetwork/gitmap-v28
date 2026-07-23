// Package cmd — clonenextcrossdir.go implements the cross-dir
// `gitmap cn <repo> <version>` form: chdir into the named repo, run the
// existing clone-next pipeline, then chdir back.
//
// Backward compatibility: `gitmap cn vX.Y.Z` (single positional) keeps
// operating on the current directory's repo as before.
package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// tryCrossDirCloneNext intercepts `cn <repo> <version>` and chdirs first.
// Returns true when it handled the invocation, false to let the normal
// in-place flow continue.
func tryCrossDirCloneNext(args []string) bool {
	positional := extractPositionalArgs(args)
	if len(positional) != 2 {
		return false
	}
	if looksLikeVersion(positional[0]) {
		return false
	}
	if !looksLikeVersion(positional[1]) {
		return false
	}

	alias := positional[0]
	version := positional[1]
	target := resolveReleaseAliasPath(alias)
	performCrossDirCloneNext(target, alias, version, args)

	return true
}

// performCrossDirCloneNext does the chdir + clone-next + chdir-back dance.
func performCrossDirCloneNext(target, alias, version string, originalArgs []string) {
	originalDir, _ := os.Getwd()
	if err := os.Chdir(target); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrRAChdirFailedFmt, target, err)
		fmt.Fprintln(os.Stderr)
		os.Exit(1)
	}
	defer func() { _ = os.Chdir(originalDir) }()

	fmt.Printf(constants.MsgCNXStartingFmt, alias, target, version)

	// Build the in-place clone-next args: just the version + any forwarded flags.
	cnArgs := []string{version}
	cnArgs = append(cnArgs, extractFlagArgs(originalArgs)...)
	runCloneNext(cnArgs)

	fmt.Printf(constants.MsgCNXReturnedFmt, originalDir)
}
