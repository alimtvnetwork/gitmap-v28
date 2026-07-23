// Package cmd — releaserebase.go implements the cross-dir `gitmap r <repo> <version>`
// form: pull --rebase the named repo, then run the standard release pipeline,
// then chdir back to the original directory.
//
// Backward compatibility: `gitmap r vX.Y.Z` (single positional arg) keeps
// running an in-place release of the current repo. The new behavior only
// triggers when TWO positional args are given AND the first does NOT look
// like a version string (e.g. v3.31.0, 3.31.0).
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// versionPattern matches the three version-shaped tokens any cross-dir
// dispatcher in this package needs to recognize:
//
//  1. Full SemVer with optional pre-release: "1.2.3", "v1.2.3", "v1.2.3-rc1"
//  2. clone-next bump shortcut:              "v++"
//  3. clone-next explicit-bump shortcut:     "v+N" (N >= 1)
//
// The bump shortcuts (#2 and #3) are clone-next-specific but live here
// because `looksLikeVersion` is shared with the release dispatcher
// (releaserebase.go) and clonenextcrossdir.go. Keeping one regex
// avoids the redeclaration footgun we hardened against in v3.113.0.
//
// Letters after `+` (e.g. "v+abc") deliberately do NOT match — those
// fall through to the alias path.
var versionPattern = regexp.MustCompile(`^(v?\d+\.\d+\.\d+([.\-+].+)?|v\+\+|v\+\d+)$`)

// looksLikeVersion returns true if s is a version-shaped token, false if it's
// more likely a repo alias (e.g. "gitmap", "my-app") or a folder path.
func looksLikeVersion(s string) bool {
	return versionPattern.MatchString(s)
}

// tryCrossDirRelease intercepts `r <repo> <version>` and delegates to the
// existing release-alias machinery with rebase pull mode. Returns true when
// it handled the invocation, false to let the normal in-place flow continue.
func tryCrossDirRelease(args []string) bool {
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
	performCrossDirRelease(target, alias, version, args)

	return true
}

// extractPositionalArgs returns args that are not flags (don't start with -).
func extractPositionalArgs(args []string) []string {
	out := make([]string, 0, len(args))
	skipNext := false
	for _, a := range args {
		if skipNext {
			skipNext = false
			continue
		}
		if len(a) > 0 && a[0] == '-' {
			continue
		}
		out = append(out, a)
	}

	return out
}

// performCrossDirRelease runs the chdir + rebase-pull + stash + release sequence.
func performCrossDirRelease(target, alias, version string, originalArgs []string) {
	originalDir, _ := os.Getwd()
	if err := os.Chdir(target); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrRAChdirFailedFmt, target, err)
		fmt.Fprintln(os.Stderr)
		os.Exit(1)
	}
	defer func() { _ = os.Chdir(originalDir) }()

	fmt.Printf(constants.MsgRRStartingFmt, alias, target, version)

	rebasePull(target)

	stashLabel := autoStashIfDirty(target, alias, version)
	if stashLabel != "" {
		defer popAutoStash(target, stashLabel)
	}

	// Forward any flags the user passed (e.g. -y, --dry-run) but replace the
	// positional args with just the version.
	releaseArgs := []string{version}
	releaseArgs = append(releaseArgs, extractFlagArgs(originalArgs)...)
	executeReleaseFromCrossDir(releaseArgs)

	fmt.Printf(constants.MsgRRReturnedFmt, originalDir)
}

// executeReleaseFromCrossDir runs the standard runRelease but bypasses the
// cross-dir interceptor (we already handled the chdir).
func executeReleaseFromCrossDir(releaseArgs []string) {
	// Recursive call is safe because releaseArgs has exactly 1 positional
	// (the version), which fails the 2-positional check in tryCrossDirRelease.
	runRelease(releaseArgs)
}

// extractFlagArgs returns only the flag-shaped tokens (-x, --foo, --foo=bar).
func extractFlagArgs(args []string) []string {
	out := make([]string, 0, len(args))
	for _, a := range args {
		if len(a) > 0 && a[0] == '-' {
			out = append(out, a)
		}
	}

	return out
}

// rebasePull runs `git fetch && git pull --rebase` in target.
// Aborts the release on failure — the user should resolve the conflict
// before retrying. This is the documented behavior contract for `r`.
func rebasePull(target string) {
	fmt.Printf(constants.MsgRRFetchingFmt, target)
	fetch := exec.Command(constants.GitBin, "fetch", "--all", "--prune")
	fetch.Dir = target
	fetch.Stdout = os.Stdout
	fetch.Stderr = os.Stderr
	if err := fetch.Run(); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrRRFetchFailedFmt, target, err)
		fmt.Fprintln(os.Stderr)
		os.Exit(1)
	}

	fmt.Printf(constants.MsgRRRebasingFmt, target)
	pull := exec.Command(constants.GitBin, "pull", "--rebase")
	pull.Dir = target
	pull.Stdout = os.Stdout
	pull.Stderr = os.Stderr
	if err := pull.Run(); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrRRRebaseFailedFmt, target, err)
		fmt.Fprintln(os.Stderr)
		os.Exit(1)
	}
}
